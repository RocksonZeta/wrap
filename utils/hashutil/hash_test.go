package hashutil_test
import (
	"testing"

	"github.com/RocksonZeta/wrap/utils/hashutil"
	"github.com/bmizerany/assert"
)

func TestRandomStr(t *testing.T) {
	md5 := hashutil.Md5("111111")
	assert.Equal(t, "96e79218965eb72c92a549dd5a330112", md5)
}
