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
