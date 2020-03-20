package sutil

import (
	"reflect"

	"github.com/fatih/structs"
)

func FieldNames(t reflect.Type) []string {
	n := t.NumField()
	r := make([]string, n)
	for i := 0; i < n; i++ {
		r[i] = t.Field(i).Name
	}
	return r
}
func Copy(src, dst interface{}) {
	s := structs.New(src)
	d := structs.New(dst)
	sm := make(map[string]*structs.Field)
	for _, f := range s.Fields() {
		sm[f.Name()] = f
	}
	for _, d := range d.Fields() {
		if f, ok := sm[d.Name()]; ok {
			d.Set(f.Value())
		}
	}
}

func Get(v interface{}, field string) interface{} {
	vV := PtrValue(reflect.ValueOf(v))
	return vV.FieldByName(field).Interface()
}
func PtrType(t reflect.Type) reflect.Type {
	if reflect.Ptr == t.Kind() {
		return PtrType(t.Elem())
	}
	return t
}
func PtrValue(v reflect.Value) reflect.Value {
	if reflect.Ptr == v.Kind() {
		return PtrValue(v.Elem())
	}
	return v
}
func Ptr(a interface{}) interface{} {
	return PtrValue(reflect.ValueOf(a)).Interface()
}

// func Struct2Map(a interface{}) map[string]interface{} {
// 	return structs.Map(a)
// }

func MapGet(m interface{}, key interface{}) interface{} {
	mV := reflect.ValueOf(m)
	r := mV.MapIndex(reflect.ValueOf(key))
	if !r.IsValid() {
		return nil
	}
	return r.Interface()
}
