// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const mockKey = "mock_arg_key"

// MockArg 是一個用於測試的 Arg 介面模擬實作
type MockArg string

func (m MockArg) Get(c *gin.Context) string {
	return string(m)
}

func (m MockArg) Setup() gin.HandlersChain {
	return []gin.HandlerFunc{func(c *gin.Context) {
		c.Set(mockKey, string(m))
		c.Next()
	}}
}

func TestHandlerWith(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockValue := "test_injected_value"
	mockArg := MockArg(mockValue)

	testHandler := HandlerWith[string](func(c *gin.Context, arg string) {
		assert.Equal(t, mockValue, arg, "Handler should receive the correct injected argument")
		c.String(http.StatusOK, "Handler executed with: %s", arg)
	})

	ginHandler := testHandler.With(mockArg)
	for _, mw := range mockArg.Setup() {
		mw(c)
	}
	ginHandler(c)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status OK")
	v, _ := c.Get(mockKey)
	assert.Equal(t, mockValue, v, "Expected context to contain the mock value")
	assert.Contains(t, w.Body.String(), "Handler executed with: test_injected_value", "Expected correct response body")
}

func TestWithAndRoutesGET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockValue := "injected_service"
	mockArg := MockArg(mockValue)
	routes := With(router, mockArg)

	routes.GET("/test", func(c *gin.Context, arg string) {
		assert.Equal(t, mockValue, arg, "Handler should receive the correct injected argument")
		c.String(http.StatusOK, "GET handler called with: %s", arg)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status OK")
	assert.Contains(t, w.Body.String(), "GET handler called with: injected_service", "Expected correct response body")
}
