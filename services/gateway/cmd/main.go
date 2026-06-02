// Package cmd
package cmd

import (
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request \n")
	io.WriteString(w, "This is my website!\n")
}
