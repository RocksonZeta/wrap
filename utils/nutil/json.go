package nutil

import (
	"encoding/json"
	"sort"

	"github.com/fatih/structs"
)

type Map map[string]interface{}

type List []interface{}

func NewMap() Map {
	return make(Map)
}
func NewMapLen(l int) Map {
	return make(Map, l)
}
func NewList(l int) List {
	return make(List, l)
}

func (m Map) Get(key string) Value {
	return ValueOf(m[key])
}
func (m Map) Has(key string) bool {
	_, ok := m[key]
	return ok
}
func (m Map) Len() int {
	return len(m)
}
func (m Map) IsNil() bool {
	return m == nil
}
func (m Map) IsEmpty() bool {
	return m == nil || len(m) == 0
}
func (m Map) Keys() []string {
	r := make([]string, len(m))
	i := 0
	for k := range m {
		r[i] = k
		i++
	}
	return r
}
func (m Map) Values() []interface{} {
	r := make([]interface{}, len(m))
	i := 0
	for _, v := range m {
		r[i] = v
		i++
	}
	return r
}
func (m Map) Struct(s interface{}) error {
	obj := structs.New(s)
	var err error
	for k, v := range m {
		f, ok := obj.FieldOk(k)
		if !ok {
			continue
		}
		err = f.Set(v)
	}
	return err
}

// List

//Get by index
func (l List) Get(i int) Value {
	return ValueOf(l[i])
}
func (l List) Len() int {
	return len(l)
}
func (l List) IsNil() bool {
	return l == nil
}
func (l List) IsEmpty() bool {
	return l == nil || len(l) == 0
}
func (l List) Sort(asc bool) {
	l.SortBy(func(i, j int) bool {
		r := l.Get(i).Compare(l.Get(j)) < 0
		if asc {
			return r
		}
		return !r
	})
}
func (l List) SortBy(less func(i, j int) bool) {
	sort.Slice(l, less)
}
func (l List) SortByField(field string, asc bool) {
	l.SortBy(func(i, j int) bool {
		v1 := l.Get(i).AsMap().Get(field)
		v2 := l.Get(j).AsMap().Get(field)
		r := v1.Compare(v2) < 0
		if asc {
			return r
		}
		return !r
	})
}

type SortOrder struct {
	Field string
	Asc   bool
}

func (l List) SortByFields(sortOrders []SortOrder) {
	l.SortBy(func(i, j int) bool {
		for _, s := range sortOrders {
			v1 := l.Get(i).AsMap().Get(s.Field)
			v2 := l.Get(j).AsMap().Get(s.Field)
			r := v1.Compare(v2)
			if r == 0 {
				continue
			}
			if s.Asc {
				return r < 0
			}
			return r > 0
		}
		return false
	})
}

func (m Map) MarshalJSON() (string, error) {
	bs, err := json.Marshal(m)
	return string(bs), err
}
func (m Map) MarshalJSONMust() string {
	r, _ := m.MarshalJSON()
	return r
}

func (l List) MarshalJSON() (string, error) {
	bs, err := json.Marshal(l)
	return string(bs), err
}
func (l List) MarshalJSONMust() string {
	r, _ := l.MarshalJSON()
	return r
}

func UnmarshalJSON(str string) (Value, error) {
	var v interface{}
	err := json.Unmarshal([]byte(str), &v)
	return ValueOf(v), err
}
func UnmarshalJSONMust(str string) Value {
	v, _ := UnmarshalJSON(str)
	return v
}
