## Upgrade from v0.3.x to v1.x
In v1.x the detect functions no longer return the MIME type and extension as
strings. Instead they return a pointer to a
[MIME](https://godoc.org/github.com/gabriel-vasile/mimetype#MIME) struct.
The returned MIME pointer is never nil, even when a non-nil error is returned too.
To get the string value of the MIME and the extension, call the
`String()` and the `Extension()` methods.

In order to play better with the stdlib `mime` package, v1.x file extensions
include the leading dot, as in ".html".

In v1.x the `text/plain` MIME type is `text/plain; charset=utf-8`.
