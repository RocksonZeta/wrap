package sutil

import (
	"reflect"

	"github.com/RocksonZeta/wrap/errs"
	"github.com/RocksonZeta/wrap/wraplog"
)

var logSlice = wraplog.Logger.Fork(pkg, "Slice")

func checkSlice(err error, state int, msg string) error {
	if err == nil {
		return nil
	}
	if msg != "" {
		msg = err.Error()
	}

	panic(errs.Err{Err: err, Module: "sutil", Pkg: pkg, State: state, Message: msg})
}

type Slice struct {
	sliceType  reflect.Type
	sliceValue reflect.Value
	addr       interface{}
	Len        int
	Cap        int
}

//FromSlice to fill slice.  var a []int; FromSlice(*a).Make()
func FromSlice(arr interface{}) *Slice {
	sliceType := reflect.TypeOf(arr)
	if reflect.Ptr != sliceType.Kind() {
		msg := "FromSlice's arg should be address of slice."
		log.Error().Func("FromSlice").Stack().Err(errs.Err{State: ErrorType, Message: msg, Pkg: pkg, Module: "sutil"}).Msg(msg)
		check(nil, ErrorType, msg)
	}
	return &Slice{sliceType: reflect.TypeOf(arr), sliceValue: reflect.ValueOf(arr), addr: arr}
}

func (s *Slice) Make(length int) *Slice {
	s.sliceValue.Elem().Set(reflect.MakeSlice(s.sliceType.Elem(), length, length))
	s.Len = length
	s.Cap = length
	return s
}
func (s *Slice) Put(i int, value interface{}) {
	s.sliceValue.Elem().Index(i).Set(reflect.ValueOf(value))
}
func (s *Slice) Get(i int) interface{} {
	return s.sliceValue.Elem().Index(i).Interface()
}
func (s *Slice) GetRef(i int) interface{} {
	return s.sliceValue.Elem().Index(i).Addr().Interface()
}
func (s *Slice) Each(fn func(i int, v interface{})) {
	for i := 0; i < s.Len; i++ {
		fn(i, s.Get(i))
	}
}
