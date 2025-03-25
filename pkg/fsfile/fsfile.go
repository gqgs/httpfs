package fsfile

import (
	"context"
	"log/slog"
	"syscall"
	"time"

	"github.com/gqgs/httpfs/pkg/client"
	"github.com/gqgs/httpfs/pkg/meta"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var _ = (fileInterface)((*file)(nil))

type fileInterface interface {
	fs.NodeOpener
	fs.NodeGetattrer
	fs.NodeReader
}

type file struct {
	fs.Inode
	client client.Client
	info   *meta.FileInfo
	logger *slog.Logger
}

func New(client client.Client, info *meta.FileInfo) *file {
	logger := slog.Default().WithGroup("fsfile")
	logger.Debug("creating new file", "info", info)
	return &file{
		logger: logger,
		client: client,
		info:   info,
	}
}

func (f *file) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	f.logger.Debug("file open call", "name", f.info.Name)
	return nil, fuse.FOPEN_KEEP_CACHE, fs.OK
}

func (f *file) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	size := min(int(int64(f.info.Size)-off), len(dest))
	f.logger.Debug("file read call", "name", f.info.Name, "offset", off, "len(dest)", len(dest), "object_size", f.info.Size, "size", size)

	if len(dest) > size {
		dest = dest[:size]
	}

	n, err := f.client.DownloadRange(f.info.Name, dest, int(off), size)
	if err != nil {
		f.logger.Error("file download error", "name", f.info.Name, "err", err, "read_bytes", n, "len(dest)", len(dest))
		return nil, fs.ToErrno(err)
	}

	f.logger.Debug("file read executed", "name", f.info.Name, "read_bytes", n, "len(dest)", len(dest))
	return fuse.ReadResultData(dest[:n]), fs.OK
}

func (f *file) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	f.logger.Debug("file getattr call", "name", f.info.Name)
	out.Mode = uint32(f.info.Mode)
	infoTime := uint64(f.info.ModTime.Unix())
	out.Nlink = 1
	out.Mtime = infoTime
	out.Atime = infoTime
	out.Ctime = infoTime
	out.Size = uint64(f.info.Size)
	out.SetTimeout(time.Minute)
	return fs.OK
}
