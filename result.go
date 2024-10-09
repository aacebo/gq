package gq

import "encoding/json"

type Result struct {
	Meta  Meta  `json:"$meta,omitempty"`
	Data  any   `json:"data,omitempty"`
	Error error `json:"error,omitempty"`
}

func (self Result) Merge(result Result) Result {
	if self.Meta == nil {
		self.Meta = result.Meta
	} else if result.Meta != nil {
		self.Meta = self.Meta.Merge(result.Meta)
	}

	if result.Data != nil {
		self.Data = result.Data
	}

	return self
}

func (self Result) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
