package execution

import (
	"fmt"
	"testing"

	"github.com/chaordic-io/gbuild/internal/config"
)

func TestSimpleExecution(t *testing.T) {

	targets := []config.Target{
		{"foo", nil, nil, "ls -ltr", nil},
		{"bar", nil, nil, "ls", nil},
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
		{"baz", nil, nil, "pwd", &[]string{"foo", "bar"}},
		{"foo", nil, nil, "ls", nil},
		{"bar", nil, nil, "echo 'bar'", &[]string{"foo"}},
	}
	for i := 1; i < 100; i++ {

		res, err := RunPlan(targets)

		if err != nil {
			t.Fatalf("Did not expect error %v", err)
		}

		if len(res) != 3 {
			t.Fatalf("Expected 3 results, got %v", res)
		}

		fmt.Println("Order of tasks finishing: " + res[0].Target.Name + " " + res[1].Target.Name + " " + res[2].Target.Name)

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
