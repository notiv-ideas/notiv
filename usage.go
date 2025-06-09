package main

import "fmt"

func printUsageShort() {
	fmt.Println(`notiv ` + appVersion + ` — Cryptographically versioned key-value CLI (powered by graviton)

Usage: notiv [global flags] <command> [subcommand] [arguments...]

Run 'notiv help' or 'notiv --help' for detailed usage.`)
}

func printUsage() {
	fmt.Println(`notiv ` + appVersion + ` — Encrypted, versioned, and taggable key-value CLI  
Built on graviton: a Merkle-tree engine by CaptainDero for the DERO project  
https://github.com/deroproject/graviton

notiv is a Git-like interface to authenticated key-value, append-only trees with full snapshot history,  
tag support, and optional encryption. Designed for programmable data workflows,  
backed by cryptographic proofs and an append-only structure.

USAGE:
  notiv [global flags] <command> [subcommand] [arguments...]

GLOBAL FLAGS:
  -n                       Run in non-interactive mode (used for password input)
  --help                   Show this usage message
  --ram                    Run with an in-memory database (ephemeral; useful for testing)
  --share=PATH             Use a specific share directory (overrides config)
  --tags="TAG1 TAG2"       Apply one or more tags to the commit
  --snapshot-version=N     Operate on a specific snapshot version
  --tree-version=N         Operate on a specific tree version
  --encrypt                Encrypt values when writing
  --decrypt                Decrypt values when reading
  --password="PASSWORD"    Supply encryption/decryption password

COMMANDS:

  help  
    Show general help or detailed help for a specific command

  config  
    list  
      Show current configuration values  
    edit <KEY> <VALUE>  
      Update a configuration key (only known keys may be modified)

  share  
    new <NAME | /full/path>  
      Create a new share (notiv database root) by name or at a specific path

  put <KEY> <VALUE> <TREE>  
      Insert or update a key-value pair in the specified tree(s)

  get <KEY> <TREE>  
      Retrieve a value by key from the specified tree(s)

  list  
    trees  
      List all trees in the current snapshot (coming soon)  
    keys <TREE>  
      List all keys from the specified tree(s)  
    pairs <TREE>  
      List all key-value pairs from the specified tree(s)

  version  
    snapshot  
      Show the current snapshot version  
    current <TREE>  
      Show the current version of the specified tree(s)  
    parent <TREE>  
      Show the parent version of the specified tree(s)

See 'notiv help' or 'notiv --help' for more details.  
Inspired by ZFS. Powered by notiv and graviton.`)
}

func printConfigUsage() {
	fmt.Println(`Usage:
  notiv config list
      Show current configuration values

  notiv config edit <KEY> <VALUE>
      Edit an existing config key (only known keys can be modified)

Examples:
  notiv config list
  notiv config edit defaultEncryptDecrypt true`)
}

func printShareUsage() {
	fmt.Println(`Usage:
  notiv share new <NAME>
      Create a new share with the given name in the default notiv directory:
	  ` + cfg["defaultnotivDir"] + `

  notiv share new </full/path/to/dir>
      Create a new share at the specified absolute path

Examples:
  notiv share new demo
  notiv share new /tmp/my-notiv-store
  notiv -ram share new`)
}

func printPutUsage() {
	fmt.Println(`Usage:
  notiv put <KEY> <VALUE> <TREE>
      Insert or update a key-value pair in one or more trees

Optional Behavior:
  - If encryption is enabled (via config or -encrypt), the value will be encrypted

Examples:
  notiv put username alice users
  notiv -encrypt put token secret123 sessions
  notiv --tags="v1.0 stable" put config '{"debug":true}' appsettings logs`)
}

func printGetUsage() {
	fmt.Println(`Usage:
  notiv get <KEY> <TREE>
      Retrieve the value for a given key from one or more trees

Optional Behavior:
  - If decryption is enabled (via config or -decrypt), values will be decrypted
  - If the key is not found, a "Key not found" message will be shown

Examples:
  notiv get username users
  notiv --decrypt get token sessions`)
}

func printListUsage() {
	fmt.Println(`Usage:
  notiv list trees
      TODO Show all trees found in the current snapshot

  notiv list keys <TREE>
      List all keys from one or more trees

  notiv list pairs <TREE>
      List all key-value pairs from one or more trees

Examples:
  TODO notiv list trees
  notiv list keys users
  notiv list pairs users products`)
}

func printVersionUsage() {
	fmt.Println(`Usage:
  notiv version
      Show the CLI version of notiv

  notiv version snapshot
      Show the current snapshot version

  notiv version current <TREE>
      Show the current version number of specified tree

  notiv version parent <TREE>
      Show the parent version number of specified tree

Examples:
  notiv version
  notiv version snapshot
  notiv version current users posts
  notiv version parent users`)
}
