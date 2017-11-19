package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hanwen/go-fuse/fuse"
	"log"
	"strings"
	"time"
)

var db *sql.DB

func mysql3Mount() {
	var err error

	log.Println("Opening database")

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		"godbfs", "FbLeX26D8X8VCdqI", "192.168.1.203", 3306, "godbfs"))
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	if _, err = db.Exec(sqlCreateMetaTable); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	if _, err = db.Exec(sqlCreateFileTable); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	if _, err = db.Exec(sqlCreateDirTable); err != nil {
		log.Fatalln("Could not create table:", err)
	}

	epoch := time.Now().UnixNano()
	if _, err = db.Exec(sqlInsertMeta, epoch, epoch, epoch, 0, 0, 0, fuse.S_IFDIR|0777, "", 0); err != nil {
		if !strings.HasPrefix(err.Error(), "Error 1062") {
			log.Fatalln("Could not initialise meta table:", err)
		}
	}

	if _, err = db.Exec(sqlInsertDir, 1, ""); err != nil {
		if !strings.HasPrefix(err.Error(), "Error 1062") {
			log.Fatalln("Could not initialise dir table:", err)
		}
	}
}
