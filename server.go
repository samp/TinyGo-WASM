// Fetch the gzip package. This is skippable when go.mod is set up correctly.
// go get github.com/NYTimes/gziphandler
// Set environment/build target. Both of my machines are windows/amd64.
// $Env:GOOS = "windows"; $Env:GOARCH = "amd64"
// Start the server
// go run server.go
// http://127.0.0.1:5501/

// This is from https://dev.bitolog.com/minimizing-go-webassembly-binary-size/
// It's also totally optional - Live Server works but will not gzip

package main

import (
	"log"
	"net/http"
	"github.com/NYTimes/gziphandler"
)

func main() {
	err := http.ListenAndServe(":5501",
		gziphandler.GzipHandler(http.FileServer(http.Dir("."))))
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
