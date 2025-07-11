// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewArg(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testValue := "hello_arg"
	contextKey := "my_arg_key"

	setupFunc := func() gin.HandlersChain {
		return []gin.HandlerFunc{func(ctx *gin.Context) {
			ctx.Set(contextKey, testValue)
			ctx.Next()
		}}
	}

	getFunc := func(ctx *gin.Context) string {
		val, exists := ctx.Get(contextKey)
		assert.True(t, exists, "Context key should exist")
		strVal, ok := val.(string)
		assert.True(t, ok, "Value should be a string")
		return strVal
	}

	arg := NewArg(setupFunc, getFunc)

	handlers := arg.Setup()
	assert.Len(t, handlers, 1, "Setup should return a HandlersChain with one handler")
	handlers[0](c) // Execute the setup handler
	val, exists := c.Get(contextKey)
	assert.True(t, exists, "Setup should set the value in context")
	assert.Equal(t, testValue, val, "Setup should set the correct value")

	retrievedValue := arg.Get(c)
	assert.Equal(t, testValue, retrievedValue, "Get should retrieve the correct value from context")
}

func TestFromDep(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testValue := "dep_value"
	depContextKey := "dep_key"

	mockDep := CreateDep(depContextKey, func(ctx *gin.Context) (string, bool) {
		return testValue, true
	})
	arg := FromDep(mockDep)

	handlers := arg.Setup()
	assert.Len(t, handlers, 1, "FromDep Setup should return a HandlersChain with one handler")
	handlers[0](c) // Execute the setup handler
	val, exists := c.Get(depContextKey)
	assert.True(t, exists, "Setup from Dep should set the value in context")
	assert.Equal(t, testValue, val, "Setup from Dep should set the correct value")

	retrievedValue := arg.Get(c)
	assert.Equal(t, testValue, retrievedValue, "Get from Dep should retrieve the correct value")
}
