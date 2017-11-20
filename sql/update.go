package sql

// Update:
const (
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
