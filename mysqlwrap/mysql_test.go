package mysqlwrap_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/RocksonZeta/wrap/mysqlwrap"
	"github.com/stretchr/testify/suite"
)

type MysqlSuite struct {
	suite.Suite
	mysql *mysqlwrap.Mysql
}

func (s *MysqlSuite) SetupTest() {
	s.mysql = mysqlwrap.NewFromUrl("root:123456@tcp(localhost)/power?charset=utf8mb4&MaxIdle=2&MaxOpen=10")
}
func (s *MysqlSuite) AfterTest() {
	// s.mysql.Close()
}

type User struct {
	Id   int
	Name string
}

func (s *MysqlSuite) TestSelect() {
	var result []User
	s.mysql.Select(&result, "select * from User limit 2")
}
func (s *MysqlSuite) TestSelectOne() {
	var result User
	s.mysql.SelectOne(&result, "select * from User where id=:Id limit 1", User{Id: 1})
}
func (s *MysqlSuite) TestPatch() {
	s.mysql.Patch("User", "Id", User{Id: 1, Name: "jim"})
}
func (s *MysqlSuite) TestPatchMap() {
	s.mysql.Patch("User", "Id", map[string]interface{}{"Id": 1, "Name": "tom"})
}
func TestMysqlSuite(t *testing.T) {
	suite.Run(t, new(MysqlSuite))
}
func TestConn(t *testing.T) {
	option := mysqlwrap.Options{
		MaxOpen: 10,
		MaxIdle: 3,
	}
	fmt.Println("hello")
	for i := 0; i < 100; i++ {
		mysqlwrap.New(option)
	}
	mysqlwrap.New(option).SelectInt("select count(*) from User")
	fmt.Println("hello123")
	time.Sleep(10 * time.Second)
}
