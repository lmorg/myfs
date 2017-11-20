package filesystem

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/lmorg/godbfs/fusefile"
	"github.com/lmorg/godbfs/sql"
	"log"
)

// OpenDir scans a directory for contents
func (fs *Fs) OpenDir(path string, context *fuse.Context) (dir []fuse.DirEntry, code fuse.Status) {
	log.Println("OpenDir:", path)

	// Get inode of directory
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
	}

	// Get directory contents
	rows, err := Db.Query(sql.GetDirContents, inode)
	if err != nil || rows == nil {
		log.Println("Error querying sqlGetDirContents:", err)
		return nil, fuse.ENOENT
	}

	for rows.Next() {
		var attr = new(fuse.Attr)
		var name string
		err := rows.Scan(&attr.Mode, &attr.Ino, &name)
		if err != nil {
			log.Println("Error scanning sqlGetDirContents:", err)
		}

		dir = append(dir, fuse.DirEntry{
			Mode: attr.Mode,
			Ino:  attr.Ino,
			Name: name,
		})

	}

	return dir, fuse.OK
}

// Open a file
func (fs *Fs) Open(pathfile string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	//if flags&fuse.O_ANYWRITE != 0 {
	//	return nil, fuse.EPERM
	//}

	path, file := SplitPath(pathfile)

	row := Db.QueryRow(sql.GetFileContents, path, file)
	if row == nil {
		log.Println("sqlGetFileContents returned nothing")
		return nil, fuse.ENOENT
	}

	var b []byte
	var inode uint64
	err := row.Scan(&inode, &b)
	if err != nil {
		log.Println("Error scanning sqlGetFileContents:", err, pathfile)
		return nil, fuse.ENOENT
	}

	//f := nodefs.NewDataFile(b)
	//log.Println("content:", f.String())
	f := fusefile.NewFile(inode)
	return f, fuse.OK
}
