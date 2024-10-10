package gq

import "encoding/json"

type Result struct {
	Meta  Meta  `json:"$meta,omitempty"`
	Data  any   `json:"data,omitempty"`
	Error error `json:"error,omitempty"`
}

func (self *Result) SetMeta(key string, value any) *Result {
	if self.Meta == nil {
		self.Meta = Meta{}
	}

	self.Meta[key] = value
	return self
}

func (self Result) Merge(result Result) Result {
	if result.Meta != nil {
		if self.Meta == nil {
			self.Meta = result.Meta
		} else {
			self.Meta = self.Meta.Merge(result.Meta)
		}
	}

	if result.Data != nil {
		self.Data = result.Data
	}

	if result.Error != nil {
		self.Error = result.Error
	}

	return self
}

func (self Result) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
