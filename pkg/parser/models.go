package parser

import (
	"gopkg.in/yaml.v3"
)

type Raw struct {
	Version string    `yaml:"version"`
	Data    yaml.Node `yaml:"data"`
}

type Spec struct {
	Version string      `yaml:"version"`
	Data    interface{} `yaml:"data"`
}

type V1 struct {
	Hello string `yaml:"hello"`
}
