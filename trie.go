package main

import "errors"

// TODO(ym): critbit trees and qp trees to get rid of the map
// rewrite... who the fuck wrote this :^)

type Trie struct {
	children map[rune]*Trie
	items    []string
	myrune   rune // Do I need this??
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
	for i, c := range text {
		if child := trie.children[c]; child != nil {
			trie = child
		} else {
			return trie, i, errors.New("Failed to find items")
		}
	}
	return trie, len(text) - 1, nil
}

// not the most efficient, but it does its job
func (trie *Trie) recurse() []string {
	list := append(make([]string, 0, len(trie.items)), trie.items...)
	for _, elem := range trie.children {
		list = append(list, elem.recurse()...)
	}
	return list
}

func (trie *Trie) GetRecurse(text string) ([]string, error) {
	tree, _, err := trie.getTree(text)
	if err == nil {
		return tree.recurse(), nil
	}
	return nil, err
}

// Should I copy the slice?
func (trie *Trie) Get(text string) ([]string, error) {
	tree, _, err := trie.getTree(text)
	if err == nil {
		return tree.items, nil
	}
	return nil, err
}

func (trie *Trie) Add(name, item string) error {
	tree, _, err := trie.getTree(name)
	if err == nil {
		tree.items = append(tree.items, item)
		return nil
	}
	return err
}

func (trie *Trie) Clear(name string) error {
	tree, _, err := trie.getTree(name)
	if err == nil {
		tree.items = nil
		return nil
	}
	return err
}

func (trie *Trie) AddBulk(name string, items []string) error {
	tree, _, err := trie.getTree(name)
	if err == nil {
		tree.items = append(tree.items, items...)
		return nil
	}
	return err
}
