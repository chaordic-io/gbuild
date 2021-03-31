package execution

import (
	"testing"

	"github.com/chaordic-io/gbuild/internal/common"
	"github.com/chaordic-io/gbuild/internal/config"
)

var log = common.NoLog{}

func TestSimpleExecution(t *testing.T) {

	targets := []config.Target{
		{"foo", nil, nil, "cd .", nil},
		{"bar", nil, nil, "cd .", nil},
	}

	res, err := RunPlan(targets, log)

	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected 2 results, got %v", res)
	}
}

// This test should pass in only a few seconds if it works, because a failed other process should cancell all others
func TestFailedExecutionWithCancelOfOthers(t *testing.T) {

	targets := []config.Target{
		{"fast", nil, nil, "asdfasdf", nil},
		{"slow", nil, nil, "sleep 30", nil},
	}

	res, err := RunPlan(targets, log)

	if err == nil {
		t.Fatalf("Did not expect error %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected 2 results, got %v", res)
	}
}

func TestDependentExecution(t *testing.T) {
	targets := []config.Target{
		{"baz", nil, nil, "cd .", &[]string{"foo", "bar"}},
		{"foo", nil, nil, "cd .", nil},
		{"bar", nil, nil, "cd .", &[]string{"foo"}},
	}
	for i := 1; i < 100; i++ {

		res, err := RunPlan(targets, log)

		if err != nil {
			t.Fatalf("Did not expect error %v", err)
		}

		if len(res) != 3 {
			t.Fatalf("Expected 3 results, got %v", res)
		}

		if res[2].Target.Name != "baz" {
			t.Fatalf("baz should always be last, but got %v", res[2].Target.Name)
		}
		if res[1].Target.Name != "bar" {
			t.Fatalf("bar should always be second, but got %v", res[2].Target.Name)
		}
		if res[0].Target.Name != "foo" {
			t.Fatalf("foo should always be first, but got %v", res[2].Target.Name)
		}

	}
}
