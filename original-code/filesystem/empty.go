// +build ignore

package filesystem

import (
	"github.com/hanwen/go-fuse/fuse"
	"log"
)

func (fs *Fs) Access(pathfile string, mode uint32, context *fuse.Context) fuse.Status {
	log.Println("Access: ", pathfile)
	return fuse.OK
}

func (fs *Fs) Truncate(pathfile string, size uint64, context *fuse.Context) fuse.Status {
	log.Println("Truncate: ", pathfile)
	return fuse.OK
}

func (fs *Fs) GetXAttr(pathfile string, attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
	log.Println("GetXAttr: ", pathfile)
	return []byte{}, fuse.OK
}

func (fs *Fs) ListXAttr(pathfile string, context *fuse.Context) (attributes []string, code fuse.Status) {
	log.Println("ListXAttr: ", pathfile)
	return []string{}, fuse.OK
}

func (fs *Fs) RemoveXAttr(pathfile string, attr string, context *fuse.Context) fuse.Status {
	log.Println("RemoveXAttr: ", pathfile)
	return fuse.OK
}

func (fs *Fs) Chmod(pathfile string, mode uint32, context *fuse.Context) (code fuse.Status) {
	log.Println("Chmod: ", pathfile)
	return fuse.OK
}

func (fs *Fs) Chown(pathfile string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	log.Println("Chown: ", pathfile)
	return fuse.OK
}

func (fs *Fs) Mknod(pathfile string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	log.Println("Mknod: ", pathfile)
	return fuse.OK
}
