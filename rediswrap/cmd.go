package rediswrap

import (
	"time"

	"github.com/RocksonZeta/wrap/wraplog"
	"github.com/go-redis/redis/v7"
)

var logSc = wraplog.Logger.Fork("github.com/RocksonZeta/wrap/rediswrap", "StringCmd")

type StringCmd struct {
	*redis.StringCmd
}

func (s *StringCmd) Int() int {
	r, err := s.StringCmd.Int()
	if err != nil {
		logSc.Error().Func("Int").Stack().Err(err).Msg(err.Error())
		check(err, ErrorTypeCast, err.Error())
	}
	return r
}
func (s *StringCmd) Int64() int64 {
	r, err := s.StringCmd.Int64()
	if err != nil {
		logSc.Error().Func("Int64").Stack().Err(err).Msg(err.Error())
		check(err, ErrorTypeCast, err.Error())
	}
	return r
}
func (s *StringCmd) Float64() float64 {
	r, err := s.StringCmd.Float64()
	if err != nil {
		logSc.Error().Func("Float64").Stack().Err(err).Msg(err.Error())
		check(err, ErrorTypeCast, err.Error())
	}
	return r
}
func (s *StringCmd) Float32() float32 {
	r, err := s.StringCmd.Float32()
	if err != nil {
		logSc.Error().Func("Float32").Stack().Err(err).Msg(err.Error())
		check(err, ErrorTypeCast, err.Error())
	}
	return r
}
func (s *StringCmd) Bytes() []byte {
	r, err := s.StringCmd.Bytes()
	if err != nil {
		logSc.Error().Func("Bytes").Stack().Err(err).Msg(err.Error())
		check(err, ErrorTypeCast, err.Error())
	}
	return r
}
func (s *StringCmd) Time() time.Time {
	r, err := s.StringCmd.Time()
	if err != nil {
		logSc.Error().Func("Time").Stack().Err(err).Msg(err.Error())
		check(err, ErrorTypeCast, err.Error())
	}
	return r
}
func (s *StringCmd) Uint64() uint64 {
	r, err := s.StringCmd.Uint64()
	if err != nil {
		logSc.Error().Func("Uint64").Stack().Err(err).Msg(err.Error())
		check(err, ErrorTypeCast, err.Error())
	}
	return r
}
