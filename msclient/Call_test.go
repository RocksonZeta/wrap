package msclient_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/msclient"
	"github.com/stretchr/testify/suite"
)

type CallSuite struct {
	suite.Suite
}
type U struct {
	Id       int
	RealName string
}

func (c *CallSuite) TestCall() {
	cookies := map[string]string{
		"sessionid": "7qd1yjwvy4icn34v1jg4qs7p4x48lkti",
	}
	call := msclient.New("https://test.iqidao.com", nil, cookies, 3)
	var r interface{}
	err := call.Get(&r, "/json/user/info?a=c", map[string]string{"a": "b"})
	c.Nil(err)
	fmt.Println(r, err)
}

func TestMysqlSuite(t *testing.T) {
	suite.Run(t, new(CallSuite))
}
