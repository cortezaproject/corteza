# Corteza Server integration and e2e tests

Main goal of these integration tests is to test entire system with the database
and without any other external services.

What, how and why we test:

 - Handling input:
 
   How HTTP traffic is handled (standard HTTP requests, websocket communication),
   to ensure proper response on valid, invalid or harmful input.
   
 - Configuration:
 
   How system behaves under different settings to cover as much different 
   configuration scenarios as possible.
   
 - Security:
 
   Are we handling access control and log events (audit log) properly?
   All services and data should be protected to prevent unwanted access
   and modifications.
   
 - Scenarios:
 
   Are complex scenarios executed as designed (e.g. is password recovery email 
   sent and can link from the email be used)
      
 - Integration with external services
 
   All external services (with exception to database) are mocked but we do test
   if communication to these services take place and if 

 - Database schema & data migrations
 
   Database migration is executed as one of the first steps in automation tests.
   This can be turned off to enable running tests on live/testing database.
   
   Entire test suite is ran inside transaction so we can rollback all changes
   that occurred while testing 


# Running tests

To run integration tests once:
```shell script
make test.integration
```

For development, you can watch file-system changes (with `nodemon` utility) and 
rerun-tests everytime:
```shell script
make watch.test.integration
```

This can be combined with any of the testing suits and flavours described below.

## Testing suites:

 - `all`: shared + services + integration
 - `integration` integration tests (from API to the DB)
 - `pkg`, `internal`: shared packages (ran as one)
 - `federation`, `system`, `compose`: services

See `Makefile` internals for details.
 

## Testing flavours

 - `test.<suite>` runs simple tests on a specific suite
 - `test.cover.<suite>` run tests with -cover and -covermode=$COVER_MODE
 - `test.coverprofile.<suite>` run cover tests with -coverprofile=$COVER_PROFILE
 
See `Makefile` internals for details.

## Environmental variables you can sue:

### Change test utility wih `GOTEST`

If you want some colors in your CLI, you can use [rakyll/gotest](https://github.com/rakyll/gotest).

```shell script
GOTEST=$GOPATH/bin/gotest make watch.test.integration
```

### Fine-tune test execution with `TEST_FLAGS`
Examples:
 - `TEST_FLAGS="-v" make ....`
 - `TEST_FLAGS="-v -run testName" make ....`
