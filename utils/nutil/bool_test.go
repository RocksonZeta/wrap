package nutil_test

import (
	"testing"

	"github.com/RocksonZeta/wrap/utils/nutil"
	"github.com/stretchr/testify/assert"
)

func TestFalse(t *testing.T) {
	var a *int
	var i interface{}
	assert.False(t, nutil.Bool(a))
	assert.False(t, nutil.Bool(i))
	assert.False(t, nutil.Bool(nil))
	assert.False(t, nutil.Bool(0))
	assert.False(t, nutil.Bool(""))
	assert.False(t, nutil.Bool([]int{}))
	assert.False(t, nutil.Bool(map[int]int{}))
}
