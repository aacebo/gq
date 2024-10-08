package query

import (
	"encoding/json"
)

// http://localhost:8080/product?query={product(id:1){name,info,price}}
type Query struct {
	Args   QueryArgs        `json:"args"`
	Fields map[string]Query `json:"fields"`
}

func (self Query) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
