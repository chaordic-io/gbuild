package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Task struct {
	Name       string    `yaml:"name"`
	MaxRetries *int      `yaml:"max_retries"`
	Path       *string   `yaml:"path"`
	Run        string    `yaml:"run"`
	DependsOn  *[]string `yaml:"depends_on"`
}

type ExecutionPlan struct {
	Name  string   `yaml:"name"`
	Tasks []string `yaml:"tasks"`
}

type Config struct {
	Tasks          []Task          `yaml:"tasks"`
	ExecutionPlans []ExecutionPlan `yaml:"execution_plans"`
}

func LoadConfig(filename string) (*Config, error) {
	c := &Config{}
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}
	err = validate(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func GetTasksForPlan(config *Config, planName string) ([]Task, error) {
	var tasks []Task
	cfg := *config

	for _, pl := range cfg.ExecutionPlans {
		if pl.Name == planName {
			for _, taskName := range pl.Tasks {
				for _, task := range cfg.Tasks {
					if task.Name == taskName {
						tasks = append(tasks, task)
					}
				}
			}
		}
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks found for plan %v, does the plan exist?", planName)
	}

	return tasks, nil
}

func validate(c *Config) error {
	conf := *c
	var taskNames []string
	for _, task := range conf.Tasks {
		if containsString(task.Name, taskNames) {
			return fmt.Errorf("a task with name %v is defined twice, names must be unique", task.Name)
		}
		taskNames = append(taskNames, task.Name)
		if task.DependsOn != nil && containsString(task.Name, *task.DependsOn) {
			return fmt.Errorf("the %v depends on itself, this is not permissible", task.Name)
		}
	}

	var planNames []string
	for _, plan := range conf.ExecutionPlans {
		if containsString(plan.Name, planNames) {
			return fmt.Errorf("an execution plan with name %v is defined twice, names must be unique", plan.Name)
		}
		var planTasks []string
		planNames = append(planNames, plan.Name)
		for _, task := range plan.Tasks {
			if !containsString(task, taskNames) {
				return fmt.Errorf("the task %v in the execution plan %v is not defined among the tasks", task, plan.Name)
			}
			if containsString(task, planTasks) {
				return fmt.Errorf("the task %v in the execution plan %v is defined twice, a task may only be run once per plan", task, plan.Name)
			}
			planNames = append(planTasks, task)
		}
	}

	return nil

}

func containsString(s string, values []string) bool {
	for _, value := range values {
		if value == s {
			return true
		}
	}
	return false
}
