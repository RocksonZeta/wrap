package stringutil_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/utils/stringutil"

	"github.com/stretchr/testify/assert"
)

func TestFileNameAppend(t *testing.T) {
	assert.Equal(t, "1_1.txt", stringutil.FileNameAppend("1.txt", "_1"))
}
func TestTemplate(t *testing.T) {
	r := stringutil.Template("hello {k1} ok{k2}", map[string]interface{}{"k1": "world", "k2": 1})
	fmt.Println(r)
	a := struct {
		A int
	}{
		12,
	}
	ra := stringutil.Template("hello {A} ok{k2}", a)
	fmt.Println(ra)
}
