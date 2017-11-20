package sql

// Query:
const (
	GetFileInode = `SELECT
							meta.inode
						FROM
							meta
						WHERE
							meta.parent = ?
							AND meta.name = ?`

	GetFileAttr = `SELECT
							meta.atime,
							meta.ctime,
							meta.mtime,
							meta.uid,
							meta.gid,
							meta.size,
							meta.mode
						FROM
							meta
						WHERE
							meta.inode = ?`

	GetFileContents = `SELECT
							file.contents
						FROM
							file
						WHERE
							file.inode = ?`

	//GetDirInode    = `SELECT inode FROM dir WHERE path = ?`
	GetDirContents = `SELECT inode, name FROM meta WHERE parent = ?`
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
