package sql

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
)
