package router

import "testing"

func TestCommandTreeHasRoot(t *testing.T) {
	tree := NewCommandTree([]string{
		"workspace list",
		"workspace read",
		"organization token delete",
		"whoami",
	})

	if !tree.HasRoot("workspace") {
		t.Fatal("expected workspace to be known root")
	}
	if !tree.HasRoot("organization") {
		t.Fatal("expected organization to be known root")
	}
	if !tree.HasRoot("whoami") {
		t.Fatal("expected whoami to be known root")
	}
	if tree.HasRoot("notacommand") {
		t.Fatal("expected notacommand to be unknown root")
	}
}
