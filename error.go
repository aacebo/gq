package gq

import (
	"encoding/json"
)

type Error struct {
	Key     string  `json:"key,omitempty"`
	Message string  `json:"message,omitempty"`
	Errors  []error `json:"errors,omitempty"`
}

func NewError(key string, message string) Error {
	return Error{
		Key:     key,
		Message: message,
		Errors:  []error{},
	}
}

func NewEmptyError(key string) Error {
	return Error{
		Key:    key,
		Errors: []error{},
	}
}

func (self Error) Add(err error) Error {
	if err, ok := err.(Error); ok {
		self.Errors = append(self.Errors, err)
		return self
	}

	self.Errors = append(self.Errors, NewError("", err.Error()))
	return self
}

func (self Error) Error() string {
	b, _ := json.Marshal(self)
	return string(b)
}

func (self Error) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
