# Changelog

## v0.1.1

- fix: Check for initialized Client in AddBreadcrumbs (#20)
- build: Bump version when releasing with Craft (#19)

## v0.1.0

- First stable release! \o/

## v0.0.1-beta.5

- feat: **[breaking]** Add `NewHTTPTransport` and `NewHTTPSyncTransport` which accepts all transport options
- feat: New `HTTPSyncTransport` that blocks after each call
- feat: New `Echo` integration
- ref: **[breaking]** Remove `BufferSize` option from `ClientOptions` and move it to `HTTPTransport` instead
- ref: Export default `HTTPTransport`
- ref: Export `net/http` integration handler
- ref: Set `Request` instantly in the package handlers, not in `recoverWithSentry` so it can be accessed later on
- ci: Add craft config

## v0.0.1-beta.4

- feat: `IgnoreErrors` client option and corresponding integration
- ref: Reworked `net/http` integration, wrote better example and complete readme
- ref: Reworked `Gin` integration, wrote better example and complete readme
- ref: Reworked `Iris` integration, wrote better example and complete readme
- ref: Reworked `Negroni` integration, wrote better example and complete readme
- ref: Reworked `Martini` integration, wrote better example and complete readme
- ref: Remove `Handle()` from frameworks handlers and return it directly from New

## v0.0.1-beta.3

- feat: `Iris` framework support with `sentryiris` package
- feat: `Gin` framework support with `sentrygin` package
- feat: `Martini` framework support with `sentrymartini` package
- feat: `Negroni` framework support with `sentrynegroni` package
- feat: Add `Hub.Clone()` for easier frameworks integration
- feat: Return `EventID` from `Recovery` methods
- feat: Add `NewScope` and `NewEvent` functions and use them in the whole codebase
- feat: Add `AddEventProcessor` to the `Client`
- fix: Operate on requests body copy instead of the original
- ref: Try to read source files from the root directory, based on the filename as well, to make it work on AWS Lambda
- ref: Remove `gocertifi` dependence and document how to provide your own certificates
- ref: **[breaking]** Remove `Decorate` and `DecorateFunc` methods in favor of `sentryhttp` package
- ref: **[breaking]** Allow for integrations to live on the client, by passing client instance in `SetupOnce` method
- ref: **[breaking]** Remove `GetIntegration` from the `Hub`
- ref: **[breaking]** Remove `GlobalEventProcessors` getter from the public API

## v0.0.1-beta.2

- feat: Add `AttachStacktrace` client option to include stacktrace for messages
- feat: Add `BufferSize` client option to configure transport buffer size
- feat: Add `SetRequest` method on a `Scope` to control `Request` context data
- feat: Add `FromHTTPRequest` for `Request` type for easier extraction
- ref: Extract `Request` information more accurately
- fix: Attach `ServerName`, `Release`, `Dist`, `Environment` options to the event
- fix: Don't log events dropped due to full transport buffer as sent
- fix: Don't panic and create an appropriate event when called `CaptureException` or `Recover` with `nil` value

## v0.0.1-beta

- Initial release
