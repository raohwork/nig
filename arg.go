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

// Arg represents an injectable argument for Gin handlers in the static approach.
//
// Unlike Manager's runtime approach using struct tags, Arg provides compile-time
// type safety by explicitly declaring dependencies. It encapsulates:
//   - A list of dependencies (Dep instances)
//   - A create function that constructs the argument value from gin.Context
//
// Arg instances are created with NewArg or convenience functions (Arg1, Arg2, Arg3),
// and are used with Wrapper to register type-safe handlers.
//
// Example:
//
//	type MyHandlerArgs struct {
//		Logger zerolog.Logger
//		ReqID  string
//	}
//	myArg := NewArg(func(c *gin.Context) MyHandlerArgs {
//		return MyHandlerArgs{
//			Logger: loggerDep.Get(c),
//			ReqID:  reqIDDep.Get(c),
//		}
//	}, loggerDep, reqIDDep)
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

// Arg1 is a convenience function to create an Arg from a single Dep.
//
// This is useful when your handler only needs one dependency and you want to
// pass it directly without wrapping it in a struct.
//
// Example:
//
//	loggerArg := Arg1(loggerDep)
//	wrapper := Use(loggerArg, router)
//	wrapper.GET("/hello", func(c *gin.Context, logger zerolog.Logger) {
//		logger.Info().Msg("hello")
//		c.JSON(200, gin.H{"message": "hello"})
//	})
func Arg1[A any](a *Dep[A]) Arg[A] {
	return NewArg(a.Get, a)
}

// Arg2 is a convenience function to create an Arg from two Deps.
//
// The f function combines the two dependency values into your desired argument type.
// This is useful for creating a struct or tuple from multiple dependencies.
//
// Example:
//
//	type HandlerArgs struct {
//		Logger zerolog.Logger
//		ReqID  string
//	}
//	arg := Arg2(func(logger zerolog.Logger, reqID string) HandlerArgs {
//		return HandlerArgs{Logger: logger, ReqID: reqID}
//	}, loggerDep, reqIDDep)
func Arg2[A, B, T any](f func(A, B) T, a *Dep[A], b *Dep[B]) Arg[T] {
	return NewArg(func(ctx *gin.Context) T {
		return f(a.Get(ctx), b.Get(ctx))
	}, a, b)
}

// Arg3 is a convenience function to create an Arg from three Deps.
//
// The f function combines the three dependency values into your desired argument type.
// For handlers requiring more than three dependencies, consider using NewArg with a
// struct to bundle them together.
//
// Example:
//
//	type HandlerArgs struct {
//		Logger zerolog.Logger
//		ReqID  string
//		User   UserInfo
//	}
//	arg := Arg3(
//		func(logger zerolog.Logger, reqID string, user UserInfo) HandlerArgs {
//			return HandlerArgs{Logger: logger, ReqID: reqID, User: user}
//		},
//		loggerDep, reqIDDep, userDep,
//	)
func Arg3[A, B, C, T any](f func(A, B, C) T, a *Dep[A], b *Dep[B], c *Dep[C]) Arg[T] {
	return NewArg(func(ctx *gin.Context) T {
		return f(a.Get(ctx), b.Get(ctx), c.Get(ctx))
	}, a, b, c)
}
