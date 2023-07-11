package util

import (
	"bytes"
	"fmt"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	// Map from service name to map from process name to URL.
	Targets map[string]map[string]string `yaml:"targets"`
}

func ReadConfig(path string) (Config, error) {
	r, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	dec := yaml.NewDecoder(bytes.NewReader(r))
	if err := dec.Decode(&c); err != nil {
		return Config{}, err
	}
	pretty.Print(c)

	if len(c.Targets) > 1 {
		return Config{}, fmt.Errorf("expected exactly one service in config")
	}

	return c, nil
}
