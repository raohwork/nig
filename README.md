# NIG - Simple Middleware Management and DI for Gin-Gonic

[![Go Reference](https://pkg.go.dev/badge/github.com/raohwork/nig.svg)](https://pkg.go.dev/github.com/raohwork/nig)

NIG is a lightweight dependency injection library for Gin-Gonic that helps mid-sized projects manage middleware and handler dependencies elegantly.

## Installation

```bash
go get github.com/raohwork/nig
```

## Quick Start

```go
// Define dependencies
reqIDDep := nig.NewDep(func(c *gin.Context) string {
    return uuid.New().String()
})

loggerDep := nig.NewDep(func(c *gin.Context) zerolog.Logger {
    return log.With().Str("request_id", reqIDDep.Get(c)).Logger()
})
nig.DependsOn(reqIDDep, loggerDep)

// Option 1: Runtime approach with Manager
mgr := nig.New(router)
mgr.Register("logger", loggerDep)

type HelloArgs struct {
    Logger zerolog.Logger `nig:"logger"`
}
mgr.GET("/hello", func(c *gin.Context, args HelloArgs) {
    args.Logger.Info().Msg("hello world")
    c.JSON(200, gin.H{"message": "hello"})
})

// Option 2: Static approach with Wrapper
arg := nig.Arg1(loggerDep)
wrapper := nig.Use(arg, router)
wrapper.GET("/hello", func(c *gin.Context, logger zerolog.Logger) {
    logger.Info().Msg("hello world")
    c.JSON(200, gin.H{"message": "hello"})
})
```

For complete examples, see `example_test.go` and `manager_integral_test.go` / `wrapper_integral_test.go`.

## Core Concepts

### Dep[T]

The core of NIG is `Dep[T]`, a factory that creates instances of type `T` for each request. The instance is extracted or computed from `gin.Context` via a user-defined function and optional middleware.

A `Dep[T]` is composed of:

- **A value of type T**: Extracted from gin.Context (like JWT claims) or generated on-the-fly (like request ID or logger)
- **Optional gin middlewares**: Run before the value is created
- **Optional dependencies**: Other `Dep[T]` instances it depends on or updates

**Important**: If every request shares the same instance of `T` (e.g., database connection), it should NOT be a `Dep[T]`. Put shared resources in a struct and define handlers as methods of that struct instead.

#### Good Examples of Dep[T]

- **Request ID**: `Dep[string]` or `Dep[uuid.UUID]`
- **Structured logger**: `Dep[zerolog.Logger]` (can be enriched with request context)
- **User info**: `Dep[UserInfo]` extracted from JWT
- **JWT token**: `Dep[[]byte]` for raw JWT
- **Session data**: `Dep[SessionData]` or `Dep[*SessionData]` if not using JWT

### Middleware Cooperation

When you define a `Dep[T]` with middleware, the middleware MUST cooperate properly with handlers:

- For required dependencies (e.g., `Dep[MyJWTClaim]`): The middleware should return HTTP 403 and call `c.Abort()` on failure
- For optional dependencies (e.g., `Dep[*MyJWTClaim]`): The middleware can set nil and let the handler decide how to handle missing values

### Dependency Relationships

NIG automatically manages dependency order using topological sorting:

```go
// Logger depends on request ID
reqIDDep := nig.NewDep(func(c *gin.Context) string {
    return uuid.New().String()
})

loggerDep := nig.NewDep(func(c *gin.Context) zerolog.Logger {
    return log.With().Str("request_id", reqIDDep.Get(c)).Logger()
})
nig.DependsOn(reqIDDep, loggerDep)

// User info updates logger with user details
userDep := nig.NewDep(func(c *gin.Context) UserInfo {
    return extractUserFromJWT(c)
})
nig.Updates(loggerDep, userDep, func(logger zerolog.Logger, user UserInfo) zerolog.Logger {
    return logger.With().Str("user_id", user.ID).Logger()
})
```

## Two Approaches

NIG offers two approaches for dependency injection. **You cannot mix them** in the same application, as they use different mechanisms for storing values in `gin.Context`.

### Approach 1: Runtime Validation (Manager)

Use `Manager` for a runtime, reflection-based approach with minimal boilerplate.

**How it works:**

1. Define your `Dep[T]` instances
2. Register them with string keys to `Manager`
3. Define handler argument structs with `nig` tags
4. Register handlers - `Manager` validates everything and wires dependencies automatically

**Example:**

```go
// Setup
mgr := nig.New(router)
mgr.Register("logger", loggerDep).
    Register("reqid", reqIDDep).
    Register("user", userDep)

// Define handler arguments
type UserHandlerArgs struct {
    Logger zerolog.Logger `nig:"logger"`
    User   UserInfo       `nig:"user"`
}

// Register handler
mgr.GET("/profile", func(c *gin.Context, args UserHandlerArgs) {
    args.Logger.Info().Msg("fetching profile")
    c.JSON(200, args.User)
})
```

**Pros:**
- Easy to use with minimal boilerplate
- Flexible - handlers only declare what they need
- Less repetitive code

**Cons:**
- Typos in struct tags cause runtime panics
- Uses reflection (slower, but typically negligible compared to network latency)
- Errors only caught at runtime

### Approach 2: Static Validation (Wrapper + Arg)

Use `Wrapper` and `Arg` for compile-time type safety with no reflection.

**How it works:**

1. Define your `Dep[T]` instances
2. Create an `Arg[T]` that combines dependencies
3. Use `Use()` to create a `Wrapper`
4. Register handlers - everything is type-checked at compile time

**Example:**

```go
// Define argument type
type UserHandlerArgs struct {
    Logger zerolog.Logger
    User   UserInfo
}

// Create Arg
arg := nig.NewArg(func(c *gin.Context) UserHandlerArgs {
    return UserHandlerArgs{
        Logger: loggerDep.Get(c),
        User:   userDep.Get(c),
    }
}, loggerDep, userDep)

// Register handler
wrapper := nig.Use(arg, router)
wrapper.GET("/profile", func(c *gin.Context, args UserHandlerArgs) {
    args.Logger.Info().Msg("fetching profile")
    c.JSON(200, args.User)
})
```

**Convenience functions for simple cases:**

```go
// Single dependency
arg1 := nig.Arg1(loggerDep)
nig.Use(arg1, router).GET("/hello", func(c *gin.Context, logger zerolog.Logger) {
    // handler code
})

// Two dependencies
arg2 := nig.Arg2(
    func(logger zerolog.Logger, reqID string) MyArgs {
        return MyArgs{Logger: logger, ReqID: reqID}
    },
    loggerDep, reqIDDep,
)

// Three dependencies
arg3 := nig.Arg3(
    func(logger zerolog.Logger, reqID string, user UserInfo) MyArgs {
        return MyArgs{Logger: logger, ReqID: reqID, User: user}
    },
    loggerDep, reqIDDep, userDep,
)
```

**Pros:**
- Full compile-time type safety
- Fast - no reflection
- If it compiles, dependencies are correctly wired (barring logical errors)

**Cons:**
- More repetitive - must write `Arg` creation code for each handler type
- More boilerplate compared to Manager approach

## Choosing an Approach

| Criterion | Manager (Runtime) | Wrapper (Static) |
|-----------|------------------|------------------|
| Type safety | Runtime | Compile-time |
| Performance | Slower (reflection) | Faster (no reflection) |
| Boilerplate | Less | More |
| Error detection | Runtime panics | Compile errors |
| Flexibility | High | Medium |

**Recommendation:**
- Use **Manager** for rapid development and when flexibility is important
- Use **Wrapper** for production code where type safety and performance are critical

## License

This project is licensed under the Mozilla Public License 2.0. See the license header in source files for details.
