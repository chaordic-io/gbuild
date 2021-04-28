package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CacheState struct {
	RootDir     *string
	WorkDir     *string
	Cache       Cache
	InChecksum  string
	GitRevs     []string
	OutChecksum *string
}

type CacheIndex struct {
	Hashes    map[string]string
	GitHashes map[string]string
}

// Probably change this interface
type CacheProvider interface {
	GetIndex() (*CacheIndex, error)
	PutIndex(CacheIndex) error
	GetCache(string) (*io.ReadCloser, error)
	PutCache(string, io.ReadCloser) error
}

type LocalFileCacheProvider struct {
	Directory string
}

func (cache *LocalFileCacheProvider) GetIndex() (*CacheIndex, error) {
	return nil, nil
}

func (cache *LocalFileCacheProvider) PutIndex(index CacheIndex) error {
	return nil
}

func (cache *LocalFileCacheProvider) GetCache(hash string) (*io.ReadCloser, error) {
	return nil, nil
}

func (cache *LocalFileCacheProvider) PutCache(hash string, reader io.ReadCloser) error {
	return nil
}

// func (state *CacheState) inputs() CacheLocation {
// 	return CacheLocation{state.RootDir, state.WorkDir, state.Cache.Inputs}
// }

// func (state *CacheState) outputs() CacheLocation {
// 	return CacheLocation{state.RootDir, state.WorkDir, state.Cache.Outputs}
// }

// Calculate in + GitRevs - DONE
// Calculate current out - DONE
// run in against index
// retrieve first priority that matches
// on way out write to index, put cache items
// .gbuild_cache/index/hash
// .gbuild_cache/index/githash

func calculateCacheState(rootDir *string, target *Target) (*[]CacheState, error) {
	if target.Caches != nil && len(*target.Caches) > 0 {
		var caches []CacheState
		for _, cache := range *target.Caches {
			gitRevs, err := GetGitHashes(rootDir, target.WorkDir, cache.Inputs)
			if err != nil {
				return nil, err
			}
			checksum, err := CheckSumWithGitIgnoreWithRelative(rootDir, target.WorkDir, cache.Inputs, true)
			if err != nil {
				return nil, err
			}
			outSum, err := CheckSumWithGitIgnoreWithRelative(rootDir, target.WorkDir, cache.Outputs, false)
			if err != nil {
				return nil, err
			}
			state := CacheState{rootDir, target.WorkDir, cache, *checksum, *gitRevs, outSum}
			caches = append(caches, state)
		}
		return &caches, nil
	}

	return nil, nil
}

func calculateCacheStates(rootDir *string, targets *[]Target) (*[]CacheState, error) {
	if targets != nil {
		var states []CacheState
		for _, target := range *targets {
			newStates, err := calculateCacheState(rootDir, &target)
			if err != nil {
				return nil, err
			}
			if newStates != nil {
				states = append(states, *newStates...)
			}
		}
		return &states, nil
	}

	return nil, nil
}

func (index *CacheIndex) getCacheFile(state *CacheState) *string {
	result := index.Hashes[state.InChecksum]
	if len(result) == 0 {
		for _, hash := range state.GitRevs {
			result = index.GitHashes[hash]
			if len(result) > 0 {
				return &result
			}
		}
		return nil
	} else {
		return &result
	}
}

func LoadCache(rootDir *string, targets *[]Target, provider CacheProvider) error {
	if provider == nil || targets == nil {
		return nil
	}
	cacheDir := prependPath(rootDir, filepath.Join(".gbuild_cache", "cache"))
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, os.ModePerm)
	}
	states, err := calculateCacheStates(nil, targets)
	if err != nil || states == nil {
		return err
	}
	index, err := provider.GetIndex()
	if err != nil {
		return err
	}
	for _, state := range *states {
		cache := index.getCacheFile(&state)
		if cache != nil && *cache != *state.OutChecksum {
			// check if we already downloaded the cache here? -
			// "has built locally with list" to avoid unpacking same cache multiple times
			hitDir := filepath.Join(cacheDir, *cache)
			if _, err := os.Stat(hitDir); os.IsNotExist(err) {
				os.MkdirAll(hitDir, os.ModePerm)
				_, err := provider.GetCache(*cache)
				if err != nil {
					return err
				}
				// unzip reader in cache-hit-dir
				// move files into target locations of state
			}
		}
	}

	return nil
}

func PutCache(rootDir *string, targets *[]Target) error {

	return nil
}

func zipTarget(toLocation string, target *Target) error {

	return nil
}

func unzip(dest string, src string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func zipWriter(src string, target string) error {

	// Get a Buffer to Write To
	outFile, err := os.Create(target)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	err = addFiles(w, src, "")

	if err != nil {
		fmt.Println(err)
		return err
	}
	// Make sure to check the error on Close.
	err = w.Close()
	return err
}

func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range files {
		// fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(filepath.Join(basePath, file.Name()))
			if err != nil {
				fmt.Println("ff " + err.Error())
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			// fmt.Println("Recursing and Adding SubDir: " + file.Name())
			// fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
	return nil
}
