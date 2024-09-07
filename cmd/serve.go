package cmd

import (
	"AmixCrypto/controllers"
	"net/http"
)

func Serve() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", controllers.Index)

	return mux, nil
}
