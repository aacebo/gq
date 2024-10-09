package gq

import (
	"context"
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Args interface {
	Validate(value any) error
}

type Field struct {
	Type        Schema                           `json:"type,omitempty"`
	Description string                           `json:"description,omitempty"`
	Args        Args                             `json:"args,omitempty"`
	DependsOn   []string                         `json:"depends_on,omitempty"`
	Use         []Middleware                     `json:"-"`
	Resolver    func(params Params) (any, error) `json:"-"`
}

func (self Field) Do(ctx context.Context, q string, value any) (any, error) {
	parser := query.Parser([]byte(q))
	query, err := parser.Parse()

	if err != nil {
		return nil, err
	}

	return self.Resolve(Params{
		Query:   query,
		Value:   value,
		Context: ctx,
	})
}

func (self Field) Resolve(params Params) (any, error) {
	if self.Use != nil {
		for _, use := range self.Use {
			if err := use(params); err != nil {
				return params.Value, err
			}
		}
	}

	if self.Args != nil {
		if err := self.Args.Validate(params.Value); err != nil {
			return params.Value, NewError(params.Key, err.Error())
		}
	}

	if self.Resolver != nil {
		value, err := self.Resolver(params)

		if err != nil {
			return value, NewError(params.Key, err.Error())
		}

		params.Value = value
	}

	if self.Type != nil {
		value, err := self.Type.Resolve(params)

		if err != nil {
			return value, NewEmptyError(params.Key).Add(err)
		}

		params.Value = value
	}

	return params.Value, nil
}

func (self Field) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
