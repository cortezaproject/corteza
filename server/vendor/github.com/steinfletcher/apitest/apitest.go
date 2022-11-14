package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

// SystemUnderTestDefaultName default name for system under test
const SystemUnderTestDefaultName = "sut"

// ConsumerDefaultName default consumer name
const ConsumerDefaultName = "cli"

var divider = strings.Repeat("-", 10)
var requestDebugPrefix = fmt.Sprintf("%s>", divider)
var responseDebugPrefix = fmt.Sprintf("<%s", divider)

// APITest is the top level struct holding the test spec
type APITest struct {
	debugEnabled             bool
	mockResponseDelayEnabled bool
	networkingEnabled        bool
	networkingHTTPClient     *http.Client
	reporter                 ReportFormatter
	verifier                 Verifier
	recorder                 *Recorder
	handler                  http.Handler
	name                     string
	request                  *Request
	response                 *Response
	observers                []Observe
	mocksObservers           []Observe
	recorderHook             RecorderHook
	mocks                    []*Mock
	t                        TestingT
	httpClient               *http.Client
	httpRequest              *http.Request
	transport                *Transport
	meta                     map[string]interface{}
	started                  time.Time
	finished                 time.Time
}

// InboundRequest used to wrap the incoming request with a timestamp
type InboundRequest struct {
	request   *http.Request
	timestamp time.Time
}

// FinalResponse used to wrap the final response with a timestamp
type FinalResponse struct {
	response  *http.Response
	timestamp time.Time
}

// Observe will be called by with the request and response on completion
type Observe func(*http.Response, *http.Request, *APITest)

// RecorderHook used to implement a custom interaction recorder
type RecorderHook func(*Recorder)

// New creates a new api test. The name is optional and will appear in test reports
func New(name ...string) *APITest {
	apiTest := &APITest{
		meta: map[string]interface{}{},
	}

	request := &Request{
		apiTest:  apiTest,
		headers:  map[string][]string{},
		query:    map[string][]string{},
		formData: map[string][]string{},
	}
	response := &Response{
		apiTest: apiTest,
		headers: map[string][]string{},
	}
	apiTest.request = request
	apiTest.response = response

	if len(name) > 0 {
		apiTest.name = name[0]
	}

	return apiTest
}

// Handler is a convenience method for creating a new apitest with a handler
func Handler(handler http.Handler) *APITest {
	return New().Handler(handler)
}

// HandlerFunc is a convenience method for creating a new apitest with a handler func
func HandlerFunc(handlerFunc http.HandlerFunc) *APITest {
	return New().HandlerFunc(handlerFunc)
}

// EnableNetworking will enable networking for provided clients
func (a *APITest) EnableNetworking(cli ...*http.Client) *APITest {
	a.networkingEnabled = true
	if len(cli) == 1 {
		a.networkingHTTPClient = cli[0]
		return a
	}
	a.networkingHTTPClient = http.DefaultClient
	return a
}

// EnableMockResponseDelay turns on mock response delays (defaults to OFF)
func (a *APITest) EnableMockResponseDelay() *APITest {
	a.mockResponseDelayEnabled = true
	return a
}

// Debug logs to the console the http wire representation of all http interactions that are intercepted by apitest. This includes the inbound request to the application under test, the response returned by the application and any interactions that are intercepted by the mock server.
func (a *APITest) Debug() *APITest {
	a.debugEnabled = true
	return a
}

// Report provides a hook to add custom formatting to the output of the test
func (a *APITest) Report(reporter ReportFormatter) *APITest {
	a.reporter = reporter
	return a
}

// Recorder provides a hook to add a recorder to the test
func (a *APITest) Recorder(recorder *Recorder) *APITest {
	a.recorder = recorder
	return a
}

// Meta provides a hook to add custom meta data to the test which can be picked up when defining a custom reporter
func (a *APITest) Meta(meta map[string]interface{}) *APITest {
	a.meta = meta
	return a
}

// Handler defines the http handler that is invoked when the test is run
func (a *APITest) Handler(handler http.Handler) *APITest {
	a.handler = handler
	return a
}

// HandlerFunc defines the http handler that is invoked when the test is run
func (a *APITest) HandlerFunc(handlerFunc http.HandlerFunc) *APITest {
	a.handler = handlerFunc
	return a
}

// Mocks is a builder method for setting the mocks
func (a *APITest) Mocks(mocks ...*Mock) *APITest {
	var m []*Mock
	for i := range mocks {
		times := mocks[i].response.mock.times
		for j := 1; j <= times; j++ {
			mockCpy := mocks[i].copy()
			mockCpy.times = 1
			m = append(m, mockCpy)
		}
	}
	a.mocks = m
	return a
}

// HttpClient allows the developer to provide a custom http client when using mocks
func (a *APITest) HttpClient(cli *http.Client) *APITest {
	a.httpClient = cli
	return a
}

// Observe is a builder method for setting the observers
func (a *APITest) Observe(observers ...Observe) *APITest {
	a.observers = observers
	return a
}

// ObserveMocks is a builder method for setting the mocks observers
func (a *APITest) ObserveMocks(observer Observe) *APITest {
	a.mocksObservers = append(a.mocksObservers, observer)
	return a
}

// RecorderHook allows the consumer to provider a function that will receive the recorder instance before the
// test runs. This can be used to inject custom events which can then be rendered in diagrams
// Deprecated: use Recorder() instead
func (a *APITest) RecorderHook(hook RecorderHook) *APITest {
	a.recorderHook = hook
	return a
}

// Request returns the request spec
func (a *APITest) Request() *Request {
	return a.request
}

// Response returns the expected response
func (a *APITest) Response() *Response {
	return a.response
}

// Request is the user defined request that will be invoked on the handler under test
type Request struct {
	interceptor     Intercept
	method          string
	url             string
	body            string
	query           map[string][]string
	queryCollection map[string][]string
	headers         map[string][]string
	formData        map[string][]string
	cookies         []*Cookie
	basicAuth       string
	apiTest         *APITest
}

// Intercept will be called before the request is made. Updates to the request will be reflected in the test
type Intercept func(*http.Request)

type pair struct {
	l string
	r string
}

// Intercept is a builder method for setting the request interceptor
func (a *APITest) Intercept(interceptor Intercept) *APITest {
	a.request.interceptor = interceptor
	return a
}

// Verifier allows consumers to override the verification implementation.
func (a *APITest) Verifier(v Verifier) *APITest {
	a.verifier = v
	return a
}

// Method is a builder method for setting the http method of the request
func (a *APITest) Method(method string) *Request {
	a.request.method = method
	return a.request
}

// HttpRequest defines the native `http.Request`
func (a *APITest) HttpRequest(req *http.Request) *Request {
	a.httpRequest = req
	return a.request
}

// Get is a convenience method for setting the request as http.MethodGet
func (a *APITest) Get(url string) *Request {
	a.request.method = http.MethodGet
	a.request.url = url
	return a.request
}

// Getf is a convenience method that adds formatting support to Get
func (a *APITest) Getf(format string, args ...interface{}) *Request {
	return a.Get(fmt.Sprintf(format, args...))
}

// Post is a convenience method for setting the request as http.MethodPost
func (a *APITest) Post(url string) *Request {
	r := a.request
	r.method = http.MethodPost
	r.url = url
	return r
}

// Postf is a convenience method that adds formatting support to Post
func (a *APITest) Postf(format string, args ...interface{}) *Request {
	return a.Post(fmt.Sprintf(format, args...))
}

// Put is a convenience method for setting the request as http.MethodPut
func (a *APITest) Put(url string) *Request {
	r := a.request
	r.method = http.MethodPut
	r.url = url
	return r
}

// Putf is a convenience method that adds formatting support to Put
func (a *APITest) Putf(format string, args ...interface{}) *Request {
	return a.Put(fmt.Sprintf(format, args...))
}

// Delete is a convenience method for setting the request as http.MethodDelete
func (a *APITest) Delete(url string) *Request {
	a.request.method = http.MethodDelete
	a.request.url = url
	return a.request
}

// Deletef is a convenience method that adds formatting support to Delete
func (a *APITest) Deletef(format string, args ...interface{}) *Request {
	return a.Delete(fmt.Sprintf(format, args...))
}

// Patch is a convenience method for setting the request as http.MethodPatch
func (a *APITest) Patch(url string) *Request {
	a.request.method = http.MethodPatch
	a.request.url = url
	return a.request
}

// Patchf is a convenience method that adds formatting support to Patch
func (a *APITest) Patchf(format string, args ...interface{}) *Request {
	return a.Patch(fmt.Sprintf(format, args...))
}

// URL is a builder method for setting the url of the request
func (r *Request) URL(url string) *Request {
	r.url = url
	return r
}

// URLf is a builder method for setting the url of the request and supports a formatter
func (r *Request) URLf(format string, args ...interface{}) *Request {
	r.url = fmt.Sprintf(format, args...)
	return r
}

// Body is a builder method to set the request body
func (r *Request) Body(b string) *Request {
	r.body = b
	return r
}

// Bodyf sets the request body and supports a formatter
func (r *Request) Bodyf(format string, args ...interface{}) *Request {
	r.body = fmt.Sprintf(format, args...)
	return r
}

// BodyFromFile is a builder method to set the request body
func (r *Request) BodyFromFile(f string) *Request {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		r.apiTest.t.Fatal(err)
	}
	r.body = string(b)
	return r
}

// JSON is a convenience method for setting the request body and content type header as "application/json".
// If v is not a string or []byte it will marshall the provided variable as json
func (r *Request) JSON(v interface{}) *Request {
	switch x := v.(type) {
	case string:
		r.body = x
	case []byte:
		r.body = string(x)
	default:
		asJSON, err := json.Marshal(x)
		if err != nil {
			r.apiTest.t.Fatal(err)
			return nil
		}
		r.body = string(asJSON)
	}
	r.ContentType("application/json")
	return r
}

// JSONFromFile is a convenience method for setting the request body and content type header as "application/json"
func (r *Request) JSONFromFile(f string) *Request {
	r.BodyFromFile(f)
	r.ContentType("application/json")
	return r
}

// GraphQLQuery is a convenience method for building a graphql POST request
func (r *Request) GraphQLQuery(query string, variables ...map[string]interface{}) *Request {
	q := GraphQLRequestBody{
		Query: query,
	}

	if len(variables) > 0 {
		q.Variables = variables[0]
	}

	return r.GraphQLRequest(q)
}

// GraphQLRequest builds a graphql POST request
func (r *Request) GraphQLRequest(body GraphQLRequestBody) *Request {
	r.ContentType("application/json")

	data, err := json.Marshal(body)
	if err != nil {
		r.apiTest.t.Fatal(err)
	}

	r.body = string(data)

	return r
}

// GraphQLRequestBody represents the POST request body as per the GraphQL spec
type GraphQLRequestBody struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

// Query is a convenience method to add a query parameter to the request.
func (r *Request) Query(key, value string) *Request {
	r.query[key] = append(r.query[key], value)
	return r
}

// QueryParams is a builder method to set the request query parameters.
// This can be used in combination with request.QueryCollection
func (r *Request) QueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.query[k] = append(r.query[k], v)
	}
	return r
}

// QueryCollection is a builder method to set the request query parameters
// This can be used in combination with request.Query
func (r *Request) QueryCollection(q map[string][]string) *Request {
	r.queryCollection = q
	return r
}

// Header is a builder method to set the request headers
func (r *Request) Header(key, value string) *Request {
	normalizedKey := textproto.CanonicalMIMEHeaderKey(key)
	r.headers[normalizedKey] = append(r.headers[normalizedKey], value)
	return r
}

// Headers is a builder method to set the request headers
func (r *Request) Headers(headers map[string]string) *Request {
	for k, v := range headers {
		normalizedKey := textproto.CanonicalMIMEHeaderKey(k)
		r.headers[normalizedKey] = append(r.headers[normalizedKey], v)
	}
	return r
}

// ContentType is a builder method to set the Content-Type header of the request
func (r *Request) ContentType(contentType string) *Request {
	normalizedKey := textproto.CanonicalMIMEHeaderKey("Content-Type")
	r.headers[normalizedKey] = []string{contentType}
	return r
}

// Cookie is a convenience method for setting a single request cookies by name and value
func (r *Request) Cookie(name, value string) *Request {
	r.cookies = append(r.cookies, &Cookie{name: &name, value: &value})
	return r
}

// Cookies is a builder method to set the request cookies
func (r *Request) Cookies(c ...*Cookie) *Request {
	r.cookies = append(r.cookies, c...)
	return r
}

// BasicAuth is a builder method to sets basic auth on the request.
func (r *Request) BasicAuth(username, password string) *Request {
	r.basicAuth = fmt.Sprintf("%s:%s", username, password)
	return r
}

// FormData is a builder method to set the body form data
// Also sets the content type of the request to application/x-www-form-urlencoded
func (r *Request) FormData(name string, values ...string) *Request {
	r.ContentType("application/x-www-form-urlencoded")
	r.formData[name] = append(r.formData[name], values...)
	return r
}

// Expect marks the request spec as complete and following code will define the expected response
func (r *Request) Expect(t TestingT) *Response {
	r.apiTest.t = t
	return r.apiTest.response
}

// Response is the user defined expected response from the application under test
type Response struct {
	status            int
	body              string
	headers           map[string][]string
	headersPresent    []string
	headersNotPresent []string
	cookies           []*Cookie
	cookiesPresent    []string
	cookiesNotPresent []string
	apiTest           *APITest
	assert            []Assert
}

// Assert is a user defined custom assertion function
type Assert func(*http.Response, *http.Request) error

// Body is the expected response body
func (r *Response) Body(b string) *Response {
	r.body = b
	return r
}

// Bodyf is the expected response body that supports a formatter
func (r *Response) Bodyf(format string, args ...interface{}) *Response {
	r.body = fmt.Sprintf(format, args...)
	return r
}

// BodyFromFile reads the given file and uses the content as the expected response body
func (r *Response) BodyFromFile(f string) *Response {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		r.apiTest.t.Fatal(err)
	}
	r.body = string(b)
	return r
}

// Cookies is the expected response cookies
func (r *Response) Cookies(cookies ...*Cookie) *Response {
	r.cookies = append(r.cookies, cookies...)
	return r
}

// Cookie is used to match on an individual cookie name/value pair in the expected response cookies
func (r *Response) Cookie(name, value string) *Response {
	r.cookies = append(r.cookies, NewCookie(name).Value(value))
	return r
}

// CookiePresent is used to assert that a cookie is present in the response,
// regardless of its value
func (r *Response) CookiePresent(cookieName string) *Response {
	r.cookiesPresent = append(r.cookiesPresent, cookieName)
	return r
}

// CookieNotPresent is used to assert that a cookie is not present in the response
func (r *Response) CookieNotPresent(cookieName string) *Response {
	r.cookiesNotPresent = append(r.cookiesNotPresent, cookieName)
	return r
}

// Header is a builder method to set the request headers
func (r *Response) Header(key, value string) *Response {
	normalizedName := textproto.CanonicalMIMEHeaderKey(key)
	r.headers[normalizedName] = append(r.headers[normalizedName], value)
	return r
}

// HeaderPresent is a builder method to set the request headers that should be present in the response
func (r *Response) HeaderPresent(name string) *Response {
	normalizedName := textproto.CanonicalMIMEHeaderKey(name)
	r.headersPresent = append(r.headersPresent, normalizedName)
	return r
}

// HeaderNotPresent is a builder method to set the request headers that should not be present in the response
func (r *Response) HeaderNotPresent(name string) *Response {
	normalizedName := textproto.CanonicalMIMEHeaderKey(name)
	r.headersNotPresent = append(r.headersNotPresent, normalizedName)
	return r
}

// Headers is a builder method to set the request headers
func (r *Response) Headers(headers map[string]string) *Response {
	for name, value := range headers {
		normalizedName := textproto.CanonicalMIMEHeaderKey(name)
		r.headers[normalizedName] = append(r.headers[textproto.CanonicalMIMEHeaderKey(normalizedName)], value)
	}
	return r
}

// Status is the expected response http status code
func (r *Response) Status(s int) *Response {
	r.status = s
	return r
}

// Assert allows the consumer to provide a user defined function containing their own
// custom assertions
func (r *Response) Assert(fn func(*http.Response, *http.Request) error) *Response {
	r.assert = append(r.assert, fn)
	return r.apiTest.response
}

// End runs the test returning the result to the caller
func (r *Response) End() Result {
	apiTest := r.apiTest
	defer func() {
		if apiTest.debugEnabled {
			apiTest.finished = time.Now()
			fmt.Println(fmt.Sprintf("Duration: %s\n", apiTest.finished.Sub(apiTest.started)))
		}
	}()

	if apiTest.handler == nil && !apiTest.networkingEnabled {
		apiTest.t.Fatal("either define a http.Handler or enable networking")
	}

	apiTest.started = time.Now()
	var res *http.Response
	if apiTest.reporter != nil {
		res = apiTest.report()
	} else {
		res = r.runTest()
	}

	var unmatchedMocks []UnmatchedMock
	for _, m := range r.apiTest.mocks {
		if m.isUsed == false {
			unmatchedMocks = append(unmatchedMocks, UnmatchedMock{
				URL: *m.request.url,
			})
			break
		}
	}

	return Result{
		Response:       res,
		unmatchedMocks: unmatchedMocks,
	}
}

// Result provides the final result
type Result struct {
	Response       *http.Response
	unmatchedMocks []UnmatchedMock
}

// UnmatchedMocks returns any mocks that were not used, e.g. there was not a matching http Request for the mock
func (r Result) UnmatchedMocks() []UnmatchedMock {
	return r.unmatchedMocks
}

// JSON unmarshal the result response body to a valid struct
func (r Result) JSON(t interface{}) {
	data, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, t)
	if err != nil {
		panic(err)
	}
}

type mockInteraction struct {
	request   *http.Request
	response  *http.Response
	timestamp time.Time
}

func (r *mockInteraction) GetRequestHost() string {
	host := r.request.Host
	if host == "" {
		host = r.request.URL.Host
	}
	return host
}

func (a *APITest) report() *http.Response {
	var capturedInboundReq *http.Request
	var capturedFinalRes *http.Response
	var capturedMockInteractions []*mockInteraction

	a.observers = append(a.observers, func(finalRes *http.Response, inboundReq *http.Request, a *APITest) {
		capturedFinalRes = copyHttpResponse(finalRes)
		capturedInboundReq = copyHttpRequest(inboundReq)
	})

	a.mocksObservers = append(a.mocksObservers, func(mockRes *http.Response, mockReq *http.Request, a *APITest) {
		capturedMockInteractions = append(capturedMockInteractions, &mockInteraction{
			request:   copyHttpRequest(mockReq),
			response:  copyHttpResponse(mockRes),
			timestamp: time.Now().UTC(),
		})
	})

	if a.recorder == nil {
		a.recorder = NewTestRecorder()
	}
	defer a.recorder.Reset()

	if a.recorderHook != nil {
		a.recorderHook(a.recorder)
	}

	a.started = time.Now()
	res := a.response.runTest()
	a.finished = time.Now()

	a.recorder.
		AddTitle(fmt.Sprintf("%s %s", capturedInboundReq.Method, capturedInboundReq.URL.String())).
		AddSubTitle(a.name).
		AddHttpRequest(HttpRequest{
			Source:    ConsumerDefaultName,
			Target:    SystemUnderTestDefaultName,
			Value:     capturedInboundReq,
			Timestamp: a.started,
		})

	for _, interaction := range capturedMockInteractions {
		a.recorder.AddHttpRequest(HttpRequest{
			Source:    SystemUnderTestDefaultName,
			Target:    interaction.GetRequestHost(),
			Value:     interaction.request,
			Timestamp: interaction.timestamp,
		})
		if interaction.response != nil {
			a.recorder.AddHttpResponse(HttpResponse{
				Source:    interaction.GetRequestHost(),
				Target:    SystemUnderTestDefaultName,
				Value:     interaction.response,
				Timestamp: interaction.timestamp,
			})
		}
	}

	a.recorder.AddHttpResponse(HttpResponse{
		Source:    SystemUnderTestDefaultName,
		Target:    ConsumerDefaultName,
		Value:     capturedFinalRes,
		Timestamp: a.finished,
	})

	sort.Slice(a.recorder.Events, func(i, j int) bool {
		return a.recorder.Events[i].GetTime().Before(a.recorder.Events[j].GetTime())
	})

	meta := map[string]interface{}{}

	for k, v := range a.meta {
		meta[k] = v
	}

	meta["status_code"] = capturedFinalRes.StatusCode
	meta["path"] = capturedInboundReq.URL.String()
	meta["method"] = capturedInboundReq.Method
	meta["name"] = a.name
	meta["hash"] = createHash(meta)
	meta["duration"] = a.finished.Sub(a.started).Nanoseconds()

	a.recorder.AddMeta(meta)
	a.reporter.Format(a.recorder)

	return res
}

func createHash(meta map[string]interface{}) string {
	path := meta["path"]
	method := meta["method"]
	name := meta["name"]
	app := meta["app"]

	prefix := fnv.New32a()
	_, err := prefix.Write([]byte(fmt.Sprintf("%s%s%s", app, strings.ToUpper(method.(string)), path)))
	if err != nil {
		panic(err)
	}

	suffix := fnv.New32a()
	_, err = suffix.Write([]byte(name.(string)))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d_%d", prefix.Sum32(), suffix.Sum32())
}

func (r *Response) runTest() *http.Response {
	a := r.apiTest
	if len(a.mocks) > 0 {
		a.transport = newTransport(
			a.mocks,
			a.httpClient,
			a.debugEnabled,
			a.mockResponseDelayEnabled,
			a.mocksObservers,
			r.apiTest,
		)
		defer a.transport.Reset()
		a.transport.Hijack()
	}
	res, req := a.doRequest()

	defer func() {
		if len(a.observers) > 0 {
			for _, observe := range a.observers {
				observe(res, req, a)
			}
		}
	}()

	if a.verifier == nil {
		a.verifier = DefaultVerifier{}
	}

	a.assertMocks()
	a.assertResponse(res)
	a.assertHeaders(res)
	a.assertCookies(res)
	a.assertFunc(res, req)

	return copyHttpResponse(res)
}

func (a *APITest) assertMocks() {
	for _, mock := range a.mocks {
		if mock.isUsed == false && mock.timesSet {
			a.verifier.Fail(a.t, "mock was not invoked expected times")
		}
	}
}

func (a *APITest) assertFunc(res *http.Response, req *http.Request) {
	if len(a.response.assert) > 0 {
		for _, assertFn := range a.response.assert {
			err := assertFn(copyHttpResponse(res), copyHttpRequest(req))
			if err != nil {
				a.verifier.NoError(a.t, err)
			}
		}
	}
}

func (a *APITest) doRequest() (*http.Response, *http.Request) {
	req := a.buildRequest()
	if a.request.interceptor != nil {
		a.request.interceptor(req)
	}
	resRecorder := httptest.NewRecorder()

	if a.debugEnabled {
		requestDump, err := httputil.DumpRequest(req, true)
		if err == nil {
			debugLog(requestDebugPrefix, "inbound http request", string(requestDump))
		}
	}

	var res *http.Response
	var err error
	if !a.networkingEnabled {
		a.serveHttp(resRecorder, copyHttpRequest(req))
		res = resRecorder.Result()
	} else {
		res, err = a.networkingHTTPClient.Do(copyHttpRequest(req))
		if err != nil {
			a.t.Fatal(err)
		}
	}

	if a.debugEnabled {
		responseDump, err := httputil.DumpResponse(res, true)
		if err == nil {
			debugLog(responseDebugPrefix, "final response", string(responseDump))
		}
	}

	return res, req
}

func (a *APITest) serveHttp(res *httptest.ResponseRecorder, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			a.t.Fatalf("%s: %s", err, debug.Stack())
		}
	}()

	a.handler.ServeHTTP(res, req)
}

func (a *APITest) buildRequest() *http.Request {
	if a.httpRequest != nil {
		return a.httpRequest
	}

	if len(a.request.formData) > 0 {
		form := url.Values{}
		for k := range a.request.formData {
			for _, value := range a.request.formData[k] {
				form.Add(k, value)
			}
		}
		a.request.body = form.Encode()
	}

	req, _ := http.NewRequest(a.request.method, a.request.url, bytes.NewBufferString(a.request.body))
	req.URL.RawQuery = formatQuery(a.request)
	req.Host = SystemUnderTestDefaultName
	if a.networkingEnabled {
		req.Host = req.URL.Host
	}

	for k, v := range a.request.headers {
		for _, headerValue := range v {
			req.Header.Add(k, headerValue)
		}
	}

	for _, cookie := range a.request.cookies {
		req.AddCookie(cookie.ToHttpCookie())
	}

	if a.request.basicAuth != "" {
		parts := strings.Split(a.request.basicAuth, ":")
		req.SetBasicAuth(parts[0], parts[1])
	}

	return req
}

func formatQuery(request *Request) string {
	var out url.Values = map[string][]string{}

	if request.queryCollection != nil {
		for _, param := range buildQueryCollection(request.queryCollection) {
			out.Add(param.l, param.r)
		}
	}

	if request.query != nil {
		for k, v := range request.query {
			for _, p := range v {
				out.Add(k, p)
			}
		}
	}

	if len(out) > 0 {
		return out.Encode()
	}

	return ""
}

func buildQueryCollection(params map[string][]string) []pair {
	if len(params) == 0 {
		return []pair{}
	}

	var pairs []pair
	for k, v := range params {
		for _, paramValue := range v {
			pairs = append(pairs, pair{l: k, r: paramValue})
		}
	}
	return pairs
}

func (a *APITest) assertResponse(res *http.Response) {
	if a.response.status != 0 {
		a.verifier.Equal(a.t, a.response.status, res.StatusCode, fmt.Sprintf("Status code %d not equal to %d", res.StatusCode, a.response.status))
	}

	if a.response.body != "" {
		var resBodyBytes []byte
		if res.Body != nil {
			resBodyBytes, _ = ioutil.ReadAll(res.Body)
			res.Body = ioutil.NopCloser(bytes.NewBuffer(resBodyBytes))
		}
		if json.Valid([]byte(a.response.body)) {
			a.verifier.JSONEq(a.t, a.response.body, string(resBodyBytes))
		} else {
			a.verifier.Equal(a.t, a.response.body, string(resBodyBytes))
		}
	}
}

func (a *APITest) assertCookies(response *http.Response) {
	if len(a.response.cookies) > 0 {
		for _, expectedCookie := range a.response.cookies {
			var mismatchedFields []string
			foundCookie := false
			for _, actualCookie := range response.Cookies() {
				cookieFound, errors := compareCookies(expectedCookie, actualCookie)
				if cookieFound {
					foundCookie = true
					mismatchedFields = append(mismatchedFields, errors...)
				}
			}
			a.verifier.Equal(a.t, true, foundCookie, "ExpectedCookie not found - "+*expectedCookie.name)
			a.verifier.Equal(a.t, 0, len(mismatchedFields), strings.Join(mismatchedFields, ","))
		}
	}

	if len(a.response.cookiesPresent) > 0 {
		for _, cookieName := range a.response.cookiesPresent {
			foundCookie := false
			for _, cookie := range response.Cookies() {
				if cookie.Name == cookieName {
					foundCookie = true
				}
			}
			a.verifier.Equal(a.t, true, foundCookie, "ExpectedCookie not found - "+cookieName)
		}
	}

	if len(a.response.cookiesNotPresent) > 0 {
		for _, cookieName := range a.response.cookiesNotPresent {
			foundCookie := false
			for _, cookie := range response.Cookies() {
				if cookie.Name == cookieName {
					foundCookie = true
				}
			}
			a.verifier.Equal(a.t, false, foundCookie, "ExpectedCookie found - "+cookieName)
		}
	}
}

func (a *APITest) assertHeaders(res *http.Response) {
	for expectedHeader, expectedValues := range a.response.headers {
		resHeaderValues, foundHeader := res.Header[expectedHeader]
		a.verifier.Equal(a.t, true, foundHeader, fmt.Sprintf("expected header '%s' not present in response", expectedHeader))

		if foundHeader {
			for _, expectedValue := range expectedValues {
				foundValue := false
				for _, resValue := range resHeaderValues {
					if expectedValue == resValue {
						foundValue = true
						break
					}
				}
				a.verifier.Equal(a.t, true, foundValue, fmt.Sprintf("mismatched values for header '%s'. Expected %s but received %s", expectedHeader, expectedValue, strings.Join(resHeaderValues, ",")))
			}
		}
	}

	if len(a.response.headersPresent) > 0 {
		for _, expectedName := range a.response.headersPresent {
			if res.Header.Get(expectedName) == "" {
				a.t.Fatalf("expected header '%s' not present in response", expectedName)
			}
		}
	}

	if len(a.response.headersNotPresent) > 0 {
		for _, name := range a.response.headersNotPresent {
			if res.Header.Get(name) != "" {
				a.t.Fatalf("did not expect header '%s' in response", name)
			}
		}
	}
}

func debugLog(prefix, header, msg string) {
	fmt.Printf("\n%s %s\n%s\n", prefix, header, msg)
}

func copyHttpResponse(response *http.Response) *http.Response {
	if response == nil {
		return nil
	}

	var resBodyBytes []byte
	if response.Body != nil {
		resBodyBytes, _ = ioutil.ReadAll(response.Body)
		response.Body = ioutil.NopCloser(bytes.NewBuffer(resBodyBytes))
	}

	resCopy := &http.Response{
		Header:        map[string][]string{},
		StatusCode:    response.StatusCode,
		Status:        response.Status,
		Body:          ioutil.NopCloser(bytes.NewBuffer(resBodyBytes)),
		Proto:         response.Proto,
		ProtoMinor:    response.ProtoMinor,
		ProtoMajor:    response.ProtoMajor,
		ContentLength: response.ContentLength,
	}

	for name, values := range response.Header {
		resCopy.Header[name] = values
	}

	return resCopy
}

func copyHttpRequest(request *http.Request) *http.Request {
	resCopy := &http.Request{
		Method:        request.Method,
		Host:          request.Host,
		Proto:         request.Proto,
		ProtoMinor:    request.ProtoMinor,
		ProtoMajor:    request.ProtoMajor,
		ContentLength: request.ContentLength,
		RemoteAddr:    request.RemoteAddr,
	}
	resCopy = resCopy.WithContext(request.Context())

	if request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(request.Body)
		resCopy.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if request.URL != nil {
		r2URL := new(url.URL)
		*r2URL = *request.URL
		resCopy.URL = r2URL
	}

	headers := make(http.Header)
	for k, values := range request.Header {
		for _, hValue := range values {
			headers.Add(k, hValue)
		}
	}
	resCopy.Header = headers

	return resCopy
}
