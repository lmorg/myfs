package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"strings"
	"time"
)

type filesystem struct {
	pathfs.FileSystem
}

// SplitPath separates file name from path
func SplitPath(pathfile string) (path, file string) {
	split := strings.Split(pathfile, "/")

	switch {
	case len(split) == 1:
		file = split[0]
	case len(split) > 1:
		file = split[len(split)-1]
		path = strings.Join(split[:len(split)-1], "/")
	}

	return
}

func (fs *filesystem) GetAttr(pathfile string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	path, file := SplitPath(pathfile)

	// Get file metadata
	row := db.QueryRow(sqlGetMetaAttr, path, file)
	if row == nil {
		log.Println("sqlGetMetaAttr returned nothing")
		return nil, fuse.ENOENT
	}

	attr := new(fuse.Attr)
	err := row.Scan(&attr.Ino, &attr.Atime, &attr.Ctime, &attr.Mtime, &attr.Uid, &attr.Gid, &attr.Size, &attr.Mode)
	if err != nil {
		log.Println("Error scanning sqlGetMetaAttr:", err, pathfile)
		return nil, fuse.ENOENT
	}

	//log.Println(attr)

	return attr, fuse.OK
}

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

func (fs *filesystem) Open(path string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}

	return nil, fuse.ENOENT
}

func (fs *filesystem) Create(pathfile string, flags uint32, mode uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	path, file := SplitPath(pathfile)

	log.Println("Create flags:", flags)

	// Get inode of parent directory
	row := db.QueryRow(sqlGetDirInode, path)
	if row == nil {
		log.Println("Nothing returned from sqlGetDirInode")
		return nil, fuse.ENOENT
	}

	var parent uint64
	err := row.Scan(&parent)
	if err != nil {
		log.Println("Error scanning sqlGetDirInode:", err)
		return nil, fuse.ENOENT
	}

	epoch := time.Now().Unix()
	result, err := db.Exec(sqlInsertMeta, epoch, epoch, epoch, context.Uid, context.Gid, 0, mode, file, parent)
	if err != nil {
		log.Println("Error sqlInsertMeta:", err)
		return nil, fuse.ENOENT
	}

	inode, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting inode of sqlInsertMeta:", err)
		return nil, fuse.ENOENT
	}

	_, err = db.Exec(sqlInsertFile, inode, "")
	if err != nil {
		log.Println("Could not initialise dir table:", err)
		return nil, fuse.ENOENT
	}

	f := nodefs.NewDataFile([]byte{})
	return f, fuse.OK
}

func (fs *filesystem) Utimens(pathfile string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	path, file := SplitPath(pathfile)

	_, err := db.Exec(sqlUpdateTime, Atime.Unix(), Mtime.Unix(), path, file)
	if err != nil {
		return fuse.ENOENT
	}

	return fuse.OK
}
