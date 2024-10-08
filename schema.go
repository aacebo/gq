package gq

import (
	"context"
	"encoding/json"
)

type Params struct {
	Query   Query           `json:"query"`
	Parent  any             `json:"parent,omitempty"`
	Key     string          `json:"key"`
	Value   any             `json:"value,omitempty"`
	Context context.Context `json:"context"`
}

func (self Params) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}

type Schema interface {
	Do(ctx context.Context, q string, value any) (any, error)
	Resolve(params Params) (any, error)
}
