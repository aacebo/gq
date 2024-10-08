package query

import (
	"encoding/json"
)

type QueryArgs map[string]any

func (self QueryArgs) Get(key string) any {
	if value, exists := self[key]; exists {
		return value
	}

	return nil
}

func (self QueryArgs) TryGet(key string, defaultValue any) any {
	if value, exists := self[key]; exists {
		return value
	}

	return defaultValue
}

func (self QueryArgs) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
