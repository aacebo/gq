package gq

import (
	"context"
	"errors"
	"reflect"

	"github.com/aacebo/gq/query"
)

type List struct {
	Type Schema `json:"type,omitempty"`
}

func (self List) Do(ctx context.Context, q string, value any) (any, error) {
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

func (self List) Resolve(params Params) (any, error) {
	value := reflect.Indirect(reflect.ValueOf(params.Value))

	if !value.IsValid() {
		return []any{}, nil
	}

	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return nil, errors.New("must be an array/slice")
	}

	for i := 0; i < value.Len(); i++ {
		index := value.Index(i)
		res, err := self.Type.Resolve(Params{
			Query:   params.Query,
			Parent:  params.Value,
			Value:   index.Interface(),
			Context: params.Context,
		})

		if err != nil {
			return nil, err
		}

		index.Set(reflect.ValueOf(res))
	}

	return value.Interface(), nil
}
