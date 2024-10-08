package gq

import (
	"context"
)

type Params struct {
	Query   Query           `json:"query"`
	Parent  any             `json:"parent,omitempty"`
	Value   any             `json:"value,omitempty"`
	Context context.Context `json:"context"`
}

type Schema interface {
	Do(ctx context.Context, q string, value any) (any, error)
	Resolve(params Params) (any, error)
}
