package internal

import (
	"strings"
	"testing"
)

func TestCheckSumDir(t *testing.T) {
	dotGit := ".git"
	target := "target"
	node := "node_modules"
	check := func(str string) bool {
		return strings.Contains(str, dotGit) || strings.Contains(str, target) || strings.Contains(str, node)
	}

	sum, err := MD5Dir("../", check)
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}
	if len(*sum) != 32 {
		t.Fatalf("Expected length to be %d, but was ", len(*sum))
	}
}

func TestCheckSumDirFail(t *testing.T) {
	dotGit := ".git"
	target := "target"
	node := "node_modules"
	check := func(str string) bool {
		return strings.Contains(str, dotGit) || strings.Contains(str, target) || strings.Contains(str, node)
	}

	_, err := MD5Dir("asdfasdfs", check)
	if err == nil {
		t.Fatalf("Expected no error, found %v", err)
	}
}

func TestIgnoreGeneration(t *testing.T) {
	fn, err := GenShouldIgnoreFn("../")
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}
	if !fn(".vscode/settings.json") {
		t.Fatal("Should be true for .vscode/settings.json")
	}

	if fn("node_modules/foo") {
		t.Fatal("Should be false for node_modules in root")
	}

	if !fn("test/node_modules/foo") {
		t.Fatal("Should be true for node_modules in root")
	}

}
