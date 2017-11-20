package sql

// Initialisation:
const (
	CreateMetaTable = `CREATE TABLE IF NOT EXISTS meta (
							inode       INT AUTO_INCREMENT PRIMARY KEY,
							atime       INT,
							ctime       INT,
							mtime       INT,
							uid			INT,
							gid			INT,
							size		INT,
							mode		INT,
							name		VARCHAR(191) UNIQUE,
							parent		INT
						);`

	CreateFileTable = `CREATE TABLE IF NOT EXISTS file (
							inode       INT PRIMARY KEY,
							contents	BLOB
						);`

	CreateDirTable = `CREATE TABLE IF NOT EXISTS dir (
							inode       INT PRIMARY KEY,
							path		VARCHAR(10000)
						);`
)
