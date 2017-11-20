package sql

// Get:
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
