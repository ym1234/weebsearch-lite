package main

import "errors"
// import "fmt"

// Trie implementation, can only handle LOWERCASE alphabetic characters
type Trie struct {
	children map[rune]*Trie
	items []string
	myrune rune // Do I need this??
}

func New() *Trie {
	return &Trie{children: make(map[rune]*Trie)}
}

func (trie *Trie) Insert(text string, items []string) {
	trie, index, err := trie.getTree(text)
	if err != nil {
		text = text[index:]
		for _, char := range text {
			trie.children[char] = &Trie{make(map[rune]*Trie), nil, char}
			trie = trie.children[char]
		}
	}
	trie.items = append(trie.items, items...)
}

func (trie *Trie) getTree(text string) (*Trie, int, error) {
	for i, c := range text{
		if child := trie.children[c]; child != nil {
			trie = child
		} else { return trie, i, errors.New("Failed to find items") }
	}
	return trie, len(text)-1, nil
}

// not the most efficient, but it does its job
func (trie *Trie) recurse() []string {
	list := make([]string, 0, len(trie.items))
	list = append(list, trie.items...)
	for _, elem := range trie.children {
		list = append(list, elem.recurse()...)
	}
	return list
}

func (trie *Trie) GetRecurse(text string) ([]string, error) {
	if tree, _, err := trie.getTree(text); err == nil {
		return tree.recurse(), nil
	} else { return nil, err }
}

// Should I copy the slice?
func (trie *Trie) Get(text string) ([]string, error) {
	if tree, _, err := trie.getTree(text); err == nil {
		return tree.items, nil
	} else { return nil, err }
}

func (trie *Trie) Add(name, item string) error {
	if tree, _, err := trie.getTree(name); err == nil {
		tree.items = append(tree.items, item)
		return nil
	} else { return err }
}

func (trie *Trie) Clear(name string) error {
	if tree, _, err := trie.getTree(name); err == nil {
		tree.items = []string{}
		return nil
	} else { return err }
}

func (trie *Trie) AddBulk(name string, items []string) error {
	if tree, _, err := trie.getTree(name); err == nil {
		tree.items = append(tree.items, items...)
		return nil
	} else { return err }
}
