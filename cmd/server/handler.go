package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gqgs/httpfs/pkg/meta"
)

func handler(opts options) error {
	root, err := os.OpenRoot(opts.root)
	if err != nil {
		return fmt.Errorf("failed to open root: %w", err)
	}

	http.HandleFunc("GET /", listFolder(root.FS()))
	http.HandleFunc("GET /{file}", getFile(root.FS()))

	log.Printf("Server is running. Visit http://localhost:%d/", opts.port)

	return http.ListenAndServe(":"+fmt.Sprint(opts.port), nil)
}

func listFolder(rootFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := fs.ReadDir(rootFS, ".")
		if err != nil {
			slog.Error("failed reading dir", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		files := make([]*meta.FileInfo, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				slog.Error("failed reading dir", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// skip hidden files
			if strings.HasPrefix(info.Name(), ".") {
				continue
			}

			files = append(files, &meta.FileInfo{
				Name:    info.Name(),
				Size:    info.Size(),
				Mode:    info.Mode(),
				ModTime: info.ModTime(),
			})
		}

		if err = json.NewEncoder(w).Encode(files); err != nil {
			slog.Error("error encoding files", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func getFile(rootFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, rootFS, r.PathValue("file"))
	}
}
