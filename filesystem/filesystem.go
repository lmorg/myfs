package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/lmorg/godbfs/file"
	"github.com/lmorg/godbfs/sql"
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
	log.Println("Lookup:", d.inode, name)

	// Get file metadata
	row := Db.QueryRow(sql.GetFileInode, d.inode, name)
	if row == nil {
		log.Println("GetFileInode returned nothing")
		return nil, fuse.ENOENT
	}

	var inode uint64
	err := row.Scan(&inode)
	if err != nil {
		log.Println("Error scanning GetFileInode:", err, name)
		return nil, fuse.ENOENT
	}

	return file.New(inode), nil
	//return nil, fuse.ENOENT
}

//var dirDirs = []fuse.Dirent{
//	{Inode: 2, Name: "hello", Type: fuse.DT_File},
//	}

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("OpenDir:", d.inode)

	/*// Get inode of directory
	row := Db.QueryRow(sql.GetDirInode, path)
	if row == nil {
		log.Println("Nothing returned from sqlGetDirInode")
		return nil, fuse.ENOENT
	}

	var inode uint64
	err := row.Scan(&inode)
	if err != nil {
		log.Println("Error scanning sqlGetDirInode:", err)
		return nil, fuse.ENOENT
	}*/

	dirs := make([]fuse.Dirent, 0)

	// Get directory contents
	rows, err := Db.Query(sql.GetDirContents, d.inode)
	if err != nil || rows == nil {
		log.Println("Error querying sqlGetDirContents:", err)
		return nil, err
	}

	for rows.Next() {
		var inode uint64
		var name string
		err := rows.Scan(&inode, &name)
		if err != nil {
			log.Println("Error scanning sqlGetDirContents:", err)
		}

		dirs = append(dirs, fuse.Dirent{
			Inode: inode,
			Name:  name,
			Type:  fuse.DT_File,
		})

	}

	log.Println("sqlGetDirContents:", dirs)
	return dirs, nil
}
