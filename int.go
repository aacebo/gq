package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Int struct{}

func (self Int) Key() string {
	return "int"
}

func (self Int) Do(params *DoParams) Result {
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

func (self Int) Resolve(params *ResolveParams) Result {
	switch value := params.Value.(type) {
	case int:
		return Result{Data: value}
	case *int:
		return Result{Data: value}
	}

	return Result{Error: NewError("", "must be an integer")}
}

func (self Int) MarshalJSON() ([]byte, error) {
	return json.Marshal("int")
}
