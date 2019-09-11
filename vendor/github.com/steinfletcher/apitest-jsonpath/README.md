[![Build Status](https://travis-ci.org/steinfletcher/apitest-jsonpath.svg?branch=master)](https://travis-ci.org/steinfletcher/apitest-jsonpath)

# apitest-jsonpath

This library provides jsonpath assertions for [apitest](https://github.com/steinfletcher/apitest).

# Installation

```bash
go get -u github.com/steinfletcher/apitest-jsonpath
```

## Examples

### Equals

`Equals` checks for value equality when the json path expression returns a single result. Given the response is `{"a": 12345}`, the result can be asserted as follows

```go
apitest.New(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Equal(`$.a`, float64(12345))).
	End()
```

we can also provide more complex expected values

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Equal(`$`, map[string]interface{}{"a": "hello", "b": float64(12345)})).
	End()
```

given the response is `{"a": "hello", "b": 12345}` 

### Contains

When the jsonpath expression returns an array, use `jsonpath.Contains` to assert the expected value is contained in the result. Given the response is `{"a": 12345, "b": [{"key": "c", "value": "result"}]}`, we can assert on the result like so

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Contains(`$.b[? @.key=="c"].value`, "result")).
	End()
```

### Len

Use `Len` to check to the length of the returned value.

```go
apitest.New().
	Handler(handler).
	Get("/articles?category=golang").
	Expect(t).
	Assert(jsonpath.Len(`$.items`, 3).
	End()
```

### Present / NotPresent

Use `Present` and `NotPresent` to check the presence of a field in the response without evaluating its value

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(Present(`$.a`)).
	Assert(NotPresent(`$.password`)).
	End()
```

### Matches

Use `Matches` to check that a single path element of type string, number or bool matches a regular expression.

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(Matches(`$.a`, `^[abc]{1,3}$`)).
	End()
```
