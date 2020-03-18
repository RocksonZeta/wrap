package sutil

type Structs []interface{}

//Col structList []struct
func Col(structList interface{}, field string) Structs {
	if nil == structList {
		return nil
	}
	rlist := NewSlice(structList)
	r := make(Structs, len(rlist))
	for i, v := range rlist {
		r[i] = v.Field(field).Value()
	}
	return r
}

func (s Structs) Ints() []int {
	r := make([]int, len(s))
	for i, v := range s {
		r[i] = v.(int)
	}
	return r
}
func (s Structs) Strings() []string {
	r := make([]string, len(s))
	for i, v := range s {
		r[i] = v.(string)
	}
	return r
}
func (s Structs) Int64s() []int64 {
	r := make([]int64, len(s))
	for i, v := range s {
		r[i] = v.(int64)
	}
	return r
}
func (s Structs) Floats() []float64 {
	r := make([]float64, len(s))
	for i, v := range s {
		r[i] = v.(float64)
	}
	return r
}

func IndexByField(structList interface{}, field string, value interface{}) int {
	list := NewSlice(structList)
	for i, v := range list {
		if v.Field(field).Value() == value {
			return i
		}
	}
	return -1
}
