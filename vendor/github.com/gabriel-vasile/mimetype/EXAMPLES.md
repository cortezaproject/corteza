## Examples
 - [Detect MIME type](#detect-mime-type)
 - [Test against a MIME type](#test-against-a-mime-type-go-playground)
 - [Whitelist](#whitelist-go-playground)
 - [Binary file vs text file](#binary-file-vs-text-file-go-playground)

### Detect MIME type
Get the MIME type from a path to a file. [<kbd>Go Playground</kbd>](https://play.golang.org/p/AkZOfZffk76)
```go
file := "testdata/pdf.pdf"
mime, err := mimetype.DetectFile(file)
fmt.Println(mime.String(), mime.Extension(), err)
// Output: application/pdf .pdf nil
```
Get the MIME type from a reader. [<kbd>Go Playground</kbd>](https://play.golang.org/p/vuftk8TPJwX)
```go
reader, _ := os.Open(file) // ignoring error for brevity's sake
mime, err := mimetype.DetectReader(reader)
fmt.Println(mime.String(), mime.Extension(), err)
// Output: application/pdf .pdf nil
```

Get the MIME type from a byte slice. [<kbd>Go Playground</kbd>](https://play.golang.org/p/uriZsOaEUI1)
```go
data, _ := ioutil.ReadFile(file) // ignoring error for brevity's sake
mime := mimetype.Detect(data)
fmt.Println(mime.String(), mime.Extension())
// Output: application/pdf .pdf
```

### Test against a MIME type [<kbd>Go Playground</kbd>](https://play.golang.org/p/H0ooIXD2N3-)
Test if a file has a specific MIME type. Different from the string comparison,
e.g.: `mime.String() == "application/zip"`, `mime.Is("application/zip")` method
has the following advantages:
 - handles MIME aliases,
 - is case insensitive,
 - ignores optional MIME parameters,
 - ignores any leading and trailing whitespace.
```go
mime, err := mimetype.DetectFile("testdata/zip.zip")
// application/x-zip is an alias of application/zip,
// therefore Is returns true both times.
fmt.Println(mime.Is("application/zip"), mime.Is("application/x-zip"), err)
// Output: true true <nil>
```

### Whitelist [<kbd>Go Playground</kbd>](https://play.golang.org/p/js_1nmWGG3j)
Test if a MIME type is in a list of allowed MIME types.
```go
allowed := []string{"text/plain", "text/html", "text/csv"}
mime, _ := mimetype.DetectFile("/etc/passwd")

if mimetype.EqualsAny(mime.String(), allowed...) {
    fmt.Printf("%s is allowed\n", mime)
} else {
    fmt.Printf("%s is not allowed\n", mime)
}
// Output: text/plain; charset=utf-8 is allowed
```

### Binary file vs text file [<kbd>Go Playground</kbd>](https://play.golang.org/p/CHEFnkn5LQp)
Considering the definition of a binary file as "a computer file that is not
a text file", they can be differentiated by searching for the `text/plain` MIME
in it's MIME hierarchy.
```go
detectedMIME, err := mimetype.DetectFile("testdata/xml.xml")

isBinary := true
for mime := detectedMIME; mime != nil; mime = mime.Parent() {
    if mime.Is("text/plain") {
        isBinary = false
    }
}

fmt.Println(isBinary, detectedMIME, err)
// Output: false text/xml; charset=utf-8 <nil>
```
