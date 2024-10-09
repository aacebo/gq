package gq

import (
	"context"
	"encoding/json"
)

type DoParams struct {
	Query   string          `json:"query"`
	Value   any             `json:"value,omitempty"`
	Context context.Context `json:"context,omitempty"`
}

func (self DoParams) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}

type ResolveParams struct {
	Query   Query           `json:"query"`
	Parent  any             `json:"parent,omitempty"`
	Key     string          `json:"key,omitempty"`
	Value   any             `json:"value,omitempty"`
	Context context.Context `json:"context,omitempty"`
}

func (self ResolveParams) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}

type Schema interface {
	Do(params DoParams) Result
	Resolve(params ResolveParams) Result
}
