package gq

import (
	"encoding/json"
	"time"

	"github.com/aacebo/gq/query"
)

type Date struct{}

func (self Date) Do(params *DoParams) Result {
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

func (self Date) Resolve(params *ResolveParams) Result {
	switch value := params.Value.(type) {
	case time.Time:
		return Result{Data: value}
	case *time.Time:
		return Result{Data: value}
	}

	return Result{Error: NewError("", "must be a Date")}
}

func (self Date) MarshalJSON() ([]byte, error) {
	return json.Marshal("Date")
}
