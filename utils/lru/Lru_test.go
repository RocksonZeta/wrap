package lru_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/RocksonZeta/wrap/utils/lru"
	"github.com/stretchr/testify/suite"
)

type LruSuite struct {
	suite.Suite
	lru *lru.Lru
}

func (s *LruSuite) BeforeTest() {
	fmt.Println("BeforeTest")
}
func (s *LruSuite) TestAdd() {
	s.lru = lru.New(lru.Options{Ttl: 2, CleanupInterval: 1})
	s.lru.Add("k1", 1)
	s.Equal(1, s.lru.Get("k1"))
	time.Sleep(3 * time.Second)
	s.Nil(s.lru.Get("k1"))

}

func TestLru(t *testing.T) {
	suite.Run(t, new(LruSuite))
}
