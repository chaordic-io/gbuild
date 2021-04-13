package internal

type CacheState struct {
	RootDir     *string
	WorkDir     *string
	Cache       Cache
	InChecksum  string
	GitRevs     []string
	OutChecksum *string
}

type CacheLocation struct {
	RootDir   *string
	WorkDir   *string
	Locations []string
}

// Probably change this interface
type CacheIndex struct {
	Hashes    map[string]string
	GitHashes map[string]string
}

// Probably change this interface
type CacheProvider interface {
	GetIndex() (CacheIndex, error)
	PutIndex(CacheIndex) error
	GetCache(string, CacheLocation) error
	PutCache(CacheLocation) error
}

func (state *CacheState) Inputs() CacheLocation {
	return CacheLocation{state.RootDir, state.WorkDir, state.Cache.Inputs}
}

func (state *CacheState) Outputs() CacheLocation {
	return CacheLocation{state.RootDir, state.WorkDir, state.Cache.Outputs}
}

// Calculate in + GitRevs - DONE
// Calculate current out - DONE
// run in against index
// retrieve first priority that matches
// on way out write to index, put cache items
// .gbuild/index/hash
// .gbuild/index/githash

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
			states = append(states, *newStates...)
		}
		return &states, nil
	}

	return nil, nil
}
