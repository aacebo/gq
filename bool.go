package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Bool struct{}

func (self Bool) Key() string {
	return "bool"
}

func (self Bool) Do(params *DoParams) Result {
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

func (self Bool) Resolve(params *ResolveParams) Result {
	switch value := params.Value.(type) {
	case bool:
		return Result{Data: value}
	case *bool:
		return Result{Data: value}
	}

	return Result{Error: NewError("", "must be a boolean")}
}

func (self Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal("bool")
}
