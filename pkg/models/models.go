package models

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/a8m/envsubst/parse"
	"github.com/getnoops/ops/pkg/queries"
	"gopkg.in/yaml.v3"

	"github.com/suzuki-shunsuke/go-convmap/convmap"
)

type NoOpsCode struct {
	Code string `json:"code"`
}

type NoOpsConfig struct {
	Name      string                     `json:"name"`
	Code      string                     `json:"code"`
	Class     queries.ConfigClass        `json:"class"`
	Resources []*queries.ResourceInput   `json:"resources"`
	Access    *queries.ConfigAccessInput `json:"access"`
}

func (rev *NoOpsConfig) Validate() error {
	if len(rev.Resources) == 0 {
		return errors.New("no resources found")
	}
	if rev.Code == "" {
		return errors.New("no code found")
	}
	return nil
}

func Replace(in []byte, envs []string) ([]byte, error) {
	parser := parse.New("bytes", envs, parse.Strict)

	replaced, err := parser.Parse(string(in))
	if err != nil {
		return nil, err
	}

	return []byte(replaced), nil
}

type LoadOptions struct {
	Env         []string
	VarFiles    []string
	ReplaceEnvs bool
}

type LoadOption func(*LoadOptions)

func WithVarFiles(varFiles []string) LoadOption {
	return func(opts *LoadOptions) {
		opts.VarFiles = varFiles
		opts.ReplaceEnvs = true
	}
}

func WithOsEnv() LoadOption {
	return func(opts *LoadOptions) {
		opts.Env = os.Environ()
		opts.ReplaceEnvs = true
	}
}

func LoadFile[T any](file string, options ...LoadOption) (*T, error) {
	opts := &LoadOptions{}
	for _, opt := range options {
		opt(opts)
	}

	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if opts.ReplaceEnvs {
		envs := opts.Env

		for _, file := range opts.VarFiles {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			r := bufio.NewScanner(f)
			r.Split(bufio.ScanLines)

			lines := []string{}
			for r.Scan() {
				lines = append(lines, r.Text())
			}
			envs = append(lines, envs...)
		}

		raw, err = Replace(raw, envs)
		if err != nil {
			return nil, err
		}
	}

	var out T
	switch filepath.Ext(file) {
	case ".json":
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(raw, &out); err != nil {
			return nil, err
		}
		a, err := convmap.Convert(&out, nil)
		if err != nil {
			return nil, err
		}
		return a.(*T), nil

	default:
		return nil, errors.New("unsupported file type")
	}
	return &out, nil
}
