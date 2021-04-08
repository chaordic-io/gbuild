package internal

import (
	"testing"
)

var log = NoLog{}

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig("../test/test.yml", log)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(c.Targets) != 3 {
		t.Fatalf("Expected 3 targets, got %v", c.Targets)
	}

	if len(*c.Targets[2].DependsOn) != 2 {
		t.Fatalf("Expected 2 dependencies, got %v", c.Targets)
	}

	if c.Targets[0].MaxRetries == nil {
		t.Fatalf("Expected Max Retries to be set, got %v", c.Targets[0].MaxRetries)
	}
	if c.Targets[0].WorkDir == nil {
		t.Fatalf("Expected Max Retries to be set, got %v", c.Targets[0].MaxRetries)
	}

	if len(c.ExecutionPlans) != 2 {
		t.Fatalf("Expected 2 targets, got %v", c.ExecutionPlans)
	}

	if len(c.ExecutionPlans[0].Targets) != 3 {
		t.Fatalf("Expected 3 targets, got %v", c.ExecutionPlans)
	}

	if len(c.ExecutionPlans[1].Targets) != 2 {
		t.Fatalf("Expected 3 targets, got %v", c.ExecutionPlans)
	}

	if c.Targets[2].Caches == nil && len(*c.Targets[2].Caches) == 2 {
		t.Fatalf("Expected Caches to be set on 3rd target, got %v", c.Targets[0].MaxRetries)
	}
}

func TestLoadConfigNoFile(t *testing.T) {
	_, err := LoadConfig("Not-there.yml", log)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}

func TestLoadInvalidFile(t *testing.T) {
	_, err := LoadConfig("../test/invalid.yaml", log)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}

func TestLoadFileThatFailsValidation(t *testing.T) {
	_, err := LoadConfig("../test/invalid2.yaml", log)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}

func TestTargetDefinedTwiceValidation(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", nil, nil},
			{"foo", nil, nil, "bar", nil, nil},
		},
		[]ExecutionPlan{},
	}

	err := validate(c, log)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestTargetSelfDependentValidations(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", &[]string{"foo"}, nil},
		},
		[]ExecutionPlan{},
	}

	err := validate(c, log)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestTargetNotDefined(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", nil, nil},
		},
		[]ExecutionPlan{{"foo", []string{"bar"}}, {"bar", []string{}}},
	}

	err := validate(c, log)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestDuplicatePlanName(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", nil, nil},
		},
		[]ExecutionPlan{
			{"foo", []string{"foo"}},
			{"foo", []string{"foo"}},
		},
	}

	err := validate(c, log)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestDuplicateTargetInPlan(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", nil, nil},
		},
		[]ExecutionPlan{{"bar", []string{"foo", "foo"}}},
	}

	err := validate(c, log)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}

}

func TestGetTargetsForPlan(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", nil, nil},
		},
		[]ExecutionPlan{{"foo", []string{"foo"}}},
	}

	targets, err := GetTargetsForPlan(c, "foo", log)
	if err != nil || len(targets) != 1 {
		t.Fatal("Did not expect an error here, expected 1 target")
	}
}

func TestGetTargetsForPlanFailure(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", &[]string{"foo"}, nil},
		},
		[]ExecutionPlan{{"foo", []string{"foo"}}},
	}

	_, err := GetTargetsForPlan(c, "bar", log)
	if err == nil {
		t.Fatal("Did not expect an error here, expected 1 target")
	}
}

func TestGetTargetsForPlanFailure2(t *testing.T) {
	c := &Config{
		[]Target{
			{"foo", nil, nil, "bar", nil, nil},
		},
		[]ExecutionPlan{{"foo", []string{}}},
	}

	_, err := GetTargetsForPlan(c, "foo", log)
	if err == nil {
		t.Fatal("Did not expect an error here, expected 1 target")
	}
}
