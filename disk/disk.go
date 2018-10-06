package disk

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"roggl-server/trie"
)

func Load(fname string) (*trie.Trie, error) {
	if fname == "" {
		return trie.MakeTrie(), nil
	}

	bytes, err := ioutil.ReadFile(fname)
	if err == nil {
		var t *trie.Trie
		err = json.Unmarshal(bytes, &t)
		return t, err
	}

	if os.IsNotExist(err) {
		return trie.MakeTrie(), nil
	}

	return nil, err
}

func Save(trie *trie.Trie, fname string) error {
	if fname == "" {
		return nil
	}

	bytes, err := json.Marshal(trie)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fname, bytes, 0644)
}
