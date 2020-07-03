package mysqlwrap

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/RocksonZeta/wrap/errs"
	"github.com/RocksonZeta/wrap/utils/sqlutil"
	"github.com/RocksonZeta/wrap/utils/sutil"
	"github.com/RocksonZeta/wrap/wraplog"
	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/guregu/null.v3"
)

var pkg = reflect.TypeOf(Mysql{}).PkgPath()
var log = wraplog.Logger.Fork(pkg, "Mysql")

const (
	ErrorInit = 1 + iota
	ErrorQuery
	ErrorInsert
	ErrorUpdate
	ErrorDelete
	ErrorClose
	ErrorTxBegin
	ErrorTxCommit
	ErrorTxRollback
)

func check(err error, state int, msg string) {
	if err == nil {
		return
	}
	if msg != "" {
		msg = err.Error()
	}
	fmt.Println("err type:", reflect.TypeOf(err))
	panic(errs.Err{Err: err, Module: "Mysql", Pkg: pkg, State: state, Message: msg})
}

type Options struct {
	Url         string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime int
}

// var mysqls sync.Map

//NewFromUrl mysqlUrl:root:6plzHiJKdUMlFZ@tcp(test.iqidao.com:43122)/good?charset=utf8mb4&MaxOpen=2000&MaxIdle=10&MaxLifetime=60
func NewFromUrl(mysqlUrl string) *Mysql {
	log.Trace().Func("NewFromUrl").Interface("mysqlUrl", mysqlUrl).Send()
	parts, err := url.Parse(mysqlUrl)
	if err != nil {
		check(err, ErrorInit, err.Error())
	}
	q := parts.Query()
	var options Options
	options.MaxOpen, _ = strconv.Atoi(q.Get("MaxOpen"))
	options.MaxIdle, _ = strconv.Atoi(q.Get("MaxIdle"))
	options.MaxLifetime, _ = strconv.Atoi(q.Get("MaxLifetime"))
	query := make(url.Values)
	for k := range q {
		if k[0] >= 'a' && k[0] <= 'z' {
			query.Add(k, q.Get(k))
		}
	}
	parts.RawQuery = query.Encode()
	options.Url = parts.String()
	return New(options)
}

func New(options Options) *Mysql {
	log.Trace().Func("New").Interface("options", options).Send()
	var err error
	var db *sql.DB
	// old, ok := mysqls.Load(options.Url)
	// if !ok {
	if options.MaxOpen <= 0 {
		options.MaxOpen = 2000
	}
	if options.MaxIdle <= 0 {
		options.MaxIdle = 10
	}
	db, err = sql.Open("mysql", options.Url)
	db.SetMaxOpenConns(options.MaxOpen)
	db.SetMaxIdleConns(options.MaxIdle)
	if options.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(options.MaxLifetime) * time.Second)
	}
	// mysqls.Store(options.Url, db)
	// } else {
	// 	db = old.(*sql.DB)
	// }
	if err != nil {
		log.Error().Func("New").Stack().Err(err).Interface("options", options).Send()
		check(err, ErrorInit, err.Error())
	}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	if log.DebugEnabled() {
		dbMap.TraceOn("[mysql]", new(mysqlLogger))
	}
	return &Mysql{DbMap: dbMap, MysqlExecutor: WrapExecutor(dbMap), options: options}
}

type mysqlLogger struct {
}

func (l mysqlLogger) Printf(format string, v ...interface{}) {
	log.Debug().Str("sql", fmt.Sprintf(format, v...)).Send()
}

func WrapExecutor(executor gorp.SqlExecutor) *MysqlExecutor {
	return &MysqlExecutor{SqlExecutor: executor}
}

func fatalError(err error) bool {
	if err == nil {
		return false
	}
	if err == sql.ErrNoRows || strings.Index(err.Error(), "sql: expected") == 0 {
		return false
	}
	return !gorp.NonFatalError(err)
}

type Mysql struct {
	*MysqlExecutor
	DbMap       *gorp.DbMap
	tmp         *MysqlExecutor
	Transaction *gorp.Transaction
	options     Options
}

func (m *Mysql) AddTable(t interface{}) {
	m.DbMap.AddTable(t).SetKeys(true, "Id")
}
func (m *Mysql) WithTx(fn func()) {
	log.Trace().Func("WithTx").Send()
	tx, err := m.DbMap.Begin()
	check(err, ErrorTxBegin, err.Error())
	defer func() {
		m.MysqlExecutor = m.tmp
		m.Transaction = nil
		m.tmp = nil
		if err := recover(); nil != err {
			log.Error().Func("WithTx").Stack().Err(err.(error)).Msg("rollback transaction")
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Error().Func("WithTx").Stack().Err(errRollback).Msg(errRollback.Error())
				check(errRollback, ErrorTxRollback, "")
			}
		} else {
			errCommit := tx.Commit()
			log.Trace().Func("WithTx").Msg("commit transaction")
			if errCommit != nil {
				log.Error().Func("WithTx").Stack().Err(errCommit).Msg(errCommit.Error())
				check(errCommit, ErrorTxRollback, "")
			}
		}
	}()
	m.tmp = m.MysqlExecutor
	m.MysqlExecutor = WrapExecutor(tx)
	m.Transaction = tx
	fn()
}
func (m *Mysql) Close() {
	m.DbMap.Db.Close()
}

type MysqlExecutor struct {
	gorp.SqlExecutor
}

func (m *MysqlExecutor) With(ctx context.Context) {
	m.SqlExecutor = m.SqlExecutor.WithContext(ctx)
}
func (m *MysqlExecutor) Select(result interface{}, query string, args ...interface{}) {
	log.Trace().Func("Select").Str("sql", query).Interface("args", args).Send()
	_, err := m.SqlExecutor.Select(result, query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("Select").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
}
func (m *MysqlExecutor) SelectOne(result interface{}, query string, args ...interface{}) {
	log.Trace().Func("SelectOne").Str("sql", query).Interface("args", args).Send()
	err := m.SqlExecutor.SelectOne(result, query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("SelectOne").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
}
func (m *MysqlExecutor) SelectInt(query string, args ...interface{}) int64 {
	log.Trace().Func("SelectInt").Str("sql", query).Interface("args", args).Send()
	r, err := m.SqlExecutor.SelectInt(query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("SelectInt").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
	return r
}
func (m *MysqlExecutor) SelectNullInt(query string, args ...interface{}) null.Int {
	log.Trace().Func("SelectNullInt").Str("sql", query).Interface("args", args).Send()
	r, err := m.SqlExecutor.SelectNullInt(query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("SelectNullInt").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
	return null.Int{r}
}
func (m *MysqlExecutor) SelectFloat(query string, args ...interface{}) int64 {
	log.Trace().Func("SelectFloat").Str("sql", query).Interface("args", args).Send()
	r, err := m.SqlExecutor.SelectInt(query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("SelectFloat").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
	return r
}
func (m *MysqlExecutor) SelectNullFloat(query string, args ...interface{}) null.Float {
	log.Trace().Func("SelectNullFloat").Str("sql", query).Interface("args", args).Send()
	r, err := m.SqlExecutor.SelectNullFloat(query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("SelectNullFloat").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
	return null.Float{r}
}
func (m *MysqlExecutor) SelectStr(query string, args ...interface{}) string {
	log.Trace().Func("SelectStr").Str("sql", query).Interface("args", args).Send()
	r, err := m.SqlExecutor.SelectStr(query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("SelectStr").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
	return r
}
func (m *MysqlExecutor) SelectNullStr(query string, args ...interface{}) null.String {
	log.Trace().Func("SelectNullFloat").Str("sql", query).Interface("args", args).Send()
	r, err := m.SqlExecutor.SelectNullStr(query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("SelectNullFloat").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
	return null.String{r}
}
func (m *MysqlExecutor) Exec(query string, args ...interface{}) sql.Result {
	log.Trace().Func("Exec").Str("sql", query).Interface("args", args).Send()
	r, err := m.SqlExecutor.Exec(query, args...)
	if err != nil && fatalError(err) {
		log.Error().Func("Exec").Stack().Err(err).Str("sql", query).Interface("args", args).Msg(err.Error())
		check(err, ErrorQuery, err.Error())
	}
	return r
}
func (m *MysqlExecutor) Get(result interface{}, table string, idField string, id interface{}) {
	log.Trace().Func("GetTable").Str("table", table).Interface("id", id).Str("idField", idField).Send()
	if idField == "" {
		idField = "id"
	}
	query := "select * from `" + table + "` where `" + idField + "`=? limit 1"
	m.SelectOne(result, query, id)
}
func (m *MysqlExecutor) GetTable(result interface{}, id interface{}) {
	log.Trace().Func("Get").Interface("id", id).Send()
	resultT := reflect.TypeOf(result)
	m.Get(result, resultT.Elem().Name(), "id", id)
}
func (m *MysqlExecutor) GetBy(result interface{}, table string, kvs map[string]interface{}) {
	log.Trace().Func("GetBy").Str("table", table).Interface("kvs", kvs).Send()
	query := "select * from `" + table + "` where 1=1 "
	for k := range kvs {
		query += " and `" + k + "`=:" + k
	}
	query += " limit 1"
	m.SelectOne(result, query, kvs)
}
func (m *MysqlExecutor) List(result interface{}, table string, idField string, ids []int) {
	log.Trace().Func("List").Str("table", table).Ints("ids", ids).Send()
	if idField == "" {
		idField = "id"
	}
	query := "select * from `" + table + "` where `" + idField + "` in (" + sqlutil.JoinInts(ids) + ") limit " + strconv.Itoa(len(ids))
	m.Select(result, query)
}
func (m *MysqlExecutor) ListByStrIds(result interface{}, table string, idField string, ids []string) {
	log.Trace().Func("List").Str("table", table).Strs("ids", ids).Send()
	if idField == "" {
		idField = "id"
	}
	query := "select * from `" + table + "` where `" + idField + "` in (" + sqlutil.JoinStrs(ids) + ") limit " + strconv.Itoa(len(ids))
	m.SelectOne(result, query)
}

func (m *MysqlExecutor) ListBy(result interface{}, table string, kvs map[string]interface{}) {
	log.Trace().Func("ListBy").Str("table", table).Interface("kvs", kvs).Send()
	query := "select * from `" + table + "` where 1=1 "
	for k := range kvs {
		query += " and `" + k + "`=:" + k
	}
	m.Select(result, query, kvs)
}

func (m *MysqlExecutor) Insert(obj interface{}) {
	log.Trace().Func("Insert").Interface("obj", obj).Send()
	err := m.SqlExecutor.Insert(obj)
	if err != nil {
		log.Error().Func("Insert").Stack().Err(err).Interface("obj", obj).Msg(err.Error())
		check(err, ErrorInsert, err.Error())
	}
}
func (m *MysqlExecutor) Inserts(list ...interface{}) {
	log.Trace().Func("Inserts").Interface("list", list).Send()
	if len(list) <= 0 {
		return
	}
	err := m.SqlExecutor.Insert(list...)
	if err != nil {
		log.Error().Func("Insert").Stack().Err(err).Interface("list", list).Msg(err.Error())
		check(err, ErrorInsert, err.Error())
	}
}
func (m *MysqlExecutor) Updates(list ...interface{}) int64 {
	log.Trace().Func("Updates").Interface("list", list).Send()
	r, err := m.SqlExecutor.Update(list...)
	if err != nil {
		log.Error().Func("Insert").Stack().Err(err).Interface("list", list).Msg(err.Error())
		check(err, ErrorInsert, err.Error())
	}
	return r
}

func (m *MysqlExecutor) Delete(table string, idField string, id int) sql.Result {
	log.Trace().Func("Delete").Str("table", table).Int("id", id).Str("idField", idField).Send()
	query := "delete from `" + table + "` where `" + idField + "`=?"
	r, err := m.SqlExecutor.Exec(query, id)
	if err != nil {
		log.Error().Func("Delete").Stack().Err(err).Str("table", table).Int("id", id).Str("idField", idField).Msg(err.Error())
		check(err, ErrorDelete, err.Error())
	}
	return r
}
func (m *MysqlExecutor) Deletes(table string, idField string, ids []int) sql.Result {
	log.Trace().Func("Delete").Str("table", table).Ints("ids", ids).Str("idField", idField).Send()
	query := "delete from `" + table + "` where `" + idField + "`in (" + sqlutil.JoinInts(ids) + ")"
	r, err := m.SqlExecutor.Exec(query)
	if err != nil {
		log.Error().Func("Delete").Stack().Err(err).Str("table", table).Ints("ids", ids).Str("idField", idField).Msg(err.Error())
		check(err, ErrorDelete, err.Error())
	}
	return r
}

//Patch update table by idFieldidField must in params ,update `table` set k1=v1,k2=v2 where `idField`=params[idField]
func (m *MysqlExecutor) Patch(table string, idField string, params interface{}) sql.Result {
	log.Trace().Func("Patch").Str("table", table).Interface("params", params).Str("idField", idField).Send()
	if idField == "" {
		idField = "id"
	}
	query := "update `" + table + "` set "
	if kvs, ok := params.(map[string]interface{}); ok {
		for field := range kvs {
			if field == idField {
				continue
			}
			query += "`" + field + "`=:" + field + ","
		}
	} else {
		s := sutil.New(params)
		for _, field := range s.Fields() {
			name := field.Name()
			// if strings.ToLower(name) == strings.ToLower(idField) {
			if name == idField {
				continue
			}
			query += "`" + name + "`=:" + name + ","
		}
	}
	query = query[0:len(query)-1] + " where `" + idField + "`=:" + idField
	return m.Exec(query, params)
}

//PatchArgs .update `table` set k1=v1,k2=v2 where `idField`=id ,
func (m *MysqlExecutor) PatchArgs(table string, idField string, idValue interface{}, kvs ...interface{}) sql.Result {
	params := sutil.Kv2Map(kvs)
	params[idField] = idValue
	return m.Patch(table, idField, params)
}
