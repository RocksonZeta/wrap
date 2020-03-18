package iriswrap

import (
	"sync"
	"time"

	"github.com/RocksonZeta/wrap/rediswrap"
	"github.com/RocksonZeta/wrap/utils/lru"
	"github.com/RocksonZeta/wrap/utils/mathutil"
	"github.com/RocksonZeta/wrap/utils/osutil"
	"github.com/RocksonZeta/wrap/wraplog"
	"github.com/go-redis/redis/v7"
	"github.com/kataras/iris/v12"
)

var slog = wraplog.Logger.Fork(pkg, "Session")

var sessionPrefix string = "session/"
var SessionUidKey string = "UID"
var SessionCookieId string = "sessionid"
var SessionCookieTtl int = 30 * 60
var SessionCookieTokenId string = "sessiontoken" //网页记住我
var SessionCookieTokenTtl int = 30 * 24 * 3600   //网页记住我的时长
var SessionCookieDomain string
var SessionHeaderId = "X-USER-TOKEN"
var RedisClient *rediswrap.Redis
var GetUidByToken func(token string, tokenType int) int

var localSessionUids *lru.Lru = lru.New(lru.Options{Ttl: SessionCookieTtl, MaxAge: 3600})
var sidNeedRefresh *cacheRefresh = newCacheRefresh()

type SessionOptions struct {
	// SessionCookieTtl    int //seconds
	// SessionCookieId     string
	// SessionUidKey       string
	// SessionHeaderId     string
	// SessionCookieDomain string
	// Redis               *rediswrap.Redis
	// GetUidByToken       func(token string) int
}

// func InitSession(options SessionOptions) {
// red = options.Redis
// GetUidByToken = options.GetUidByToken
// SessionCookieDomain = options.SessionCookieDomain
// if options.SessionCookieTtl > 0 {
// 	SessionCookieTtl = options.SessionCookieTtl
// }
// if options.SessionCookieId != "" {
// 	SessionCookieId = options.SessionCookieId
// }
// if options.SessionUidKey != "" {
// 	SessionUidKey = options.SessionUidKey
// }
// if options.SessionHeaderId != "" {
// 	SessionHeaderId = options.SessionHeaderId
// }
// localUids =
// sidNeedRefresh =

// }

func SessionFilter(ctx iris.Context) {
	log.Trace().Func("SessionFilter").Str("method", ctx.Method()).Str("path", ctx.Path()).Send()
	c := ctx.(*Context)
	headerToken := c.GetHeader(SessionHeaderId)
	var cookieToken string
	var sid string
	if headerToken == "" {
		cookieToken = c.GetCookie(SessionCookieTokenId)
		sid = c.GetCookie(SessionCookieId)
		if sid == "" {
			sid = mathutil.RandomStr(32, false)
		}
	}
	c.SetCookieLocal(SessionCookieId, sid, SessionCookieTtl, true)
	c.Session = &Session{Sid: sid, HeaderToken: headerToken, CookieToken: cookieToken}
	c.Next()
}

type Session struct {
	Sid         string
	HeaderToken string
	CookieToken string
	uid         int
	// isToken bool //是临时生成的sessionid，还是标识用户的token
}

func (s Session) cacheKey() string {
	return sessionPrefix + s.Sid
}
func (s Session) GetString(key string) string {
	var r string
	s.Get(key, &r)
	return r
}
func (s Session) GetInt(key string) int {
	var r int
	s.Get(key, &r)
	return r
}
func (s Session) Get(key string, result interface{}) {
	log.Trace().Func("Get").Send()

	// red := s.RedisFactory()
	// defer red.Close()
	// if key == SessionUidKey {
	// 	localUids.
	// }

	RedisClient.HGetJson(s.cacheKey(), key, result)
	sidNeedRefresh.Add(s.Sid)
	// red.Expire(s.cacheKey(), time.Duration(s.Ttl)*time.Second)
}
func (s *Session) Set(key string, value interface{}) {
	log.Trace().Func("Set").Send()
	// red := s.RedisFactory()
	// defer red.Close()
	RedisClient.HSetJson(s.cacheKey(), key, value)
	RedisClient.Expire(s.cacheKey(), time.Duration(SessionCookieTtl)*time.Second)
}
func (s *Session) SetUid(uid int) {
	s.uid = uid
	localSessionUids.Add(s.Sid, uid)
	s.Set(SessionUidKey, uid)
}
func (s *Session) Uid() int {
	log.Trace().Func("Uid").Send()
	if s.uid > 0 {
		sidNeedRefresh.Add(s.Sid)
		return s.uid
	}
	var uid int
	//从本地缓存中获取uid
	if s.Sid != "" {
		uid = localSessionUids.GetInt(s.Sid)
		if uid > 0 {
			s.uid = uid
			sidNeedRefresh.Add(s.Sid)
			return uid
		}
		//从redis中获取uid
		uid = s.GetInt(SessionUidKey)
		if uid > 0 {
			s.uid = uid
			localSessionUids.Add(s.Sid, uid)
		}
	}
	//从数据库中用token换取uid
	if s.CookieToken != "" && uid <= 0 {
		uid = GetUidByToken(s.Sid, 0)
		s.SetUid(uid)
	}
	//从数据库中用token换取uid
	if s.HeaderToken != "" && uid <= 0 {
		uid = GetUidByToken(s.Sid, 1)
		s.SetUid(uid)
	}

	return uid
}
func (s *Session) Refresh() {
	// red := s.RedisFactory()
	// defer red.Close()
	RedisClient.Expire(s.cacheKey(), time.Duration(SessionCookieTtl)*time.Second)
}
func (s *Session) Remove(key string) {
	// red := s.RedisFactory()
	// defer red.Close()
	RedisClient.HDel(s.cacheKey(), key)
}
func (s *Session) Destroy() {
	// red := s.RedisFactory()
	// defer red.Close()
	RedisClient.Del(s.cacheKey())
	localSessionUids.Delete(s.Sid)
}

// session cache refresher
type cacheRefresh struct {
	sids   map[string]bool
	ticker *time.Ticker
	lock   sync.Mutex
	ttl    int64 //seconds
	// red    *rediswrap.Redis
}

func newCacheRefresh() *cacheRefresh {
	r := &cacheRefresh{
		sids: make(map[string]bool),
		// red:    red,
	}
	osutil.Go(func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			r.refresh()
		}
	})
	return r
}

func (c *cacheRefresh) Add(sid string) {
	c.lock.Lock()
	c.sids[sid] = true
	c.lock.Unlock()
}
func (c *cacheRefresh) refresh() {
	if len(c.sids) == 0 {
		return
	}
	c.lock.Lock()
	arr := make([]string, len(c.sids))
	i := 0
	for sid := range c.sids {
		arr[i] = sid
		i++
		delete(c.sids, sid)
	}
	c.lock.Unlock()
	RedisClient.Pipelined(func(p redis.Pipeliner) error {
		for _, sid := range arr {
			p.Expire(sid, time.Duration(SessionCookieTtl)*time.Second)
		}
		return nil
	})
}
