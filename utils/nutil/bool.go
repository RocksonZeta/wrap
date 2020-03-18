package nutil

import (
	"fmt"
	"reflect"
)

// convert value to bool type
func Bool(v interface{}) bool {
	fmt.Println("v", v, v == nil)
	if v == nil {
		return false
	}
	switch n := v.(type) {
	case uint8:
		return n != 0
	case uint16:
		return n != 0
	case uint32:
		return n != 0
	case uint64:
		return n != 0
	case int8:
		return n != 0
	case int16:
		return n != 0
	case int32:
		return n != 0
	case int64:
		return n != 0
	case int:
		return n != 0
	case uint:
		return n != 0
	case uintptr:
		return n != 0
	case float32:
		return n != 0
	case float64:
		return n != 0
	case string:
		return n != ""
	case complex64:
		return real(n) != 0 || imag(n) != 0
	case complex128:
		return real(n) != 0 || imag(n) != 0
	}
	vV := reflect.ValueOf(v)
	switch vV.Kind() {
	case reflect.Map:
		return vV.Len() > 0
	case reflect.Slice:
		return vV.Len() > 0
	case reflect.Array:
		return vV.Len() > 0
	case reflect.Ptr:
		return !vV.IsNil()
	}
	return true
}
