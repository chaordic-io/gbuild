package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chaordic-io/gbuild/internal/common"
)

func main() {
	dotGit := ".git"
	target := "target"
	node := "node_modules"
	check := func(str string) bool {
		return strings.Contains(str, dotGit) || strings.Contains(str, target) || strings.Contains(str, node)
	}

	start := time.Now()
	m, err := common.MD5Dir(os.Args[1], check)
	if err != nil {
		fmt.Println(err)
		return
	}
	since := time.Since(start)
	fmt.Println(*m)
	fmt.Println("")
	fmt.Printf("Calculated in %v", since)
}
