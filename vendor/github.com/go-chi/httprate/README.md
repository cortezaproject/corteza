# httprate

![](https://github.com/go-chi/httprate/workflows/build/badge.svg?branch=master)

net/http request rate limiter based on the Sliding Window Counter pattern inspired by
CloudFlare https://blog.cloudflare.com/counting-things-a-lot-of-different-things/.

The sliding window counter pattern is accurate, smooths traffic and offers a simple counter
design to share a rate-limit amoung a cluster of servers. For example, if you'd like
to use redis to coordinate a rate-limit across a group of microservices you just need
to implement the httprate.LimitCounter interface to support an atomic increment
and get. 


## Example

```go
package main

import (
  "net/http"

  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
  "github.com/go-chi/httprate"
)

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Logger)

  // Enable httprate request limiter of 100 requests per minute.
  //
  // In the code example below, rate-limiting is bound to the request IP address
  // via the LimitByIP middleware handler.
  //
  // To have a single rate-limiter for all requests, use httprate.LimitAll(..).
  //
  // Please see _example/main.go for other more, or read the library code.
  r.Use(httprate.LimitByIP(100, 1*time.Minute))

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("."))
  })

  http.ListenAndServe(":3333", r)
}
```

## LICENSE

MIT
