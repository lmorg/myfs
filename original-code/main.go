package original_code

import (
	"flag"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/lmorg/godbfs/filesystem"
	"github.com/lmorg/godbfs/fusefile"
	"github.com/lmorg/godbfs/sql"
	"log"
	"os"
)

var (
	fMountPoint string
)

func main() {
	Flags()
	Mount()
}

func Flags() {
	flag.StringVar(&fMountPoint, "m", "", "mount point")

	flag.Parse()

	if fMountPoint == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func Mount() error {
	db := sql.InitDb()
	filesystem.Db = db
	fusefile.Db = db

	vfs := pathfs.NewPathNodeFs(
		&filesystem.Fs{FileSystem: pathfs.NewDefaultFileSystem()}, nil,
	)

	server, _, err := nodefs.MountRoot(fMountPoint, vfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount failed (%s): %v\n", fMountPoint, err)
	}
	server.Serve()

	return nil
}
