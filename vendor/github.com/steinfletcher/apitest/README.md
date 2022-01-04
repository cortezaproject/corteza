<p align="center">
    <img src="https://apitest.dev/static/images/dummy.svg" width="150">
</p>

<p align="center">
<a href="https://godoc.org/github.com/steinfletcher/apitest"><img src="https://godoc.org/github.com/steinfletcher/apitest?status.svg" alt="Godoc" /></a>
<a href="https://coveralls.io/github/steinfletcher/apitest?branch=master&service=github"><img src="https://coveralls.io/repos/github/steinfletcher/apitest/badge.svg?branch=master" alt="Coverage Status"/></a>
<a href="https://circleci.com/gh/steinfletcher/apitest"><img src="https://circleci.com/gh/steinfletcher/apitest.svg?style=shield" alt="Build Status" /></a>
<a href="https://goreportcard.com/report/github.com/steinfletcher/apitest"><img src="https://goreportcard.com/badge/github.com/steinfletcher/apitest" alt="Go Report Card" /></a>
<a href="https://github.com/avelino/awesome-go/#testing"><img src="https://awesome.re/mentioned-badge.svg" alt="Mentioned in Awesome Go" /></a>
</p>

# apitest

A simple and extensible behavioural testing library. Supports mocking external http calls and renders sequence diagrams on completion.

In behavioural tests the internal structure of the app is not known by the tests. Data is input to the system and the outputs are expected to meet certain conditions.

Join the conversation at #apitest on [https://gophers.slack.com](https://gophers.slack.com).

<span>Logo by <a target="_blank" href="https://twitter.com/egonelbre">@egonelbre</a><span>

## Documentation

Please visit [https://apitest.dev](https://apitest.dev) for the latest documentation.

## Installation

```bash
go get -u github.com/steinfletcher/apitest
```

## Examples

### Framework and library integration examples

| Example                                                                                              | Comment                                                                                                    |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| [gin](https://github.com/steinfletcher/apitest/tree/master/examples/gin)                             | popular martini-like web framework                                                                         |
| [graphql](https://github.com/steinfletcher/apitest/tree/master/examples/graphql)                     | using gqlgen.com to generate a graphql server                                                              |
| [gorilla](https://github.com/steinfletcher/apitest/tree/master/examples/gorilla)                     | the gorilla web toolkit                                                                                    |
| [iris](https://github.com/steinfletcher/apitest/tree/master/examples/iris)                           | iris web framework                                                                                         |
| [echo](https://github.com/steinfletcher/apitest/tree/master/examples/echo)                           | High performance, extensible, minimalist Go web framework                                                  |
| [fiber](https://github.com/steinfletcher/apitest/tree/master/examples/fiber)                         | Express inspired web framework written in Go                                                               |
| [httprouter](https://github.com/steinfletcher/apitest/tree/master/examples/httprouter)               | High performance HTTP request router that scales well                                                      |
| [mocks](https://github.com/steinfletcher/apitest/tree/master/examples/mocks)                         | example mocking out external http calls                                                                    |
| [sequence diagrams](https://github.com/steinfletcher/apitest/tree/master/examples/sequence-diagrams) | generate sequence diagrams from tests. See the [demo](http://demo-html.apitest.dev.s3-website-eu-west-1.amazonaws.com/)  |
| [Ginkgo](https://github.com/steinfletcher/apitest/tree/master/examples/ginkgo) | Ginkgo BDD test framework|

### Companion libraries

| Library                                                                 | Comment                                        |
| ----------------------------------------------------------------------- | -----------------------------------------------|
| [JSONPath](https://github.com/steinfletcher/apitest-jsonpath)           | JSONPath assertion addons                      |
| [CSS Selectors](https://github.com/steinfletcher/apitest-css-selector)  | CSS selector assertion addons                  |
| [PlantUML](https://github.com/steinfletcher/apitest-plantuml)           | Export sequence diagrams as plantUML           |
| [DynamoDB](https://github.com/steinfletcher/apitest-dynamodb)           | Add DynamoDB interactions to sequence diagrams |

### Credits

This library was influenced by the following software packages:

* [YatSpec](https://github.com/bodar/yatspec) for creating sequence diagrams from tests
* [MockMVC](https://spring.io) and [superagent](https://github.com/visionmedia/superagent) for the concept and behavioural testing approach
* [Gock](https://github.com/h2non/gock) for the approach to mocking HTTP services in Go
* [Baloo](https://github.com/h2non/baloo) for API design

### Code snippets

#### JSON body matcher

```go
func TestApi(t *testing.T) {
	apitest.New().
		Handler(handler).
		Get("/user/1234").
		Expect(t).
		Body(`{"id": "1234", "name": "Tate"}`).
		Status(http.StatusCreated).
		End()
}
```

#### JSONPath

For asserting on parts of the response body JSONPath may be used. A separate module must be installed which provides these assertions - `go get -u github.com/steinfletcher/apitest-jsonpath`. This is packaged separately to keep this library dependency free.

Given the response is `{"a": 12345, "b": [{"key": "c", "value": "result"}]}`

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Contains(`$.b[? @.key=="c"].value`, "result")).
		End()
}
```

and `jsonpath.Equals` checks for value equality

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Equal(`$.a`, float64(12345))).
		End()
}
```

#### Custom assert functions

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			return nil
		}).
		End()
}
```

#### Assert cookies

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Patch("/hello").
		Expect(t).
		Status(http.StatusOK).
		Cookies(apitest.Cookie"ABC").Value("12345")).
		CookiePresent("Session-Token").
		CookieNotPresent("XXX").
		Cookies(
			apitest.Cookie("ABC").Value("12345"),
			apitest.Cookie("DEF").Value("67890"),
		).
		End()
}
```

#### Assert headers

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		Expect(t).
		Status(http.StatusOK).
		Headers(map[string]string{"ABC": "12345"}).
		End()
}
```

#### Mocking external http calls

```go
var getUser = apitest.NewMock().
	Get("/user/12345").
	RespondWith().
	Body(`{"name": "jon", "id": "1234"}`).
	Status(http.StatusOK).
	End()

var getPreferences = apitest.NewMock().
	Get("/preferences/12345").
	RespondWith().
	Body(`{"is_contactable": true}`).
	Status(http.StatusOK).
	End()

func TestApi(t *testing.T) {
	apitest.New().
		Mocks(getUser, getPreferences).
		Handler(handler).
		Get("/hello").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "jon", "id": "1234"}`).
		End()
}
```

#### Generating sequence diagrams from tests

```go

func TestApi(t *testing.T) {
	apitest.New().
		Report(apitest.SequenceDiagram()).
		Mocks(getUser, getPreferences).
		Handler(handler).
		Get("/hello").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "jon", "id": "1234"}`).
		End()
}
```

It is possible to override the default storage location by passing the formatter instance `Report(apitest.NewSequenceDiagramFormatter(".sequence-diagrams"))`.
You can bring your own formatter too if you want to produce custom output. By default a sequence diagram is rendered on a html page. See the [demo](http://demo-html.apitest.dev.s3-website-eu-west-1.amazonaws.com/)

#### Debugging http requests and responses generated by api test and any mocks

```go
func TestApi(t *testing.T) {
	apitest.New().
		Debug().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

#### Provide basic auth in the request

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		BasicAuth("username", "password").
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

#### Provide cookies in the request

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		Cookies(apitest.Cookie("ABC").Value("12345")).
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

#### Provide headers in the request

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Delete("/hello").
		Headers(map[string]string{"My-Header": "12345"}).
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

#### Provide query parameters in the request

`Query`, `QueryParams` and `QueryCollection` can all be used in combination 

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		QueryParams(map[string]string{"a": "1", "b": "2"}).
		Query("c", "d").
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

Providing `{"a": {"b", "c", "d"}` results in parameters encoded as `a=b&a=c&a=d`.
`QueryCollection` can be used in combination with `Query`

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Get("/hello").
		QueryCollection(map[string][]string{"a": {"b", "c", "d"}}).
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

#### Provide a url encoded form body in the request

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Post("/hello").
		FormData("a", "1").
		FormData("b", "2").
		FormData("b", "3").
		FormData("c", "4", "5", "6").
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

#### Capture the request and response data

```go
func TestApi(t *testing.T) {
	apitest.New().
		Observe(func(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
			// do something with res and req
		}).
		Handler(handler).
		Get("/hello").
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

#### Intercept the request

This is useful for mutating the request before it is sent to the system under test.

```go
func TestApi(t *testing.T) {
	apitest.Handler(handler).
		Intercept(func(req *http.Request) {
			req.URL.RawQuery = "a[]=xxx&a[]=yyy"
		}).
		Get("/hello").
		Expect(t).
		Status(http.StatusOK).
		End()
}
```

## Contributing

View the [contributing guide](CONTRIBUTING.md).
