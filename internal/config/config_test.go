package config

import "testing"

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig("test.yml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(c.Tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %v", c.Tasks)
	}

	if len(*c.Tasks[2].DependsOn) != 2 {
		t.Fatalf("Expected 2 dependencies, got %v", c.Tasks)
	}

	if c.Tasks[0].MaxRetries == nil {
		t.Fatalf("Expected Max Retries to be set, got %v", c.Tasks[0].MaxRetries)
	}
	if c.Tasks[0].Path == nil {
		t.Fatalf("Expected Max Retries to be set, got %v", c.Tasks[0].MaxRetries)
	}

	if len(c.ExecutionPlans) != 2 {
		t.Fatalf("Expected 2 tasks, got %v", c.ExecutionPlans)
	}

	if len(c.ExecutionPlans[0].Tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %v", c.ExecutionPlans)
	}

	if len(c.ExecutionPlans[1].Tasks) != 2 {
		t.Fatalf("Expected 3 tasks, got %v", c.ExecutionPlans)
	}
}

func TestTaskDefinedTwiceValidation(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", nil},
			{"foo", nil, nil, "bar", nil},
		},
		[]ExecutionPlan{},
	}

	err := validate(c)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestTaskSelfDependentValidations(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", &[]string{"foo"}},
		},
		[]ExecutionPlan{},
	}

	err := validate(c)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestTaskNotDefined(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", &[]string{"foo"}},
		},
		[]ExecutionPlan{{"foo", []string{"bar"}}, {"bar", []string{}}},
	}

	err := validate(c)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestDuplicatePlanName(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", &[]string{"foo"}},
		},
		[]ExecutionPlan{{"foo", []string{}}, {"foo", []string{}}},
	}

	err := validate(c)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestDuplicateTaskInPlan(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", &[]string{"foo"}},
		},
		[]ExecutionPlan{{"foo", []string{"foo", "foo"}}},
	}

	err := validate(c)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestGetTasksForPlan(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", &[]string{"foo"}},
		},
		[]ExecutionPlan{{"foo", []string{"foo"}}},
	}

	tasks, err := GetTasksForPlan(c, "foo")
	if err != nil || len(tasks) != 1 {
		t.Fatal("Did not expect an error here, expected 1 task")
	}
}

func TestGetTasksForPlanFailure(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", &[]string{"foo"}},
		},
		[]ExecutionPlan{{"foo", []string{"foo"}}},
	}

	_, err := GetTasksForPlan(c, "bar")
	if err == nil {
		t.Fatal("Did not expect an error here, expected 1 task")
	}
}

func TestGetTasksForPlanFailure2(t *testing.T) {
	c := &Config{
		[]Task{
			{"foo", nil, nil, "bar", &[]string{"foo"}},
		},
		[]ExecutionPlan{{"foo", []string{}}},
	}

	_, err := GetTasksForPlan(c, "foo")
	if err == nil {
		t.Fatal("Did not expect an error here, expected 1 task")
	}
}
