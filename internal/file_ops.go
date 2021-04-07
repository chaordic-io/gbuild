package internal

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
	if len(paths) == 1 {
		return String(fmt.Sprintf("%x", m[paths[0]])), nil
	}
	sort.Strings(paths)
	toChecksum := ""
	for _, path := range paths {
		toChecksum = toChecksum + fmt.Sprintf("%x", m[path])
	}
	hash := md5.Sum([]byte(toChecksum))

	return String(hex.EncodeToString(hash[:])), nil
}

func Gitignores(root *string) ([]gitignoreFile, error) {

	var ignore []gitignoreFile
	if root == nil {
		root = String(".")
	}

	err := filepath.WalkDir(*root, func(path string, info fs.DirEntry, err error) error {
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
			path = strings.Replace(strings.ReplaceAll(path, "../", ""), ".gitignore", "", 1)
			ignore = append(ignore, gitignoreFile{path, gitignore.NewMatcher(patterns)})
		}

		return nil
	})

	return ignore, err
}

func GenShouldIgnoreFn(projectRoot *string, relativePath *string) (func(string) bool, error) {

	files, err := Gitignores(projectRoot)

	if err != nil {
		return nil, err
	}
	ignoreFn := func(file string) bool {
		file = prependPath(relativePath, file)
		for _, gitignore := range files {
			if strings.HasPrefix(file, gitignore.path) && gitignore.matcher.Match(strings.Split(file, "/"), false) {
				return true
			}
		}
		return false
	}

	return ignoreFn, nil
}

func prependPath(relativePath *string, file string) string {
	if relativePath != nil {
		if strings.HasSuffix(*relativePath, "/") {
			return *relativePath + file
		} else {
			return *relativePath + "/" + file
		}
	} else {
		return file
	}
}

func CheckSumWithGitIgnoreWithRelative(projectRoot *string, relativePath *string, inputs []string) (*string, error) {
	var calculatedInputs []string
	for _, file := range inputs {
		calculatedInputs = append(calculatedInputs, prependPath(projectRoot, prependPath(relativePath, file)))
	}

	fn, err := GenShouldIgnoreFn(projectRoot, relativePath)
	if err != nil {
		return nil, err
	}

	toChecksum := ""
	if len(calculatedInputs) == 1 {
		return MD5Dir(calculatedInputs[0], fn)
	}
	for _, input := range calculatedInputs {
		sum, err := MD5Dir(input, fn)
		if err != nil {
			return nil, err
		}
		toChecksum = toChecksum + *sum
	}
	hash := md5.Sum([]byte(toChecksum))

	return String(hex.EncodeToString(hash[:])), nil
}
