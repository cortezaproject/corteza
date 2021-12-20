# Corteza Server automated testing

Corteza server has three types of automated tests: unit, store and integration tests.

## Unit tests

Unit tests follow common Golang testing patterns with `_test.go` files inside packages. 
Test files use the same package name as the non-test files to allow testing of unexported functions

## Store tests

[Store tests](store/tests) are intended to verify implementation of the persistence layer across multiple store backends.

## Integration tests

[integration tests](tests/README.md) are intended to verify correct exectution of API request to storage layer and back.
