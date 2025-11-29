// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// wrapperAbortSuite tests abort behavior in Wrapper approach
type wrapperAbortSuite struct {
	suite.Suite
}

func TestWrapperIntegral(t *testing.T) {
	suite.Run(t, new(wrapperAbortSuite))
}

func (s *wrapperAbortSuite) Test_Normal() {
	// Track middleware execution order
	var callOrder []string

	// Create depB (base dependency)
	// depB has an initial value of 10
	depB := NewDep(
		func(c *gin.Context) int { return 10 },
		func(c *gin.Context) {
			callOrder = append(callOrder, "B")
			c.Next()
		},
	)

	// Create depA (depends on depB)
	// depA has a value of 5 and updates depB's value to B + A
	depA := Updates(
		depB,
		NewDep(
			func(c *gin.Context) int { return 5 },
			func(c *gin.Context) {
				callOrder = append(callOrder, "A")
				c.Next()
			},
		),
		func(b, a int) int { return b + a }, // depB will be updated to 10 + 5 = 15
	)

	// Define handler argument structure
	type TestArgs struct {
		ValueA int
		ValueB int
	}

	// Create Arg using Arg2 helper
	testArg := Arg2(
		func(a, b int) TestArgs {
			return TestArgs{
				ValueA: a,
				ValueB: b,
			}
		},
		depA,
		depB,
	)

	// Variables to verify handler execution and received values
	var handlerCalled bool
	var gotA, gotB int

	// Define the handler
	handler := func(c *gin.Context, args TestArgs) {
		handlerCalled = true
		gotA = args.ValueA
		gotB = args.ValueB
		c.JSON(200, gin.H{
			"value_a": args.ValueA,
			"value_b": args.ValueB,
			"sum":     args.ValueA + args.ValueB,
		})
	}

	// Set up test server
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Use Wrapper API
	wrapper := Use(testArg, router)
	wrapper.GET("/test", handler)

	// Send test request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assertion 1: Middleware execution order (B must be before A)
	s.Equal([]string{"B", "A"}, callOrder, "middleware should be called in correct order (B before A)")

	// Assertion 2: Handler was called
	s.True(handlerCalled, "handler should be called")

	// Assertion 3: Values obtained in handler are correct
	s.Equal(5, gotA, "value A should be 5")
	s.Equal(15, gotB, "value B should be updated to 15 (10 + 5)")

	// Assertion 4: HTTP response status code is correct
	s.Equal(200, w.Code, "response status should be 200")

	// Assertion 5: HTTP response content is correct
	var response struct {
		A   int `json:"value_a"`
		B   int `json:"value_b"`
		Sum int `json:"sum"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err, "response should be valid JSON")
	s.Equal(5, response.A, "response value_a should be 5")
	s.Equal(15, response.B, "response value_b should be 15")
	s.Equal(20, response.Sum, "response sum should be 20")
}

func (s *wrapperAbortSuite) Test_AbortInMiddlewareParameter() {
	// Track middleware execution order
	var callOrder []string
	var handlerCalled bool

	// Create depB (base dependency) that aborts in middleware without calling c.Next()
	depB := NewDep(
		func(c *gin.Context) int { return 10 },
		func(c *gin.Context) {
			callOrder = append(callOrder, "B")
			c.Abort()
			c.Status(403)
			// NOT calling c.Next() - this prevents setValue and subsequent deps from running
		},
	)

	// Create depA (depends on depB) - should never execute
	depA := Updates(
		depB,
		NewDep(
			func(c *gin.Context) int { return 5 },
			func(c *gin.Context) {
				callOrder = append(callOrder, "A")
				c.Next()
			},
		),
		func(b, a int) int { return b + a },
	)

	// Define handler argument structure
	type TestArgs struct {
		ValueA int
		ValueB int
	}

	// Create Arg using Arg2 helper
	testArg := Arg2(
		func(a, b int) TestArgs {
			return TestArgs{
				ValueA: a,
				ValueB: b,
			}
		},
		depA,
		depB,
	)

	// Define the handler
	handler := func(c *gin.Context, args TestArgs) {
		handlerCalled = true
		c.JSON(200, gin.H{"message": "should not reach here"})
	}

	// Set up test server
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Use Wrapper API
	wrapper := Use(testArg, router)
	wrapper.GET("/test", handler)

	// Send test request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assertions
	s.Equal([]string{"B"}, callOrder, "only depB middleware should execute")
	s.False(handlerCalled, "handler should not be called")
	s.Equal(403, w.Code, "response status should be 403")
}

func (s *wrapperAbortSuite) Test_AbortInCreateFunction() {
	// Track middleware execution order
	var callOrder []string
	var handlerCalled bool

	// Create depB (base dependency) that aborts in create function
	depB := NewDep(
		func(c *gin.Context) int {
			c.Abort()
			c.Status(403)
			return 10 // Value still returned and stored
		},
		func(c *gin.Context) {
			callOrder = append(callOrder, "B")
			c.Next() // Allows setValue to run
		},
	)

	// Create depA (depends on depB) - should never execute due to abort flag
	depA := Updates(
		depB,
		NewDep(
			func(c *gin.Context) int { return 5 },
			func(c *gin.Context) {
				callOrder = append(callOrder, "A")
				c.Next()
			},
		),
		func(b, a int) int { return b + a },
	)

	// Define handler argument structure
	type TestArgs struct {
		ValueA int
		ValueB int
	}

	// Create Arg using Arg2 helper
	testArg := Arg2(
		func(a, b int) TestArgs {
			return TestArgs{
				ValueA: a,
				ValueB: b,
			}
		},
		depA,
		depB,
	)

	// Define the handler
	handler := func(c *gin.Context, args TestArgs) {
		handlerCalled = true
		c.JSON(200, gin.H{"message": "should not reach here"})
	}

	// Set up test server
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Use Wrapper API
	wrapper := Use(testArg, router)
	wrapper.GET("/test", handler)

	// Send test request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assertions
	s.Equal([]string{"B"}, callOrder, "only depB middleware should execute")
	s.False(handlerCalled, "handler should not be called")
	s.Equal(403, w.Code, "response status should be 403")
}
