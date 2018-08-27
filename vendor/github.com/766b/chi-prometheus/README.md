# chi-prometheus

[Prometheus](http://prometheus.io) middleware for [chi](https://github.com/go-chi/chi).

This is a port of [negroni-prometheus](https://github.com/zbindenren/negroni-prometheus) middleware (written by [Rene Zbinden](https://github.com/zbindenren)).

## Why

[Logging v. instrumentation](http://peter.bourgon.org/blog/2016/02/07/logging-v-instrumentation.html)

Instead of logging request times, it is considered best practice to provide an endpoint for instrumentation tools (like prometheus).

## Installation

    go get github.com/766b/chi-prometheus

## Usage

Take a look at the [example](./example/main.go).

## What do you get

An endpoint with the following information (stripped output):

    # HELP chi_request_duration_milliseconds How long it took to process the request, partitioned by status code, method and HTTP path.
    # TYPE chi_request_duration_milliseconds histogram
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/metrics",service="serviceName",le="300"} 1
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/metrics",service="serviceName",le="1200"} 1
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/metrics",service="serviceName",le="5000"} 1
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/metrics",service="serviceName",le="+Inf"} 1
    chi_request_duration_milliseconds_sum{code="OK",method="GET",path="/metrics",service="serviceName"} 2.003123
    chi_request_duration_milliseconds_count{code="OK",method="GET",path="/metrics",service="serviceName"} 1
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/ok",service="serviceName",le="300"} 0
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/ok",service="serviceName",le="1200"} 0
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/ok",service="serviceName",le="5000"} 2
    chi_request_duration_milliseconds_bucke{code="OK",method="GET",path="/ok",service="serviceName",le="+Inf"} 2
    chi_request_duration_milliseconds_sum{code="OK",method="GET",path="/ok",service="serviceName"} 4747.529026
    chi_request_duration_milliseconds_count{code="OK",method="GET",path="/ok",service="serviceName"} 2
    # HELP chi_requests_total How many HTTP requests processed, partitioned by status code, method and HTTP path.
    # TYPE chi_requests_total counter
    chi_requests_total{code="OK",method="GET",path="/metrics",service="serviceName"} 1
    chi_requests_total{code="OK",method="GET",path="/ok",service="serviceName"} 2
