package osswrap_test

import (
	"testing"

	"github.com/RocksonZeta/wrap/osswrap"
	"github.com/stretchr/testify/suite"
)

type OssSuite struct {
	suite.Suite
	client *osswrap.Oss
}

func (s *OssSuite) SetupTest() {
	var options osswrap.Options
	s.client = osswrap.New(options, "sy-video-test")
}

func (s *OssSuite) TestPut() {
	s.client.PutFile("oss.go", "/test/wrap/osswrap/oss.go")
}
func (s *OssSuite) TestPutDir() {
	s.client.PutDir(".", "/test/wrap/osswrap", nil)
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(OssSuite))
}
