package nutil

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

func ValueOf(v interface{}) Value {
	return Value{V: v}
}

type Value struct {
	V interface{}
}

func (v Value) Float64() (float64, error) {
	return Float64(v.V)
}
func (v Value) Float32() (float32, error) {
	x, err := v.Float64()
	return float32(x), err
}
func (v Value) Int64() (int64, error) {
	return Int64(v.V)
}
func (v Value) Int32() (int32, error) {
	x, err := v.Int64()
	return int32(x), err
}
func (v Value) Int() (int, error) {
	x, err := v.Int64()
	return int(x), err
}
func (v Value) Int16() (int16, error) {
	x, err := v.Int64()
	return int16(x), err
}
func (v Value) Int8() (int8, error) {
	x, err := v.Int64()
	return int8(x), err
}
func (v Value) Byte() (byte, error) {
	x, err := v.Int64()
	return byte(x), err
}
func (v Value) Uint64() (uint64, error) {
	return Uint64(v.V)
}
func (v Value) Uint32() (uint32, error) {
	x, err := v.Int64()
	return uint32(x), err
}
func (v Value) Uint() (uint, error) {
	x, err := v.Int64()
	return uint(x), err
}
func (v Value) Uint16() (uint16, error) {
	x, err := v.Int64()
	return uint16(x), err
}
func (v Value) Uint8() (uint8, error) {
	x, err := v.Uint64()
	return uint8(x), err
}
func (v Value) Bool() bool {
	if IsInt(v.V) {
		return v.Int64Must() > 0
	}
	if IsFloat(v.V) {
		return v.Float64Must() > 0
	}
	if IsString(v.V) {
		return v.String() != ""
	}
	return v.V != nil
}

func (v Value) String() string {
	return String(v.V)
}

func (v Value) Float64Must() float64 {
	x, err := Float64(v.V)
	if err != nil {
		panic(err)
	}
	return x
}
func (v Value) Float32Must() float32 {
	return float32(v.Float64Must())
}
func (v Value) Int64Must() int64 {
	x, err := Int64(v.V)
	if err != nil {
		panic(err)
	}
	return x
}
func (v Value) Int32Must() int32 {
	return int32(v.Int64Must())
}
func (v Value) IntMust() int {
	return int(v.Int64Must())
}
func (v Value) Int16Must() int16 {
	return int16(v.Int64Must())
}
func (v Value) Int8Must() int8 {
	return int8(v.Int64Must())
}
func (v Value) ByteMust() byte {
	return byte(v.Int64Must())
}
func (v Value) Uint64Must() uint64 {
	x, err := Uint64(v.V)
	if err != nil {
		panic(err)
	}
	return x
}
func (v Value) Uint32Must() uint32 {
	return uint32(v.Uint64Must())
}
func (v Value) UintMust() uint {
	return uint(v.Uint64Must())
}
func (v Value) Uint16Must() uint16 {
	return uint16(v.Uint64Must())
}
func (v Value) Uint8Must() uint8 {
	return uint8(v.Uint64Must())
}

func (v Value) AsFloat64() float64 {
	x, _ := Float64(v.V)
	return x
}
func (v Value) AsInt64() int64 {
	return int64(v.AsFloat64())
}
func (v Value) AsInt32() int32 {
	return int32(v.AsFloat64())
}
func (v Value) AsInt() int {
	return int(v.AsFloat64())
}
func (v Value) AsInt16() int16 {
	return int16(v.AsFloat64())
}
func (v Value) AsInt8() int8 {
	return int8(v.AsFloat64())
}
func (v Value) AsByte() byte {
	return byte(v.AsFloat64())
}
func (v Value) AsMap() Map {
	switch x := v.V.(type) {
	case map[string]interface{}:
		return x
	case Map:
		return x
	}
	if v.IsMap() {
		vr := reflect.ValueOf(v.V)
		r := NewMapLen(vr.Len())
		iter := vr.MapRange()
		for iter.Next() {
			r[ValueOf(iter.Key().Interface()).String()] = iter.Value().Interface()
		}
		return r
	}
	if v.IsSlice() {
		vr := reflect.ValueOf(v.V)
		n := vr.Len()
		r := NewMapLen(n)
		for i := 0; i < n; i++ {
			r[strconv.Itoa(i)] = vr.Index(i).Interface()
		}
		return r
	}
	if v.IsStruct() {
		return structs.New(v.V).Map()
	}
	return nil
}

func (v Value) AsList() List {
	switch v.V.(type) {
	case List, []interface{}:
		return v.V.(List)
	}
	if v.IsSlice() || v.IsArray() {
		vr := reflect.ValueOf(v.V)
		n := vr.Len()
		r := NewList(n)
		for i := 0; i < n; i++ {
			r[i] = vr.Index(i).Interface()
		}
		return r
	}
	return nil
}

func (v Value) IsString() bool {
	return IsString(v.V)
}
func (v Value) IsNumber() bool {
	return IsNumber(v.V)
}
func (v Value) IsNumberStr() bool {
	return IsNumberStr(v.V)
}
func (v Value) IsNumberable() bool {
	return IsNumberable(v.V)
}
func (v Value) IsFloat() bool {
	return IsFloat(v.V)
}
func (v Value) IsInteger() bool {
	return IsInteger(v.V)
}

func (v Value) IsArray() bool {
	return reflect.TypeOf(v.V).Kind().String() == "array"
}
func (v Value) IsSlice() bool {
	return reflect.TypeOf(v.V).Kind().String() == "slice"
}
func (v Value) IsMap() bool {
	return reflect.TypeOf(v.V).Kind().String() == "map"
}
func (v Value) IsStruct() bool {
	return reflect.TypeOf(v.V).Kind().String() == "struct"
}
func (v Value) IsPtr() bool {
	return reflect.TypeOf(v.V).Kind().String() == "ptr"
}

//Compare 0:v ==a ,1:v>a,-1:v<a
func (v Value) Compare(a Value) int {
	if v.IsString() || a.IsString() {
		return strings.Compare(v.String(), a.String())
	}
	v1, err1 := v.Float64()
	v2, err2 := a.Float64()
	if err1 == nil && err2 == nil {
		r := v1 - v2
		if r == 0 {
			return 0
		} else if r > 0 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

func (v Value) MarshalJSON() (string, error) {
	bs, err := json.Marshal(v.V)
	return string(bs), err
}
func (v Value) MarshalJSONMust() string {
	r, _ := v.MarshalJSON()
	return r
}
