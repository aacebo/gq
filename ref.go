package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Ref struct {
	Type Schema `json:"type"`
}

func (self Ref) Key() string {
	return self.Type.Key()
}

func (self Ref) Do(params *DoParams) Result {
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

func (self Ref) Resolve(params *ResolveParams) Result {
	return self.Type.Resolve(params)
}

func (self Ref) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Key())
}
