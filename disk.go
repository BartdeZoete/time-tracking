package main

import "roggl-server/disk"

var fname = ""

func loadTrie() {
	trie, err := disk.Load(fname)
	if err != nil {
		panic(err)
	}

	t = trie
}

func saveTrie() {
	if err := disk.Save(t, fname); err != nil {
		panic(err)
	}
}
