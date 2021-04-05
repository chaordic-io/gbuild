package common

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

type gitignoreFile struct {
	path    string
	matcher gitignore.Matcher
}

func sumFiles(done <-chan struct{}, root string, shouldIgnoreFn func(string) bool) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)
	go func() {
		var wg sync.WaitGroup
		err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if shouldIgnoreFn(path) {
				return nil
			}
			if !info.Type().IsRegular() {
				return nil
			}
			wg.Add(1)
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case c <- result{path, md5.Sum(data), err}:
				case <-done:
				}
				wg.Done()
			}()
			select {
			case <-done:
				return errors.New("walk canceled")
			default:
				return nil
			}
		})
		go func() {
			wg.Wait()
			close(c)
		}()
		errc <- err
	}()
	return c, errc
}

func MD5All(root string, shouldIgnoreFn func(string) bool) (map[string][md5.Size]byte, error) {

	done := make(chan struct{})
	defer close(done)

	c, errc := sumFiles(done, root, shouldIgnoreFn)

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func MD5Dir(root string, shouldIgnoreFn func(string) bool) (*string, error) {
	m, err := MD5All(root, shouldIgnoreFn)
	if err != nil {
		return nil, err
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	toChecksum := ""
	for _, path := range paths {
		toChecksum = toChecksum + fmt.Sprintf("%x", m[path])
	}
	hash := md5.Sum([]byte(toChecksum))
	dirHash := hex.EncodeToString(hash[:])

	return &dirHash, nil
}
func Gitignores(root string) ([]gitignoreFile, error) {

	var ignore []gitignoreFile

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.Type().IsRegular() {
			return nil
		}

		// could do this with goroutines and channels - make ignore into a channel
		if info.Name() == ".gitignore" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			var patterns []gitignore.Pattern

			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if len(line) > 0 && !strings.HasPrefix(line, "#") {
					line = strings.TrimSpace(strings.Split(line, " #")[0])
					patterns = append(patterns, gitignore.ParsePattern(line, nil))
				}
			}
			ignore = append(ignore, gitignoreFile{path, gitignore.NewMatcher(patterns)})
		}

		return nil
	})

	return ignore, err
}
