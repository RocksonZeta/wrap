package nutil_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/utils/nutil"
	"github.com/stretchr/testify/suite"
)

type JsonSuite struct {
	suite.Suite
}

func (s *JsonSuite) SetupSuite() {
}
func (s *JsonSuite) SetupTest() {
}
func (s *JsonSuite) TearDownTest() {
}

type obj1 struct {
	Id   int
	Name string
}
type obj struct {
	Id     int
	Name   string
	Parent obj1
}

func (s *JsonSuite) TestMarshal() {
	a := obj{1, "jim", obj1{2, "tom"}}
	v := nutil.ValueOf(a)
	js := nutil.UnmarshalJSONMust(v.MarshalJSONMust())
	fmt.Println(js.AsMap().Get("Parent").AsMap())
}
func (s *JsonSuite) TestStruct() {
	var a obj
	b := nutil.NewMap()
	b["Id"] = "1"
	b["Id1"] = 1
	b.Struct(&a)
	fmt.Println(a)
}
func (s *JsonSuite) TestSortFields() {
	a := make([]obj, 4)
	a[0] = obj{Id: 1, Name: "jim"}
	a[1] = obj{Id: 1, Name: "tom"}
	a[2] = obj{Id: 4}
	v := nutil.ValueOf(a)
	l := v.AsList()
	l.SortByFields([]nutil.SortOrder{{"Id", false}, {"Name", true}})
	fmt.Println(l)
}
func (s *JsonSuite) TestSortField() {
	a := make([]obj, 4)
	a[0] = obj{Id: 1, Name: "jim"}
	a[1] = obj{Id: 8, Name: "tome"}
	a[2] = obj{Id: 4}
	v := nutil.ValueOf(a)
	l := v.AsList()
	l.SortByField("Id", false)
	fmt.Println(l)
}
func (s *JsonSuite) TestSort() {
	a := []int{5, 1, 2, 3}
	v := nutil.ValueOf(a)
	l := v.AsList()
	l.Sort(true)
	fmt.Println(l)
}
func TestJsonSuite(t *testing.T) {
	suite.Run(t, new(JsonSuite))
}
