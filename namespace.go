package gq

import (
	"encoding/json"
	"fmt"
)

type Namespace struct {
	schemas map[string]Schema
}

func New() *Namespace {
	return &Namespace{schemas: map[string]Schema{}}
}

func (self *Namespace) Register(schema Schema) *Namespace {
	self.schemas[schema.Key()] = schema
	return self
}

func (self Namespace) Get(key string) Schema {
	return self.schemas[key]
}

func (self Namespace) Ref(key string) Ref {
	return Ref{resolve: func() Schema {
		return self.schemas[key]
	}}
}

func (self Namespace) Do(key string, params *DoParams) Result {
	schema, exists := self.schemas[key]

	if !exists {
		return Result{Error: NewError("", fmt.Sprintf("schema '%s' not found", key))}
	}

	return schema.Do(params)
}

func (self Namespace) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.schemas)
}
