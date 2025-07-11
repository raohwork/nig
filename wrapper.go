// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlerWith is a type alias for a Gin handler function that takes one injected argument.
// This simplifies the signature of handlers that use nig for dependency injection.
type HandlerWith[T any] func(c *gin.Context, arg T)

// With wraps a HandlerWith[T] into a standard Gin handler function (gin.HandlerFunc).
// It retrieves the injected argument from the Gin context using the provided Arg[T]
// and passes it to the underlying HandlerWith[T].
func (h HandlerWith[T]) With(a Arg[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c, a.Get(c))
	}
}

// With creates a new Routes for handlers that require one injected argument.
// It takes a Gin router group and an Arg[T] instance. The Arg's Setup middleware
// is automatically applied to the router group, ensuring the dependency is ready
// before the handler executes.
func With[T any](r gin.IRouter, arg Arg[T]) Routes[T] {
	return Routes[T]{arg: arg, r: r.Group("", arg.Setup()...)}
}

// Routes provides a fluent API for defining Gin routes for handlers that take
// one injected argument. It simplifies the process of attaching handlers to routes
// by automatically handling the argument injection.
type Routes[T any] struct {
	arg Arg[T]
	r   gin.IRoutes
}

// Use registers middleware for the router group.
func (w Routes[T]) Use(middlewate ...HandlerWith[T]) Routes[T] {
	mw := make([]gin.HandlerFunc, len(middlewate))
	for i, m := range middlewate {
		mw[i] = m.With(w.arg)
	}
	w.r.Use(mw...)
	return w
}

// Any registers a route that matches any HTTP method.
func (w Routes[T]) Any(path string, handler HandlerWith[T]) Routes[T] {
	w.r.Any(path, handler.With(w.arg))
	return w
}

// GET registers a route that matches the GET HTTP method.
func (w Routes[T]) GET(path string, handler HandlerWith[T]) Routes[T] {
	w.r.GET(path, handler.With(w.arg))
	return w
}

// POST registers a route that matches the POST HTTP method.
func (w Routes[T]) POST(path string, handler HandlerWith[T]) Routes[T] {
	w.r.POST(path, handler.With(w.arg))
	return w
}

// DELETE registers a route that matches the DELETE HTTP method.
func (w Routes[T]) DELETE(path string, handler HandlerWith[T]) Routes[T] {
	w.r.DELETE(path, handler.With(w.arg))
	return w
}

// PATCH registers a route that matches the PATCH HTTP method.
func (w Routes[T]) PATCH(path string, handler HandlerWith[T]) Routes[T] {
	w.r.PATCH(path, handler.With(w.arg))
	return w
}

// PUT registers a route that matches the PUT HTTP method.
func (w Routes[T]) PUT(path string, handler HandlerWith[T]) Routes[T] {
	w.r.PUT(path, handler.With(w.arg))
	return w
}

// OPTIONS registers a route that matches the OPTIONS HTTP method.
func (w Routes[T]) OPTIONS(path string, handler HandlerWith[T]) Routes[T] {
	w.r.OPTIONS(path, handler.With(w.arg))
	return w
}

// HEAD registers a route that matches the HEAD HTTP method.
func (w Routes[T]) HEAD(path string, handler HandlerWith[T]) Routes[T] {
	w.r.HEAD(path, handler.With(w.arg))
	return w
}

// Match registers a route that matches the specified HTTP methods.
func (w Routes[T]) Match(methods []string, path string, handler HandlerWith[T]) Routes[T] {
	w.r.Match(methods, path, handler.With(w.arg))
	return w
}

// StaticFile registers a single route to serve a static file from the local filesystem.
func (w Routes[T]) StaticFile(relativePath, filepath string) Routes[T] {
	w.r.StaticFile(relativePath, filepath)
	return w
}

func (w Routes[T]) StaticFileFS(relativePath, filepath string, fs http.FileSystem) Routes[T] {
	w.r.StaticFileFS(relativePath, filepath, fs)
	return w
}

// Static registers a route to serve static files from a specified directory.
func (w Routes[T]) Static(relativePath, root string) Routes[T] {
	w.r.Static(relativePath, root)
	return w
}

// StaticFS registers a route to serve static files from a specified file system.
func (w Routes[T]) StaticFS(relativePath string, fs http.FileSystem) Routes[T] {
	w.r.StaticFS(relativePath, fs)
	return w
}

// HandlerWith2 is a type alias for a Gin handler function that takes two injected arguments.
type HandlerWith2[A, B any] func(c *gin.Context, argA A, argB B)

// With wraps a HandlerWith2[A, B] into a standard Gin handler function.
func (h HandlerWith2[A, B]) With(argA Arg[A], argB Arg[B]) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c, argA.Get(c), argB.Get(c))
	}
}

// With2Args creates a new Routes2 for handlers that require two injected arguments.
// It takes a Gin router group and two Arg instances. Their Setup middlewares
// are automatically applied to the router group.
func With2Args[A, B any](r gin.IRouter, argA Arg[A], argB Arg[B]) Routes2[A, B] {
	h := append(argA.Setup(), argB.Setup()...)
	return Routes2[A, B]{argA: argA, argB: argB, r: r.Group("", h...)}
}

// Routes2 provides a fluent API for defining Gin routes for handlers that take
// two injected arguments.
type Routes2[A, B any] struct {
	argA Arg[A]
	argB Arg[B]
	r    gin.IRoutes
}

// Use registers middleware for the router group.
func (w Routes2[A, B]) Use(middlewate ...HandlerWith2[A, B]) Routes2[A, B] {
	mw := make([]gin.HandlerFunc, len(middlewate))
	for i, m := range middlewate {
		mw[i] = m.With(w.argA, w.argB)
	}
	w.r.Use(mw...)
	return w
}

// Any registers a route that matches any HTTP method.
func (w Routes2[A, B]) Any(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.Any(path, handler.With(w.argA, w.argB))
	return w
}

// GET registers a route that matches the GET HTTP method.
func (w Routes2[A, B]) GET(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.GET(path, handler.With(w.argA, w.argB))
	return w
}

// POST registers a route that matches the POST HTTP method.
func (w Routes2[A, B]) POST(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.POST(path, handler.With(w.argA, w.argB))
	return w
}

// DELETE registers a route that matches the DELETE HTTP method.
func (w Routes2[A, B]) DELETE(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.DELETE(path, handler.With(w.argA, w.argB))
	return w
}

// PATCH registers a route that matches the PATCH HTTP method.
func (w Routes2[A, B]) PATCH(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.PATCH(path, handler.With(w.argA, w.argB))
	return w
}

// PUT registers a route that matches the PUT HTTP method.
func (w Routes2[A, B]) PUT(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.PUT(path, handler.With(w.argA, w.argB))
	return w
}

// OPTIONS registers a route that matches the OPTIONS HTTP method.
func (w Routes2[A, B]) OPTIONS(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.OPTIONS(path, handler.With(w.argA, w.argB))
	return w
}

// HEAD registers a route that matches the HEAD HTTP method.
func (w Routes2[A, B]) HEAD(path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.HEAD(path, handler.With(w.argA, w.argB))
	return w
}

// Match registers a route that matches the specified HTTP methods.
func (w Routes2[A, B]) Match(methods []string, path string, handler HandlerWith2[A, B]) Routes2[A, B] {
	w.r.Match(methods, path, handler.With(w.argA, w.argB))
	return w
}

// StaticFile registers a single route to serve a static file from the local filesystem.
func (w Routes2[A, B]) StaticFile(relativePath, filepath string) Routes2[A, B] {
	w.r.StaticFile(relativePath, filepath)
	return w
}

func (w Routes2[A, B]) StaticFileFS(relativePath, filepath string, fs http.FileSystem) Routes2[A, B] {
	w.r.StaticFileFS(relativePath, filepath, fs)
	return w
}

// Static registers a route to serve static files from a specified directory.
func (w Routes2[A, B]) Static(relativePath, root string) Routes2[A, B] {
	w.r.Static(relativePath, root)
	return w
}

// StaticFS registers a route to serve static files from a specified file system.
func (w Routes2[A, B]) StaticFS(relativePath string, fs http.FileSystem) Routes2[A, B] {
	w.r.StaticFS(relativePath, fs)
	return w
}

// HandlerWith3 is a type alias for a Gin handler function that takes three injected arguments.
type HandlerWith3[A, B, C any] func(c *gin.Context, argA A, argB B, argC C)

// With wraps a HandlerWith3[A, B, C] into a standard Gin handler function.
func (h HandlerWith3[A, B, C]) With(argA Arg[A], argB Arg[B], argC Arg[C]) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c, argA.Get(c), argB.Get(c), argC.Get(c))
	}
}

// With3Args creates a new Routes3 for handlers that require three injected arguments.
// It takes a Gin router group and three Arg instances. Their Setup middlewares
// are automatically applied to the router group.
func With3Args[A, B, C any](r gin.IRouter, argA Arg[A], argB Arg[B], argC Arg[C]) Routes3[A, B, C] {
	h := append(append(argA.Setup(), argB.Setup()...), argC.Setup()...)
	return Routes3[A, B, C]{argA: argA, argB: argB, argC: argC, r: r.Group("", h...)}
}

// Routes3 provides a fluent API for defining Gin routes for handlers that take
// three injected arguments.
type Routes3[A, B, C any] struct {
	argA Arg[A]
	argB Arg[B]
	argC Arg[C]
	r    gin.IRoutes
}

// Use registers middleware for the router group.
func (w Routes3[A, B, C]) Use(middlewate ...HandlerWith3[A, B, C]) Routes3[A, B, C] {
	mw := make([]gin.HandlerFunc, len(middlewate))
	for i, m := range middlewate {
		mw[i] = m.With(w.argA, w.argB, w.argC)
	}
	w.r.Use(mw...)
	return w
}

// Any registers a route that matches any HTTP method.
func (w Routes3[A, B, C]) Any(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.Any(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// GET registers a route that matches the GET HTTP method.
func (w Routes3[A, B, C]) GET(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.GET(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// POST registers a route that matches the POST HTTP method.
func (w Routes3[A, B, C]) POST(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.POST(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// DELETE registers a route that matches the DELETE HTTP method.
func (w Routes3[A, B, C]) DELETE(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.DELETE(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// PATCH registers a route that matches the PATCH HTTP method.
func (w Routes3[A, B, C]) PATCH(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.PATCH(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// PUT registers a route that matches the PUT HTTP method.
func (w Routes3[A, B, C]) PUT(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.PUT(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// OPTIONS registers a route that matches the OPTIONS HTTP method.
func (w Routes3[A, B, C]) OPTIONS(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.OPTIONS(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// HEAD registers a route that matches the HEAD HTTP method.
func (w Routes3[A, B, C]) HEAD(path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.HEAD(path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// Match registers a route that matches the specified HTTP methods.
func (w Routes3[A, B, C]) Match(methods []string, path string, handler HandlerWith3[A, B, C]) Routes3[A, B, C] {
	w.r.Match(methods, path, handler.With(w.argA, w.argB, w.argC))
	return w
}

// StaticFile registers a single route to serve a static file from the local filesystem.
func (w Routes3[A, B, C]) StaticFile(relativePath, filepath string) Routes3[A, B, C] {
	w.r.StaticFile(relativePath, filepath)
	return w
}

func (w Routes3[A, B, C]) StaticFileFS(relativePath, filepath string, fs http.FileSystem) Routes3[A, B, C] {
	w.r.StaticFileFS(relativePath, filepath, fs)
	return w
}

// Static registers a route to serve static files from a specified directory.
func (w Routes3[A, B, C]) Static(relativePath, root string) Routes3[A, B, C] {
	w.r.Static(relativePath, root)
	return w
}

// StaticFS registers a route to serve static files from a specified file system.
func (w Routes3[A, B, C]) StaticFS(relativePath string, fs http.FileSystem) Routes3[A, B, C] {
	w.r.StaticFS(relativePath, fs)
	return w
}
