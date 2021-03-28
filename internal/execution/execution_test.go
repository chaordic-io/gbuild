package execution

import (
	"testing"

	"github.com/chaordic-io/gbuild/internal/config"
)

func TestSimpleExecution(t *testing.T) {

	targets := []config.Target{
		{"foo", nil, nil, "foo", nil, nil, nil},
		{"bar", nil, nil, "bar", nil, nil, nil},
	}

	res, err := RunPlan(targets)

	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected 2 results, got %v", res)
	}

}

func TestDependentExecution(t *testing.T) {
	targets := []config.Target{
		{"baz", nil, nil, "baz", &[]string{"foo", "bar"}, nil, nil},
		{"foo", nil, nil, "foo", nil, nil, nil},
		{"bar", nil, nil, "bar", nil, nil, nil},
	}
	for i := 1; i < 100; i++ {

		res, err := RunPlan(targets)

		if err != nil {
			t.Fatalf("Did not expect error %v", err)
		}

		if len(res) != 3 {
			t.Fatalf("Expected 3 results, got %v", res)
		}

		if res[2].Target.Name != "baz" {
			t.Fatalf("baz should always be last, but got %v", res)
		}

	}
}
