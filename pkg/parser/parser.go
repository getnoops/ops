package parser

import (
	"context"
	"io"

	"gopkg.in/yaml.v3"
)

type Parser interface {
	Parse(ctx context.Context, data io.Reader) (*Spec, error)
}

type parser struct {
}

func (p *parser) Parse(ctx context.Context, data io.Reader) (*Spec, error) {
	d := yaml.NewDecoder(data)

	var raw Raw
	if err := d.Decode(&raw); err != nil {
		return nil, err
	}

	switch raw.Version {
	case "v1":
		v1 := new(V1)
		if err := raw.Data.Decode(v1); err != nil {
			return nil, err
		}

		return &Spec{
			Version: raw.Version,
			Data:    v1,
		}, nil
	}

	return nil, nil
}

func NewParser() Parser {
	return &parser{}
}
