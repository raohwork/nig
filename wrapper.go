// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import "github.com/gin-gonic/gin"

// Use creates a Wrapper for registering type-safe handlers with the given Arg.
//
// This is the entry point for the static, compile-time approach to dependency
// injection. Unlike Manager which uses reflection and struct tags, Wrapper
// provides full type safety at compile time.
//
// The arg parameter defines what dependencies are needed and how to construct
// the handler argument from gin.Context.
// The r parameter is the gin router (or router group) where handlers will be registered.
//
// Returns a Wrapper that can be used to register handlers via GET, POST, etc.
//
// Example:
//
//	loggerArg := Arg1(loggerDep)
//	wrapper := Use(loggerArg, router)
//	wrapper.GET("/hello", func(c *gin.Context, logger zerolog.Logger) {
//		logger.Info().Msg("handling request")
//		c.JSON(200, gin.H{"message": "hello"})
//	})
func Use[T any](arg Arg[T], r gin.IRouter) Wrapper[T] {
	return Wrapper[T]{
		r:   r,
		arg: arg,
	}
}

// Wrapper provides a type-safe way to register Gin handlers with dependency injection.
//
// It wraps a gin router and an Arg[T], providing methods to register handlers that
// accept *gin.Context and a typed argument T. The wrapper ensures all required
// dependencies are set up before calling your handler.
//
// This approach offers compile-time type safety: if your code compiles, your
// dependencies are correctly wired. This is in contrast to Manager's runtime
// approach which uses reflection and may panic at runtime.
//
// Wrapper is created with the Use function and provides HTTP method registration
// functions (GET, POST, etc.) that accept handlers of type func(*gin.Context, T).
//
// Example:
//
//	type MyArgs struct {
//		Logger zerolog.Logger
//		ReqID  string
//	}
//	arg := NewArg(func(c *gin.Context) MyArgs {
//		return MyArgs{
//			Logger: loggerDep.Get(c),
//			ReqID:  reqIDDep.Get(c),
//		}
//	}, loggerDep, reqIDDep)
//
//	wrapper := Use(arg, router)
//	wrapper.GET("/hello", func(c *gin.Context, args MyArgs) {
//		args.Logger.Info().Str("request_id", args.ReqID).Msg("hello")
//		c.JSON(200, gin.H{"message": "hello"})
//	})
type Wrapper[T any] struct {
	r   gin.IRouter
	arg Arg[T]
}

func (w Wrapper[T]) mws() []gin.HandlerFunc {
	mws, err := w.arg.middlewares()
	if err != nil {
		panic(err)
	}
	return mws
}

func (w Wrapper[T]) h(h func(*gin.Context, T)) gin.HandlerFunc {
	return func(c *gin.Context) { h(c, w.arg.create(c)) }
}

func (w Wrapper[T]) GET(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).GET(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) POST(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).POST(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) PUT(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).PUT(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) PATCH(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).PATCH(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) DELETE(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).DELETE(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) HEAD(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).HEAD(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) OPTIONS(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).OPTIONS(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) Any(pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).Any(pattern, w.h(handler))
	return w
}

func (w Wrapper[T]) Handle(method, pattern string, handler func(*gin.Context, T)) Wrapper[T] {
	w.r.Group("", w.mws()...).Handle(method, pattern, w.h(handler))
	return w
}
