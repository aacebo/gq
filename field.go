package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Args interface {
	Validate(value any) error
}

type Field struct {
	Type        Schema                                   `json:"type,omitempty"`
	Description string                                   `json:"description,omitempty"`
	Args        Args                                     `json:"args,omitempty"`
	DependsOn   []string                                 `json:"depends_on,omitempty"`
	Use         []Middleware                             `json:"-"`
	Resolver    func(params *ResolveParams) (any, error) `json:"-"`
}

func (self Field) Do(params *DoParams) Result {
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

func (self Field) Resolve(params *ResolveParams) Result {
	res := Result{Meta: Meta{}}
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
		Key:     params.Key,
		Value:   params.Value,
		Context: params.Context,
	})

	if result.Error != nil {
		res.Error = result.Error
		return res
	}

	res.Meta = result.Meta
	res.Data = result.Data
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
		value, err := self.Resolver(params)

		if err != nil {
			res.Error = NewError(params.Key, err.Error())
			return res
		}

		params.Value = value
	}

	if self.Type != nil {
		result := self.Type.Resolve(params)

		if result.Error != nil {
			res.Error = NewEmptyError(params.Key).Add(result.Error)
			return res
		}

		if result.Meta != nil {
			res.Meta = res.Meta.Merge(result.Meta)
		}

		params.Value = result.Data
	}

	res.Data = params.Value
	return res
}

func (self Field) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
