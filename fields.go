package gq

import (
	"encoding/json"
)

type Fields map[string]Schema

func (self Fields) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
