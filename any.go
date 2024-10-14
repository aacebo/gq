package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Any struct{}

func (self Any) Key() string {
	return "any"
}

func (self Any) Do(params *DoParams) Result {
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

func (self Any) Resolve(params *ResolveParams) Result {
	return Result{Data: params.Value}
}

func (self Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Key())
}
