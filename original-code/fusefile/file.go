package fusefile

import (
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"time"
)

type File struct {
	nodefs.File
	inode uint64
}

func NewFile(inode uint64) nodefs.File {
	f := new(File)
	f.inode = inode
	return f
}

func (f *File) SetInode(inode *nodefs.Inode)                                     {}
func (f *File) String() string                                                   { return fmt.Sprint("inode:", f.inode) }
func (f *File) InnerFile() nodefs.File                                           { return f }
func (f *File) Release()                                                         {}
func (f *File) Flock(flags int) fuse.Status                                      { return fuse.OK }
func (f *File) Flush() fuse.Status                                               { return fuse.OK }
func (f *File) Fsync(flags int) (code fuse.Status)                               { return fuse.OK }
func (f *File) Truncate(size uint64) fuse.Status                                 { return fuse.OK }
func (f *File) GetAttr(out *fuse.Attr) fuse.Status                               { return fuse.OK }
func (f *File) Chown(uid uint32, gid uint32) fuse.Status                         { return fuse.OK }
func (f *File) Chmod(perms uint32) fuse.Status                                   { return fuse.OK }
func (f *File) Utimens(atime *time.Time, mtime *time.Time) fuse.Status           { return fuse.OK }
func (f *File) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) { return fuse.OK }

func (f *File) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	rr := fuse.ReadResultData([]byte{})
	return rr, fuse.OK
}
