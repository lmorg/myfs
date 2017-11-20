package sql

// Query:
const (
	GetMetaAttr = `SELECT
							meta.inode,
							meta.atime,
							meta.ctime,
							meta.mtime,
							meta.uid,
							meta.gid,
							meta.size,
							meta.mode
						FROM
							meta,
							dir
						WHERE
							meta.parent = dir.inode
							AND dir.path = ?
							AND meta.name = ?`

	GetFileContents = `SELECT
							file.inode,
							file.contents
						FROM
							meta,
							dir,
							file
						WHERE
							meta.parent = dir.inode
							AND file.inode = meta.inode
							AND dir.path = ?
							AND meta.name = ?`

	GetDirInode    = `SELECT inode FROM dir WHERE path = ?`
	GetDirContents = `SELECT mode, inode, name FROM meta WHERE parent = ?`
)

// Modify:
const (
	InsertMeta = `INSERT INTO
                            meta
                                (
										atime,
										ctime,
										mtime,
										uid,
										gid,
										size,
										mode,
										name,
										parent
                                )
                            values
                                (
                                    ?, ?, ?, ?, ?, ?, ?, ?, ?
                                )`

	InsertFile = `INSERT INTO
                            file
                                (
										inode,
										contents
                                )
                            values
                                (
                                    ?, ?
                                )`

	InsertDir = `INSERT INTO
                            dir
                                (
										inode,
										path
                                )
                            values
                                (
                                    ?, ?
                                )`

	UpdateTime = `UPDATE
							meta,
							dir
						SET
							meta.atime = ?,
							meta.mtime = ?
						WHERE
							meta.parent = dir.inode
							AND dir.path = ?
							AND meta.name = ?`
)
