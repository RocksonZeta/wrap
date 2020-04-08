package errs

import (
	"encoding/json"
	"fmt"
)

// const (
// 	LayerWrap    = "wrap"
// 	LayerDao     = "dao"
// 	LayerService = "service"
// 	LayerRoute   = "route"
// )

type Err struct {
	State     int         //异常的状态码， 0：无异常
	Data      interface{} //异常相关的参数
	Message   string      //消息
	Pkg       string
	Module    string //异常的模块
	Err       error  //被wrap的异常
	UserError bool   //是否是用户异常，否则是系统异常
}

// type JsonResult Err

func (e Err) Error() string {
	bs, _ := json.Marshal(e)
	return string(bs)
}
func (e Err) String() string {
	return fmt.Sprintf("Status:%d,Message:%s", e.State, e.Message)
}

type Result struct {
	State   int
	Data    interface{} `json:"Data,omitempty"`
	Message string      `json:"Message,omitempty"`
	Pkg     string      `json:"Pkg,omitempty"`
	Module  string      `json:"Module,omitempty"`
}

func (e Err) Result() Result {
	var r Result
	r.State = e.State
	r.Data = e.Data
	r.Message = e.Message
	r.Pkg = e.Pkg
	r.Module = e.Module
	return r
}

// type UserError struct {
// 	Err
// }

func NewUserError(state int, msg string) *Err {
	return &Err{State: state, Message: msg, UserError: true}
}

// func NewUserError1(state int, msg string, data interface{}, err error) *UserError {
// 	return &UserError{Err{State: state, Message: msg, Data: data, Err: err}}
// }
