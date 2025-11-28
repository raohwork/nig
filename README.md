# NIG - Simple middleware management and DI for Gin-Gonic

[![Go Reference](https://pkg.go.dev/badge/github.com/raohwork/nig.svg)](https://pkg.go.dev/github.com/raohwork/nig)

NIG helps mid-sized project to manage gin middlewares.

# TL; DR

Take a look at `example_test.go` and `integral_test.go`.

# Concept and usage

The core of NIG is `Dep[T]`, a factory to create instance of `T` which is needed by handler, and the instance is extracted/computed from gin.Context via a user defined function or middleware.

A `Dep[T]` is composed by:

- An instance of `T`: extracted from gin.Context (like jwt), generated on-the-fly (like request id or logger).
- Set of gin middlewares: optional.
- Another `Dep[T]`: To compute from, or to be updated.

To be clear, if every request shares same instance of `T`, db connection for example, it is NOT a `Dep[T]`. You probably should put it in a struct, and define handler as method of the struct.

Here are some good examples of `Dep[T]`:

- Request ID: likely Dep[string] or Dep[uuid.UUID].
- Structured logger: Dep[zerolog.Logger] is my favorite.
- User info: Dep[UserInfo] which is extracted from jwt.
- JWT: Dep[[]byte] if you need raw JWT.
- Session data: Dep[SessionData] or Dep[*SessionData] if you don't use JWT.

If you defines `Dep[T]` properly, NIG ensures your handler will get working instance of `T`, by running middlewares bundled with `Dep[T]` and it's dependencies. The middleware bundled with `Dep[T]` MUST cooperate with handler well. For example, a `Dep[MyJWTClaim]` should return HTTP 403 and call `c.Abort()` if it failed to extract claim data from JWT, unless you want to handle it in your handler (and you probably want `Dep[*MyJWTClaim]` instead).

There will be two different approach: runtime and statically. You cannot mix the two approaches together: the key used to store value in gin.Context is filled by `Manager` to ensure it is aligned with struct tag, while `Wrapper` (`Arg[T]` actually) also fills the key if not filled.

## Ensure at runtime

`Manager` helps you to manage your `Dep[T]` and handlers. 

1. Define your `Dep[T]`s.
2. Register it with a key to `Manager`.
3. Define a struct `helloArgs`, write field tags so `Manager` knows what you need.
4. Register your `func handleHello(c *gin.Context, arg helloArgs)` to `Manager`.
5. `Manager` validates if everything looks okay before actually register it to gin router.
6. `Manager` warps your handler and register it to router:
   - list required middlewares in order.
   - fill the struct with instances created by middleware.
   - pass gin.Context and struct to your handler.
   
Pros:

- Easy to use
- Less repeative

Cons:

- Typo causes panic
- Reflection is slow, but should be acceptable comparing to network latency

## Ensure statically

The `Wrapper` and `Arg` pair helps you to ensure your dependencies at compile time.

1. Define your `Dep[T]`.
2. Define a struct `helloArgs`.
3. Create an `Arg[helloArgs]` with `NewArg` and required `Dep[T]`s.
4. Register your handler `func(*gin.Context, helloArgs)` with `Use`.

Pros:

- Compile ok = program ok unless logicall error
- Fast, no reflection at all

Cons:

- Quite repeative as you have to write similar code for each struct

## License

This project is licensed under the Mozilla Public License 2.0. See the license header in source files for details.
