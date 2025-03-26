package server

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gqgs/httpfs/pkg/meta"
)

type server struct {
	root fs.FS
}

func New(root fs.FS) *server {
	return &server{
		root: root,
	}
}

func (s *server) ListFolder(w http.ResponseWriter, r *http.Request) {
	entries, err := fs.ReadDir(s.root, ".")
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

func (s *server) GetFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, s.root, r.PathValue("file"))
}

func (s *server) RandomImage(w http.ResponseWriter, r *http.Request) {
	matches, _ := fs.Glob(s.root, "*")
	randomFile := matches[rand.Intn(len(matches))]
	http.ServeFileFS(w, r, s.root, randomFile)
}
