// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import "github.com/gin-gonic/gin"

// Arg represents an injectable argument for Gin handlers. It defines how to
// set up the necessary Gin middleware (e.g., for dependency resolution) and
// how to retrieve the actual argument value from the context.
//
// Implementations of Arg are typically created using NewArg.
type Arg[T any] interface {
	// Setup returns a Gin HandlersChain that includes middleware responsible
	// for preparing the argument's dependencies in the Gin context.
	Setup() gin.HandlersChain

	// Get retrieves the argument's value from the Gin context. This method
	// is called by the nig.Wrapper to inject the argument into the handler.
	Get(*gin.Context) T
}

type argImpl[T any] struct {
	setup func() gin.HandlersChain
	get   func(*gin.Context) T
}

func (a *argImpl[T]) Setup() gin.HandlersChain { return a.setup() }
func (a *argImpl[T]) Get(c *gin.Context) T     { return a.get(c) }

// NewArg creates a new Arg instance. It takes two functions:
//   - setup: A function that returns a Gin HandlersChain. These handlers
//     are executed before the main handler and are responsible for
//     setting up any necessary dependencies in the Gin context.
//   - get:   A function that extracts the argument's value from the Gin context.
//     This function is called by the nig.Wrapper to provide the argument
//     to the Gin handler.
func NewArg[T any](setup func() gin.HandlersChain, get func(*gin.Context) T) Arg[T] {
	return &argImpl[T]{setup: setup, get: get}
}

// FromDep creates an Arg from an existing Dep. This is useful when a handler
// directly depends on a single dependency managed by a Dep.
func FromDep[T any](dep Dep[T]) Arg[T] {
	return &argImpl[T]{
		setup: func() gin.HandlersChain { return gin.HandlersChain{dep.Setup} },
		get:   dep.Get,
	}
}
