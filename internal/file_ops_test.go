package internal

import (
	"crypto/md5"
	"encoding/hex"
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
	fn, err := genShouldIgnoreFn(String("../"), nil)
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
		t.Fatal("Should be true for node_modules in sub-folder")
	}
}

func TestRelativePath(t *testing.T) {
	fn, err := genShouldIgnoreFn(String("../"), String("test"))
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}
	if !fn(".vscode/settings.json") {
		t.Fatal("Should be true for .vscode/settings.json")
	}

	if !fn("node_modules/bar") {
		t.Fatal("Should be false for node_modules in root")
	}
}

func TestSingleInput(t *testing.T) {
	check := func(str string) bool {
		return false
	}
	res2, err := MD5Dir(".", check)
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}
	res, err := CheckSumWithGitIgnoreWithRelative(String("../"), nil, []string{"internal"})
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}

	if *res != *res2 {
		t.Fatalf("Expected checksums to match, found %v and %v", *res, *res2)
	}
}

func TestDoubleInput(t *testing.T) {
	check := func(str string) bool {
		return false
	}
	res2, err := MD5Dir(".", check)
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}
	res3, err := MD5Dir("../cmd", check)
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}
	hash := md5.Sum([]byte(*res2 + *res3))

	md5Str := hex.EncodeToString(hash[:])

	res, err := CheckSumWithGitIgnoreWithRelative(String("../"), nil, []string{"internal", "cmd"})
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}

	if *res != md5Str {
		t.Fatalf("Expected checksums to match, found %v and %v", *res, *res2)
	}
}

func TestAccountingForGitIgnore(t *testing.T) {
	check := func(str string) bool {
		return false
	}
	res2, err := MD5Dir("../", check)
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}

	res, err := CheckSumWithGitIgnoreWithRelative(String("../"), nil, []string{"."})
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}

	if *res == *res2 {
		t.Fatalf("Expected checksums to be different, found %v and %v", *res, *res2)
	}
}
