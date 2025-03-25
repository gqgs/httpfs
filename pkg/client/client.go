package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gqgs/httpfs/pkg/meta"
)

type client struct {
	server *url.URL
	logger *slog.Logger
}

type Client interface {
	ListDir() ([]*meta.FileInfo, error)
	DownloadRange(name string, dest []byte, off, size int) (n int, err error)
}

func New(server string) (*client, error) {
	logger := slog.Default().WithGroup("client")
	parsed, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("failed parsing url: %w", err)
	}
	return &client{
		server: parsed,
		logger: logger,
	}, nil
}

func (c *client) ListDir() ([]*meta.FileInfo, error) {
	c.logger.Debug("list dir request")
	result, err := http.Get(c.server.String())
	if err != nil {
		return nil, fmt.Errorf("failed to list dir: %w", err)
	}
	defer result.Body.Close()

	var files []*meta.FileInfo
	err = json.NewDecoder(result.Body).Decode(&files)
	return files, err
}

func (c *client) DownloadRange(name string, dest []byte, off, size int) (n int, err error) {
	c.logger.Debug("download range call", "name", name, "len(dest)", len(dest), "off", off, "size", size)
	req, err := http.NewRequest(http.MethodGet, c.server.JoinPath(name).String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed creating request: %w", err)
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", off, off+size-1))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed executing request: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAtLeast(resp.Body, dest, size)
}
