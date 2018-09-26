package main

import (
	"fmt"
	"time"
	//"errors"
	"bufio"
	"os"
	"strings"
	"encoding/json"
	//"io/ioutil"
	"log"
	"net/http"
)

type Trie struct {
	Children map[string]*Trie
	Value []Entry
}

type Entry struct {
	Begin time.Time
	End time.Time
}

func MakeTrie() *Trie {
	return &Trie {
		Children: make(map[string]*Trie),
	}
}

func (t *Trie) Get(path []string) []Entry {
	curr := t

	for len(path) > 0 {
		child, has := curr.Children[path[0]]
		if !has {
			return []Entry{}
		}

		curr = child
		path = path[1:]
	}

	return curr.Value
}

func (t *Trie) Add(path []string) { // TODO add message when project already exists
	curr := t

	for len(path) > 0 {
		if curr.Children[path[0]] == nil {
			curr.Children[path[0]] = MakeTrie()
		}

		curr, _ = curr.Children[path[0]]
		path = path[1:]
	}
}

/*
func (t *Trie) IsEmpty() bool {
	return len(t.Children) == 0
}
*/

func (t *Trie) Record(path []string) bool {
	timeNow := time.Now()

	if t.IsRecording() {
		return false
	}

	curr := t

	for len(path) > 0 {
		child, _ := curr.Children[path[0]]
		curr = child

		if child != nil && len(path) == 1 {
			curr.Value = append(curr.Value, Entry{Begin: timeNow, End: time.Time{}})
			return true
		}

		path = path[1:]
	}

	return false
}

func (t *Trie) Stop() bool { // TODO fix not recording message caused by recursion
	curr := t

	for k, _ := range curr.Children {
		child, _ := curr.Children[k]

		if child != nil { // child is a contender
			entryCount := len(child.Value)

			if entryCount > 0 { // still a contender
				i := entryCount - 1
				beginSet := child.Value[i].Begin != time.Time{}
				endSet := child.Value[i].End != time.Time{}
				if beginSet && !endSet { // child is being recorded
					child.Value[i].End = time.Now()
					return true
				}
			}

			child.Stop()
		} // else check brothers: next k
	}

	return false
}

func (t *Trie) IsRecording() bool {
	curr := t

	for k, _ := range curr.Children {
		child, _ := curr.Children[k]

		if child != nil { // child is a contender
			entryCount := len(child.Value)

			if entryCount > 0 { // still a contender
				i := entryCount - 1
				beginSet := child.Value[i].Begin != time.Time{}
				endSet := child.Value[i].End != time.Time{}
				if beginSet && !endSet { // child is being recorded
					return true
				}
			}

			child.IsRecording()
		} // else check brothers: next k
	}

	return false
}

func printTrie(root *Trie) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	encoder.Encode(root)
}

func parsePath(path string) []string{
	return strings.Split(path, ";")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[0:])
}

func entriesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Oke %s!" + " 1", r.URL.Path[0:])
}

func main() {
	http.HandleFunc("/entries", handler)
	http.HandleFunc("/entries/tik", entriesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	t := MakeTrie()
	t.Children["gitaar"] = MakeTrie()
	t.Children["gitaar"].Children["oefenen"] = MakeTrie()
	t.Children["lezen"] = MakeTrie()
	reader := bufio.NewReader(os.Stdin)
	//path1 := []string{"Gitaar", "Oefenen", "Canco del Lladre"}
	//path2 := []string{"Gitaar", "Oefenen"}

	for {
		str, _ := reader.ReadString('\n')
		str = strings.TrimSpace(str)
		fields := strings.Fields(str)
		var path []string

		if len(fields) > 1 {
			path = parsePath(fields[1])
			fmt.Println(path)
		}

		switch fields[0] {
		case "a":
			t.Add(path)
		case "r":
			if !t.Record(path) {
				fmt.Println("Project not found.")
			}
		case "s":
			if !t.Stop() {
				fmt.Println("Not recording.")
			}
		case "print":
			printTrie(t)
		case "kaas":
			fmt.Println("Nice.")
		case "zaad":
			fmt.Println("Bah.")
		case "quit": // code for saving
			return
		default:
			fmt.Println("Default")
		}
	}
}
