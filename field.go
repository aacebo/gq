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
	if self.Resolver != nil {
		return self.Resolver(params)
	}

	return params.Value, nil
}

func (self Field) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
