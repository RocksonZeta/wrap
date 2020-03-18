package osswrap

import (
	"io"
	"io/ioutil"
	syslog "log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/RocksonZeta/wrap/errs"
	"github.com/RocksonZeta/wrap/wraplog"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	OssErrorInit = 1 + iota
	OssErrorUse
	OssErrorPut
	OssErrorGet
)

var pkg = reflect.TypeOf(Oss{}).PkgPath()
var log = wraplog.Logger.Fork(pkg, "Oss")

func check(err error, state int, msg string) error {
	if err == nil {
		return nil
	}
	if msg != "" {
		msg = err.Error()
	}
	panic(errs.Err{Err: err, Module: "Oss", Pkg: pkg, State: state, Message: msg})
}

type Oss struct {
	client  *oss.Client
	bucket  *oss.Bucket
	options Options
}

type Options struct {
	Endpoint, AccessKeyId, AccessKeySecret string
	ConnectTimeoutSec, ReadWriteTimeout    int64
}

func New(options Options, bucketName string) *Oss {
	log.Trace().Func("New").Interface("options", options).Str("bucketName", bucketName).Send()
	if options.ConnectTimeoutSec <= 0 {
		options.ConnectTimeoutSec = 30
	}
	if options.ReadWriteTimeout <= 0 {
		options.ReadWriteTimeout = 10 * 60
	}
	client, err := oss.New(options.Endpoint, options.AccessKeyId, options.AccessKeySecret, oss.Timeout(options.ConnectTimeoutSec, options.ReadWriteTimeout))
	if err != nil {
		log.Error().Func("New").Stack().Err(err).Interface("options", options).Str("bucketName", bucketName).Msg(err.Error())
		check(err, OssErrorInit, "")
	}
	if log.DebugEnabled() {
		ossLogger := syslog.New(log, "oss", 0644)
		oss.SetLogger(ossLogger)
	}

	r := &Oss{client: client, options: options}
	if bucketName != "" {
		r.Use(bucketName)

	}
	return r
}

func (o *Oss) Use(bucketName string) {
	log.Trace().Func("Use").Str("bucketName", bucketName).Send()
	var err error
	o.bucket, err = o.client.Bucket(bucketName)
	if err != nil {
		log.Error().Func("Use").Stack().Err(err).Str("bucketName", bucketName).Msg(err.Error())
		check(err, OssErrorUse, "")
	}
}
func (o *Oss) Bucket() *oss.Bucket {
	return o.bucket
}
func (o *Oss) Client() *oss.Client {
	return o.client
}
func (o *Oss) Put(src io.Reader, dst string, filename ...string) {
	log.Trace().Func("Put").Str("dst", dst).Strs("filename", filename).Send()
	var options []oss.Option
	if len(filename) > 0 {
		options = append(options, oss.ContentDisposition("attachment;filename=\""+filepath.Base(filename[0])+"\""))
	}
	err := o.bucket.PutObject(formatDst(dst), src, options...)
	if err != nil {
		log.Error().Func("Put").Err(err).Str("dst", dst).Strs("filename", filename).Msg(err.Error())
		check(err, OssErrorPut, "")
	}
}

func (o *Oss) PutFile(src, dst string, filename ...string) {
	log.Trace().Func("PutFile").Str("src", src).Str("dst", dst).Strs("filename", filename).Send()
	var options []oss.Option
	if len(filename) > 0 {
		options = append(options, oss.ContentDisposition("attachment;filename=\""+filepath.Base(filename[0])+"\""))
	}
	err := o.bucket.PutObjectFromFile(formatDst(dst), src, options...)
	if err != nil {
		log.Error().Func("PutFile").Stack().Err(err).Str("src", src).Str("dst", dst).Strs("filename", filename).Msg(err.Error())
		check(err, OssErrorPut, "")
	}
}
func (o *Oss) Get(src string) io.ReadCloser {
	log.Trace().Func("Get").Str("src", src).Send()
	r, err := o.bucket.GetObject(src)
	if err != nil {
		log.Error().Func("Get").Stack().Err(err).Str("src", src).Msg(err.Error())
		check(err, OssErrorGet, "")
	}
	return r
}
func (o *Oss) GetFile(src, dst string) {
	log.Trace().Func("GetFile").Str("src", src).Str("dst", dst).Send()
	err := o.bucket.GetObjectToFile(src, formatDst(dst))
	if err != nil {
		log.Error().Func("GetFile").Stack().Err(err).Str("src", src).Str("dst", dst).Msg(err.Error())
		check(err, OssErrorGet, "")
	}
}
func formatDst(dst string) string {
	if '/' == dst[0] {
		return dst[1:]
	}
	return dst
}
func (o *Oss) Proxy() string {
	log.Trace().Func("GetProxyFile").Send()
	return "https://" + o.bucket.BucketName + "." + o.client.Config.Endpoint
}

//PutDir copy src/* -> dst/*
func (o *Oss) PutDir(src, dst string, filter func(path string)) {
	log.Trace().Func("PutDir").Str("src", src).Str("dst", dst).Send()
	if filter != nil {
		filter(src)
	}
	stat, err := os.Stat(src)
	if err != nil {
		log.Error().Func("PutDir").Stack().Err(err).Str("src", src).Str("dst", dst).Msg(err.Error())
		check(err, OssErrorPut, err.Error())
	}
	if !stat.IsDir() {
		o.PutFile(src, dst)
		return
	}
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Error().Func("PutDir").Stack().Err(err).Str("src", src).Str("dst", dst).Msg(err.Error())
		check(err, OssErrorPut, err.Error())
	}
	for _, file := range files {
		o.PutDir(filepath.Join(src, file.Name()), filepath.Join(dst, file.Name()), filter)
	}
}
