package fileutil_test

import (
	"fmt"
	"testing"

	"github.com/RocksonZeta/wrap/utils/fileutil"
)

func TestFindFileDir(t *testing.T) {
	dir := fileutil.FindFileDir("main.go")
	fmt.Println(dir)

}
