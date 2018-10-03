# Dockertest [![Build Status](https://travis-ci.org/fortytw2/dockertest.svg?branch=master)](https://travis-ci.org/fortytw2/dockertest)

`dockertest` allows you to quickly and easily test database interactions by
creating and destroying databases within your tests using `docker`.

It works by controlling the docker daemon running locally with `exec.Command`.
The flow is as follows ->

1. find a free port on the local machine
2. launch docker container and bind that port
3. wait until the container needs to be shutdown

`dockertest` is inspired by https://divan.github.io/posts/integration_testing/
and https://github.com/ory-am/dockertest - however, it does not add 300k loc of dependencies (guesstimated) to your project. See https://github.com/fsouza/go-dockerclient/issues/599 for more info on this.

# Installation

```sh
go get -u github.com/fortytw2/dockertest
```

currently the tests depend on `github.com/lib/pq`

# How good is it?

Currently the manipulation of the docker daemon is somewhat fragile, as it depends on `exec.Command` and a well placed `time.Sleep` for shutdown. In an
ideal world, this would use the docker api via the docker socket directly,
but it currently works well enough for now. Contributions welcome

# Usage

Postgres example copied from `github.com/fortytw2/hydrocarbon`

```go
func TestDBBits(t *testing.T) {
	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	})
	defer container.Shutdown()
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	db, err := sql.Open("postgres", "postgres://postgres:postgres@" + container.Addr + "?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// run tests on the db, etc
}
```

It should be trivial to adapt the above bits to work just as well from a
`TestMain` function, if you want to avoid running a new container for each
individual test function - subtests also help here.

### Docker-machine host

Note that dockertest will give priority to the `DOCKER_MACHINE_NAME` when looking for your container address and will fallback to `localhost` if it fails to find it.

# License

MIT, see LICENSE

