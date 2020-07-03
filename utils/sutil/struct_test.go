package sutil_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/utils/sutil"
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
func (s *StructSuite) TestKv2Map() {
	m := sutil.Kv2Map("k1", 1, "k2", "2")
	fmt.Println(m)
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(StructSuite))

}
