package internal

import (
	"testing"
)

var l = NoLog{}

func TestSimpleExecution(t *testing.T) {

	targets := []Target{
		{"foo", nil, nil, "cd .", nil, nil},
		{"bar", nil, nil, "cd .", nil, nil},
	}

	res, err := RunPlan(targets, l)

	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected 2 results, got %v", res)
	}
}

// This test should pass in only a few seconds if it works, because a failed other process should cancell all others
func TestFailedExecutionWithCancelOfOthers(t *testing.T) {

	targets := []Target{
		{"fast", nil, nil, "asdfasdf", nil, nil},
		{"slow", nil, nil, "sleep 5", nil, nil},
	}

	res, err := RunPlan(targets, l)

	if err == nil {
		t.Fatalf("Did not expect error %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected 2 results, got %v", res)
	}
}

func TestFailedExecutionWithCancelOfOthersRetries(t *testing.T) {

	targets := []Target{
		{"fast", Int(2), nil, "asdfasdf", nil, nil},
		{"slow", nil, nil, "sleep 5", nil, nil},
	}

	res, err := RunPlan(targets, l)

	if err == nil {
		t.Fatalf("Did not expect error %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected 2 results, got %v", res)
	}
}

func TestDependentExecution(t *testing.T) {
	targets := []Target{
		{"baz", nil, nil, "cd .", &[]string{"foo", "bar"}, nil},
		{"foo", nil, nil, "cd .", nil, nil},
		{"bar", nil, nil, "cd .", &[]string{"foo"}, nil},
	}
	for i := 1; i < 100; i++ {

		res, err := RunPlan(targets, l)

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

func TestNonExistentExecutionDir(t *testing.T) {
	targets := []Target{
		{"foo", nil, String("foobar"), "cd .", nil, nil},
	}

	_, err := RunPlan(targets, l)

	if err == nil {
		t.Fatalf("Did not expect lack of error")
	}
}

func TestExistingDir(t *testing.T) {
	targets := []Target{
		{"foo", nil, String("../internal"), "cd .", nil, nil},
	}

	_, err := RunPlan(targets, l)

	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}
}
