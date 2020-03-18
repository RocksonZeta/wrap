package mathutil_test

import (
	"testing"

	"github.com/RocksonZeta/wrap/utils/mathutil"
	"github.com/bmizerany/assert"
)

func TestRandomStr(t *testing.T) {
	str := mathutil.RandomStr(32, true)
	assert.Equal(t, 32, len(str))
}
