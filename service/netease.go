package service

import (
	"fmt"
	"net/http"
)

func NeteaseSearchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
