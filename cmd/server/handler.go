package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gqgs/httpfs/pkg/server"
)

func handler(opts options) error {
	root, err := os.OpenRoot(opts.root)
	if err != nil {
		return fmt.Errorf("failed to open root: %w", err)
	}
	srv := server.New(root.FS())

	http.HandleFunc("GET /", srv.ListFolder)
	http.HandleFunc("GET /{file}", srv.GetFile)
	http.HandleFunc("GET /random", srv.RandomImage)

	log.Printf("Server is running. Visit http://localhost:%d/", opts.port)

	return http.ListenAndServe(":"+fmt.Sprint(opts.port), nil)
}
