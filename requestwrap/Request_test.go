package requestwrap_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/requestwrap"
	"github.com/stretchr/testify/suite"
)

type RequestSuite struct {
	suite.Suite
	req *requestwrap.Request
}

func (s *RequestSuite) SetupSuite() {
	fmt.Println("SetupSuite")
}
func (s *RequestSuite) SetupTest() {
	fmt.Println("SetupTest")
}
func (s *RequestSuite) TearDownTest() {
	fmt.Println("TearDownTest")
}
func (s *RequestSuite) TestGet() {
	s.req = requestwrap.New("https://test.iqidao.com", nil, nil, 3)
	bs, err := s.req.Get("/json/user/info", map[string]string{"a": "b"})
	fmt.Println(string(bs), err)
	// bs, err = s.req.Get("/json/user/info", nil)
	res, err := s.req.Request.Get("https://test.iqidao.com/json/user/info")
	fmt.Println(string(res.StatusCode), err)

}
func (s *RequestSuite) TestPost() {
	req := requestwrap.New("http://localhost:9000", nil, map[string]string{"sid": "xx"}, 3)
	bs, err := req.Post("/form", nil, map[string]string{"name": "jim"}, nil)
	s.Nil(err)
	fmt.Println(string(bs))
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}
