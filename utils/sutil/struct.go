package sutil

import (
	"reflect"

	"github.com/fatih/structs"
)

type M map[string]interface{}
type A []interface{}
type MS []M

//New struct
func New(s interface{}) *structs.Struct {
	return structs.New(s)
}

//NewSlice structList : []struct
func NewSlice(structList interface{}) []*structs.Struct {
	sV := reflect.ValueOf(structList)
	l := sV.Len()
	r := make([]*structs.Struct, l)
	for i := 0; i < l; i++ {
		r[i] = New(sV.Index(i).Interface())
	}
	return r
}

//NewMap m : map[?]struct
func NewMap(m interface{}) map[interface{}]*structs.Struct {
	sV := reflect.ValueOf(m)
	keys := sV.MapKeys()
	r := make(map[interface{}]*structs.Struct, len(keys))
	for _, v := range keys {
		r[v] = New(sV.MapIndex(v).Interface())
	}
	return r
}

func Struct2Map(s interface{}) map[string]interface{} {
	switch t := s.(type) {
	case map[string]interface{}:
		return t
	case M:
		return t
	}
	return structs.New(s).Map()
}
func Pick(s interface{}, keys ...string) map[string]interface{} {
	m := Struct2Map(s)
	if len(keys) == 0 {
		return m
	}
	r := make(map[string]interface{})
	for _, k := range keys {
		if v, ok := m[k]; ok {
			r[k] = v
		}
	}
	return r
}
func Unpick(s interface{}, keys ...string) map[string]interface{} {
	m := Struct2Map(s)
	if len(keys) == 0 {
		return m
	}
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}
	r := make(map[string]interface{})
	for k := range m {
		if !keyMap[k] {
			r[k] = m[k]
		}
	}
	return r
}
func Kv2Map(s ...interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	for i := 0; i < len(s); i += 2 {
		if i >= len(s)-1 {
			r[s[i].(string)] = nil
			return r
		}
		k := s[i].(string)
		if i+1 < len(s) {
			r[k] = s[i+1]
		}
	}
	return r
}
