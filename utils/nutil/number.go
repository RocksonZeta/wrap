package nutil

import (
	"errors"
	"fmt"
	"strconv"
)

//Add a+b
func Add(a, b interface{}) interface{} {
	if IsString(a) || IsString(b) {
		return String(a) + String(b)
	}
	if IsInt(a) && IsInt(b) {
		return Int64Must(a) + Int64Must(b)
	}
	return Float64Must(a) + Float64Must(b)
}

//AddInt a+b
func AddInt(a, b interface{}) int {
	return IntMust(a) + IntMust(b)
}

//Sub  a-b
func Sub(a, b interface{}) interface{} {
	if IsInt(a) && IsInt(b) {
		return Int64Must(a) - Int64Must(b)
	}
	return Float64Must(a) - Float64Must(b)
}

//Mul  a*b
func Mul(a, b interface{}) interface{} {
	if IsInt(a) && IsInt(b) {
		return Int64Must(a) * Int64Must(b)
	}
	return Float64Must(a) * Float64Must(b)
}

//Div  a/b
func Div(a, b interface{}) interface{} {
	if IsInt(a) && IsInt(b) {
		return Int64Must(a) / Int64Must(b)
	}
	return Float64Must(a) / Float64Must(b)
}

//Div  a/b
func DivFloat(a, b interface{}) float64 {
	return Float64Must(a) / Float64Must(b)
}

//Mod  a%b
func Mod(a, b interface{}) interface{} {
	return Int64Must(a) % Int64Must(b)
}

func IsInt(v interface{}) bool {
	switch n := v.(type) {
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case int:
		return true
	case uint:
		return true
	case uintptr:
		return true
	case string:
		_, err := strconv.ParseInt(n, 10, 64)
		return err == nil
	}
	return false
}
func IsFloat(v interface{}) bool {
	switch n := v.(type) {
	case float32:
		return true
	case float64:
		return true
	case string:
		_, err := strconv.ParseFloat(n, 64)
		return err == nil
	}
	return false
}
func IsString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}
func String(v interface{}) string {
	if a, ok := v.(string); ok {
		return a
	}
	return fmt.Sprintf("%v", v)
}
func Int64Must(v interface{}) int64 {
	a, _ := Int64(v)
	return a
}
func Float64Must(v interface{}) float64 {
	a, _ := Float64(v)
	return a
}
func IntMust(v interface{}) int {
	return int(Int64Must(v))
}
func Int(v interface{}) (int, error) {
	a, err := Int64(v)
	return int(a), err
}
func Float32Must(v interface{}) float32 {
	return float32(Float64Must(v))
}
func Float32(v interface{}) (float32, error) {
	a, err := Float64(v)
	return float32(a), err
}
func Int64(v interface{}) (int64, error) {
	switch n := v.(type) {
	case bool:
		if n {
			return 1, nil
		}
		return 0, nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return int64(n), nil
	case int:
		return int64(n), nil
	case uint:
		return int64(n), nil
	case uintptr:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case float64:
		return int64(n), nil
	case string:
		a, err := strconv.ParseFloat(n, 64)
		return int64(a), err
	}
	return 0, errors.New("bad integer format.")
}
func Float64(v interface{}) (float64, error) {
	switch n := v.(type) {
	case bool:
		if n {
			return 1, nil
		}
		return 0, nil
	case uint8:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil

	case int:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uintptr:
		return float64(n), nil
	case float32:
		return float64(n), nil
	case float64:
		return float64(n), nil
	case string:
		return strconv.ParseFloat(n, 64)
	}
	return 0, errors.New("bad float format.")
}
func IsNumber(v interface{}) bool {
	switch n := v.(type) {
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case int:
		return true
	case uint:
		return true
	case uintptr:
		return true
	case float32:
		return true
	case float64:
		return true
	case string:
		_, err := strconv.ParseFloat(n, 64)
		return err == nil
	}
	return false
}
