package main

import (
	"log/slog"
	"os"
)

//go:generate go tool argsgen

type options struct {
	mountpoint string `arg:"mountpoint for fuse root,required"`
	server     string `arg:"server serving files,required"`
	debug      bool   `arg:"enable debug mode"`
	verbose    bool   `arg:"enable verbose debug mode"`
}

func main() {
	o := options{
		mountpoint: os.Getenv("HTTPFS_MOUNTPOINT"),
		server:     os.Getenv("HTTPFS_SERVER"),
	}
	o.MustParse()

	if o.debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if err := handler(o); err != nil {
		slog.Error(err.Error())
	}
}
