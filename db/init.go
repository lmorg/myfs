package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hanwen/go-fuse/fuse"
	sequal "github.com/lmorg/myfs/sql"
	"log"
	"strings"
	"time"
)

var db *sql.DB

func Initalise() {
	var err error

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		"myfs", "qwerty", "localhost", 3306, "myfs")

	log.Println("Opening database:", conn)

	db, err = sql.Open("mysql", conn)
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	if _, err = db.Exec(sequal.CreateMetaTable); err != nil {
		log.Fatalln("(CreateMetaTable):", err, sequal.CreateMetaTable)
	}

	if _, err = db.Exec(sequal.CreateFileTable); err != nil {
		log.Fatalln("(CreateFileTable):", err)
	}

	if _, err = db.Exec(sequal.CreateDirTable); err != nil {
		log.Fatalln("(CreateDirTable):", err)
	}

	epoch := time.Now().Unix()
	if _, err = db.Exec(sequal.InsertMeta, epoch, epoch, epoch, 0, 0, 0, fuse.S_IFDIR|0777, "", 0); err != nil {
		if !strings.HasPrefix(err.Error(), "Error 1062") {
			log.Fatalln("Could not initialise meta table:", err)
		}
	}

	if _, err = db.Exec(sequal.InsertDir, 1, ""); err != nil {
		if !strings.HasPrefix(err.Error(), "Error 1062") {
			log.Fatalln("Could not initialise dir table:", err)
		}
	}

	return
}
