package cmd

import (
	"fmt"
	"net/http"
)

func Serve() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", Index)

	return mux, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello From Server!!")
}
