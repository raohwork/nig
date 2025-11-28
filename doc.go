// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package nig provides a lightweight dependency injection library for Gin-Gonic
// that helps mid-sized projects manage middleware and handler dependencies elegantly.
//
// It is not meant to solve complex DI problems. You should use dedicated DI libraries
// like Fx or Wire for complex dependency graphs.
//
// # Core Concepts
//
// Dep[T] represents a per-request dependency that can be created, set, and retrieved
// from the Gin context. It manages the lifecycle of the dependency within a request,
// along with its associated middleware and dependencies on other Dep instances.
//
// # Two Approaches
//
// NIG offers two approaches for dependency injection. You cannot mix them in the
// same application, as they use different mechanisms for storing values in gin.Context.
//
// Runtime Validation (Manager):
//
//   - Manager uses reflection and struct tags for dependency injection
//   - Define handlers accepting structs with "nig" tags
//   - Easy to use with minimal boilerplate
//   - Errors caught at runtime (panics on misconfiguration)
//   - Best for rapid development and flexibility
//
// Example:
//
//	mgr := nig.New(router)
//	mgr.Register("logger", loggerDep)
//
//	type HelloArgs struct {
//	    Logger zerolog.Logger `nig:"logger"`
//	}
//	mgr.GET("/hello", func(c *gin.Context, args HelloArgs) {
//	    args.Logger.Info().Msg("hello")
//	    c.JSON(200, gin.H{"message": "hello"})
//	})
//
// Static Validation (Wrapper + Arg):
//
//   - Wrapper and Arg provide compile-time type safety with no reflection
//   - Use NewArg or convenience functions (Arg1, Arg2, Arg3) to define arguments
//   - Full compile-time type checking
//   - Faster (no reflection overhead)
//   - Best for production code where type safety is critical
//
// Example:
//
//	arg := nig.Arg1(loggerDep)
//	wrapper := nig.Use(arg, router)
//	wrapper.GET("/hello", func(c *gin.Context, logger zerolog.Logger) {
//	    logger.Info().Msg("hello")
//	    c.JSON(200, gin.H{"message": "hello"})
//	})
//
// For complete examples, see the package examples and test files.
package nig
