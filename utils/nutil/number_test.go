package nutil_test

import (
	"testing"

	"github.com/RocksonZeta/wrap/utils/nutil"
	"github.com/stretchr/testify/assert"
)

func TestInt64(t *testing.T) {
	r, err := nutil.Int64("12.2")
	assert.Nil(t, err)
	assert.Equal(t, int64(12), r)
}
func TestAdd(t *testing.T) {
	r := nutil.Add("12", 2)
	assert.Equal(t, "122", r)
	r1 := nutil.Add(1.1, 1)
	assert.Equal(t, 2.1, r1)
}
func TestMul(t *testing.T) {
	r := nutil.Mul("12", 1)
	assert.Equal(t, int64(12), r)
}
func TestSub(t *testing.T) {
	r := nutil.Sub("12", 1)
	assert.Equal(t, int64(11), r)
}
func TestDiv(t *testing.T) {
	r := nutil.Div("12", "1.2")
	assert.Equal(t, float64(10), r)
}
func TestString(t *testing.T) {
	r := nutil.String(1.2)
	assert.Equal(t, "1.2", r)

}
