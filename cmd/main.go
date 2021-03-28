package main

import (
	"github.com/chaordic-io/gbuild/internal/config"
	"github.com/chaordic-io/gbuild/internal/execution"
)

func main() {
	targets := []config.Target{
		{"baz", nil, nil, "baz", &[]string{"foo", "bar"}, nil, nil},
		{"foo", nil, nil, "foo", nil, nil, nil},
		{"bar", nil, nil, "bar", nil, nil, nil},
	}
	execution.RunPlan(targets)
}
