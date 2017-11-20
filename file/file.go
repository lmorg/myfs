package file

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/lmorg/myfs/db"
	"github.com/lmorg/myfs/sql"
	"golang.org/x/net/context"
	"log"
	"time"
)

type File struct {
	inode uint64
}

func New(inode uint64) File { return File{inode: inode} }

func (f File) Attr(ctx context.Context, a *fuse.Attr) (err error) {
	var atime, ctime, mtime int64

	// Get file metadata
	err = db.ScanRec(
		db.QueryRec(sql.GetFileAttr, f.inode),
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

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Println("open ###############")

	//if !req.Flags.IsReadOnly() {
	//	return nil, fuse.Errno(syscall.EACCES)
	//}

	switch {
	case req.Flags&fuse.OpenCreate != 0:
		q := db.Insert(sql.InsertMeta)
		if q.Err != nil {
			return nil, q.Err
		}
		db.Insert(sql.InsertFile, q.Inode)
	}

	//resp.Flags |= fuse.OpenKeepCache

	return f, nil
}

func (f File) ReadAll(ctx context.Context) (b []byte, err error) {
	// Get file contents
	err = db.ScanRec(
		db.QueryRec(sql.GetFileContents, f.inode),
		&b,
	)

	return
}
