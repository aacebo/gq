package gq

import "encoding/json"

type Meta map[string]any

func (self Meta) Merge(meta Meta) Meta {
	for key, value := range meta {
		self[key] = value
	}

	return self
}

func (self Meta) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
