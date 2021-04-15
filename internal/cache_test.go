package internal

import (
	"testing"
)

func cacheToTarget(caches *[]Cache) Target {
	return Target{"foo", nil, String("internal"), "bla", nil, caches}
}
func TestCalculateCacheStatesWithOne(t *testing.T) {
	cache := cacheToTarget(&[]Cache{
		Cache{[]string{"config.go", "execution.go"}, []string{"config_test.go", "execution_test.go"}},
	})

	state, err := calculateCacheStates(String("../"), &[]Target{cache})
	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}
	if state == nil || len(*state) != 1 {
		t.Fatalf("Did not expect state of nil or empty, %v", state)
	}
}
