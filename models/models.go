package models

import (
	"fmt"
	. "io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name   string  `yaml:"name"`   // Supporting both JSON and YAML.
	Envs   []Env   `yaml:"envs"`   // Supporting both JSON and YAML.
	Groups []Group `yaml:"groups"` // Supporting both JSON and YAML.
}

type Env struct {
	Name      string                      `yaml:"name"`       // Supporting both JSON and YAML.
	DependsOn []string                    `yaml:"depends_on"` // Supporting both JSON and YAML.
	Features  map[interface{}]interface{} `yaml:"features"`
}

type Group struct {
	Name string `yaml:"name"` // Supporting both JSON and YAML.
}

func (f *Env) GetDependsOn() string {
	return fmt.Sprintf("[%s]", strings.Join(f.DependsOn, ","))
}

func (f *Env) GetDependsOnArray() []string {
	return f.DependsOn
}

func (e *Env) HasDependencies() bool {
	if e.DependsOn == nil || len(e.DependsOn) == 0 {
		return false
	} else {
		return true
	}
}

func (e *Env) HasFeature(feature string) bool {
	_, ok := e.Features[feature]
	return ok
}

func (e *Env) Feature(feature string) string {
	if v, ok := e.Features[feature]; ok {
		if v == nil {
			return ""
		}
		if _, ok := v.(map[interface{}]interface{}); ok {
			panic(fmt.Sprintf("Using a map as the value of a feature is currently unsupported (on feature '%s')\n", feature))
		}
		if _, ok := v.([]interface{}); ok {
			panic(fmt.Sprintf("Using an array as the value of a feature is currently unsupported (on feature '%s')\n", feature))
		}
		return fmt.Sprintf("%v", v)
	} else {
		return ""
	}
}

func Load(y []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(y, &config)

	for index, env := range config.Envs {
		if len(env.DependsOn) == 0 {
			config.Envs[index].DependsOn = nil
		}
		if env.Features == nil {
			config.Envs[index].Features = map[interface{}]interface{}{}
		}
	}

	return config, err
}

func LoadFromFile(path string) (Config, error) {
	y, err := ReadFile(path)
	if err != nil {
		panic(err)
	}

	return Load(y)
}
