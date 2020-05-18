package rediswrap

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/RocksonZeta/wrap/errs"
	"github.com/RocksonZeta/wrap/utils/sutil"
	"github.com/RocksonZeta/wrap/wraplog"
	"github.com/go-redis/redis/v7"
)

var pkg = reflect.TypeOf(Redis{}).PkgPath()

var log = wraplog.Logger.Fork(pkg, "Redis")

const (
	ErrorInit = 1 + iota
	ErrorCmd
	ErrorClose
	ErrorMarshal
	ErrorTypeCast
	ErrorPipeline
)

func isNilError(err error) bool {
	return "redis: nil" == err.Error()
}

func check(err error, state int, msg string) {
	if err == nil {
		return
	}
	if msg != "" {
		msg = err.Error()
	}
	panic(errs.Err{Err: err, Module: "Redis", Pkg: pkg, State: state, Message: msg})
}

// var redises sync.Map

type Options struct {
	//Url eg. redis://:qwerty@localhost:6379/1
	Url          string
	PoolSize     int
	MinIdleConns int
	IdleTimeout  int
}

func ParseUrl(redisUrl string) (*redis.Options, error) {
	urlParts, err := url.Parse(redisUrl)
	if err != nil {
		return nil, err
	}
	q := urlParts.Query()
	urlParts.RawQuery = ""
	opt, err := redis.ParseURL(urlParts.String())
	if v := q.Get("PoolSize"); v != "" {
		opt.PoolSize, _ = strconv.Atoi(v)
	}
	if v := q.Get("MinIdleConns"); v != "" {
		opt.MinIdleConns, _ = strconv.Atoi(v)
	}
	if v := q.Get("IdleTimeout"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.IdleTimeout = time.Duration(vi) * time.Second
	}
	if v := q.Get("PoolTimeout"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.PoolTimeout = time.Duration(vi) * time.Second
	}
	if v := q.Get("ReadTimeout"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.ReadTimeout = time.Duration(vi) * time.Second
	}
	if v := q.Get("DialTimeout"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.DialTimeout = time.Duration(vi) * time.Second
	}
	if v := q.Get("WriteTimeout"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.WriteTimeout = time.Duration(vi) * time.Second
	}
	if v := q.Get("IdleCheckFrequency"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.IdleCheckFrequency = time.Duration(vi) * time.Second
	}
	if v := q.Get("MaxConnAge"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.MaxConnAge = time.Duration(vi) * time.Second
	}
	if v := q.Get("MaxRetries"); v != "" {
		opt.MaxRetries, _ = strconv.Atoi(v)
	}
	if v := q.Get("MaxRetryBackoff"); v != "" {
		vi, _ := strconv.Atoi(v)
		opt.MaxRetryBackoff = time.Duration(vi) * time.Second
	}
	return opt, nil
}
func NewFromUrl(redisUrl string) *Redis {
	option, err := ParseUrl(redisUrl)
	if err != nil {
		check(err, ErrorInit, err.Error())
	}
	return NewFromOption(option)
}
func NewFromOption(option *redis.Options) *Redis {
	client := redis.NewClient(option)
	client.AddHook(redisHook{})
	newRed := &Redis{Client: client}
	return newRed
}
func New(options Options) *Redis {
	log.Trace().Interface("options", options).Send()
	// red, ok := redises.Load(options.Url)
	// if ok {
	// 	return red.(*Redis)
	// }

	opt, err := redis.ParseURL(options.Url)
	if err != nil {
		log.Error().Func("New").Stack().Err(err).Interface("options", options).Msg(err.Error())
		check(err, ErrorInit, err.Error())
	}
	if options.PoolSize > 0 {
		opt.PoolSize = options.PoolSize
	}
	if options.MinIdleConns > 0 {
		opt.MinIdleConns = options.MinIdleConns
	}
	if options.IdleTimeout > 0 {
		opt.IdleTimeout = time.Duration(options.IdleTimeout) * time.Second
	}
	return NewFromOption(opt)
}

var logHoot redis.Hook = redisHook{}

type redisHook struct {
	redis.Hook
}

func (redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	log.Debug().Func("BeforeProcess").Interface("cmd", cmd.Args()).Send()
	return ctx, nil
}

func (redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if err := cmd.Err(); err != nil && !isNilError(err) {
		log.Error().Func("AfterProcess").Interface("cmd", cmd.Args()).Stack().Err(err).Msg(err.Error())
		check(err, ErrorCmd, err.Error())
	}
	return nil
}
func (redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if log.DebugEnabled() {
		args := make([][]interface{}, len(cmds))
		for i, cmd := range cmds {
			args[i] = cmd.Args()
		}
		log.Debug().Func("BeforeProcessPipeline").Interface("cmds", args).Send()
	}
	return ctx, nil
}
func (redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			log.Error().Func("AfterProcessPipeline").Interface("cmd", cmd.Args()).Stack().Err(err).Msg(err.Error())
			check(err, ErrorPipeline, err.Error())
		}
	}
	return nil
}

type Redis struct {
	*redis.Client
	// options Options
}

func (r *Redis) Close() {
	log.Trace().Func("Close").Send()

	// err := r.Client.Close()
	// if err != nil {
	// 	log.Error().Func("Close").Stack().Err(err).Msg(err.Error())
	// 	check(err, ErrorMarshal, err.Error())
	// }
}
func (r *Redis) Marshal(v interface{}) string {
	log.Trace().Func("Marshal").Interface("v", v).Send()
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error().Func("Marshal").Stack().Err(err).Interface("v", v).Msg(err.Error())
		check(err, ErrorMarshal, err.Error())
	}
	return string(bs)
}

func (r *Redis) Unmarshal(str string, result interface{}) {
	log.Trace().Func("Unmarshal").Str("str", str).Interface("result", result).Send()
	if str == "" {
		return
	}
	err := json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Error().Func("Unmarshal").Stack().Err(err).Str("str", str).Msg(err.Error())
		check(err, ErrorMarshal, err.Error())
	}
}

func (r *Redis) GetJson(key string, result interface{}) bool {
	log.Trace().Func("GetJson").Str("key", key).Interface("result", result).Send()
	cmd := r.Get(key)
	err := cmd.Err()
	if err != nil && isNilError(err) {
		return false
	}
	str := cmd.Val()

	r.Unmarshal(str, result)
	return true
}
func (r *Redis) SetJson(key string, result interface{}, secs int) {
	log.Trace().Func("SetJson").Str("key", key).Interface("result", result).Int("secs", secs).Send()
	r.Set(key, r.Marshal(result), secs)
}
func (r *Redis) GetJsons(keys []string, result interface{}) bool {
	log.Trace().Func("GetJsons").Strs("keys", keys).Interface("result", result).Send()
	cmd := r.MGet(keys...)
	err := cmd.Err()
	if err != nil && isNilError(err) {
		return false
	}
	strs := cmd.Val()
	results := sutil.FromSlice(result).Make(len(keys))
	for i, v := range strs {
		if v == nil {
			continue
		}
		r.Unmarshal(v.(string), results.GetRef(i))
	}
	return true
}

func (r *Redis) SetJsons(keys []string, values interface{}, secs int) {
	if len(keys) <= 0 {
		return
	}
	r.Pipelined(func(p redis.Pipeliner) error {
		args := make([]interface{}, len(keys)*2)
		valuesV := sutil.PtrValue(reflect.ValueOf(values))

		for i := 0; i < len(args); i += 2 {
			args[i] = keys[i/2]
			args[i+1] = r.Marshal(valuesV.Index(i / 2).Interface())
		}
		p.MSet(args...)
		for _, k := range keys {
			p.Expire(k, time.Duration(secs)*time.Second)
		}
		return nil
	})
}
func (r *Redis) SetJsonsByKv(kvs map[string]interface{}, secs int) {
	log.Trace().Func("SetJsons").Interface("kvs", kvs).Send()
	if len(kvs) == 0 {
		return
	}
	r.Pipelined(func(p redis.Pipeliner) error {
		args := make([]string, len(kvs)*2)
		i := 0
		for k, v := range kvs {
			args[i] = k
			args[i+1] = r.Marshal(v)
			i += 2
		}
		p.MSet(args)
		for k := range kvs {
			p.Expire(k, time.Duration(secs)*time.Second)
		}
		return nil
	})
}

func (r *Redis) Get(key string) *redis.StringCmd {
	log.Trace().Func("Get").Str("key", key).Send()
	cmd := r.Client.Get(key)
	return cmd
}
func (r *Redis) Set(key string, value interface{}, secs int) {
	log.Trace().Func("Get").Str("key", key).Send()
	r.Client.Set(key, value, time.Duration(secs)*time.Second)
}

func (r *Redis) HSetJson(key, field string, value interface{}) {
	log.Trace().Func("HSetJson").Str("key", key).Send()
	r.Client.HSet(key, field, r.Marshal(value))
}
func (r *Redis) HGetJson(key, field string, result interface{}) bool {
	log.Trace().Func("HGetJson").Str("key", key).Str("field", field).Interface("result", result).Send()
	cmd := r.Client.HGet(key, field)
	str := cmd.Val()
	err := cmd.Err()
	if err != nil && isNilError(err) {
		return false
	}
	r.Unmarshal(str, result)
	return true
}
func (r *Redis) HMGetAllJson(key string, result interface{}) {
	log.Trace().Func("HMGetAllJson").Str("key", key).Send()
	m := r.Client.HGetAll(key).Val()
	resultMap := sutil.FromMap(result).MakeWithSize(len(m))
	for key, value := range m {
		err := resultMap.PutJson(key, []byte(value))
		if err != nil {
			check(err, ErrorMarshal, err.Error())
		}
	}
}
