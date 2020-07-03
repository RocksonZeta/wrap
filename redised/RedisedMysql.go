package redised

import (
	"reflect"
	"sort"

	"github.com/RocksonZeta/wrap/mysqlwrap"
	"github.com/RocksonZeta/wrap/rediswrap"
	"github.com/RocksonZeta/wrap/utils/cutil"
	"github.com/RocksonZeta/wrap/utils/nutil"
	"github.com/RocksonZeta/wrap/utils/sutil"
	"github.com/RocksonZeta/wrap/wraplog"
)

var pkg = reflect.TypeOf(RedisedMysql{}).PkgPath()
var log = wraplog.Logger.Fork(pkg, "RedisedMysql")

const (
	ErrorInit = 1 + iota
)

func DefaultKeyFn(table string, id interface{}) string {
	return table + "/" + nutil.String(id)
}

// map (table key value ...)-> redis key : table/k1/v1/k2/v2  order by keys
func DefaultKVFn(table string, kvs map[string]interface{}) string {
	keys := make([]string, len(kvs))
	i := 0
	for k := range kvs {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	result := table
	for _, k := range keys {
		result += "/" + nutil.String(k) + "/" + nutil.String(kvs[k])
	}
	return result
}

var DefaultTtl = 3600

func NewRedisedMysql(redis *rediswrap.Redis, mysql *mysqlwrap.Mysql) *RedisedMysql {
	r := &RedisedMysql{
		Redis: redis,
		Mysql: mysql,
		KeyFn: DefaultKeyFn,
		KVFn:  DefaultKVFn,
		Ttl:   DefaultTtl,
	}

	return r
}

type RedisedMysql struct {
	Redis *rediswrap.Redis
	Mysql *mysqlwrap.Mysql
	KeyFn func(table string, id interface{}) string
	KVFn  func(table string, kvs map[string]interface{}) string
	Ttl   int
}

func (r *RedisedMysql) Close() {
	if r.Redis != nil {
		r.Redis.Close()
	}
	if r.Mysql != nil {
		r.Mysql.Close()
	}
}

func (r *RedisedMysql) ClearCache(table string, ids ...int) {
	keys := make([]string, len(ids))
	for i, v := range ids {
		keys[i] = r.KeyFn(table, v)
	}
	r.Redis.Del(keys...)
}

func (r *RedisedMysql) Get(result interface{}, table, idField string, id interface{}) {
	key := r.KeyFn(table, id)
	exists := r.Redis.GetJson(key, result)
	if exists {
		return
	}
	r.Mysql.Get(result, table, idField, id)
	r.Redis.SetJson(key, result, r.Ttl)
}

//GetByKvs kv key -> id key -> result
func (r *RedisedMysql) GetByKvs(result interface{}, table, idField string, kvs map[string]interface{}) {
	key := r.KVFn(table, kvs)
	cachedId, _ := r.Redis.Get(key).Int()
	if cachedId > 0 {
		r.Get(result, table, idField, cachedId)
		return
	}
	r.Mysql.GetBy(result, table, kvs)
	r.Redis.Set(key, sutil.Get(result, idField), r.Ttl)
}

//ListByKvs kvs keys -> list -> values
func (r *RedisedMysql) ListByKvs(result interface{}, table, idField string, kvs map[string]interface{}) {
	key := r.KVFn(table, kvs)
	var ids []int
	exists := r.Redis.GetJson(key, &ids)
	if exists {
		r.List(result, table, idField, ids)
		return
	}
	r.Mysql.ListBy(result, table, kvs)
	idsFromMysql := cutil.ColInt(result, idField)
	r.Redis.SetJson(key, idsFromMysql, r.Ttl)
}

func (r *RedisedMysql) List(result interface{}, table, idField string, ids []int) {
	keys := make([]string, len(ids))
	for i, v := range ids {
		keys[i] = r.KeyFn(table, v)
	}
	r.Redis.GetJsons(keys, result)
	resultV := reflect.ValueOf(result).Elem()
	allCached := cutil.All(result, func(i int, v interface{}) bool {
		return sutil.New(v).Field(idField).Value() != 0
	})
	if allCached {
		return
	}
	resultV.Set(reflect.Zero(reflect.TypeOf(result).Elem()))
	r.Mysql.List(result, table, idField, ids)
	var mkeys []string
	cutil.Map(result, &mkeys, func(i int, v interface{}) interface{} {
		return r.KeyFn(table, sutil.New(v).Field(idField).Value())
	})
	r.Redis.SetJsons(mkeys, result, r.Ttl)
}

//List  result: eg. &[]table.User . redis -> set key result
func (r *RedisedMysql) ListBySql(result interface{}, key string, query string, args ...interface{}) {
	exists := r.Redis.GetJson(key, &result)
	if exists {
		return
	}
	r.Mysql.Select(result, query, args...)
	r.Redis.SetJson(key, result, r.Ttl)
}

//ListBySqlHash  set list result to redis map . redis -> hset key field result
func (r *RedisedMysql) ListBySqlIntoMap(result interface{}, key, field string, query string, args ...interface{}) {
	exists := r.Redis.HGetJson(key, field, result)
	if exists {
		return
	}
	r.Mysql.Select(result, query, args...)
	r.Redis.HSet(key, field, result, r.Ttl)
}

//DelCacheList if we use listByXX method,we should clear cache in delete method
func (r *RedisedMysql) DelCacheList(table string, kvs map[string]interface{}) {
	r.Redis.Del(r.KVFn(table, kvs))
}

//DelCacheList if we use listByXX method,we should clear cache in delete method
func (r *RedisedMysql) DelCacheListArgs(table string, kvs ...interface{}) {
	r.Redis.Del(r.KVFn(table, sutil.Kv2Map(kvs...)))
}
