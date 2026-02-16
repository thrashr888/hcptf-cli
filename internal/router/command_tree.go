package router

import "strings"

type commandNode struct {
	children map[string]*commandNode
	terminal bool
}

// CommandTree stores command paths as a tokenized tree generated from
// the command registry.
type CommandTree struct {
	root  *commandNode
	roots map[string]struct{}
}

// NewCommandTree builds a command tree from full command paths
// (for example: "workspace list", "organization token delete").
func NewCommandTree(commandPaths []string) *CommandTree {
	t := &CommandTree{
		root: &commandNode{
			children: make(map[string]*commandNode),
		},
		roots: make(map[string]struct{}),
	}

	for _, path := range commandPaths {
		tokens := strings.Fields(strings.TrimSpace(path))
		if len(tokens) == 0 {
			continue
		}

		t.roots[tokens[0]] = struct{}{}

		node := t.root
		for _, token := range tokens {
			child, ok := node.children[token]
			if !ok {
				child = &commandNode{
					children: make(map[string]*commandNode),
				}
				node.children[token] = child
			}
			node = child
		}
		node.terminal = true
	}

	return t
}

func (t *CommandTree) HasRoot(token string) bool {
	if t == nil {
		return false
	}
	_, ok := t.roots[token]
	return ok
}
