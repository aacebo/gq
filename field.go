package gq

import (
	"encoding/json"
	"time"

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
	now := time.Now()
	res := Result{Meta: Meta{}}

	defer func() {
		res.Meta["$elapse"] = time.Now().Sub(now).Milliseconds()
	}()

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
