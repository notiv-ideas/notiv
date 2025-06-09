package main

import (
	"fmt"
	"os"

	"github.com/deroproject/graviton"
)

func listTrees() {
	// graviton does not currently support list trees...
	// this is poor construction
	// TODO , duh

	// snapshot, err := store.LoadSnapshot(*snapshot_version_flag)
	// if err != nil {
	// 	panic(err)
	// }
	// file, err := os.ReadFile(shareDir + "/0/0/0/0.dfs") // I can imainge that this isn't a very healthy way of doing this.
	// if err != nil {
	// 	panic(err)
	// }

	// var result []string
	// var buf bytes.Buffer
	// for _, b := range file {
	// 	if b >= 32 && b <= 126 {
	// 		buf.WriteByte(b)
	// 	} else {
	// 		if buf.Len() >= 3 { // I guess we would need to minLimit tree names to 3
	// 			result = append(result, buf.String())
	// 		}
	// 		buf.Reset()
	// 	}
	// }

	// if buf.Len() >= 3 { // I guess we would need to minLimit tree names to 3
	// 	result = append(result, buf.String())
	// }

	// reg := regexp.MustCompile(`:[a-zA-Z0-9_:@-]+`) // yeah, this doesn't work as expected
	// found := make(map[string]struct{})
	// for _, str := range result {
	// 	matches := reg.FindAllString(str, -1)
	// 	for _, match := range matches {
	// 		found[match] = struct{}{}
	// 	}
	// }
	// var names []string
	// for name := range found {
	// 	tree, err := snapshot.GetTree(name[1:])
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	names = append(names, tree.GetName())
	// }

	// sort.Strings(names)

	// for _, name := range names {
	// 	fmt.Println(name)
	// }
	fmt.Println("currently unsupported")
}

func listKeys(treeName string) {
	listWith(treeName, func(k, _ []byte) {
		fmt.Printf("%s\n", k)
	})
}

func listPairs(treeName string) {
	listWith(treeName, func(k, v []byte) {
		fmt.Printf("%s %s\n", k, v)
	})
}

func listWith(treeName string, handle func(k, v []byte)) {
	snapshot, err := store.LoadSnapshot(*snapshot_version_flag)
	if err != nil {
		panic(err)
	}
	var pass_hash string
	useDecryption := cfg["defaultEncryptDecrypt"] == "true" || (decryption_flag != nil && *decryption_flag)
	if useDecryption {
		pass_hash = readPassword()
		if checksum != decrypt(cfg["defaultEncryptedCHECKSUM"], pass_hash) {
			fmt.Println("password incorrect, checksum failed to decrypt")
			os.Exit(1)
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

	fmt.Println(tree.GetName())
	c := tree.Cursor()
	for k, v, err := c.First(); err == nil && k != nil; k, v, err = c.Next() {
		var value string
		if useDecryption {
			value = decrypt(string(v), pass_hash)
		} else {
			value = string(v)
		}
		handle(k, []byte(value))
	}
}
