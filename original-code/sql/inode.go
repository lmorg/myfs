package sql

const (
	UpdateContentsByInode = `UPDATE
							file
						SET
							file.contents = ?
						WHERE
							file.inode = ?`
)
