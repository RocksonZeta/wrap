package nutil

import "reflect"

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
func (v Value) IsString() bool {
	return IsString(v.V)
}
func (v Value) IsNumber() bool {
	return IsNumber(v.V)
}
func (v Value) IsFloat() bool {
	return IsFloat(v.V)
}
func (v Value) IsInteger() bool {
	return IsInteger(v.V)
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
