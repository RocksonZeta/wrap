package cutil

import (
	"errors"
	"reflect"
	"sort"
	"time"

	"github.com/RocksonZeta/wrap/utils/sutil"
	"gopkg.in/guregu/null.v3"
)

//Group group []T -> map[field][]T
func Group(arr interface{}, m interface{}, field string) {
	GroupInto(arr, m, field, "")
	// mt := reflect.TypeOf(m)
	// arrV := reflect.ValueOf(arr)
	// arrT := reflect.TypeOf(arr)
	// r := reflect.MakeMap(mt.Elem())
	// mLen := make(map[interface{}]int)
	// iLen := make(map[interface{}]int)
	// for i := 0; i < arrV.Len(); i++ {
	// 	ele := arrV.Index(i)
	// 	key := ele.FieldByName(field).Interface()
	// 	if l, ok := mLen[key]; ok {
	// 		mLen[key] = l + 1
	// 	} else {
	// 		mLen[key] = 1
	// 		iLen[key] = 0
	// 	}
	// }
	// for i := 0; i < arrV.Len(); i++ {
	// 	ele := arrV.Index(i)
	// 	key := ele.FieldByName(field)
	// 	kv := key.Interface()
	// 	vs := r.MapIndex(key)
	// 	if !vs.IsValid() {
	// 		vs = reflect.MakeSlice(arrT, mLen[kv], mLen[kv])
	// 		vs.Index(0).Set(ele)
	// 		r.SetMapIndex(key, vs)
	// 		iLen[kv] = 1
	// 	} else {
	// 		vs.Index(iLen[kv]).Set(ele)
	// 		iLen[kv]++
	// 	}
	// }
	// reflect.ValueOf(m).Elem().Set(r)
}

//GroupOne []T -> map[field]T , m should be map ref *map
func GroupOne(arr interface{}, m interface{}, field string) {
	if arr == nil {
		return
	}
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	arrLen := arrV.Len()
	if 0 == arrLen {
		return
	}
	mt := reflect.TypeOf(m)
	r := reflect.MakeMap(mt.Elem())
	for i := 0; i < arrLen; i++ {
		ele := arrV.Index(i)
		r.SetMapIndex(ele.FieldByName(field), ele)
	}
	reflect.ValueOf(m).Elem().Set(r)
}
func GroupOneFunc(arr interface{}, m interface{}, cb func(i int) interface{}) {
	if arr == nil {
		return
	}
	mt := reflect.TypeOf(m)
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	r := reflect.MakeMap(mt.Elem())
	arrLen := arrV.Len()
	for i := 0; i < arrLen; i++ {
		ele := arrV.Index(i)
		r.SetMapIndex(reflect.ValueOf(cb(i)), ele)
	}
	reflect.ValueOf(m).Elem().Set(r)
}

func GroupInto(arr interface{}, m interface{}, field, intoField string) {
	if arr == nil {
		return
	}
	mt := reflect.TypeOf(m)
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	arrLen := arrV.Len()
	if arrLen <= 0 {
		return
	}
	// arrT := reflect.TypeOf(arr)
	// targetEleType.
	arrT := mt.Elem().Elem()
	r := reflect.MakeMap(mt.Elem())
	mLen := make(map[interface{}]int)
	iLen := make(map[interface{}]int)
	iKey := make(map[int]interface{})
	for i := 0; i < arrLen; i++ {
		ele := arrV.Index(i)
		key := ele.FieldByName(field).Interface()
		iKey[i] = key
		if l, ok := mLen[key]; ok {
			mLen[key] = l + 1
		} else {
			mLen[key] = 1
			iLen[key] = 0
		}
	}
	for i := 0; i < arrLen; i++ {
		ele := arrV.Index(i)
		key := ele.FieldByName(field)
		kv := iKey[i]
		vs := r.MapIndex(key)
		if !vs.IsValid() {
			vs = reflect.MakeSlice(arrT, mLen[kv], mLen[kv])
			if intoField != "" {
				vs.Index(0).FieldByName(intoField).Set(ele)
			} else {
				vs.Index(0).Set(ele)
			}
			r.SetMapIndex(key, vs)
			iLen[kv] = 1
		} else {
			if intoField != "" {
				vs.Index(iLen[kv]).FieldByName(intoField).Set(ele)
			} else {
				vs.Index(iLen[kv]).Set(ele)
			}
			iLen[kv]++
		}
	}
	reflect.ValueOf(m).Elem().Set(r)
}

func GroupFunc(arr interface{}, m interface{}, cb func(i int) interface{}) {
	GroupIntoFunc(arr, m, cb, "")
}
func GroupIntoFunc(arr interface{}, m interface{}, cb func(i int) interface{}, intoField string) {
	if arr == nil {
		return
	}
	mt := reflect.TypeOf(m)
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	arrLen := arrV.Len()
	// arrT := reflect.TypeOf(arr)
	// targetEleType.
	arrT := mt.Elem().Elem()
	r := reflect.MakeMap(mt.Elem())
	mLen := make(map[interface{}]int)
	iLen := make(map[interface{}]int)
	iKey := make(map[int]interface{}, arrLen)
	for i := 0; i < arrLen; i++ {
		// ele := arrV.Index(i)
		// key := ele.FieldByName(field).Interface()
		key := cb(i)
		iKey[i] = key
		if l, ok := mLen[key]; ok {
			mLen[key] = l + 1
		} else {
			mLen[key] = 1
			iLen[key] = 0
		}
	}
	for i := 0; i < arrV.Len(); i++ {
		ele := arrV.Index(i)
		// key := ele.FieldByName(field)
		// key := cb(i)
		key := iKey[i]
		// kv := key.Interface()
		keyValue := reflect.ValueOf(key)
		vs := r.MapIndex(keyValue)
		if !vs.IsValid() {
			vs = reflect.MakeSlice(arrT, mLen[key], mLen[key])
			if intoField != "" {
				vs.Index(0).FieldByName(intoField).Set(ele)
			} else {
				vs.Index(0).Set(ele)
			}
			r.SetMapIndex(keyValue, vs)
			iLen[key] = 1
		} else {
			if intoField != "" {
				vs.Index(iLen[key]).FieldByName(intoField).Set(ele)
			} else {
				vs.Index(iLen[key]).Set(ele)
			}
			iLen[key]++
		}
	}
	reflect.ValueOf(m).Elem().Set(r)
}

//Col []T->[]field
func Col(arr interface{}, result interface{}, field string) {
	arrV := sutil.PtrValue(reflect.ValueOf(arr))
	arrLen := arrV.Len()
	colType := reflect.TypeOf(result).Elem()
	col := reflect.MakeSlice(colType, arrLen, arrLen)
	for i := 0; i < arrLen; i++ {
		ele := arrV.Index(i)
		col.Index(i).Set(ele.FieldByName(field))
	}
	reflect.ValueOf(result).Elem().Set(col)
}
func ColInt2Str(arr interface{}, field string, fn func(id int) string) []string {
	ints := ColInt(arr, field)
	r := make([]string, len(ints))
	for i, v := range ints {
		r[i] = fn(v)
	}
	return r
}

//ColInt 获取结构体数组的整数列。arr:T[]
func ColInt(arr interface{}, field string) []int {
	var r []int
	Col(arr, &r, field)
	return r
}
func ColIntUnique(arr interface{}, field string) []int {
	l := ColInt(arr, field)
	m := make(map[int]bool, len(l))
	var r []int
	for _, v := range l {
		if !m[v] {
			r = append(r, v)
		}
	}
	return r
}

//ColInt 获取结构体数组的整数列,以map的形式返回。arr:T[]
func ColIntMap(arr interface{}, field string) map[int]bool {
	c := ColInt(arr, field)
	r := make(map[int]bool, len(c))
	for _, v := range c {
		r[v] = true
	}
	return r
}
func ColNullInt(arr interface{}, field string) []null.Int {
	var r []null.Int
	Col(arr, &r, field)
	return r
}
func ColNullIntValues(arr interface{}, field string) []int {
	ints := ColNullInt(arr, field)
	r := make([]int, len(ints))
	for i, v := range ints {
		if v.Valid {
			r[i] = int(v.Int64)
		}
	}
	return r
}
func ColNullIntValuesUnique(arr interface{}, field string) []int {
	ints := ColNullInt(arr, field)
	r := make([]int, len(ints))
	m := make(map[int64]bool, len(ints))
	for i, v := range ints {
		if v.Valid && !m[v.Int64] {
			r[i] = int(v.Int64)
		}
	}
	return r
}

//Index 找出元素位置 , arr:[]T
func Index(arr interface{}, field string, value interface{}) int {
	arrV := reflect.ValueOf(arr)
	arrLen := arrV.Len()
	for i := 0; i < arrLen; i++ {
		if arrV.Index(i).FieldByName(field).Interface() == value {
			return i
		}
	}
	return -1
}

//ChangeMapKey 修改字典的键。
func ChangeMapKey(m interface{}, r interface{}, field string) {
	if m == nil {
		return
	}
	// mT := reflect.TypeOf(m)
	mV := reflect.ValueOf(m)
	rT := reflect.TypeOf(r)
	rV := reflect.MakeMap(rT.Elem())
	mKeys := mV.MapKeys()
	for _, key := range mKeys {
		ele := mV.MapIndex(key)
		rV.SetMapIndex(ele.FieldByName(field), ele)
	}
	reflect.ValueOf(r).Elem().Set(rV)
}

//ChangeIndex move i->j ,
func ChangeIndex(arr interface{}, i, j int) {
	if i == j {
		return
	}
	if i > j {
		i, j = j, i
	}
	arrV := reflect.ValueOf(arr)
	lastEle := arrV.Index(j).Interface()
	for ; j > i; j-- {
		arrV.Index(j).Set(arrV.Index(j - 1))
	}
	arrV.Index(i).Set(reflect.ValueOf(lastEle))
}

//ToMapIndex  []int ->map[v]index
func ToMapIndex(ints []int) map[int]int {
	l := len(ints)
	r := make(map[int]int, l)

	for i := l - 1; i >= 0; i-- {
		r[ints[i]] = i
	}
	return r
}

//Sort 排序结构体数组 。arr:[]T,fieldAsc : field true|false,...
func Sort(arr interface{}, fieldAsc ...interface{}) {
	if nil == arr {
		return
	}
	arrV := reflect.ValueOf(arr)
	sort.Slice(arr, func(i, j int) bool {
		for k := 0; k < len(fieldAsc); k += 2 {
			asc := true
			if k < len(fieldAsc)-1 {
				asc = fieldAsc[k+1].(bool)
			}
			field := fieldAsc[k].(string)
			ei := arrV.Index(i).FieldByName(field).Interface()
			ej := arrV.Index(j).FieldByName(field).Interface()
			if ei == ej {
				continue
			}
			le, err := LE(ei, ej)
			if err != nil {
				panic(err)
			}
			return asc && le || !asc && !le
		}
		return false
	})
}
func SortLikeInts(arr interface{}, field string, orderedInts []int) {
	if nil == arr {
		return
	}
	im := IntsIndexMap(orderedInts)
	arrV := reflect.ValueOf(arr)
	sort.Slice(arr, func(i, j int) bool {
		ei := arrV.Index(i).FieldByName(field).Interface().(int)
		ej := arrV.Index(j).FieldByName(field).Interface().(int)
		i1, ok1 := im[ei]
		i2, ok2 := im[ej]
		if !ok1 {
			return false
		}
		if !ok2 {
			return true
		}
		return i1 < i2
	})

}

//arr -> [value]index
func IntsIndexMap(arr []int) map[int]int {
	r := make(map[int]int, len(arr))
	for i, v := range arr {
		r[v] = i
	}
	return r
}

//Gather 通过arr中结构体的的某个域，按顺序挑选
func Gather(arr interface{}, field string, ids []int) interface{} {
	n := len(ids)
	arrV := reflect.ValueOf(arr)
	arrLen := arrV.Len()
	rV := reflect.MakeSlice(arrV.Type(), n, n)
	m := make(map[int]int, n)
	for i, v := range ids {
		m[v] = i
	}
	for i := 0; i < arrLen; i++ {
		f := arrV.Index(i)
		v := int(f.FieldByName(field).Int())
		if j, ok := m[v]; ok {
			rV.Index(j).Set(f)
		}
	}
	return rV.Interface()
}

//LE a<=b :true ,else false
func LE(a, b interface{}) (bool, error) {
	switch a.(type) {
	case int:
		return a.(int) <= b.(int), nil
	case int64:
		return a.(int64) <= b.(int64), nil
	case float32:
		return a.(float32) <= b.(float32), nil
	case float64:
		return a.(float64) <= b.(float64), nil
	case string:
		return a.(string) <= b.(string), nil
	case byte:
		return a.(byte) <= b.(byte), nil
	case bool:
		return !a.(bool) || !b.(bool), nil
	case int32:
		return a.(int32) <= b.(int32), nil
	case time.Time:
		return a.(time.Time).Before(b.(time.Time)), nil
	case null.Int:
		return a.(null.Int).Int64 <= b.(null.Int).Int64, nil
	case null.String:
		return a.(null.String).String <= b.(null.String).String, nil
	}
	return false, errors.New("structutil.LE not support this type compare:" + reflect.TypeOf(a).String())
}

func CopySlice(src interface{}, dst interface{}, fields ...string) {
	// srcV := reflect.ValueOf(src)
	// srcLen := srcV.Len()
	// dstT := reflect.TypeOf(dst)
	// dstEleT := dstT.Elem().Elem()
	// for i := 0; i < srcLen; i++ {
	// 	dstEle := reflect.New(dstEleT)
	// }

}

func Unique(ids []int) []int {
	m := make(map[int]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	r := make([]int, len(m))
	i := 0
	for id := range m {
		r[i] = id
		i++
	}
	return r
}

func IntsNotInArray(ids []int, arr interface{}, field string) []int {
	arrV := reflect.ValueOf(arr)
	arrLen := arrV.Len()
	arrMap := make(map[int]bool, arrLen)
	var r []int
	for i := 0; i < arrLen; i++ {
		arrMap[int(arrV.Index(i).FieldByName(field).Int())] = true
	}
	for _, id := range ids {
		if _, ok := arrMap[id]; !ok {
			r = append(r, id)
		}
	}
	return r
}
func MapIds(ids []int, fn func(id int) string) []string {
	r := make([]string, len(ids))
	for i, id := range ids {
		r[i] = fn(id)
	}
	return r
}

// func Map(arr interface{}, r interface{}, fn func(i int) interface{}) {
// 	arrV := reflect.ValueOf(arr)
// 	arrLen := arrV.Len()
// 	rT := reflect.TypeOf(r).Elem()
// 	rV := reflect.MakeSlice(rT, arrLen, arrLen)
// 	for i := 0; i < arrLen; i++ {
// 		rV.Index(i).Set(reflect.ValueOf(fn(i)))
// 	}
// 	reflect.ValueOf(r).Elem().Set(rV)
// }
// func MapString(arr interface{}, fn func(i int) interface{}) []string {
// 	var r []string
// 	Map(arr, &r, fn)
// 	return r
// }
func Reduce(arr interface{}, fn func(i int) interface{}) {
	arrV := reflect.ValueOf(arr)
	arrLen := arrV.Len()
	for i := 0; i < arrLen; i++ {
		fn(i)
	}
}

func MapKeysInt(m interface{}) []int {
	var r []int
	MapKeys(m, &r)
	return r
}
func MapKeysString(m interface{}) []string {
	var r []string
	MapKeys(m, &r)
	return r
}
func MapKeys(m interface{}, r interface{}) {
	if m == nil {
		return
	}
	// mT := reflect.TypeOf(m)
	mV := reflect.ValueOf(m)
	rT := reflect.TypeOf(r)
	mLen := mV.Len()
	rV := reflect.MakeSlice(rT.Elem(), mLen, mLen)
	mKeys := mV.MapKeys()
	for i, key := range mKeys {
		rV.Index(i).Set(key)
	}
	reflect.ValueOf(r).Elem().Set(rV)
}
func MapValues(m interface{}, r interface{}) {
	if m == nil {
		return
	}
	// mT := reflect.TypeOf(m)
	mV := reflect.ValueOf(m)
	rT := reflect.TypeOf(r)
	mLen := mV.Len()
	rV := reflect.MakeSlice(rT.Elem(), mLen, mLen)
	mKeys := mV.MapKeys()
	for i, key := range mKeys {
		rV.Index(i).Set(mV.MapIndex(key))
	}
	reflect.ValueOf(r).Elem().Set(rV)
}

func Join(a interface{}, b interface{}, r interface{}) {
	aV := sutil.PtrValue(reflect.ValueOf(a))
	bV := sutil.PtrValue(reflect.ValueOf(b))
	rV := sutil.PtrValue(reflect.ValueOf(r))
	aLen := aV.Len()
	bLen := bV.Len()
	rLen := aLen + bLen
	s := reflect.MakeSlice(rV.Type(), aLen+bLen, rLen)
	i := 0
	for ; i < aLen; i++ {
		s.Index(i).Set(aV.Index(i))
	}
	for ; i < aLen+bLen; i++ {
		s.Index(i).Set(bV.Index(i - aLen))
	}
	rV.Set(s)
}
