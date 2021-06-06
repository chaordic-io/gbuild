package api

import (
	"io"
)

type CacheIndex struct {
	Hashes    map[string]string
	GitHashes map[string]string
}

// Probably change this interface
type CacheProvider interface {
	GetIndex() (*CacheIndex, error)
	PutIndex(CacheIndex) error
	GetCache(string) (*io.Reader, error)
	PutCache(string, io.Reader) error
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

func (cache *LocalFileCacheProvider) GetCache(hash string) (*io.Reader, error) {
	return nil, nil
}

func (cache *LocalFileCacheProvider) PutCache(hash string, reader io.Reader) error {
	return nil
}
