package gq

import (
	"encoding/json"
)

type Args interface {
	Validate(value any) error
}

type Field struct {
	Type        Schema       `json:"type,omitempty"`
	Description string       `json:"description,omitempty"`
	Args        Args         `json:"args,omitempty"`
	DependsOn   []string     `json:"depends_on,omitempty"`
	Use         []Middleware `json:"-"`
	Resolver    Resolver     `json:"-"`
}

func (self Field) Resolve(params *ResolveParams) Result {
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

	res := next(&ResolveParams{
		Query:   params.Query,
		Parent:  params.Parent,
		Key:     params.Key,
		Value:   params.Value,
		Context: params.Context,
	})

	return res
}

func (self Field) resolve(params *ResolveParams, _ Resolver) Result {
	res := Result{Meta: Meta{}}

	if self.Args != nil {
		if err := self.Args.Validate(params.Value); err != nil {
			res.Error = NewError(params.Key, err.Error())
			return res
		}
	}

	if self.Resolver != nil {
		result := self.Resolver(params)

		if result.Error != nil {
			res.Error = NewError(params.Key, result.Error.Error())
			return res
		}

		if result.Meta != nil && !result.Meta.Empty() {
			res.Meta = res.Meta.Merge(result.Meta)
		}

		params.Value = result.Data
	}

	if self.Type != nil {
		result := self.Type.Resolve(params)

		if result.Error != nil {
			res.Error = result.Error
			return res
		}

		if result.Meta != nil && !result.Meta.Empty() {
			res.Meta = res.Meta.Merge(result.Meta)
		}

		params.Value = result.Data
	}

	if res.Meta.Empty() {
		res.Meta = nil
	}

	res.Data = params.Value
	return res
}

func (self Field) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
