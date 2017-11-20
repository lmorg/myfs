// +build ignore

package fusefile

import "github.com/hanwen/go-fuse/fuse"

type File interface {
	// Called upon registering the filehandle in the inode.
	SetInode(*Inode)

	// The String method is for debug printing.
	String() string

	// Wrappers around other File implementations, should return
	// the inner file here.
	InnerFile() File

	Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status)
	Write(data []byte, off int64) (written uint32, code fuse.Status)

	Flock(flags int) fuse.Status

	// Flush is called for close() call on a file descriptor. In
	// case of duplicated descriptor, it may be called more than
	// once for a file.
	Flush() fuse.Status

	// This is called to before the file handle is forgotten. This
	// method has no return value, so nothing can synchronizes on
	// the call. Any cleanup that requires specific synchronization or
	// could fail with I/O errors should happen in Flush instead.
	Release()
	Fsync(flags int) (code fuse.Status)

	// The methods below may be called on closed files, due to
	// concurrency.  In that case, you should return EBADF.
	Truncate(size uint64) fuse.Status
	GetAttr(out *fuse.Attr) fuse.Status
	Chown(uid uint32, gid uint32) fuse.Status
	Chmod(perms uint32) fuse.Status
	Utimens(atime *time.Time, mtime *time.Time) fuse.Status
	Allocate(off uint64, size uint64, mode uint32) (code fuse.Status)
}
