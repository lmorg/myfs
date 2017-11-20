package filesystem

import "github.com/hanwen/go-fuse/fuse/pathfs"

type Fs struct {
	pathfs.FileSystem
}
