package rediswrap_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/RocksonZeta/wrap/rediswrap"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/suite"
)

type RedisSuite struct {
	suite.Suite
	client *rediswrap.Redis
}

type User struct {
	Id   int
	Name string
}

func (s *RedisSuite) SetupTest() {
	var options rediswrap.Options
	// options.PoolSize = 100
	// options.MinIdleConns = 1
	s.client = rediswrap.New(options)
}
func (s *RedisSuite) TearDownTest() {
	s.client.Close()
}
func (s *RedisSuite) TestGet() {
	// u := User{Id: 1, Name: "jim"}
	// s.client.SetJson("k1", u, 3)
	var old User
	s.False(s.client.GetJson("k1", &old))
}
func (s *RedisSuite) TestGetM() {
	u := User{Id: 1, Name: "jim"}
	s.client.SetJson("k1", u, 3)
	s.client.SetJson("k2", u, 3)
	var old []User
	s.True(s.client.GetJsons([]string{"k1", "k0", "k2"}, &old))
	fmt.Println(old)
}
func (s *RedisSuite) TestHGetAll() {
	u := User{Id: 1, Name: "jim"}
	s.client.HSetJson("ha", "k1", []User{u})
	s.client.HSetJson("ha", "k2", []User{u})
	var all map[string][]User
	s.client.HMGetAllJson("ha", &all)
	s.Equal(2, len(all))
	fmt.Println("all:", all)
}
func (s *RedisSuite) TestHGet() {
	ps := s.client.Subscribe("c1")
	go func() {
		for {
			v := <-ps.Channel()
			fmt.Println("receive:", v)
			break
		}
	}()
	s.client.HSet("h", "k1", "v1")
	r := s.client.HGet("h", "k1").Val()
	s.Equal("v1", r)
	r0 := s.client.HGet("h", "k0").Val()
	s.Equal("", r0)
	incr := s.client.HIncrBy("h", "k0", 1)

	s.Equal(int64(1), incr)
	s.client.Del("h")
	s.client.Publish("c1", "hello redis")
}
func (s *RedisSuite) TestPipeline() {
	var incr *redis.IntCmd
	s.client.Pipelined(func(pipe redis.Pipeliner) error {
		key := "pipelined_counter"
		pipe.Set(key, "a", 10*time.Second)
		incr = pipe.Incr(key)
		pipe.Expire(key, 10*time.Second)
		return nil
	})
	fmt.Println(incr.Val())
}
func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(RedisSuite))

}
