package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/lmorg/myfs/db"
	"github.com/lmorg/myfs/file"
	"github.com/lmorg/myfs/sql"
	"golang.org/x/net/context"
	"log"
	"os"
)

type Fs struct{}

func (Fs) Root() (fs.Node, error) {
	return Dir{inode: 1}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	inode uint64
}

func (d Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = d.inode
	a.Mode = os.ModeDir | 0555
	return nil
}

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	// Get file metadata
	var inode uint64
	err := db.ScanRec(
		db.QueryRec(sql.GetFileInode, d.inode, name),
		&inode,
	)

	if err != nil {
		return nil, fuse.ENOENT
	}

	return file.New(inode), nil
}

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var (
		inode    uint64
		name     string
		dirs     []fuse.Dirent
		callback func()
	)

	callback = func() {
		dirs = append(dirs, fuse.Dirent{
			Inode: inode,
			Name:  name,
			Type:  fuse.DT_File,
		})
	}

	// Get directory contents
	q := db.QueryRows(sql.GetDirContents, d.inode)
	err := db.ScanRows(q, callback, &inode, &name)

	return dirs, err
}

func (d *Dir) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Println("open dir ###############")

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

	return nil, nil
}
