// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func TestValidateHandlerType(t *testing.T) {
	suite.Run(t, new(validateHandlerTypeSuite))
}

type validateHandlerTypeSuite struct {
	suite.Suite
}

// Test valid handler with one struct parameter
func (s *validateHandlerTypeSuite) Test_ValidHandler_OneStruct() {
	type Args struct {
		Field int
	}
	handler := func(c *gin.Context, args Args) {}

	info, err := validateHandlerType(handler)
	s.Require().NoError(err)
	s.NotNil(info.handlerFunc)
	s.Len(info.argStructs, 1)
}

// Test valid handler with multiple struct parameters
func (s *validateHandlerTypeSuite) Test_ValidHandler_MultipleStructs() {
	type Args1 struct {
		Field1 int
	}
	type Args2 struct {
		Field2 string
	}
	type Args3 struct {
		Field3 bool
	}
	handler := func(c *gin.Context, args1 Args1, args2 Args2, args3 Args3) {}

	info, err := validateHandlerType(handler)
	s.Require().NoError(err)
	s.NotNil(info.handlerFunc)
	s.Len(info.argStructs, 3)
}

// Test error: handler is not a function
func (s *validateHandlerTypeSuite) Test_Error_NotAFunction() {
	handler := "not a function"

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}

// Test error: handler has no parameters
func (s *validateHandlerTypeSuite) Test_Error_NoParameters() {
	handler := func() {}

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}

// Test error: handler has only context parameter
func (s *validateHandlerTypeSuite) Test_Error_OnlyContextParameter() {
	handler := func(c *gin.Context) {}

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}

// Test error: first parameter is not *gin.Context
func (s *validateHandlerTypeSuite) Test_Error_FirstParamNotContext() {
	type Args struct {
		Field int
	}
	handler := func(s string, args Args) {}

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}

// Test error: second parameter is not a struct
func (s *validateHandlerTypeSuite) Test_Error_SecondParamNotStruct() {
	handler := func(c *gin.Context, i int) {}

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}

// Test error: one of multiple parameters is not a struct
func (s *validateHandlerTypeSuite) Test_Error_NonStructInMultipleParams() {
	type Args1 struct {
		Field1 int
	}
	handler := func(c *gin.Context, args1 Args1, s string) {}

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}

// Test error: handler has return value
func (s *validateHandlerTypeSuite) Test_Error_HasReturnValue() {
	type Args struct {
		Field int
	}
	handler := func(c *gin.Context, args Args) error {
		return nil
	}

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}

// Test error: handler has multiple return values
func (s *validateHandlerTypeSuite) Test_Error_HasMultipleReturnValues() {
	type Args struct {
		Field int
	}
	handler := func(c *gin.Context, args Args) (int, error) {
		return 0, nil
	}

	info, err := validateHandlerType(handler)
	s.Require().Error(err)
	s.Empty(info.argStructs)
}
