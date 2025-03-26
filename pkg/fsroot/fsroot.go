package fsroot

import (
	"context"
	"log"
	"log/slog"
	"syscall"
	"time"

	"github.com/gqgs/httpfs/pkg/client"
	"github.com/gqgs/httpfs/pkg/fsfile"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var _ = (rootInterface)((*root)(nil))

type rootInterface interface {
	fs.InodeEmbedder
	fs.NodeOnAdder
	fs.NodeGetattrer
}

type root struct {
	fs.Inode
	client    client.Client
	logger    *slog.Logger
	createdAt uint64
}

func New(client client.Client) *root {
	logger := slog.Default().WithGroup("fsroot")
	logger.Debug("creating new root")
	return &root{
		client:    client,
		logger:    logger,
		createdAt: uint64(time.Now().Unix()),
	}
}

func (r *root) OnAdd(ctx context.Context) {
	r.logger.Debug("onAdd called")
	infos, err := r.client.ListDir()
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range infos {
		p := r.EmbeddedInode()
		file := fsfile.New(r.client, info)

		// Create the file. The Inode must be persistent,
		// because its life time is not under control of the
		// kernel.
		child := p.NewPersistentInode(ctx, file, fs.StableAttr{})

		// And add it
		p.AddChild(info.Name, child, false)
	}
}

func (r *root) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	r.logger.Debug("getattr call")
	out.Mode = 0400
	out.Nlink = 1
	out.Mtime = r.createdAt
	out.Atime = r.createdAt
	out.Ctime = r.createdAt
	out.SetTimeout(time.Minute)
	return fs.OK
}
