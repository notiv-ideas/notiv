package main

import (
	"fmt"

	"github.com/deroproject/graviton"
)

func get(key string, treeName string) {

	snapshot, err := store.LoadSnapshot(*snapshot_version_flag)
	if err != nil {
		panic(err)
	}
	var pass_hash string
	if cfg["defaultDecryptPolicy"] == "true" || (decryption_flag != nil && *decryption_flag) {
		pass_hash = readPassword()
		if checksum != decrypt(cfg["defaultEncryptedCHECKSUM"], pass_hash) {
			fmt.Println("password incorrect, checksum failed to decrypt")
			return
		}
	}
	var tree *graviton.Tree
	if tree_version_flag != nil && *tree_version_flag != 0 {
		tree, err = snapshot.GetTreeWithVersion(treeName, *tree_version_flag)
	} else {
		tree, err = snapshot.GetTree(treeName)
	}
	if err != nil {
		panic(err)
	}
	b, err := tree.Get([]byte(key))
	if err != nil {
		panic(err)
	}
	var value string
	if cfg["defaultDecryptPolicy"] == "true" || (decryption_flag != nil && *decryption_flag) {
		value = decrypt(string(b), pass_hash)
	} else {
		value = string(b)
	}
	if value == "" {
		fmt.Println("Key not found")
	} else {
		fmt.Printf("%s\n", value)
	}
}
