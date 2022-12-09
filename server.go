//go get github.com/NYTimes/gziphandler
//$Env:GOOS = "windows"; $Env:GOARCH = "amd64"
//go run server.go

// This is from https://dev.bitolog.com/minimizing-go-webassembly-binary-size/

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
