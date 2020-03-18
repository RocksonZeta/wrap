package wraps

import (
	"sync"

	"github.com/RocksonZeta/wrap/mysqlwrap"
	"github.com/RocksonZeta/wrap/osswrap"
	"github.com/RocksonZeta/wrap/rediswrap"
	"github.com/RocksonZeta/wrap/requestwrap"
)

// wrap实例管理

var redises sync.Map
var mysqls sync.Map
var https sync.Map
var osses sync.Map

func GetRedis(options rediswrap.Options) *rediswrap.Redis {
	old, ok := redises.Load(options.Url)
	if ok {
		return old.(*rediswrap.Redis)
	}
	red := rediswrap.New(options)
	redises.Store(options.Url, red)
	return red
}
func GetMysql(options mysqlwrap.Options) *mysqlwrap.Mysql {
	old, ok := mysqls.Load(options.Url)
	if ok {
		return old.(*mysqlwrap.Mysql)
	}
	n := mysqlwrap.New(options)
	mysqls.Store(options.Url, n)
	return n
}
func GetOss(options osswrap.Options, bucketName string) *osswrap.Oss {
	key := bucketName + "." + options.Endpoint
	old, ok := mysqls.Load(key)
	if ok {
		return old.(*osswrap.Oss)
	}
	n := osswrap.New(options, bucketName)
	mysqls.Store(key, n)
	return n
}
func GetRequest(baseUrl string, headers, cookies map[string]string, timeout int) *requestwrap.Request {
	old, ok := mysqls.Load(baseUrl)
	if ok {
		return old.(*requestwrap.Request)
	}
	n := requestwrap.New(baseUrl, headers, cookies, timeout)
	mysqls.Store(baseUrl, n)
	return n
}
