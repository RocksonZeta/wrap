package cutil

import (
	"reflect"

	"github.com/RocksonZeta/wrap/utils/sutil"
)

func UniqueStrs(ss []string) []string {
	m := make(map[string]bool, len(ss))
	var r []string
	for _, v := range ss {
		if !m[v] {
			r = append(r, v)
		}
		m[v] = true
	}
	return r
}

func UniqueInts(vs []int) []int {
	m := make(map[int]bool)
	var r []int
	for _, v := range vs {
		if !m[v] {
			r = append(r, v)
		}
		m[v] = true
	}
	return r
}

func Filter(arr interface{}, fn func(i int, v interface{}) bool) interface{} {
	if arr == nil {
		return nil
	}
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	n := arrV.Len()
	j := 0
	result := reflect.MakeSlice(arrV.Type(), n, n)
	for i := 0; i < n; i++ {
		ele := arrV.Index(i)
		if fn(i, ele.Interface()) {
			result.Index(i).Set(ele)
			j++
		}
	}
	return result.Slice(0, j).Interface()
}

func Map(arr interface{}, result interface{}, fn func(i int, v interface{}) interface{}) {
	if arr == nil {
		return
	}
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	n := arrV.Len()
	resultSlice := sutil.FromSlice(result).Make(n)
	for i := 0; i < n; i++ {
		mapedValue := fn(i, arrV.Index(i).Interface())
		resultSlice.Put(i, mapedValue)
	}
}

//Any of fn(ele) return true , result is true
func Any(arr interface{}, fn func(i int, v interface{}) bool) bool {
	if arr == nil {
		return false
	}
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	n := arrV.Len()
	for i := 0; i < n; i++ {
		ele := arrV.Index(i)
		if fn(i, ele.Interface()) {
			return true
		}
	}
	return false
}

//All of fn(ele) return true , result is true other false
func All(arr interface{}, fn func(i int, v interface{}) bool) bool {
	if arr == nil {
		return false
	}
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	n := arrV.Len()
	for i := 0; i < n; i++ {
		ele := arrV.Index(i)
		if !fn(i, ele.Interface()) {
			return false
		}
	}
	return true
}
