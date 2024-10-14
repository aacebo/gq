package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type String struct{}

func (self String) Key() string {
	return "string"
}

func (self String) Do(params *DoParams) Result {
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

func (self String) Resolve(params *ResolveParams) Result {
	switch value := params.Value.(type) {
	case string:
		return Result{Data: value}
	case *string:
		return Result{Data: value}
	}

	return Result{Error: NewError("", "must be a string")}
}

func (self String) MarshalJSON() ([]byte, error) {
	return json.Marshal("string")
}
