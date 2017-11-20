package file

import (
	"bazil.org/fuse"
	"github.com/lmorg/myfs/sql"
	"golang.org/x/net/context"
	"time"
)

type File struct {
	inode uint64
}

func New(inode uint64) File { return File{inode: inode} }

func (f File) Attr(ctx context.Context, a *fuse.Attr) (err error) {
	var atime, ctime, mtime int64

	// Get file metadata
	err = sql.ScanRec(
		sql.QueryRec(sql.GetFileAttr, f.inode),
		&atime, &ctime, &mtime, &a.Uid, &a.Gid, &a.Size, &a.Mode,
	)

	/*if err != nil {
		return fuse.ENOENT
	}*/

	a.Atime = time.Unix(atime, 0)
	a.Ctime = time.Unix(ctime, 0)
	a.Mtime = time.Unix(mtime, 0)

	return
}

func (f File) ReadAll(ctx context.Context) (b []byte, err error) {
	// Get file contents
	err = sql.ScanRec(
		sql.QueryRec(sql.GetFileContents, f.inode),
		&b,
	)

	return
}
