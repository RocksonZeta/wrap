package sutil

import (
	"encoding/json"
	"reflect"

	"github.com/RocksonZeta/wrap/errs"
	"github.com/RocksonZeta/wrap/wraplog"
)

const (
	ErrorType = 1
)

var log = wraplog.Logger.Fork(pkg, "Map")
var pkg = reflect.TypeOf(Map{}).PkgPath()

func check(err error, state int, msg string) error {
	if err == nil {
		return nil
	}
	if msg != "" {
		msg = err.Error()
	}

	panic(errs.Err{Err: err, Module: "sutil", Pkg: pkg, State: state, Message: msg})
}

//Map like map[string]User
type Map struct {
	mapType  reflect.Type
	mapValue reflect.Value
	addr     interface{}
}

// NewMap useage :
// var m map[string]User
// NewMap(&m)
func FromMap(emptyMapAddr interface{}) *Map {
	mapType := reflect.TypeOf(emptyMapAddr)
	if reflect.Ptr != mapType.Kind() {
		msg := "FromMap's arg should be address of map."
		log.Error().Func("FromMap").Stack().Err(errs.Err{State: ErrorType, Message: msg, Pkg: pkg, Module: "sutil"}).Msg(msg)
		check(nil, ErrorType, msg)
	}
	return &Map{mapType: reflect.TypeOf(emptyMapAddr), mapValue: reflect.ValueOf(emptyMapAddr), addr: emptyMapAddr}
}

func (m *Map) Make() *Map {
	m.mapValue.Elem().Set(reflect.MakeMap(m.mapType.Elem()))
	return m
}
func (m *Map) MakeWithSize(size int) *Map {
	m.mapValue.Elem().Set(reflect.MakeMapWithSize(m.mapType.Elem(), size))
	return m
}
func (m *Map) Put(key, value interface{}) {
	m.mapValue.Elem().SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
}
func (m *Map) Get(key interface{}) interface{} {
	return m.mapValue.Elem().MapIndex(reflect.ValueOf(key)).Interface()
}
func (m *Map) NewValue() interface{} {
	return reflect.New(m.mapType.Elem().Elem()).Interface()
}
func (m *Map) EmptyValue() interface{} {
	return reflect.New(m.mapType.Elem().Elem()).Elem().Interface()
}
func (m *Map) Each(fn func(k, v interface{})) {
	ele := m.mapValue.Elem()
	for _, key := range ele.MapKeys() {
		fn(key.Interface(), ele.MapIndex(key).Interface())
	}
}

func (m *Map) PutJson(key interface{}, jsonBs []byte) error {
	valueT := m.mapType.Elem().Elem()
	valueV := reflect.New(valueT)
	value := valueV.Interface()
	err := json.Unmarshal(jsonBs, &value)
	if err != nil {
		return err
	}
	m.mapValue.Elem().SetMapIndex(reflect.ValueOf(key), valueV.Elem())
	return nil
}
