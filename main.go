package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"flag"
	"github.com/lmorg/godbfs/file"
	"github.com/lmorg/godbfs/filesystem"
	"github.com/lmorg/godbfs/sql"
	"log"
	"os"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	db := sql.InitDb()
	filesystem.Db = db
	file.Db = db

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("godbfs"),
		fuse.Subtype("mysqlfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("volume"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, filesystem.Fs{})
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}
