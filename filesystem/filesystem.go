package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/lmorg/myfs/file"
	"github.com/lmorg/myfs/sql"
	"golang.org/x/net/context"
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
	err := sql.ScanRec(
		sql.QueryRec(sql.GetFileInode, d.inode, name),
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
	q := sql.QueryRows(sql.GetDirContents, d.inode)
	err := sql.ScanRows(q, callback, &inode, &name)

	return dirs, err
}
