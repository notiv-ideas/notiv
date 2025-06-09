package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// flags are used to bypass configuration defaults
var (
	ram_flag              *bool   // default: false
	share_flag            *string // default: cfg["defaultShareDir"]
	tags_flag             *string // default: ""
	snapshot_version_flag *uint64 // default: 0
	tree_version_flag     *uint64
	encryption_flag       *bool   // default: false
	decryption_flag       *bool   // default: false
	password_flag         *string // default: ""
	help_flag             *bool   // default: false
	non_interactive_flag  *bool   // default: false
	// file_flag       *string // default: ""
)

func rangeTagsFlag() (tags []string) {
	if tags_flag != nil && *tags_flag != "" {
		for _, tag := range strings.Split(*tags_flag, " ") {
			tags = append(tags, tag)
		}
	}
	return
}

// custom flags and args parser
func parseArgs() {
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	help_flag = flagset.Bool("help", false, "Display help information about the available flags and options")
	non_interactive_flag = flagset.Bool("n", false, "Hide interactive prompts; eg. readpassword")
	ram_flag = flagset.Bool("ram", false, "Enable in-memory operations (do not use disk storage)")
	share_flag = flagset.String("share", "", "Path to shared notiv data for synchronization (e.g., --share=/path/to/data)")
	tags_flag = flagset.String("tags", "", "Comma-separated list of tags to associate with the data (e.g., --tags=\"foo bar\")")
	snapshot_version_flag = flagset.Uint64("snapshot-version", 0, "Specify the version number (e.g., --version=123)")
	tree_version_flag = flagset.Uint64("tree-version", 0, "Specify the version number (e.g., --version=123)")
	encryption_flag = flagset.Bool("encrypt", false, "Enable encryption for the data (default is true)")
	decryption_flag = flagset.Bool("decrypt", false, "Enable decryption of encrypted data (default is false)")

	// The application never stores passwords, it uses a hashed version as a key for encryption/decryption
	password_flag = flagset.String("password", "", "Password for encryption/decryption (hashed to create a key for securing data)")

	var flags []string
	var args []string

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			args = append(args, arg)
		}
	}

	if err := flagset.Parse(flags); err != nil {
		fmt.Println("Error pasing flags", err)
		os.Exit(1)
	}
	os.Args = append([]string{os.Args[0]}, args...)
}
