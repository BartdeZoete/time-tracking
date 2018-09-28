package trie

import (
	"errors"
	"time"
)

type Value struct {
	Begin *time.Time `json:"begin"`
	End   *time.Time `json:"end"`
	Tags  []string   `json:"tags"`
}

type Trie struct {
	Children map[string]*Trie `json:"children"`
	Value    []Value          `json:"value"`
}

func MakeTrie() *Trie {
	return &Trie{
		Children: make(map[string]*Trie),
	}
}

func (t *Trie) Get(path []string) []Value {
	curr := t

	for len(path) > 0 {
		child, has := curr.Children[path[0]]
		if !has {
			return []Value{}
		}

		curr = child
		path = path[1:]
	}

	return curr.Value
}

func (t *Trie) Add(path []string) error {
	curr := t
	set := false

	for len(path) > 0 {
		if curr.Children[path[0]] == nil {
			curr.Children[path[0]] = MakeTrie()
			set = true
		}

		curr, _ = curr.Children[path[0]]
		path = path[1:]
	}

	if !set {
		return errors.New("project already exists")
	}
	return nil
}

/*
func (t *Trie) IsEmpty() bool {
	return len(t.Children) == 0
}
*/

func (t *Trie) Record(path []string) bool {
	if t.IsRecording() {
		return false
	}

	timeNow := time.Now()
	curr := t

	for len(path) > 0 {
		child, _ := curr.Children[path[0]]
		curr = child

		if child != nil && len(path) == 1 {
			curr.Value = append(curr.Value, Value{
				Begin: &timeNow,
				End:   nil,
			})
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
				begin := child.Value[i].Begin
				end := child.Value[i].End
				if begin != nil && end == nil { // child is being recorded
					now := time.Now()
					child.Value[i].End = &now
					return true
				}
			}

			if stopped := child.Stop(); stopped {
				return true
			}
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
				begin := child.Value[i].Begin
				end := child.Value[i].End
				if begin != nil && end == nil { // child is being recorded
					return true
				}
			}

			child.IsRecording()
		} // else check brothers: next k
	}

	return false
}
