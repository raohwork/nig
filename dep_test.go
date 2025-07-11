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

func TestNewDep(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	contextKey := "test_dep_key"
	expectedValue := "my_dependency_value"

	creater := func(ctx *gin.Context) (string, bool) {
		return expectedValue, true
	}

	setter := func(ctx *gin.Context, val string) {
		ctx.Set(contextKey, val)
	}

	getter := func(ctx *gin.Context) (string, bool) {
		val, exists := ctx.Get(contextKey)
		if !exists {
			return "", false
		}
		strVal, ok := val.(string)
		return strVal, ok
	}

	dep := NewDep(creater, setter, getter)

	assert.True(t, dep.Set(c), "Set should return true when dependency is created")
	val, exists := c.Get(contextKey)
	assert.True(t, exists, "Dependency should be set in context")
	assert.Equal(t, expectedValue, val, "Set should set the correct value")

	assert.True(t, dep.Set(c), "Set should return true when dependency is already present")

	retrievedValue := dep.Get(c)
	assert.Equal(t, expectedValue, retrievedValue, "Get should retrieve the correct value")

	updatedValue := "new_dependency_value"
	dep.Update(c, updatedValue)
	retrievedUpdatedValue := dep.Get(c)
	assert.Equal(t, updatedValue, retrievedUpdatedValue, "Update should change the dependency value")

	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	dep.Setup(c2)
	val2, exists2 := c2.Get(contextKey)
	assert.True(t, exists2, "Setup should ensure dependency is present")
	assert.Equal(t, expectedValue, val2, "Setup should set the correct value")
}

func TestCreateDep(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	contextKey := "created_dep_key"
	expectedValue := "created_dependency_value"

	creater := func(ctx *gin.Context) (string, bool) {
		return expectedValue, true
	}

	dep := CreateDep(contextKey, creater)

	assert.True(t, dep.Set(c), "Set should create the dependency")
	val, exists := c.Get(contextKey)
	assert.True(t, exists, "Dependency should be in context")
	assert.Equal(t, expectedValue, val, "Correct value should be set")

	retrievedValue := dep.Get(c)
	assert.Equal(t, expectedValue, retrievedValue, "Get should retrieve the correct value")

	updatedValue := "updated_created_value"
	dep.Update(c, updatedValue)
	retrievedUpdatedValue := dep.Get(c)
	assert.Equal(t, updatedValue, retrievedUpdatedValue, "Update should work for created dep")
}

func TestUse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	contextKey := "use_dep_key"
	predefinedValue := "predefined_value"

	dep := Use(predefinedValue, contextKey)

	assert.True(t, dep.Set(c), "Set should always return true for Use dep")
	val, exists := c.Get(contextKey)
	assert.True(t, exists, "Predefined value should be set in context")
	assert.Equal(t, predefinedValue, val, "Correct predefined value should be set")

	retrievedValue := dep.Get(c)
	assert.Equal(t, predefinedValue, retrievedValue, "Get should retrieve the predefined value")

	updatedValue := "updated_use_value"
	dep.Update(c, updatedValue)
	retrievedUpdatedValue := dep.Get(c)
	assert.Equal(t, updatedValue, retrievedUpdatedValue, "Update should work for Use dep")
}
func TestDepCreaterFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	contextKey := "fail_dep_key"

	creater := func(ctx *gin.Context) (string, bool) {
		// Simulate failure to create dependency
		return "", false
	}

	setter := func(ctx *gin.Context, val string) {
		ctx.Set(contextKey, val)
	}

	getter := func(ctx *gin.Context) (string, bool) {
		val, exists := ctx.Get(contextKey)
		if !exists {
			return "", false
		}
		strVal, ok := val.(string)
		return strVal, ok
	}

	dep := NewDep(creater, setter, getter)

	assert.False(t, dep.Set(c), "Set should return false when creater fails")
	_, exists := c.Get(contextKey)
	assert.False(t, exists, "Dependency should not be set if creater fails")

	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	dep.Setup(c2)
	_, exists2 := c2.Get(contextKey)
	assert.False(t, exists2, "Setup should not set dependency if creater fails")
}
