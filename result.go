package gq

import "encoding/json"

type Meta map[string]any

type Result struct {
	Meta  Meta  `json:"$meta,omitempty"`
	Data  any   `json:"data,omitempty"`
	Error error `json:"error,omitempty"`
}

func (self Result) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
