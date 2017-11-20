package fusefile

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/lmorg/godbfs/sql"
)

func (f *File) Write(data []byte, off int64) (written uint32, code fuse.Status) {
	_, err := Db.Exec(sql.UpdateContentsByInode, data, f.inode)
	if err != nil {
		return 0, fuse.ENOENT
	}

	return uint32(len(data)), fuse.OK
}
