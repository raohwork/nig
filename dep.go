// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"github.com/gin-gonic/gin"
)

// Dep represents a dependency that can be created, set, and retrieved from the
// Gin context. It manages the lifecycle of a dependency within a request.
//
// Implementations of Dep are typically created using NewDep.
type Dep[T any] interface {
	// Setup is a Gin middleware that ensures the dependency is available in the
	// Gin context. If the dependency is not already present, it will be created
	// and set.
	//
	// This is designed to be used as a middleware. The default implementation
	// in this package is just call Set and c.Abort() if it returns false.
	Setup(*gin.Context)

	// Set attempts to create and set the dependency in the Gin context if it's
	// not already present. It returns true if the dependency is successfully
	// set or was already present, false otherwise.
	//
	// Creation and error handling are done here. Typically you should only
	// use this method when you have a dependency that depends on another
	// dependency.
	Set(*gin.Context) bool

	// Get retrieves the dependency's value from the Gin context. It assumes
	// the dependency has already been set (e.g., by calling Setup or Set).
	Get(*gin.Context) T

	// Update replaces the existing dependency value in the Gin context with the
	// provided value. This is useful for modifying a dependency's state during
	// a request.
	Update(*gin.Context, T)
}

type dep[T any] struct {
	creater func(*gin.Context) (T, bool)
	setter  func(*gin.Context, T)
	getter  func(*gin.Context) (T, bool)
}

func (d *dep[T]) Setup(c *gin.Context) {
	if !d.Set(c) {
		c.Abort()
	}
}

func (d *dep[T]) Set(c *gin.Context) bool {
	v, ok := d.getter(c)
	if ok {
		return true
	}

	v, ok = d.creater(c)
	if ok {
		d.setter(c, v)
	}
	return ok
}

func (d *dep[T]) Get(c *gin.Context) T {
	v, _ := d.getter(c)
	return v
}

func (d *dep[T]) Update(c *gin.Context, val T) {
	d.setter(c, val)
}

// NewDep creates a new Dep instance. It takes three functions:
//   - creater: A function that attempts to create or retrieve the dependency's
//     value. It returns the value and a boolean indicating success.
//     This function is called only if the dependency is not already
//     present in the context. You should ensure that this function calls c.Abort
//     if it fails to create the dependency, so Dep.Setup does not continue
//     processing the request.
//   - setter:  A function that stores the dependency's value in the Gin context.
//     This is typically done using `c.Set("key", value)`.
//   - getter:  A function that retrieves the dependency's value from the Gin
//     context. It should return the value and a boolean indicating
//     if the value was found and is of the correct type.
//     This is typically done using `c.Get("key")` and type assertion.
//
// In general you should use CreateDep instead of this function, unless you need
// specific behavior like setting multiple keys in the Gin context or saving the
// dependency in a different way.
func NewDep[T any](creater func(*gin.Context) (T, bool), setter func(*gin.Context, T), getter func(*gin.Context) (T, bool)) Dep[T] {
	return &dep[T]{
		creater: creater,
		setter:  setter,
		getter:  getter,
	}
}

// CreateDep is a convenience function to create a Dep instance with a specific key
// in the Gin context. It simplifies the creation of dependencies by providing a
// common pattern for creating, setting, and getting dependencies.
func CreateDep[T any](key string, creater func(*gin.Context) (T, bool)) Dep[T] {
	return NewDep(
		creater,
		func(c *gin.Context, val T) {
			c.Set(key, val)
		},
		func(c *gin.Context) (ret T, ok bool) {
			v, ok := c.Get(key)
			if !ok {
				return ret, false
			}
			ret, ok = v.(T)
			return
		},
	)
}

// Use is a convenience function to create a Dep instance that always returns
// a predefined value. This is useful for cases where you have an external
// dependency that should be checked/created elsewhere like db connection.
func Use[T any](v T, key string) Dep[T] {
	return CreateDep(
		key,
		func(*gin.Context) (T, bool) {
			return v, true
		},
	)
}
