package file

import (
	"bazil.org/fuse"
	"github.com/lmorg/godbfs/sql"
	"golang.org/x/net/context"
	"log"
	"time"
)

type File struct {
	inode uint64
}

func New(inode uint64) File { return File{inode: inode} }

func (f File) Attr(ctx context.Context, a *fuse.Attr) error {
	// Get file metadata
	row := Db.QueryRow(sql.GetFileAttr, f.inode)
	if row == nil {
		log.Println("GetFileAttr returned nothing")
		return fuse.ENOENT
	}

	var atime, ctime, mtime int64
	err := row.Scan(&atime, &ctime, &mtime, &a.Uid, &a.Gid, &a.Size, &a.Mode)
	if err != nil {
		log.Println("Error scanning GetFileAttr:", f.inode, err)
		return fuse.ENOENT
	}
	a.Atime = time.Unix(atime, 0)
	a.Ctime = time.Unix(ctime, 0)
	a.Mtime = time.Unix(mtime, 0)

	/*	meta.atime,
		meta.ctime,
		meta.mtime,
		meta.uid,
		meta.gid,
		meta.size,
		meta.mode*/

	return nil
}

func (f File) ReadAll(ctx context.Context) ([]byte, error) {
	row := Db.QueryRow(sql.GetFileContents, f.inode)
	if row == nil {
		log.Println("sqlGetFileContents returned nothing")
		return nil, fuse.ENOENT
	}

	var b []byte
	err := row.Scan(&b)
	if err != nil {
		log.Println("Error scanning sqlGetFileContents:", f.inode, err)
		return nil, fuse.ENOENT
	}

	return b, nil
}
