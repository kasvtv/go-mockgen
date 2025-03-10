# go-mockgen

[![PkgGoDev](https://pkg.go.dev/badge/badge/github.com/derision-test/go-mockgen.svg)](https://pkg.go.dev/github.com/derision-test/go-mockgen)
[![Build status](https://github.com/derision-test/go-mockgen/actions/workflows/test.yml/badge.svg)](https://github.com/derision-test/go-mockgen/actions/workflows/test.yml)
[![Latest release](https://img.shields.io/github/release/derision-test/go-mockgen.svg)](https://github.com/derision-test/go-mockgen/releases/)

A mock interface code generator (supports generics as of [v1.2.0](https://github.com/derision-test/go-mockgen/releases/tag/v1.2.0) 🎉).

## Generating Mocks

Install with `go get -u github.com/derision-test/go-mockgen/...`.

Mocks should be generated via `go generate` and should be regenerated on each update to the target interface. For example, in `gen.go`:

```go
package mocks

//go:generate go-mockgen -f github.com/cache/user/pkg -i Cache -o mock_cache_test.go
```

Depending on how you prefer to structure your code, you can either

1. generate mocks next to the implementation (as a sibling or in a sibling `mocks` package), or
2. generate mocks as needed in test code (generating them into a `_test.go` file).

### Flags

The following flags are defined by the binary.

| Name               | Short Flag | Description |
| ------------------ | ---------- | ----------- |
| package            | p          | The name of the generated package. Is the name of target directory if dirname or filename is supplied by default. |
| prefix             |            | A prefix used in the name of each mock struct. Should be TitleCase by convention. |
| constructor-prefix |            | A prefix used in the name of each mock constructor function (after the initial `New`/`NewStrict` prefixes). Should be TitleCase by convention. |
| interfaces         | i          | A list of interfaces to generate given the import paths. |
| exclude            | e          | A list of interfaces to exclude from generation. |
| filename           | o          | The target output file. All mocks are written to this file. |
| dirname            | d          | The target output directory. Each mock will be written to a unique file. |
| force              | f          | Do not abort if a write to disk would overwrite an existing file. |
| disable-formatting |            | Do not run goimports over the rendered files (enabled by default). |
| goimports          |            | Path to the goimports binary (uses goimports on your PATH by default). |
| for-test           |            | Append _test suffix to generated package names and file names. |
| file-prefix        |            | Content that is written at the top of each generated file. |

### Configuration file

A configuration file is also supported. If no command line arguments are supplied, then the file `mockgen.yaml` in the current directory is used for input. The structure of the configuration file is as follows (where each entry in the `mocks` list can supply a value for each flag described above):

```yaml
force: true
mocks:
  - filename: foo/bar/mock_cache_test.go
    path: github.com/usr/pkg/cache
    interfaces:
      - Cache
  - filename: foo/baz/mocks_test.go
    # Supports multiple package sources in a single file
    sources:
      - path: github.com/usr/pkg/timer
        interfaces:
          - Timer
      - path: github.com/usr/pkg/stopwatch
        interfaces:
          - LapTimer
          - Stopwatch
```

The top level of the configuration file may also set the keys `exclude`, `prefix`, `constructor-prefix`, `goimports`, `file-prefix`, `force`, `disable-formatting`, and `for-tests`. Top-level excludes will also be applied to each mock generator entry. The values for interface and constructor prefixes, goimports, generated packag names, and file content prefixes will apply to each mock generator entry source(s) if a value is not set. The remaining boolean values will be true for each mock generator entry if set at the top level (regardless of the setting of each entry).

To organize long lists of mocks, multiple files can be used, as follows.

```yaml
include-config-paths:
  - foo.mockgen.yaml
  - bar.mockgen.yaml
  - baz.mockgen.yaml
mocks:
  - filename: foo/bar/mock_cache_test.go
    path: github.com/usr/pkg/cache
    interfaces:
      - Cache
```

This file results in the mocks defined in the `mockgen.yaml` file, concatenated with the mocks defined in `{foo,bar,baz}.mockgen.yaml`. The included config paths do not have global-level configuration and should encode a top-level mocks array, e.g., 

```yaml
- filename: mock_cache_test.go
  path: github.com/usr/pkg/cache
  interfaces:
    - Cache
- filename: mock_timer_test.go
  path: github.com/usr/pkg/timer
  interfaces:
    - Timer
- filename: mock_stopwatch_test.go
  path: github.com/usr/pkg/stopwatch
  interfaces:
    - LapTimer
    - Stopwatch
```

## Testing with Mocks

A mock value fulfills all of the methods of the target interface from which it was generated. Unless overridden, all methods of the mock will return zero values for everything. To override a specific method, you can set its `hook` or its `return values`.

A hook is a method that is called on each invocation and allows the test to specify complex behaviors in the mocked interface (conditionally returning values, synchronizing on external state, etc,). The default hook for a method is set with the `SetDefaultHook` method.

```go
func TestCache(t *testing.T) {
    cache := mocks.NewMockCache[string, int]()
    cache.GetFunc.SetDefaultHook(func (key string) (int, bool) {
        if key == "expected" {
            return 42, true
        }
        return nil, false
    })

    testSubject := NewThingThatNeedsCache(cache)
    // ...
}
```

In the cases where you don't need specific behaviors but just need to return some data, the setup gets a bit easier with `SetDefaultReturn`.

```go
func TestCache(t *testing.T) {
    cache := mocks.NewMockCache[string, int]()
    cache.GetFunc.SetDefaultReturn(42, true)

    testSubject := NewThingThatNeedsCache(cache)
    // ...
}
```

Hook and return values can also be _stacked_ when your test can anticipate multiple calls to the same function. Pushing a hook or a return value will set the hook or return value for _one_ invocation of the mocked method. Once this hook or return value has been spent, it will be removed from the queue. Hooks and return values can be interleaved. If the queue is empty, the default hook will be invoked (or the default return values returned).

The following example will test a cache that returns values 50, 51, and 52 in sequence, then panic if there is an unexpected fourth call.

```go
func TestCache(t *testing.T) {
    cache := mocks.NewMockCache[string, int]()
    cache.GetFunc.SetDefaultHook(func (key string) (int, bool) {
        panic("unexpected call")
    })
    cache.GetFunc.PushReturn(50, true)
    cache.GetFunc.PushReturn(51, true)
    cache.GetFunc.PushReturn(52, true)

    testSubject := NewThingThatNeedsCache(cache)
    // ...
}
```

Note that this "panic by default" behavior is given automatically when using the `NewStrictMockCache` constructor, also automatically generated for all mocks.

### Assertions

Mocks track their invocations and can be retrieved via the `History` method. Structs are generated for each method type containing fields for each argument and result type. Raw assertions can be performed on these values.

```go
allCalls := cache.GetFunc.History()
allCalls[0].Arg0 // key (type string)
allCalls[0].Result0 // value (type int)
allCalls[0].Result1 // exists flag (type bool)
```

### Testify integration

This library also contains an API that integrates with the style of [Testify](https://github.com/stretchr/testify) assertions.

To use the assertions, import the assert and require packages by name.

```go
import (
    mockassert "github.com/derision-test/go-mockgen/testutil/assert"
    mockrequire "github.com/derision-test/go-mockgen/testutil/require"
)
```

The following methods are defined in both packages.

- `Called(t, mockFn, msgAndArgs...)`
- `NotCalled(t, mockFn, msgAndArgs...)`
- `CalledOnce(t, mockFn, msgAndArgs...)`
- `CalledN(t, mockFn, n, msgAndArgs...)`
- `CalledWith(t, mockFn, msgAndArgs...)`
- `NotCalledWith(t, mockFn, msgAndArgs...)`
- `CalledOnceWith(t, mockFn, msgAndArgs...)`
- `CalledNWith(t, mockFn, n, msgAndArgs...)`
- `CalledAtNWith(t, mockFn, n, msgAndArgs...)`

These methods can be used as follows.

```go
// cache.Get called 3 times
mockassert.CalledN(t, cache.GetFunc, 3)

// Ensure cache.Set("foo", 42) was called
mockassert.CalledWith(cache.SetFunc, mockassert.Values("foo", 42))

// Ensure cache.Set("foo", _) was called
mockassert.CalledWith(cache.SetFunc, mockassert.Values("foo", mockassert.Skip))
```

### Gomega integration

This library also contains a set of [Gomega](https://onsi.github.io/gomega/) matchers which simplify assertions over a mocked method's call history.

To use the matchers, import the matchers package anonymously.

```go
import . "github.com/derision-test/go-mockgen/testutil/gomega"
```

The following matchers are defined.

- `BeCalled()`
- `BeCalledN(n)`
- `BeCalledOnce()`
- `BeCalledWith(args...)`
- `BeCalledNWith(args...)`
- `BeCalledOnceWith(args...)`
- `BeAnything()`

These matchers can be used as follows.

```go
// cache.Get called 3 times
Expect(cache.GetFunc).To(BeCalledN(3))

// Ensure cache.Set("foo", "bar") was called
Expect(cache.SetFunc).To(BeCalledWith("foo", "bar"))

// Ensure cache.Set("foo", _) was called
Expect(cache.SetFunc).To(BeCalledWith("foo", BeAnything()))
```
