package cutil_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/RocksonZeta/wrap/utils/cutil"
)

func TestUniqueStrs(t *testing.T) {
	r := cutil.UniqueStrs([]string{})
	fmt.Println(r)
}
func TestUniqueInts(t *testing.T) {
	r := cutil.UniqueInts([]int{1, 1})
	fmt.Println(r)
}
func TestFilter(t *testing.T) {
	arr := []int{1, 2, 3, 4}
	r := cutil.Filter(arr, func(i int, v interface{}) bool {
		return v.(int)%2 == 0
	})
	fmt.Println(r)
}
func TestMap(t *testing.T) {
	arr := []int{1, 2, 3, 4}
	var result []string
	cutil.Map(arr, &result, func(i int, v interface{}) interface{} {
		return "c" + strconv.Itoa(v.(int))
	})
	fmt.Println(result)
}
