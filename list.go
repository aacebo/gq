package gq

import (
	"reflect"
	"strconv"

	"github.com/aacebo/gq/query"
)

type List struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Type        Schema       `json:"type,omitempty"`
	Use         []Middleware `json:"-"`
}

func (self List) Key() string {
	return self.Name
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
	res := Result{}
	routes := []Middleware{}

	if self.Use != nil {
		for _, route := range self.Use {
			routes = append(routes, route)
		}
	}

	routes = append(routes, self.resolve)

	var next Resolver

	i := -1
	next = func(params *ResolveParams) Result {
		i++

		if i > (len(routes) - 1) {
			return Result{}
		}

		return routes[i](params, next)
	}

	result := next(&ResolveParams{
		Query:   params.Query,
		Parent:  params.Parent,
		Key:     "items",
		Value:   params.Value,
		Context: params.Context,
	})

	if result.Error != nil {
		res.Error = NewEmptyError(params.Key).Add(result.Error)
		return res
	}

	res.Meta = result.Meta
	res.Data = result.Data
	return res
}

func (self List) resolve(params *ResolveParams, _ Resolver) Result {
	value := reflect.Indirect(reflect.ValueOf(params.Value))
	res := Result{}

	if !value.IsValid() {
		return res
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

		if result.Meta != nil && !result.Meta.Empty() {
			res.SetMeta(key, result.Meta)
		}
	}

	if len(err.Errors) > 0 {
		res.Error = err
		return res
	}

	if res.Meta != nil && res.Meta.Empty() {
		res.Meta = nil
	}

	res.Data = value.Interface()
	return res
}
