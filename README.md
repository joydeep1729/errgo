# errgo

[![Go Reference](https://pkg.go.dev/badge/github.com/joydeep1729/errgo.svg)](https://pkg.go.dev/github.com/joydeep1729/errgo)

Simple error handling primitives with some extra functionalities. This project is a fork of the renowned `github.com/pkg/errors` library.

`go get github.com/joydeep1729/errgo`

The `errgo` package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error. In addition to all the standard features of `pkg/errors`, it introduces natively supported dual-unwrapping.

## Adding context to an error

The `errors.Wrap` function returns a new error that adds context to the original error. For example:
```go
import errors "github.com/joydeep1729/errgo"

_, err := ioutil.ReadAll(r)
if err != nil {
        return errors.Wrap(err, "read failed")
}
```

## Retrieving the cause of an error

Depending on the nature of the error it may be necessary to reverse the operation of `errors.Wrap` to retrieve the original error for inspection. 

With `errgo`, we have introduced dual-unwrapping features to preserve outer layer wrappers and inner errors:

### 1. Step-by-Step Single-Level Unwrapping (`UnwrapWithOuter`)
`UnwrapWithOuter` unwraps exactly one level of the error chain. It returns the immediate wrapped error as `inner`, and the outer context/message added at the current level as `outer`.
```go
import errors "github.com/joydeep1729/errgo"

// Unwraps one level (top-level wrapper context)
inner, outer := errors.UnwrapWithOuter(err)
```

### 2. Direct Root-Cause Unwrapping (`UnwrapToCauseWithOuter`)
`UnwrapToCauseWithOuter` unwraps the entire error chain all the way down to the root cause. It returns the deepest underlying error as `cause`, and the combined outer context of all wrapper layers as `outer`.
```go
import errors "github.com/joydeep1729/errgo"

// Unwraps the entire chain directly to the root cause
cause, outer := errors.UnwrapToCauseWithOuter(err)
```

Standard unwrapping and Go 1.13+ `errors.Is` and `errors.As` are also completely supported out-of-the-box.

## License

BSD-2-Clause
