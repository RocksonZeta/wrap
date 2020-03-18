package sutil_test

import (
	"testing"

	"github.com/fatih/structs"

	"github.com/stretchr/testify/suite"
)

type StructSuite struct {
	suite.Suite
}
type User struct {
	Id   int
	Name string
}

func (s *StructSuite) TestMap() {
	var u User
	m := structs.New(&u)
	m.Field("Id").Set(1)
	s.Equal(1, u.Id)
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(StructSuite))

}
