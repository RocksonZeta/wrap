package osutil

import (
	"github.com/RocksonZeta/wrap/wraplog"
)

// var pkg = reflect.TypeOf(struct A{}{}).PkgPath()
var log = wraplog.Logger.Fork("github.com/RocksonZeta/wrap/util/osutil", "")

//Go go func safely
func Go(fn func()) {
	go func() {
		if err := recover(); nil != err {
			if e, ok := err.(error); ok {
				log.Error().Func("Go").Stack().Err(err.(error)).Msg(e.Error())
			} else {
				log.Error().Func("Go").Stack().Err(err.(error)).Interface("err", err).Send()
			}
		}
		fn()
	}()
}
