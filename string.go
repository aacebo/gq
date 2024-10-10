package gq

import (
	"encoding/json"
	"reflect"

	"github.com/aacebo/gq/query"
)

type String struct{}

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
	value := reflect.Indirect(reflect.ValueOf(params.Value))

	if value.Kind() != reflect.String {
		return Result{Error: NewError("", "must be a string")}
	}

	return Result{Data: value.Interface()}
}

func (self String) MarshalJSON() ([]byte, error) {
	return json.Marshal("string")
}
