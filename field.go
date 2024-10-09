package gq

import (
	"encoding/json"

	"github.com/aacebo/gq/query"
)

type Args interface {
	Validate(value any) error
}

type Field struct {
	Type        Schema                                  `json:"type,omitempty"`
	Description string                                  `json:"description,omitempty"`
	Args        Args                                    `json:"args,omitempty"`
	DependsOn   []string                                `json:"depends_on,omitempty"`
	Use         []Middleware                            `json:"-"`
	Resolver    func(params ResolveParams) (any, error) `json:"-"`
}

func (self Field) Do(params DoParams) Result {
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

func (self Field) Resolve(params ResolveParams) Result {
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

	if self.Args != nil {
		if err := self.Args.Validate(params.Value); err != nil {
			return Result{
				Data:  params.Value,
				Error: NewError(params.Key, err.Error()),
			}
		}
	}

	if self.Resolver != nil {
		value, err := self.Resolver(params)

		if err != nil {
			return Result{
				Data:  value,
				Error: NewError(params.Key, err.Error()),
			}
		}

		params.Value = value
	}

	if self.Type != nil {
		res := self.Type.Resolve(params)

		if res.Error != nil {
			return Result{
				Data:  res.Data,
				Error: NewEmptyError(params.Key).Add(res.Error),
			}
		}

		params.Value = res.Data
	}

	return Result{Data: params.Value}
}

func (self Field) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
