package gq

import (
	"reflect"
	"strconv"
	"time"

	"github.com/aacebo/gq/query"
)

type List struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Type        Schema       `json:"type,omitempty"`
	Use         []Middleware `json:"-"`
}

func (self List) Do(params *DoParams) Result {
	parser := query.Parser([]byte(params.Query))
	query, err := parser.Parse()

	if err != nil {
		return Result{Error: err}
	}

	return self.Resolve(&ResolveParams{
		Query:   query,
		Key:     self.Name,
		Value:   params.Value,
		Context: params.Context,
	})
}

func (self List) Resolve(params *ResolveParams) Result {
	value := reflect.Indirect(reflect.ValueOf(params.Value))
	now := time.Now()
	res := Result{Meta: Meta{}}

	defer func() {
		res.Meta["$elapse"] = time.Now().Sub(now).Milliseconds()
	}()

	if !value.IsValid() {
		return res
	}

	if self.Use != nil {
		for _, use := range self.Use {
			result := use(params)

			if result.Error != nil {
				res.Error = NewError(params.Key, result.Error.Error())
				return res
			}

			res = res.Merge(result)
		}
	}

	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		res.Error = NewError(params.Key, "must be an array/slice")
		return res
	}

	err := NewEmptyError(params.Key)

	for i := 0; i < value.Len(); i++ {
		key := strconv.Itoa(i)
		index := value.Index(i)
		result := self.Type.Resolve(&ResolveParams{
			Query:   params.Query,
			Parent:  params.Value,
			Key:     key,
			Value:   index.Interface(),
			Context: params.Context,
		})

		if result.Error != nil {
			err = err.Add(result.Error)
			continue
		}

		index.Set(reflect.ValueOf(result.Data))
		res.Meta[key] = result.Meta
	}

	if len(err.Errors) > 0 {
		res.Error = err
		return res
	}

	res.Data = value.Interface()
	return res
}
