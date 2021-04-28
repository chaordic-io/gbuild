package internal

import (
	"os"
	"testing"
)

func cacheToTarget(caches *[]Cache) Target {
	return Target{"foo", nil, String("internal"), "bla", nil, caches}
}
func TestCalculateCacheStatesWithOne(t *testing.T) {
	cache := cacheToTarget(&[]Cache{
		{[]string{"config.go", "execution.go"}, []string{"config_test.go", "execution_test.go"}},
	})

	state, err := calculateCacheStates(String("../"), &[]Target{cache})
	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}
	if state == nil || len(*state) != 1 {
		t.Fatalf("Did not expect state of nil or empty, %v", state)
	}
}

func TestZipUnzip(t *testing.T) {

	os.RemoveAll("../tmp")
	res, err := CheckSumWithGitIgnoreWithRelative(String("../"), nil, []string{"."}, true)
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}

	os.MkdirAll("../tmp", os.ModePerm)
	os.MkdirAll("../tmp/unzip", os.ModePerm)
	err = zipWriter(map[string]string{"../": ""}, "../tmp/file.zip")
	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}

	unzip("../tmp/unzip", "../tmp/file.zip")

	res2, err := CheckSumWithGitIgnoreWithRelative(String("../tmp/unzip"), nil, []string{"."}, true)
	if err != nil {
		t.Fatalf("Expected no error, found %v", err)
	}

	if *res != *res2 {
		t.Fatalf("Expected checksums to match between folder zipped and folder unzipped %v, %v", *res, *res2)
	}

}
