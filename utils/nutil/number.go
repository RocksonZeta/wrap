package nutil

import (
	"errors"
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
func IsInteger(v interface{}) bool {
	switch v.(type) {
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
	}
	return false
}
func IsString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}
func String(v interface{}) string {

	switch e := v.(type) {
	case string:
		return e
	case bool:
		if e {
			return "1"
		} else {
			return "0"
		}
	case int:
		return strconv.Itoa(e)
	case byte:
		return strconv.FormatInt(int64(e), 10)
	case int16:
		return strconv.FormatInt(int64(e), 10)
	case int32:
		return strconv.FormatInt(int64(e), 10)
	case int64:
		return strconv.FormatInt(e, 10)
	case uint:
		return strconv.FormatUint(uint64(e), 10)
	case uintptr:
		return strconv.FormatUint(uint64(e), 10)
	case uint16:
		return strconv.FormatUint(uint64(e), 10)
	case uint32:
		return strconv.FormatUint(uint64(e), 10)
	case uint64:
		return strconv.FormatUint(e, 10)
	case float32:
		return strconv.FormatFloat(float64(e), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(e, 'f', -1, 64)
	}
	return ""
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
		return strconv.ParseInt(n, 10, 64)
	}
	return 0, errors.New("bad int format")
}
func Uint64(v interface{}) (uint64, error) {
	switch n := v.(type) {
	case bool:
		if n {
			return 1, nil
		}
		return 0, nil
	case uint8:
		return uint64(n), nil
	case uint16:
		return uint64(n), nil
	case uint32:
		return uint64(n), nil
	case uint64:
		return uint64(n), nil
	case int8:
		return uint64(n), nil
	case int16:
		return uint64(n), nil
	case int32:
		return uint64(n), nil
	case int64:
		return uint64(n), nil
	case int:
		return uint64(n), nil
	case uint:
		return uint64(n), nil
	case uintptr:
		return uint64(n), nil
	case float32:
		return uint64(n), nil
	case float64:
		return uint64(n), nil
	case string:
		return strconv.ParseUint(n, 10, 64)
	}
	return 0, errors.New("bad uint format")
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
	return 0, errors.New("bad float format")
}
func IsNumberable(v interface{}) bool {
	if IsNumber(v) || IsNumberStr(v) {
		return true
	}
	return false
}
func IsNumberStr(v interface{}) bool {
	if n, ok := v.(string); ok {
		_, err := strconv.ParseFloat(n, 64)
		return err == nil
	}
	return false
}
func IsNumber(v interface{}) bool {
	switch v.(type) {
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

	}
	return false
}
