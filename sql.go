package main

// Initialisation:
const (
	sqlCreateMetaTable = `CREATE TABLE IF NOT EXISTS meta (
							inode       INT AUTO_INCREMENT PRIMARY KEY,
							atime       INT,
							ctime       INT,
							mtime       INT,
							uid			INT,
							gid			INT,
							size		INT,
							mode		INT,
							name		VARCHAR(767) UNIQUE,
							parent		INT
						);`

	sqlCreateDataTable = `CREATE TABLE IF NOT EXISTS data (
							inode       INT PRIMARY KEY,
							data		BLOB
						);`

	sqlCreateDirTable = `CREATE TABLE IF NOT EXISTS dir (
							inode       INT PRIMARY KEY,
							path		VARCHAR(10000)
						);`
)

// Query:
const (
	sqlGetMetaAttr = `SELECT
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
							meta.inode = dir.inode
							AND dir.path = ?
							AND meta.name = ?`

	sqlGetDirInode    = `SELECT inode FROM dir WHERE path = ?`
	sqlGetDirContents = `SELECT mode, inode, name FROM meta WHERE inode = ?`
)

// Modify:
const (
	sqlInsertMeta = `INSERT INTO
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

	sqlInsertDir = `INSERT INTO
                            dir
                                (
										inode,
										path
                                )
                            values
                                (
                                    ?, ?
                                )`
)
