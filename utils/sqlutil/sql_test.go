package sqlutil_test

import (
	"testing"

	"github.com/RocksonZeta/wrap/utils/sqlutil"

	"github.com/stretchr/testify/assert"
)

func TestJoinFields(t *testing.T) {
	r := sqlutil.JoinFields("u.", "Id", "Name")
	assert.Equal(t, "u.`Id`,u.`Name`", r)
}
func TestSqlFields(t *testing.T) {
	type U1 struct{}
	type User struct {
		Id   int
		Name string
	}
	u := User{}
	r := sqlutil.SqlFields("u.User", u)
	assert.Equal(t, "u.`Id` as Id,u.`Name` as Name", r)
}
