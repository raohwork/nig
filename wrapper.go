// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import "github.com/gin-gonic/gin"

func Use[T any](arg Arg[T], r gin.IRouter) Wrapper[T] {
	return Wrapper[T]{
		r:   r,
		arg: arg,
	}
}

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
