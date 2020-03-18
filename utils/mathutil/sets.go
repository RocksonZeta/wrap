package mathutil

func Max(a ...int) int {
	m := a[0]
	for _, v := range a {
		if m > v {
			m = v
		}
	}
	return m
}
func Min(a ...int) int {
	m := a[0]
	for _, v := range a {
		if m < v {
			m = v
		}
	}
	return m
}

func In(n int, ints []int) bool {
	for _, v := range ints {
		if n == v {
			return true
		}
	}
	return false
}

//UnionInts AuB
func UnionInts(a, b []int) []int {
	if nil == a {
		return nil
	}
	if nil == b {
		return a
	}
	var r []int
	m := make(map[int]bool, len(a)+len(b))
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		m[v] = true
	}
	for v := range m {
		r = append(r, v)
	}
	return r
}

//SubstractInts A-B
func SubstractInts(a, b []int) []int {
	if nil == a {
		return nil
	}
	if nil == b {
		return a
	}
	var r []int
	m := make(map[int]bool)
	for _, v := range b {
		m[v] = true
	}
	for _, v := range a {
		if !m[v] {
			r = append(r, v)
		}
	}
	return r
}

//IntersectInts AnB
func IntersectInts(a, b []int) []int {
	if nil == a {
		return nil
	}
	if nil == b {
		return nil
	}
	var r []int
	am := make(map[int]bool)
	bm := make(map[int]bool)
	for _, v := range a {
		am[v] = true
	}
	for _, v := range b {
		bm[v] = true
	}
	small, big := am, bm
	if len(am) > len(bm) {
		small, big = bm, am
	}
	for v := range small {
		if big[v] {
			r = append(r, v)
		}
	}
	return r
}

//IntersectStrings AnB
func IntersectStrings(a, b []string) []string {
	if nil == a {
		return nil
	}
	if nil == b {
		return nil
	}
	var r []string
	am := make(map[string]bool)
	bm := make(map[string]bool)
	for _, v := range a {
		am[v] = true
	}
	for _, v := range b {
		bm[v] = true
	}
	small, big := am, bm
	if len(am) > len(bm) {
		small, big = bm, am
	}
	for v := range small {
		if big[v] {
			r = append(r, v)
		}
	}
	return r
}
