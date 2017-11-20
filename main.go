package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"flag"
	"github.com/lmorg/myfs/filesystem"
	"github.com/lmorg/myfs/sql"
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

	sql.InitDb()

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
