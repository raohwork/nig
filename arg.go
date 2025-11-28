// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"strconv"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var keyInt = uint64(0)

type Arg[T any] struct {
	deps   []anotherDep
	create func(*gin.Context) T
}

func (a Arg[T]) middlewares() ([]gin.HandlerFunc, error) {
	// sort deps
	deps, err := sortDeps(a.deps)
	if err != nil {
		return nil, err
	}

	// count size and fill key
	size := 0
	for _, d := range deps {
		size += len(d.setup())
		if d.getKey() == "" {
			cur := atomic.AddUint64(&keyInt, 1)
			key := "_nig:dep:" + strconv.FormatUint(cur, 10)
			d.setKey(key)
		}
	}
	ret := make([]gin.HandlerFunc, 0, size)
	for _, d := range deps {
		ret = append(ret, d.setup()...)
	}

	return ret, nil
}

// NewArg creates a Arg[T]
//
// Typical usage is:
//
//	type MyArg struct {
//		Logger zerolog.Logger
//		ReqID string
//	}
//	myarg := NewArg(func(c *gin.Context) MyArg {
//		return MyArg{
//			Logger: loggerDep.Get(c),
//			ReqID:  reqidDep.Get(c),
//		}
//	}, loggerDep, reqidDep)
//
// You don't have to define ReqID field if you don't need the value of request
// id, and you can still pass reqidDep so it can update logger.
func NewArg[T any](create func(*gin.Context) T, deps ...anotherDep) Arg[T] {
	return Arg[T]{
		deps:   deps,
		create: create,
	}
}

func Arg1[A any](a *Dep[A]) Arg[A] {
	return NewArg(a.Get, a)
}

func Arg2[A, B, T any](f func(A, B) T, a *Dep[A], b *Dep[B]) Arg[T] {
	return NewArg(func(ctx *gin.Context) T {
		return f(a.Get(ctx), b.Get(ctx))
	}, a, b)
}

func Arg3[A, B, C, T any](f func(A, B, C) T, a *Dep[A], b *Dep[B], c *Dep[C]) Arg[T] {
	return NewArg(func(ctx *gin.Context) T {
		return f(a.Get(ctx), b.Get(ctx), c.Get(ctx))
	}, a, b, c)
}
