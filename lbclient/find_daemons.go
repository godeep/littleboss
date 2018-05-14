package lbclient // import "crawshaw.io/littleboss/lbclient"

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"crawshaw.io/littleboss/lbrpc"
)

func FindDaemons() (clients []*lbrpc.Client, err error) {
	var socketpaths []string
	uid, gid := syscall.Getuid(), syscall.Getgid()
	tempDir, err := os.Open(os.TempDir())
	if err != nil {
		return nil, err
	}
	fis, err := tempDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		if !strings.HasPrefix(fi.Name(), "littleboss-") || !fi.IsDir() {
			continue
		}
		if stat := fi.Sys().(*syscall.Stat_t); int(stat.Uid) != uid || int(stat.Gid) != gid {
			log.Printf("%s has bad UID/GID %d/%d, want %d/%d", filepath.Join(os.TempDir(), fi.Name()), stat.Uid, stat.Gid, uid, gid)
			continue
		}
		if fi.Mode().Perm() != 0700 {
			log.Printf("%s has bad permissions: %s, want 0700", filepath.Join(os.TempDir(), fi.Name()), fi.Mode().Perm(), gid)
			continue
		}
		dir, err := os.Open(filepath.Join(os.TempDir(), fi.Name()))
		if err != nil {
			log.Print(err)
			continue
		}
		dirFIs, err := dir.Readdir(1)
		dir.Close()
		if err != nil {
			log.Print(err)
			continue
		}
		if len(dirFIs) == 0 {
			log.Printf("%s: no socket file", filepath.Join(os.TempDir(), fi.Name()))
			continue
		}
		ffi := dirFIs[0]
		if !strings.HasPrefix(ffi.Name(), "littleboss.") {
			continue
		}
		if ffi.Mode()&os.ModeSocket != os.ModeSocket {
			log.Printf("%s: not a socket", filepath.Join(os.TempDir(), fi.Name(), ffi.Name()))
			continue
		}
		socketpaths = append(socketpaths, filepath.Join(os.TempDir(), fi.Name(), ffi.Name()))
	}

	ch := make(chan *lbrpc.Client, len(socketpaths))
	for _, socketpath := range socketpaths {
		go func(socketpath string) {
			c, err := lbrpc.NewClient(socketpath)
			if err != nil {
				log.Printf("%s: %v", socketpath, err)
			}
			ch <- c
		}(socketpath)
	}
	for range socketpaths {
		if s := <-ch; s != nil {
			clients = append(clients, s)
		}
	}

	return clients, nil
}