package nutil_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/RocksonZeta/wrap/utils/nutil"
	"github.com/stretchr/testify/suite"
)

type ValueSuite struct {
	suite.Suite
	x *nutil.Value
}

func (s *ValueSuite) SetupSuite() {
}
func (s *ValueSuite) SetupTest() {
}
func (s *ValueSuite) TearDownTest() {
}
func (s *ValueSuite) TestInt() {
	var x interface{}
	x = "123.0"
	v := nutil.ValueOf(x)
	s.Equal(v.AsInt(), 123)
}
func (s *ValueSuite) TestString() {
	v := nutil.ValueOf(123)
	s.Equal(v.String(), "123")
	a := make([]int, 3)
	// a[0] = 2
	fmt.Println(reflect.TypeOf(a).Kind())
	fmt.Println(nutil.ValueOf(true).String())
}
func TestValueSuite(t *testing.T) {
	suite.Run(t, new(ValueSuite))
}
