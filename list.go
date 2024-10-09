package gq

import (
	"reflect"
	"strconv"

	"github.com/aacebo/gq/query"
)

type List struct {
	Type Schema       `json:"type,omitempty"`
	Use  []Middleware `json:"-"`
}

func (self List) Do(params DoParams) Result {
	parser := query.Parser([]byte(params.Query))
	query, err := parser.Parse()

	if err != nil {
		return Result{Error: err}
	}

	return self.Resolve(ResolveParams{
		Query:   query,
		Value:   params.Value,
		Context: params.Context,
	})
}

func (self List) Resolve(params ResolveParams) Result {
	value := reflect.Indirect(reflect.ValueOf(params.Value))

	if !value.IsValid() {
		return Result{}
	}

	if self.Use != nil {
		for _, use := range self.Use {
			if err := use(params); err != nil {
				return Result{
					Data:  params.Value,
					Error: err,
				}
			}
		}
	}

	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return Result{Error: NewError(params.Key, "must be an array/slice")}
	}

	err := NewEmptyError(params.Key)

	for i := 0; i < value.Len(); i++ {
		index := value.Index(i)
		res := self.Type.Resolve(ResolveParams{
			Query:   params.Query,
			Parent:  params.Value,
			Key:     strconv.Itoa(i),
			Value:   index.Interface(),
			Context: params.Context,
		})

		if res.Error != nil {
			err = err.Add(res.Error)
			continue
		}

		index.Set(reflect.ValueOf(res.Data))
	}

	if len(err.Errors) > 0 {
		return Result{
			Data:  value.Interface(),
			Error: err,
		}
	}

	return Result{Data: value.Interface()}
}
