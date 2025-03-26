package main

import (
	"fmt"
	"log"

	"github.com/gqgs/httpfs/pkg/client"
	"github.com/gqgs/httpfs/pkg/fsroot"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

func handler(o options) error {
	client, err := client.New(o.server)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	rootInode := fsroot.New(client)
	server, err := fs.Mount(o.mountpoint, rootInode, &fs.Options{
		MountOptions: fuse.MountOptions{
			Options: []string{"ro"},
			Debug:   o.verbose,
		},
	})
	if err != nil {
		return err
	}

	log.Printf("Mounted on %s", o.mountpoint)
	log.Printf("Unmount by calling 'fusermount -u %s'", o.mountpoint)

	server.Wait()
	return nil
}
