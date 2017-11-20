package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hanwen/go-fuse/fuse"
	"log"
	"strings"
	"time"
)

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
							name		VARCHAR(767) UNIQUE,
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

func InitDb() (db *sql.DB) {
	var err error

	log.Println("Opening database")

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		"godbfs", "FbLeX26D8X8VCdqI", "192.168.1.203", 3306, "godbfs"))
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	if _, err = db.Exec(CreateMetaTable); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	if _, err = db.Exec(CreateFileTable); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	if _, err = db.Exec(CreateDirTable); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	epoch := time.Now().Unix()
	if _, err = db.Exec(InsertMeta, epoch, epoch, epoch, 0, 0, 0, fuse.S_IFDIR|0777, "", 1); err != nil {
		if !strings.HasPrefix(err.Error(), "Error 1062") {
			log.Fatalln("Could not initialise meta table:", err)
		}
	}

	if _, err = db.Exec(InsertDir, 1, ""); err != nil {
		if !strings.HasPrefix(err.Error(), "Error 1062") {
			log.Fatalln("Could not initialise dir table:", err)
		}
	}

	return
}
