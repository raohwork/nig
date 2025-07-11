// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package nig provides a lightweight dependency injection and handler wrapping
// mechanism for Gin-Gonic, simplifying the management of dependencies and
// enabling cleaner, more testable Gin handlers.
//
// It is not mean to solve complex problems. You should use dedicated DI libraries
// like Fx or Wire for that.
//
// It introduces three core concepts:
//   - Dep: Represents a dependency that can be created, set, and retrieved from
//     the Gin context. It manages the lifecycle of dependency within a request.
//   - Arg: An injectable argument for Gin handlers. It defines how to set up
//     the necessary Gin middleware (e.g., for dependency resolution) and how
//     to retrieve the actual argument value from the context.
//   - Wrapper: Provides a fluent API for defining Gin routes, automatically
//     injecting the specified arguments into your handlers. Wrappers are available
//     for handlers taking one, two, or three injected arguments. If you have more
//     than three arguments, you are recommended to use structs to bundle them.
package nig
