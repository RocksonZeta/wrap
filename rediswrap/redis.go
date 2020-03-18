package rediswrap

import (
	"context"
	"encoding/json"
	"reflect"
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
	client := redis.NewClient(opt)
	// if log.DebugEnabled() {
	client.AddHook(redisHook{})
	// }
	newRed := &Redis{Client: client, options: options}
	redises.Store(options.Url, newRed)
	return newRed
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
	options Options
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
	log.Trace().Func("Unmarshal").Str("str", str).Send()
	if str == "" {
		return
	}
	err := json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Error().Func("Unmarshal").Stack().Err(err).Str("str", str).Msg(err.Error())
		check(err, ErrorMarshal, err.Error())
	}
}

func (r *Redis) GetJson(key string, result interface{}) {
	log.Trace().Func("GetJson").Str("key", key).Interface("result", result).Send()
	str := r.Get(key).Val()
	r.Unmarshal(str, result)
}
func (r *Redis) SetJson(key string, result interface{}, secs int) {
	log.Trace().Func("SetJson").Str("key", key).Interface("result", result).Int("secs", secs).Send()
	r.Set(key, r.Marshal(result), secs)
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
func (r *Redis) HGetJson(key, field string, result interface{}) {
	log.Trace().Func("HGetJson").Str("key", key).Str("field", field).Interface("result", result).Send()
	str := r.Client.HGet(key, field).Val()
	r.Unmarshal(str, result)
}
func (r *Redis) HMGetAllJson(key string, result interface{}) {
	log.Trace().Func("HMGetAllJson").Str("key", key).Send()
	m := r.Client.HGetAll(key).Val()
	resultMap := sutil.FromMap(result)
	resultMap.MakeWithSize(len(m))
	for key, value := range m {
		err := resultMap.PutJson(key, []byte(value))
		if err != nil {
			log.Error().Func("HMGetAllJson").Stack().Err(err).Str("key", key)
			check(err, ErrorMarshal, err.Error())
		}
	}
}
