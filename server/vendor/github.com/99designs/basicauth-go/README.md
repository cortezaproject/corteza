basicauth-go
=================
[![GoDoc](https://godoc.org/github.com/99designs/basicauth-go?status.svg)](https://godoc.org/github.com/99designs/basicauth-go)
[![Build Status](https://travis-ci.org/99designs/basicauth-go.svg)](https://travis-ci.org/99designs/basicauth-go)


golang middleware for HTTP basic auth.

```go
// Chi

router.Use(basicauth.New("MyRealm", map[string][]string{
    "bob": {"password1", "password2"},
}))


// Manual wrapping

middleware := basicauth.New("MyRealm", map[string][]string{
    "bob": {"password1", "password2"},
})

h := middlware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)) {
    /// do stuff
})

log.Fatal(http.ListenAndServe(":8080", h))
```

### env loading
If your environment looks like this:
```bash
SOME_PREFIX_BOB=password
SOME_PREFIX_JANE=password1,password2
```

you can load it like this:
```go
middleware := basicauth.NewFromEnv("MyRealm", "SOME_PREFIX")
```

