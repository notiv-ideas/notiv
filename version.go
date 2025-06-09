package main

import (
	"fmt"
	"log"

	"github.com/deroproject/graviton"
)

func version() {
	fmt.Println(appVersion)
}

func snapshotVersion() {
	snapshot, err := store.LoadSnapshot(*snapshot_version_flag)
	if err != nil {
		panic(err)
	}
	fmt.Println("snapshot version:", snapshot.GetVersion())
}

func treeVersion(param string, treeName string) {
	snapshot, err := store.LoadSnapshot(*snapshot_version_flag)
	if err != nil {
		panic(err)
	}
	var tree *graviton.Tree
	if tree_version_flag != nil && *tree_version_flag != 0 {
		tree, err = snapshot.GetTreeWithVersion(treeName, *tree_version_flag)
	} else {
		tree, err = snapshot.GetTree(treeName)
	}
	if err != nil {
		log.Printf("Error loading tree %s: %v\n", treeName, err)
		return
	}
	if param == "current" {
		fmt.Printf("%s current version: %d\n", treeName, tree.GetVersion())
	} else if param == "parent" {
		fmt.Printf("%s parent version: %d\n", treeName, tree.GetParentVersion())
	} else {
		fmt.Printf("unknown version param: %s\n", param)
	}
}
