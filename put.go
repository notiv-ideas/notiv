package main

import (
	"fmt"

	"github.com/deroproject/graviton"
)

func put(key, value string, treeName string) {
	snapshot, err := store.LoadSnapshot(*snapshot_version_flag)
	if err != nil {
		panic(err)
	}
	useEncryption := cfg["defaultEncryptPolicy"] == "true" || (encryption_flag != nil && *encryption_flag)
	if useEncryption {
		pass_hash := readPassword()
		ok := (checksum == decrypt(cfg["defaultEncryptedCHECKSUM"], pass_hash))
		if !ok {
			fmt.Println("password incorrect, checksum failed to decrypt")
			return
		} else {
			value = encrypt(value, pass_hash)
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
	err = tree.Put([]byte(key), []byte(value))
	if err != nil {
		panic(err)
	}

	if err := tree.Commit(rangeTagsFlag()...); err != nil {
		panic(err)
	}
	fmt.Println(tree.GetName(), key, value, "current tree version:", tree.GetVersion(), "ok")

}
