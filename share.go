package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deroproject/graviton"
)

func newShare() {
	diskStorePolicy := cfg["defaultDiskStore"] != "true" && cfg["defaultDiskStore"] == "false"

	if diskStorePolicy || *ram_flag {

		if *ram_flag {
			fmt.Println("flag --ram is present, MemStore in use; if present, NAME ignored")
		}

		if share_flag != nil {
			fmt.Println("defaultDiskStore is not \"true\"; MemStore in use, flag ignored")
		}
		_, err := graviton.NewMemStore()
		if err != nil {
			panic(err)
		}
		fmt.Println("notiv created new in-memory store")

	} else {
		var filename string
		isFlagged := share_flag != nil && *share_flag != ""

		if isFlagged {
			if len(os.Args) < 4 {
				filename = "data"
			} else {
				filename = os.Args[3]
			}
			shareDir = filepath.Join(*share_flag, ".notiv", "share")
		} else {
			if shareDir != "" {
				shareDir = filepath.Join(notivDir, "share")
			}
			if len(os.Args) < 4 {
				shareDir = cfg["defaultDataDir"]
			} else {
				if os.Args[3] == "" {
					fmt.Println("share name cannot be empty string")
					return
				}
				arg := os.Args[3]
				var err error
				if filepath.IsAbs(arg) || strings.HasPrefix(arg, ".") || strings.Contains(arg, string(os.PathSeparator)) {
					filename = "data"
					// If it's an absolute path, starts with `.` (like `./share`), or includes path separators
					shareDir, err = filepath.Abs(arg)
					if err != nil {
						fmt.Fprintf(os.Stderr, "invalid path: %v\n", err)
						panic(err)
					}
					shareDir = (filepath.Join(shareDir, notivDir, "share"))
				} else {
					filename = arg
					shareDir = filepath.Join(notivDir, "share")
				}
			}
		}

		storeDir := filepath.Join(shareDir, filename)
		// check if the version root is present
		version_root := filepath.Join(storeDir, "version_root.bin")
		if _, err := os.Stat(version_root); err == nil {
			fmt.Println("Share already created at", storeDir)
			return
		} else if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error checking file: %v\n", err)
			return
		}
		// doesn't overwrite, just opens
		_, err := graviton.NewDiskStore(storeDir)
		if err != nil {
			panic(err)
		}
		fmt.Println("Share created at", storeDir)
	}

}

func setShare() (store *graviton.Store, err error) {
	if cfg["defaultDiskStore"] != "true" &&
		cfg["defaultDiskStore"] == "false" || *ram_flag {
		// take custom shares when appropriate, else default applies
		if share_flag != nil {
			fmt.Println("defaultDiskStore is not \"true\"; MemStore in use, flag ignored")
		}
		store, err = graviton.NewMemStore()
		if err != nil {
			return nil, err
		}
		// fmt.Println("notiv created new in-memory store")
	} else {
		// take custom shares when appropriate, else default applies
		var filename string
		if share_flag != nil && *share_flag != "" {
			filename = "data"
			shareDir = *share_flag
		} else {
			filename = "data"
			shareDir = cfg["defaultShareDir"]
		}
		storeDir := filepath.Join(shareDir, filename)
		if _, err := os.Stat(storeDir); err != nil {
			fmt.Println("create new local notiv share: eg. notiv new share", shareDir)
			os.Exit(1)
		}
		store, err = graviton.NewDiskStore(storeDir)
		if err != nil {
			return nil, err
		}
		// fmt.Println("notiv created new persistent store:", shareDir)
	}
	return
}
