package internal

import (
	"io"
	"os"
	"path/filepath"
)

type CacheState struct {
	RootDir     *string
	WorkDir     *string
	Cache       Cache
	InChecksum  string
	GitRevs     []string
	OutChecksum *string
}

// type CacheLocation struct {
// 	RootDir   *string
// 	WorkDir   *string
// 	Locations []string
// }

// Probably change this interface
type CacheIndex struct {
	Hashes    map[string]string
	GitHashes map[string]string
}

// Probably change this interface
type CacheProvider interface {
	GetIndex() (CacheIndex, error)
	PutIndex(CacheIndex) error
	GetCache(string) (io.ReadCloser, error)
	PutCache(string, io.ReadCloser) error
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
	return nil
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
	for _, state := range *states {
		cache := index.getCacheFile(&state)
		if cache != nil && *cache != *state.OutChecksum {
			// check if we already downloaded the cache here?
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
