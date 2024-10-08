package query

import (
	"encoding/json"
)

type Query struct {
	Args   QueryArgs        `json:"args,omitempty"`
	Fields map[string]Query `json:"fields,omitempty"`
}

func (self Query) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
