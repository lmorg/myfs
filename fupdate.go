package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"log"
	"time"
)

// Create a file
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

	// Create file metadata
	epoch := time.Now().Unix()
	result, err := db.Exec(sqlInsertMeta, epoch, epoch, epoch, context.Uid, context.Gid, 0, mode, file, parent)
	if err != nil {
		log.Println("Error sqlInsertMeta:", err)
		return nil, fuse.ENOENT
	}

	// Get file inode
	inode, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting inode of sqlInsertMeta:", err)
		return nil, fuse.ENOENT
	}

	// Write file to table
	_, err = db.Exec(sqlInsertFile, inode, "")
	if err != nil {
		log.Println("Could not insert into file table:", err)
		return nil, fuse.ENOENT
	}

	f := nodefs.NewDataFile([]byte{})
	return f, fuse.OK
}

// Utimens updates the access and modify times
func (fs *filesystem) Utimens(pathfile string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	path, file := SplitPath(pathfile)

	_, err := db.Exec(sqlUpdateTime, Atime.Unix(), Mtime.Unix(), path, file)
	if err != nil {
		return fuse.ENOENT
	}

	return fuse.OK
}
