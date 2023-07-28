
* Install protobuf: https://github.com/protocolbuffers/protobuf
* Install the Go plugin: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
* Download the correct version of the proto file: https://github.com/sass/sass/blob/main/spec/embedded_sass.proto
* protoc --go_opt=Membedded_sass.proto=github.com/bep/godartsass/internal/embeddedsass --go_opt=paths=source_relative --go_out=. embedded_sass.proto
