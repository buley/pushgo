package goplus

import (
	"fmt"
	"net/http"
)

func init() {
	// Register a handler for /hello URLs.
	http.HandleFunc("/", hello)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}
