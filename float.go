package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Float struct{}

func (self Float) Key() string {
	return "float"
}

func (self Float) Do(params *DoParams) Result {
	parser := query.Parser([]byte(params.Query))
	query, err := parser.Parse()

	if err != nil {
		return Result{Error: err}
	}

	return self.Resolve(&ResolveParams{
		Query:   query,
		Value:   params.Value,
		Context: params.Context,
	})
}

func (self Float) Resolve(params *ResolveParams) Result {
	switch value := params.Value.(type) {
	case float32:
		return Result{Data: value}
	case *float32:
		return Result{Data: value}
	case float64:
		return Result{Data: value}
	case *float64:
		return Result{Data: value}
	}

	return Result{Error: NewError("", "must be a float")}
}

func (self Float) MarshalJSON() ([]byte, error) {
	return json.Marshal("float")
}
