package fileutil

import (
	"os"
	"path/filepath"

	"github.com/RocksonZeta/wrap/errs"
)

const (
	ErrorFile = 1 + iota
)

func check(err error, state int, msg string) error {
	if err == nil {
		return nil
	}
	if msg != "" {
		msg = err.Error()
	}
	panic(errs.Err{Err: err, Module: "fileutil", Pkg: "github.com/RocksonZeta/wrap/fileutil", State: state, Message: msg})
}

func FindFileDir(fileName string) string {
	abs, err := filepath.Abs(fileName)
	if err != nil {
		check(err, 1, err.Error())
	}
	cur := filepath.Dir(abs)

	for {
		_, err := os.Stat(filepath.Join(cur, fileName))
		if err != nil && os.IsNotExist(err) {
			if filepath.Dir(cur) == cur {
				return ""
			}
			cur = filepath.Dir(cur)
			continue
		}
		break
	}
	return cur
}
