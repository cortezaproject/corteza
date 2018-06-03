# resputil

The package provides utilities to respond with some structured JSON payloads.
By default, anything you pass to `JSON` will be encapsulated depending on the type.

~~~go
func JSON(w http.ResponseWriter, responses ...interface{}) {
~~~

The `response` variadic parameter may be any of the following:

- string
- error
- int
- struct
- nil
- interface
- func() error
- func() (interface{}, error)
- func() ([]byte, error)

It will output the first non-empty value. In the case where the result is taken
from a `interface{}, error` return, it will first output the error if it's not
empty, and then output the value, *even if empty*. For all other cases it will
skip to the next item in the responses parameter.

A parameter of the type of `func() ([]byte, error)` will not return a wrapped structure.
It's assumed the `[]byte` result is an as-is payload which should be returned. A typical
use case is to return the output of `json.Marshal` which matches this signature.

The motivation behind it is to provide more reasonable error handling, when you
want to break out of your function with less code. It's in part an attempt to get
rid of all the `if err != nil {` checks in your code, but at the same time it's
also something that changes how your code might be laid in order to fully take
advantage of what it gives you.

For example, when writing APIs, you might structure your API call into several
logical units, that have different responsibilities:

- request validation (parameters)
- request processing (issuing SQL queries based on parameters)
	- this one may be significantly broken down into many stages
- the actual response payload

Taking advantage of scope, this may look like this:

~~~go
mux.HandleFunc("/api/*", func(w http.ResponseWriter, r *http.Request) {
	owner := login.Decode(r)
	call := chi.URLParam(r, "*")

	// validate request
	validate := func() error {
		if owner == "" {
			return errors.New("Missing login info. Try to relogin")
		}
		if call == "" {
			return errors.New("Unknown API call")
		}
		if r.Method == "POST" {
			return errors.Wrap(r.ParseForm(), "Error parsing POST data")
		}
		return nil
	}

	// process request
	process := func() (interface{}, error) {
		params := map[string]interface{}{
			"owner": owner,
		}
		urlQuery := r.URL.Query()
		for name, param := range urlQuery {
			params[name] = param[0]
		}
		postVars := r.Form
		for name, param := range postVars {
			params[name] = param[0]
		}
		return sqlAPI(call, params)
	}

	// process request
	resputil.JSON(w, validate, process)
})
~~~

Since `validate` and `process` are closures, they may access anything within the scope of
their parent function. This means that you can have a `RequestParameters` struct, a response
struct, and actually extend the logic of this further. This would be one possible way:

~~~go
	// Parameters
	params := &CommentListThread{
		CommentList: &CommentList{
			NewsID:     chi.URLParam(r, "id"),
			SessionID:  r.URL.Query().Get("session_id"),
			Sort:       r.URL.Query().Get("sort"),
			Order:      r.URL.Query().Get("order"),
			PageNumber: parseInt64(r.URL.Query().Get("pageNumber")),
			PageSize:   parseInt64(r.URL.Query().Get("pageSize")),
		},
		SelfID: parseInt64(r.URL.Query().Get("self_id")),
	}

	/* steps:

	0. validate inputs
	1. with self_id=0 parameters:
		a. get comments with self_id 0 in the pagenumber/pagesize range,
		b. get all child comments with parent comment IDs,
		c. add 5 comments with date/asc to parent comments
		d. return comments data
	2. with self_id>0 parameters:
		a. get comments with self_id X in the pagenumber/pagesize range,
		b. return comments data

	*/

	// Parameters are included in the response
	result := params

	validate := func() error {
		if !is(params.Sort, "date", "rating") {
			params.Sort = "date"
		}
		if !is(params.Order, "asc", "desc") {
			params.Order = "asc"
		}
		if params.PageNumber < 0 {
			params.PageNumber = 0
		}
		if params.PageSize < 10 {
			params.PageSize = 10
		}
		if params.PageSize > 100 {
			params.PageSize = 100
		}
		return nil
	}

	// more code here ...

	resputil.JSON(w, validate, process, addReplies, result)
~~~

The parent function is broken down into closures, that represent some stage of the issued
request. Depending on what works for you, each stage individually can produce an error using
a `func() error` declaration like shown here. If a non-empty value is returned, it will
be encoded into JSON and written to the HTTP output.

This pattern of use also allows a more functional approach to what you're responding with. For example,
if you favor something closer to an ORM approach, then you could do something more similar to this:

~~~go
func (p *ProjectHTTP) create(w http.ResponseWriter, r *http.Request) {
	project := Project{}.New()
	resputil.JSON(
		w,
		project.SetName(r.PostFormValue("name")),
		project.Save(),
		project,
	)
}
~~~

In this case, both `SetName` and `Save` will be invoked, regardless of the fact if `SetName` returned
an error. Care should be taken between mixing `error` and `func() error` parameters, due to the order
in which they will be invoked. In the above case, both SetName and Save are invoked before `JSON()`, but
if we would omit `()` from `Save()`, the Save function would be invoked by `JSON()`.

## Tests

The package has 100% code coverage, but errors are possible. Due to the fact that the implementation is
aimed at generic "take anything" use, errors may occur at runtime.

## Other notes

- The package makes use of `pkg/errors`, returning a stack trace in the JSON response if configured with `SetConfig`
- There are helper functions `OK()` and `Success(string)` to format successful messages (`{ "success": { "message": "..." } }`)
- Errors will be formatted as `{ "error": { "message": "..." } }` [according to some Google conventions](https://cloud.google.com/storage/docs/json_api/v1/status-codes) which may or may not have a RFC
- HTTP response codes are not honored, you will always get a 200 OK response and a descriptive JSON payload
- valid responses are nested within `{ "response": ... }`
- if the set of all passed data to respond with it empty or all it's values are empty, `{ "response": false }` will be returned

## License

Written by [@TitPetric](https://twitter.com/TitPetric) and licensed under the permissive [WTFPL](http://www.wtfpl.net/txt/copying/).