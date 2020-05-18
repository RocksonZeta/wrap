package requestwrap

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/RocksonZeta/wrap/errs"
	"github.com/RocksonZeta/wrap/wraplog"
	"github.com/mozillazg/request"
)

const (
	ErrorBadUrl = 1 + iota
)

var pkg = reflect.TypeOf(Request{}).PkgPath()
var log = wraplog.Logger.Fork(pkg, "Request")

func check(err error, state int, msg string) {
	if err == nil {
		return
	}
	if msg != "" {
		msg = err.Error()
	}
	panic(errs.Err{Err: err, Module: "Request", Pkg: pkg, State: state, Message: msg})
}

type Request struct {
	Request *request.Request
	BaseUrl string
}

func New(baseUrl string, headers, cookies map[string]string, timeout int) *Request {
	log.Trace().Func("New").Str("baseUrl", baseUrl).Interface("headers", headers).Interface("cookies", cookies).Int("timeout", timeout).Send()
	req := request.NewRequest(new(http.Client))
	req.Client.Timeout = time.Duration(timeout) * time.Second
	req.Headers = headers
	req.Cookies = cookies
	r := &Request{Request: req, BaseUrl: strings.TrimRight(baseUrl, "/")}
	req.Hooks = append(req.Hooks, &hook{r})
	return r
}

type hook struct {
	Request *Request
}

func (h *hook) BeforeRequest(req *http.Request) (resp *http.Response, err error) {
	log.Trace().Func("BeforeRequest").Str("method", req.Method).Str("uri", req.URL.RequestURI()).Interface("header", req.Header).Interface("cookies", h.Request.Request.Cookies).Send()
	return
}
func (h *hook) AfterRequest(req *http.Request, resp *http.Response, err error) (newResp *http.Response, newErr error) {
	if err != nil {
		log.Error().Func("AfterRequest").Err(err).Stack().Str("method", req.Method).Str("uri", req.URL.RequestURI()).Msg(err.Error())
		newErr = err
	}
	newResp = resp
	h.Request.clean()
	return
}
func (c *Request) uri(path string) string {
	if '/' != path[0] {
		path = "/" + path
	}
	return c.BaseUrl + path
}
func (c *Request) clean() {
	log.Trace().Func("clean").Send()
	c.Request.Body = nil
	c.Request.Data = nil
	c.Request.Params = nil
	c.Request.Files = nil
	c.Request.Json = nil
}

func MergeQuery(p string, query map[string]string) string {
	if len(query) == 0 {
		return p
	}
	urlObject, err := url.Parse(p)
	if err != nil {
		check(err, ErrorBadUrl, "bad format url:"+p)
	}
	pathQuery := urlObject.Query()
	for k, v := range query {
		pathQuery.Add(k, v)
	}
	urlObject.RawQuery = pathQuery.Encode()
	return urlObject.String()
}
func (c *Request) BasicAuth(user, password string) {
	c.Request.BasicAuth = request.BasicAuth{Username: user, Password: password}
}
func (c *Request) Get(path string, query map[string]string) (body []byte, res *request.Response, err error) {
	log.Trace().Func("Get").Interface("uri", c.uri(path)).Interface("query", query).Send()
	c.Request.Params = query
	res, err = c.Request.Get(c.uri(path))
	if err != nil {
		log.Error().Func("Get").Stack().Err(err).Interface("path", path).Interface("query", query).Interface("headers", c.Request.Headers).Interface("cookies", c.Request.Cookies).Dur("timeout", c.Request.Client.Timeout).Send()
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}

func (c *Request) Post(path string, query, form map[string]string, files []request.FileField) (body []byte, res *request.Response, err error) {
	log.Trace().Func("Post").Str("uri", c.uri(path)).Interface("query", query).Interface("form", form).Int("files", len(files)).Send()
	c.Request.Params = query
	c.Request.Data = form
	c.Request.Files = files
	res, err = c.Request.Post(c.uri(path))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}

// PostForm send post form request.
//
// url can be string or *url.URL or ur.URL
//
// form can be map[string]string or map[string][]string or string or io.Reader
//
// 	form := map[string]string{
// 		"a": "1",
// 		"b": "2",
// 	}
//
// 	form := map[string][]string{
// 		"a": []string{"1", "2"},
// 		"b": []string{"2", "3"},
// 	}
//
// 	form : = "a=1&b=2"
//
// 	form : = strings.NewReader("a=1&b=2")
//
func (c *Request) PostForm(path string, query map[string]string, form interface{}, files []request.FileField) (body []byte, res *request.Response, err error) {
	log.Trace().Func("Post").Str("uri", c.uri(path)).Interface("query", query).Interface("form", form).Int("files", len(files)).Send()
	c.Request.Params = query
	// c.Request.Data = form
	c.Request.Files = files
	res, err = c.Request.PostForm(c.uri(path), form)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}

func (c *Request) PostFile(path string, query, form map[string]string, files map[string]string) (body []byte, res *request.Response, err error) {
	log.Trace().Func("PostFile").Str("uri", c.uri(path)).Interface("query", query).Interface("form", form).Interface("files", files).Send()
	fileFields := make([]request.FileField, len(files))
	i := 0
	for field, file := range files {
		input, ferr := os.Open(file)
		err = ferr
		if ferr != nil {
			log.Error().Func("PostFile").Stack().Err(err).Str("uri", c.uri(path)).Interface("query", query).Interface("form", form).Interface("files", files).Send()
			return
		}
		defer input.Close()
		fileFields[i] = request.FileField{FieldName: field, FileName: filepath.Base(file), File: input}
	}
	return c.Post(c.uri(path), query, form, fileFields)
}
func (c *Request) Put(path string, query, form map[string]string) (body []byte, res *request.Response, err error) {
	log.Trace().Func("Put").Str("uri", c.uri(path)).Interface("query", query).Interface("form", form).Send()
	c.Request.Params = query
	c.Request.Data = form
	res, err = c.Request.Put(c.uri(path))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}
func (c *Request) Delete(path string, query, form map[string]string) (body []byte, res *request.Response, err error) {
	log.Trace().Func("Delete").Str("uri", c.uri(path)).Interface("query", query).Interface("form", form).Send()
	c.Request.Params = query
	c.Request.Data = form
	res, err = c.Request.Delete(c.uri(path))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}
func (c *Request) Patch(path string, query, form map[string]string) (body []byte, res *request.Response, err error) {
	log.Trace().Func("Patch").Str("uri", c.uri(path)).Interface("query", query).Interface("form", form).Send()
	c.Request.Params = query
	c.Request.Data = form
	res, err = c.Request.Patch(c.uri(path))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}
func (c *Request) PostJson(path string, query map[string]string, value interface{}) (body []byte, res *request.Response, err error) {
	log.Trace().Func("PostJson").Str("uri", c.uri(path)).Interface("query", query).Interface("value", value).Send()
	c.Request.Params = query
	c.Request.Json = value
	res, err = c.Request.Post(c.uri(path))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}
