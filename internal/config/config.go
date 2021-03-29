package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Target struct {
	Name       string    `yaml:"name"`
	MaxRetries *int      `yaml:"max_retries"`
	WorkDir    *string   `yaml:"work_dir"`
	Run        string    `yaml:"run"`
	DependsOn  *[]string `yaml:"depends_on"`
}

type ExecutionPlan struct {
	Name    string   `yaml:"name"`
	Targets []string `yaml:"targets"`
}

type Config struct {
	Targets        []Target        `yaml:"targets"`
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

func GetTargetsForPlan(config *Config, planName string) ([]Target, error) {
	var targets []Target
	cfg := *config

	for _, pl := range cfg.ExecutionPlans {
		if pl.Name == planName {
			for _, targetName := range pl.Targets {
				for _, target := range cfg.Targets {
					if target.Name == targetName {
						targets = append(targets, target)
					}
				}
			}
		}
	}
	if len(targets) == 0 {
		return nil, fmt.Errorf("no targets found for plan %v, does the plan exist?", planName)
	}

	return targets, nil
}

func validate(c *Config) error {
	conf := *c
	var targetNames []string
	for _, target := range conf.Targets {
		if containsString(target.Name, targetNames) {
			return fmt.Errorf("a target with name %v is defined twice, names must be unique", target.Name)
		}
		targetNames = append(targetNames, target.Name)
		if target.DependsOn != nil && containsString(target.Name, *target.DependsOn) {
			return fmt.Errorf("the %v depends on itself, this is not permissible", target.Name)
		}
	}

	var planNames []string
	for _, plan := range conf.ExecutionPlans {
		if containsString(plan.Name, planNames) {
			return fmt.Errorf("an execution plan with name %v is defined twice, names must be unique", plan.Name)
		}
		var planTargets []string
		planNames = append(planNames, plan.Name)
		for _, target := range plan.Targets {
			if !containsString(target, targetNames) {
				return fmt.Errorf("the target %v in the execution plan %v is not defined among the targets", target, plan.Name)
			}
			if containsString(target, planTargets) {
				return fmt.Errorf("the target %v in the execution plan %v is defined twice, a target may only be run once per plan", target, plan.Name)
			}
			planNames = append(planTargets, target)
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
