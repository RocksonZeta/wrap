package sutil_test

import (
	"testing"

	"github.com/RocksonZeta/wrap/utils/sutil"
	"github.com/stretchr/testify/suite"
)

type MapSuite struct {
	suite.Suite
}

func (s *MapSuite) TestMap() {
	var m map[int]User
	fm := sutil.FromMap(&m)
	fm.Make()
	fm.Put(1, User{Id: 1})
	s.Equal(1, fm.Get(1).(User).Id)
	fm.Each(func(k, v interface{}) {
		s.Equal(1, k)
	})
}

func TestMapSuite(t *testing.T) {
	suite.Run(t, new(MapSuite))

}
