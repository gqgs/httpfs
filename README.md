# httpfs

**httpfs** is a FUSE-based filesystem that enables mounting remote directories over HTTP, providing local access to remote files in a read-only mode. This project focuses exclusively on __read-only__ access, which simplifies maintenance and reduces the attack surface compared to traditional network filesystems like SSHFS or Samba.

## Features

- 🔒 **Read-Only Access:** Ensures the integrity of remote files by preventing local modifications.
- 🚀 **FUSE Integration:** Seamlessly integrates remote files into the local filesystem.
- 🌐 **HTTP Protocol:** Communicates with an HTTP server to serve files and directory structures.

## Requirements

- Linux Kernel 2.6.9 or later (FUSE support required)
- FUSE library installed ([FUSE project](https://github.com/libfuse/libfuse))
