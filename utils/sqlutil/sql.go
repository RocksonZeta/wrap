package sqlutil

import (
	"reflect"
	"strconv"
	"strings"
)

func JoinInts(is []int) string {
	s := ""
	for _, iv := range is {
		s += strconv.Itoa(iv)
		s += ","
	}
	if "" == s {
		return "null"
	}
	return s[0 : len(s)-1]
}

func JoinInt64s(is []int64) string {
	s := ""
	for _, iv := range is {
		s += strconv.FormatInt(iv, 10)
		s += ","
	}
	if "" == s {
		return "null"
	}
	return s[0 : len(s)-1]
}
func IntArgs(is []int) []interface{} {
	r := make([]interface{}, len(is))
	for i, v := range is {
		r[i] = v
	}
	return r
}
func Args(args interface{}) []interface{} {
	argsV := reflect.ValueOf(args)
	argsLen := argsV.Len()
	r := make([]interface{}, argsLen)
	for i := 0; i < argsLen; i++ {
		r[i] = argsV.Index(i).Interface()
	}
	return r
}
func ArgsRef(args interface{}) []interface{} {
	argsV := reflect.ValueOf(args)
	argsLen := argsV.Len()
	r := make([]interface{}, argsLen)
	for i := 0; i < argsLen; i++ {
		ele := argsV.Index(i)
		if reflect.Ptr == ele.Kind() {
			r[i] = ele
		} else {
			r[i] = ele.Addr().Interface()
		}
	}
	return r
}
func trim(s, cut string) string {
	sl, cl := len(s), len(cut)
	if sl < cl {
		return s
	}
	if s[0:cl] == cut {
		return s[cl:]
	}
	return s
}
func JoinFields(prefix string, fields ...string) string {
	r := ""
	l := len(fields)
	j := strings.Index(prefix, ".")
	p := prefix
	q := ""
	if j != len(prefix)-1 {
		p = prefix[0 : j+1]
		q = prefix[j+1:]
	}
	for i, f := range fields {
		r += p + "`" + trim(f, q) + "`"
		if q != "" {
			r += " as " + f
		}
		if i != l-1 {
			r += ","
		}
	}
	return r
}
func Fields(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	l := t.NumField()
	r := make([]string, l)
	for i := 0; i < l; i++ {
		r[i] = t.Field(i).Name
	}
	return r
}

func SqlFields(prefix string, obj interface{}) string {
	return JoinFields(prefix, Fields(obj)...)
}
