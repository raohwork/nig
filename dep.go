// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"slices"

	"github.com/gin-gonic/gin"
)

type anotherDep interface {
	setup() []gin.HandlerFunc
	setKey(string)
	dependsOn() []anotherDep
	getKey() string
	_equal(anotherDep) bool
}

// Dep represents a dependency of you gin application.
//
// A dependency is composed by a value with set of gin middlewares.
// You can not create your own implementation, only [CreateDep] is
// capable to create new Dep instance. The value of Dep is mean to
// varies with different request. If the value is unchanged among
// all request, you should place it in a shared struct, as Dep passes
// value via [gin.Context.Set], which creates some overhead.
//
// There are 2 types of dependencies in general: read only and updatable.
// A common example of read only Dep is request id. Logger is a good
// example of updatable Dep. Take a look at [Updates] for more info.
//
// It's quite common that a Dep depends another Dep. Take a look at
// [DependsOn] for more info.
type Dep[T any] struct {
	key     string
	mws     []gin.HandlerFunc
	depends []anotherDep
	create  func(*gin.Context) T
}

func (i *Dep[T]) setup() []gin.HandlerFunc   { return i.mws }
func (i *Dep[T]) _equal(v anotherDep) bool   { return i.key == v.getKey() }
func (i *Dep[T]) setKey(key string)          { i.key = key }
func (i *Dep[T]) getKey() string             { return i.key }
func (i *Dep[T]) setValue(c *gin.Context)    { i.Update(c, i.create(c)) }
func (i *Dep[T]) Update(c *gin.Context, v T) { c.Set(i.key, v) }
func (i *Dep[T]) dependsOn() []anotherDep    { return i.depends }
func (i *Dep[T]) Get(c *gin.Context) T {
	v, _ := c.Get(i.key)
	t, _ := v.(T)
	return t
}
func (i *Dep[T]) addDep(b anotherDep) {
	if !slices.Contains(i.depends, b) {
		i.depends = append(i.depends, b)
	}
}

func NewDep[T any](creater func(*gin.Context) T, mw ...gin.HandlerFunc) *Dep[T] {
	ret := &Dep[T]{
		create: creater,
	}
	ret.mws = append(mw, ret.setValue)
	return ret
}

func DependsOn[D, T any](dep *Dep[D], victim *Dep[T]) *Dep[T] {
	victim.addDep(dep)
	return victim
}

func Updates[D, T any](dep *Dep[D], victim *Dep[T], f func(D, T) D) *Dep[T] {
	victim.addDep(dep)
	victim.mws = append(victim.mws, func(c *gin.Context) {
		dep.Update(c, f(dep.Get(c), victim.Get(c)))
	})
	return victim
}
