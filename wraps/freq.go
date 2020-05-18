package wraps

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/RocksonZeta/wrap/msclient"
	"github.com/RocksonZeta/wrap/mysqlwrap"
	"github.com/RocksonZeta/wrap/osswrap"
	"github.com/RocksonZeta/wrap/redised"
	"github.com/RocksonZeta/wrap/rediswrap"
	"github.com/RocksonZeta/wrap/requestwrap"
	"github.com/RocksonZeta/wrap/utils/fileutil"
	"github.com/RocksonZeta/wrap/wraplog"
	"gopkg.in/yaml.v2"
)

var log = wraplog.Logger.Fork("github.com/RocksonZeta/wrap/wraps", "")

// wrap实例管理

// all thread safe ,call be singleton
var redises sync.Map
var mysqls sync.Map
var osses sync.Map
var requests sync.Map
var calls sync.Map

func GetRedisedMysqlUrl(redisUrl, mysqlUrl string) *redised.RedisedMysql {
	return redised.NewRedisedMysql(GetRedisUrl(redisUrl), GetMysqlUrl(mysqlUrl))
}
func GetRedisedMysql(ro rediswrap.Options, mo mysqlwrap.Options) *redised.RedisedMysql {
	return redised.NewRedisedMysql(GetRedis(ro), GetMysql(mo))
}
func GetRedis(options rediswrap.Options) *rediswrap.Redis {
	old, ok := redises.Load(options.Url)
	if ok {
		return old.(*rediswrap.Redis)
	}
	red := rediswrap.New(options)
	redises.Store(options.Url, red)
	return red
}
func GetRedisUrl(redisUrl string) *rediswrap.Redis {
	old, ok := redises.Load(redisUrl)
	if ok {
		return old.(*rediswrap.Redis)
	}
	red := rediswrap.NewFromUrl(redisUrl)
	redises.Store(redisUrl, red)
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
func GetMysqlUrl(mysqlUrl string) *mysqlwrap.Mysql {
	old, ok := mysqls.Load(mysqlUrl)
	if ok {
		return old.(*mysqlwrap.Mysql)
	}
	n := mysqlwrap.NewFromUrl(mysqlUrl)
	mysqls.Store(mysqlUrl, n)
	return n
}
func GetOss(options osswrap.Options, bucketName string) *osswrap.Oss {
	key := bucketName + "." + options.Endpoint
	old, ok := osses.Load(key)
	if ok {
		return old.(*osswrap.Oss)
	}
	n := osswrap.New(options, bucketName)
	osses.Store(key, n)
	return n
}
func GetRequest(baseUrl string, headers, cookies map[string]string, timeout int) *requestwrap.Request {
	old, ok := requests.Load(baseUrl)
	if ok {
		return old.(*requestwrap.Request)
	}
	n := requestwrap.New(baseUrl, headers, cookies, timeout)
	requests.Store(baseUrl, n)
	return n
}

func GetCall(baseUrl string, headers, cookies map[string]string, timeout int) *msclient.Call {
	old, ok := calls.Load(baseUrl)
	if ok {
		return old.(*msclient.Call)
	}
	n := msclient.New(baseUrl, headers, cookies, timeout)
	calls.Store(baseUrl, n)
	return n
}

func GetConfig(config interface{}, configFileName string) error {
	cwd := fileutil.FindFileDir(configFileName)
	bs, err := ioutil.ReadFile(filepath.Join(cwd, configFileName))
	if err != nil {
		fmt.Println("read config.yml err. " + err.Error())
		return err
	} else {
		err = yaml.Unmarshal(bs, config)
		if err != nil {
			fmt.Println("Unmarshal config.yml err. " + err.Error())
			return err
		}
	}
	return nil
}
