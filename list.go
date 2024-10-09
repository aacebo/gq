package gq

import (
	"context"
	"reflect"
	"strconv"

	"github.com/aacebo/gq/query"
)

type List struct {
	Type Schema       `json:"type,omitempty"`
	Use  []Middleware `json:"-"`
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
		return nil, nil
	}

	if self.Use != nil {
		for _, use := range self.Use {
			if err := use(params); err != nil {
				return nil, err
			}
		}
	}

	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return nil, NewError(params.Key, "must be an array/slice")
	}

	err := NewEmptyError(params.Key)

	for i := 0; i < value.Len(); i++ {
		index := value.Index(i)
		res, e := self.Type.Resolve(Params{
			Query:   params.Query,
			Parent:  params.Value,
			Key:     strconv.Itoa(i),
			Value:   index.Interface(),
			Context: params.Context,
		})

		if e != nil {
			err = err.Add(e)
			continue
		}

		index.Set(reflect.ValueOf(res))
	}

	if len(err.Errors) > 0 {
		return value.Interface(), err
	}

	return value.Interface(), nil
}
