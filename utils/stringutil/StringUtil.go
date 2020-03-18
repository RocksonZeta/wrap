package stringutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/fatih/structs"

	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Ints(ss []string) []int {
	r := make([]int, len(ss))
	for i, v := range ss {
		r[i], _ = strconv.Atoi(v)
	}
	return r
}

func IntsDefault(s string, dv ...int) int {
	v, err := strconv.Atoi(s)
	if nil != err {
		if len(dv) > 0 {
			return dv[0]
		}
		return 0
	}
	return v
}
func Strings(ints []int) []string {
	r := make([]string, len(ints))
	for i, v := range ints {
		r[i] = strconv.Itoa(v)
	}
	return r
}

func IsUtf8(s []byte) bool {
	return utf8.Valid(s)
}
func GbkToUtf8(s []byte) ([]byte, error) {
	if utf8.Valid(s) {
		return s, nil
	}
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return s, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return s, e
	}
	return d, nil
}

//FileNameAppend hello.jpg,_1 -> hello_1.jpg
func FileNameAppend(filename, subname string) string {
	i := strings.LastIndex(filename, ".")
	if i == -1 {
		return filename + subname
	}
	return filename[0:i] + subname + filename[i:]
}

func Template(tmpl string, obj interface{}) string {
	var kv map[string]interface{}
	if objKv, ok := obj.(map[string]interface{}); ok {
		kv = objKv
	} else {
		kv = structs.New(obj).Map()
	}
	re := regexp.MustCompile(`\{\w+\}`)
	return re.ReplaceAllStringFunc(tmpl, func(name string) string {
		var key string
		if len(name) > 2 {
			key = name[1 : len(name)-1]
		}
		value, ok := kv[key]
		if !ok {
			return name
		}
		return fmt.Sprintf("%v", value)
	})
}
