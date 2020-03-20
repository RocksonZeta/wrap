package sutil_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/utils/sutil"
	"github.com/stretchr/testify/suite"
)

type SliceSuite struct {
	suite.Suite
}

func (s *SliceSuite) TestSlice() {
	var m []User

	arr := sutil.FromSlice(&m).Make(2)
	arr.Put(0, User{1, "jim"})
	fmt.Println(m)
	v0 := arr.GetRef(0).(*User)
	v0.Id = 2
	fmt.Println(m)

}
func (s *SliceSuite) TestWrap() {
	arr := []User{User{1, "jim"}}

	sl := sutil.FromSlice(&arr)
	fmt.Println(sl.Get(0))

}

func TestSliceSuite(t *testing.T) {
	suite.Run(t, new(SliceSuite))

}
