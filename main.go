package main

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/deroproject/graviton"
)

const appVersion = `0.0.1`
const checksum = `ok`

var (
	notivDir   string
	shareDir   string
	projectDir string
	confPath   string
	confDir    string
	cfg        map[string]string

	store *graviton.Store
)

func main() {
	// Step 1: Set up signal catching (SIGINT and SIGTERM for example)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Step 2: Optionally handle cleanup in a goroutine (for graceful shutdown)
	go func() {
		sig := <-sigs
		fmt.Printf("Received signal: %s. Exiting gracefully...\n", sig)
		// You can add any cleanup logic here
		os.Exit(0) // Or use a more graceful shutdown if needed
	}()

	parseArgs()

	if *help_flag {
		printUsage()
		return
	}

	if len(os.Args) <= 1 {
		printUsageShort()
		return
	}

	cmd := os.Args[1]

	// opening sequence
	switch {
	case !slices.Contains(
		[]string{
			// change-password
			// encrypt
			// decrypt
			"help",
			"version",
			"init",
			"config",
			"share",
			"put",
			"get",
			"list",
		}, cmd):
		fmt.Println("notiv: unknown command:", cmd)
		return
	case cmd == "version":
		if len(os.Args) == 2 {
			version()
			return
		}
	case cmd == "help":
		var params string
		if len(os.Args) > 2 {
			params = os.Args[2]
		}
		switch params {
		default:
			printUsage()
			return
		}
	default:
		//
	}

	// init sequence
	switch cmd {
	case "init":
		initCfg()
		return
	default:
		if err := loadCfg(); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("fatal: no notiv configs found, %s\n", os.ErrNotExist)
				fmt.Printf("hint: 'notiv init'\n")
				return
			} else {
				panic(err)
			}
		}
	}

	// config sequence
	switch cmd {
	case "config":
		var param string
		if len(os.Args) > 2 {
			param = os.Args[2]
		}
		switch param {
		case "list":
			listCfg()
			return
		case "edit":
			if len(os.Args) < 5 {
				fmt.Println("Usage: config edit <KEY> <VALUE>")
				return
			}
			key := os.Args[3]
			value := os.Args[4]
			editCfg(key, value)
			return
		default:
			printConfigUsage()
			return
		}
	case "share":
		var param string
		if len(os.Args) > 2 {
			param = os.Args[2]
		}
		switch param {
		case "new":
			newShare()
			return
		default:
			printShareUsage()
			return
		}
	default:
		var err error
		store, err = setShare()
		if err != nil {
			panic(err)
		}
	}

	switch cmd {
	case "put":
		if len(os.Args) < 5 {
			printPutUsage()
			return
		}
		key := os.Args[2]
		value := os.Args[3]
		treeName := os.Args[4]
		put(key, value, treeName)
	case "get":
		if len(os.Args) < 4 {
			printGetUsage()
			return
		}
		key := os.Args[2]
		treeName := os.Args[3]
		get(key, treeName)
	case "list":
		var param string
		if len(os.Args) > 2 {
			param = os.Args[2]
		}
		switch param {
		case "trees":
			listTrees()
		case "keys":
			if len(os.Args) < 4 {
				printListUsage()
				return
			}
			treeName := os.Args[3]
			listKeys(treeName)
		case "pairs":
			if len(os.Args) < 4 {
				printListUsage()
				return
			}
			treeName := os.Args[3]
			listPairs(treeName)
		default:
			fmt.Println("Unknown list subcommand:", param)
			printListUsage()
			return
		}
	case "version":
		if len(os.Args) == 2 {
			version()
			return
		}

		param := os.Args[2]
		treeName := os.Args[3]

		switch param {
		case "snapshot":
			snapshotVersion()
		case "current", "parent":
			if len(treeName) == 0 {
				fmt.Println("Please specify one or more tree names.")
				return
			}
			treeVersion(param, treeName)
		default:
			fmt.Println("Unknown version subcommand:", param)
			printVersionUsage()
			return
		}
	default:
		fmt.Println("Unknown command:", cmd)
		printUsageShort()
		return
	}
}
