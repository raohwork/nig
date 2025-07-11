# NIG - Lightweight Dependency Injection for Gin-Gonic

[![Go Reference](https://pkg.go.dev/badge/github.com/raohwork/nig.svg)](https://pkg.go.dev/github.com/raohwork/nig)

A lightweight dependency injection and handler wrapping mechanism for [Gin-Gonic](https://github.com/gin-gonic/gin), simplifying the management of dependencies and enabling cleaner, more testable Gin handlers.

> **Note**: NIG is not meant to solve complex problems. For more sophisticated dependency injection needs, consider using dedicated DI libraries like [Fx](https://github.com/uber-go/fx) or [Wire](https://github.com/google/wire).

## Core Concepts

NIG introduces three core concepts:

- **Dep**: Represents a dependency that can be created, set, and retrieved from the Gin context. It manages the lifecycle of dependencies within a request.
- **Arg**: An injectable argument for Gin handlers. It typically wraps one or more Dep instances to define the required middleware setup and value extraction from the context. Most commonly created using `FromDep()` or by combining multiple dependencies.
- **Wrapper**: Provides a fluent API for defining Gin routes, automatically injecting the specified arguments into your handlers.

## Installation

```bash
go get github.com/raohwork/nig
```

## Usage

NIG follows a simple three-step pattern:

### 1. Define Dependencies
Wrap your dependencies (logger, database, etc.) as `Dep` instances:

```go
var logger = nig.CreateDep("logger", func(c *gin.Context) (*zap.Logger, bool) {
    // Initialize your logger with request context
    return zap.NewProduction()
})

var db = nig.Use(myDBPool, "db") // Use existing connection pool
```

### 2. Create Handler Arguments
Combine related dependencies into `Arg` instances:

```go
// Single dependency
userArg := nig.FromDep(jwtDep)

// Multiple dependencies bundled together
type MyServices struct {
    Logger
    DB
}
servicesArg := nig.NewArg(func() gin.HandlersChain {
    return gin.HandlersChain{loggerDep.Setup, dbDep.Setup}
}, func(c *gin.Context) MyServices {
    return MyServices{ Logger: loggerDep.Get(c), DB: dbDep.Get(c) }
})
```

### 3. Type-Safe Handlers
Your handlers receive actual types, with dependencies visible in the function signature:

```go
nig.With(router, servicesArg).
    GET("/users", func(c *gin.Context, services MyServices) {
        // services.Logger and services.DB are ready to use
    })

nig.With2Args(router, servicesArg, userArg).
    GET("/profile", func(c *gin.Context, services MyServices, user *User) {
        // Both services and user are automatically injected
    })
```

**Key Benefits:**
- Dependencies are explicit in function signatures
- Type-safe dependency injection
- No magic strings or reflection at runtime

## Complete Example

See [example_test.go](example_test.go) for a comprehensive example showing:

- Request ID extraction and logging setup
- Database dependency management  
- JWT authentication with dependency chaining
- Multiple argument injection patterns

## License

This project is licensed under the Mozilla Public License 2.0. See the license header in source files for details.
