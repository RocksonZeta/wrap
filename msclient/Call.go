package msclient

import (
	"encoding/json"
	"reflect"

	"github.com/RocksonZeta/wrap/errs"
	"github.com/RocksonZeta/wrap/requestwrap"
	"github.com/RocksonZeta/wrap/wraplog"
	"github.com/mozillazg/request"
)

const (
	ErrorResp = 1 + iota
	ErrorMarshal
)

var pkg = reflect.TypeOf(Call{}).PkgPath()
var log = wraplog.Logger.Fork(pkg, "Call")

func check(err error, state int, msg string) {
	if err == nil {
		return
	}
	if msg != "" {
		msg = err.Error()
	}
	panic(errs.Err{Err: err, Module: "Call", Pkg: pkg, State: state, Message: msg})
}

type Call struct {
	*requestwrap.Request
	Unmarshaller func(body []byte, result interface{}) error
}

func New(baseUrl string, headers, cookies map[string]string, timeout int) *Call {
	log.Trace().Func("NewCall").Interface("headers", headers).Interface("cookies", cookies).Send()
	return &Call{Request: requestwrap.New(baseUrl, headers, cookies, timeout), Unmarshaller: DefaultUnmarshaller}
}

//DefaultUnmarshaller unmarshal to errs.Err
func DefaultUnmarshaller(body []byte, result interface{}) error {
	var re errs.Err
	re.Data = result
	err := json.Unmarshal(body, &re)
	if err != nil {
		log.Error().Func("Unmarshal").Stack().Err(err).Str("body", string(body)).Interface("result", result)
		return err
	}
	if re.State != 0 {
		return re
	}
	return nil
}

func (c *Call) Unmarshal(body []byte, result interface{}) error {
	log.Trace().Func("Unmarshal").Str("body", string(body)).Send()
	return c.Unmarshaller(body, result)
}
func (c *Call) Get(result interface{}, path string, query map[string]string) error {
	log.Trace().Func("Get").Interface("path", path).Interface("query", query)
	body, _, err := c.Request.Get(path, query)
	if err != nil {
		return err
	}
	return c.Unmarshal(body, result)
}

func (c *Call) Post(result interface{}, path string, query, form map[string]string, files []request.FileField) error {
	log.Trace().Func("Post").Str("path", path).Interface("query", query).Interface("form", form).Int("files", len(files)).Send()
	body, _, err := c.Request.Post(path, query, form, files)
	if err != nil {
		return err
	}
	return c.Unmarshal(body, result)
}
func (c *Call) PostFile(result interface{}, path string, query, form map[string]string, files map[string]string) error {
	log.Trace().Func("PostFile").Str("path", path).Interface("query", query).Interface("form", form).Interface("files", files).Send()
	body, _, err := c.Request.PostFile(path, query, form, files)
	if err != nil {
		return err
	}
	return c.Unmarshal(body, result)
}
func (c *Call) Put(result interface{}, path string, query, form map[string]string) error {
	log.Trace().Func("Put").Str("path", path).Interface("query", query).Interface("form", form).Send()
	body, _, err := c.Request.Put(path, query, form)
	if err != nil {
		return err
	}
	return c.Unmarshal(body, result)
}
func (c *Call) Delete(result interface{}, path string, query, form map[string]string) error {
	log.Trace().Func("Delete").Str("path", path).Interface("query", query).Interface("form", form).Send()
	body, _, err := c.Request.Delete(path, query, form)
	if err != nil {
		return err
	}
	return c.Unmarshal(body, result)
}
func (c *Call) Patch(result interface{}, path string, query, form map[string]string) error {
	log.Trace().Func("Patch").Str("path", path).Interface("query", query).Interface("form", form).Send()
	body, _, err := c.Request.Delete(path, query, form)
	if err != nil {
		return err
	}
	return c.Unmarshal(body, result)
}
