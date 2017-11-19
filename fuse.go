package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"strings"
)

type filesystem struct {
	pathfs.FileSystem
}

func (fs *filesystem) GetAttr(pathfile string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	// Separate file name from path
	split := strings.Split(pathfile, "/")
	var path, file string
	switch {
	case len(split) == 1:
		file = split[0]
	case len(split) > 1:
		file = split[len(split)-1]
		path = strings.Join(split[:len(split)-1], "/")
	}

	// Get file metadata
	row := db.QueryRow(sqlGetMetaAttr, path, file)
	if row == nil {
		log.Println("sqlGetMetaAttr returned nothing")
		return nil, fuse.ENOENT
	}

	attr := new(fuse.Attr)
	err := row.Scan(&attr.Ino, &attr.Atime, &attr.Ctime, &attr.Mtime, &attr.Uid, &attr.Gid, &attr.Size, &attr.Mode)
	if err != nil {
		log.Println("Error scanning sqlGetMetaAttr:", err)
		return nil, fuse.ENOENT
	}

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
	/*if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}

	// root files
	if FileHierarchy[path].Size != 0 {
		return nodefs.NewDataFile(FileHierarchy[path].Content), fuse.OK
	}

	// the only other files are in tier 3,
	// so exit if anything other than 3rd tier
	split := strings.Split(path, "/")
	if len(split) != 3 {
		return nil, fuse.ENOENT
	}

	if DirHierarchy[split[0]].Leases != nil &&
		DirHierarchy[split[0]].Rx.MatchString(split[1]) {

		leases, err := DirHierarchy[split[0]].Leases()
		if err != nil {
			return nil, fuse.ENOENT
		}

		if leases[split[1]].IP != "" &&
			LeaseHierarchy[split[2]].Size != 0 {

			f := leases[split[1]].Property(split[2])
			if f == nil {
				return nil, fuse.ENOENT
			}

			return nodefs.NewDataFile(f.Content), fuse.OK

		}
	}*/

	return nil, fuse.ENOENT
}
