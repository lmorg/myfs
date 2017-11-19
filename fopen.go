package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"log"
)

// OpenDir scans a directory for contents
func (fs *filesystem) OpenDir(path string, context *fuse.Context) (dir []fuse.DirEntry, code fuse.Status) {
	log.Println("OpenDir:", path)

	// Get inode of directory
	row := db.QueryRow(sqlGetDirInode, path)
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
	rows, err := db.Query(sqlGetDirContents, inode)
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
func (fs *filesystem) Open(pathfile string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	//if flags&fuse.O_ANYWRITE != 0 {
	//	return nil, fuse.EPERM
	//}

	path, file := SplitPath(pathfile)

	row := db.QueryRow(sqlGetFileContents, path, file)
	if row == nil {
		log.Println("sqlGetFileContents returned nothing")
		return nil, fuse.ENOENT
	}

	var b []byte
	err := row.Scan(&b)
	if err != nil {
		log.Println("Error scanning sqlGetFileContents:", err, pathfile)
		return nil, fuse.ENOENT
	}

	f := nodefs.NewDataFile(b)
	log.Println("content:", f.String())
	return f, fuse.OK
}
