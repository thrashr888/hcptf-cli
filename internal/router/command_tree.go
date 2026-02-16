package router

import "strings"

type commandNode struct {
	children map[string]*commandNode
	terminal bool
}

// CommandTree stores command paths as a tokenized tree generated from
// the command registry.
type CommandTree struct {
	root                  *commandNode
	treeRoots             map[string]struct{}
	orgCollection         map[string]string
	resourceKeywords      map[string]struct{}
}

// NewCommandTree builds a command tree from full command paths
// (for example: "workspace list", "organization token delete").
func NewCommandTree(commandPaths []string) *CommandTree {
	t := &CommandTree{
		root:             &commandNode{children: make(map[string]*commandNode)},
		treeRoots:        make(map[string]struct{}),
		orgCollection:    make(map[string]string),
		resourceKeywords: make(map[string]struct{}),
	}

	rootVerbs := make(map[string]map[string]struct{})

	for _, path := range commandPaths {
		tokens := strings.Fields(strings.TrimSpace(path))
		if len(tokens) == 0 {
			continue
		}

		t.treeRoots[tokens[0]] = struct{}{}

		if len(tokens) == 2 {
			verbs, ok := rootVerbs[tokens[0]]
			if !ok {
				verbs = make(map[string]struct{})
				rootVerbs[tokens[0]] = verbs
			}
			verbs[tokens[1]] = struct{}{}
		}

		node := t.root
		for _, token := range tokens {
			child, ok := node.children[token]
			if !ok {
				child = &commandNode{children: make(map[string]*commandNode)}
				node.children[token] = child
			}
			node = child
		}
		node.terminal = true
	}

	for _, path := range commandPaths {
		tokens := strings.Fields(strings.TrimSpace(path))
		if len(tokens) == 0 {
			continue
		}
		if tokens[len(tokens)-1] != "list" {
			continue
		}

		// Top-level list commands map to collection resources:
		//   workspace list -> workspaces
		//   workspace resource list -> resources
		if len(tokens) == 2 {
			root := tokens[0]
			keyword := resourceKeyword(root)
			t.resourceKeywords[keyword] = struct{}{}
			if shouldBeOrgCollection(root, rootVerbs[root]) {
				t.orgCollection[keyword] = root
			}
			continue
		}

		if tokens[0] != "workspace" {
			continue
		}

		if len(tokens) >= 3 {
			keyword := resourceKeyword(tokens[len(tokens)-2])
			t.resourceKeywords[keyword] = struct{}{}
		}
	}

	for root, verbs := range rootVerbs {
		if _, ok := verbs["outputs"]; ok {
			t.resourceKeywords[root] = struct{}{}
		}
	}

	return t
}

func (t *CommandTree) HasRoot(token string) bool {
	if t == nil {
		return false
	}
	_, ok := t.treeRoots[token]
	return ok
}

func (t *CommandTree) OrgCollectionNamespace(token string) (string, bool) {
	if t == nil {
		return "", false
	}
	namespace, ok := t.orgCollection[token]
	return namespace, ok
}

func (t *CommandTree) IsResourceKeyword(token string) bool {
	if t == nil {
		return false
	}
	_, ok := t.resourceKeywords[token]
	return ok
}

func shouldBeOrgCollection(root string, verbs map[string]struct{}) bool {
	if root == "" {
		return false
	}
	if _, ok := verbs["list"]; !ok {
		return false
	}
	if _, ok := verbs["delete"]; !ok {
		return false
	}

	if _, ok := verbs["read"]; ok {
		return true
	}
	if _, ok := verbs["show"]; !ok {
		return false
	}
	if _, ok := verbs["add-member"]; ok {
		return true
	}
	_, ok := verbs["remove-member"]
	return ok
}

func resourceKeyword(token string) string {
	if token == "" {
		return token
	}

	if strings.HasPrefix(token, "workspace") && token != "workspace" {
		base := strings.TrimPrefix(token, "workspace")
		if base != "" {
			return pluralize(base)
		}
	}

	if strings.HasSuffix(token, "result") {
		return strings.TrimSuffix(token, "result") + "s"
	}

	return pluralize(token)
}

func pluralize(token string) string {
	if token == "" {
		return token
	}

	if strings.HasSuffix(token, "s") || strings.HasSuffix(token, "sh") || strings.HasSuffix(token, "ch") || strings.HasSuffix(token, "x") || strings.HasSuffix(token, "z") {
		return token + "es"
	}

	if strings.HasSuffix(token, "y") {
		return strings.TrimSuffix(token, "y") + "ies"
	}

	return token + "s"
}
