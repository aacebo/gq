package gq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aacebo/gq/query"
)

type Object[T any] struct {
	Name        string `json:"name"` // must be unique
	Description string `json:"description,omitempty"`
	Fields      Fields `json:"fields,omitempty"`
}

func (self Object[T]) Do(ctx context.Context, q string, value any) (any, error) {
	parser := query.Parser([]byte(q))
	query, err := parser.Parse()

	if err != nil {
		return nil, err
	}

	return self.Resolve(Params{
		Query:   query,
		Value:   value,
		Context: ctx,
	})
}

func (self Object[T]) Resolve(params Params) (any, error) {
	if self.Fields == nil {
		return nil, nil
	}

	object := new(T)

	for key, subQuery := range params.Query.Fields {
		schema, exists := self.Fields[key]

		if !exists {
			return nil, fmt.Errorf("field %s not found", key)
		}

		if field, ok := schema.(Field); ok && field.Args != nil {
			if err := field.Args.Validate(subQuery.Args); err != nil {
				return nil, err
			}
		}

		value, err := schema.Resolve(Params{
			Query:   subQuery,
			Parent:  params.Value,
			Value:   self.getKey(key, params.Value),
			Context: params.Context,
		})

		if err != nil {
			return nil, err
		}

		if field, ok := schema.(Field); ok && field.Type != nil {
			res, err := field.Type.Resolve(Params{
				Query:   subQuery,
				Parent:  params.Value,
				Value:   value,
				Context: params.Context,
			})

			if err != nil {
				return nil, err
			}

			value = res
		}

		if err = self.setKey(key, value, object); err != nil {
			return nil, err
		}
	}

	return *object, nil
}

func (self Object[T]) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}

func (self Object[T]) getKey(key string, object any) any {
	value := reflect.Indirect(reflect.ValueOf(object))

	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Map {
		return self.getMapKey(key, value)
	}

	return self.getStructKey(key, value)
}

func (self Object[T]) setKey(key string, val any, object any) error {
	value := reflect.Indirect(reflect.ValueOf(object))

	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Map {
		return self.setMapKey(key, val, value)
	}

	return self.setStructKey(key, val, value)
}

func (self Object[T]) getMapKey(key string, object reflect.Value) any {
	value := reflect.Indirect(object.MapIndex(reflect.ValueOf(key)))

	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	if value.IsValid() && value.CanInterface() {
		return value.Interface()
	}

	return nil
}

func (self Object[T]) setMapKey(key string, val any, object reflect.Value) error {
	value := object.MapIndex(reflect.ValueOf(key))

	if value.CanSet() {
		value.Set(reflect.ValueOf(val))
	}

	return nil
}

func (self Object[T]) getStructKey(key string, object reflect.Value) any {
	name, exists := self.getStructFieldByName(key, object)

	if !exists {
		return nil
	}

	value := reflect.Indirect(object.FieldByName(name))

	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	if value.IsValid() && value.CanInterface() {
		return value.Interface()
	}

	return nil
}

func (self Object[T]) setStructKey(key string, val any, object reflect.Value) error {
	name, exists := self.getStructFieldByName(key, object)

	if !exists {
		return fmt.Errorf("[%s] => field `%s` not found", self.Name, key)
	}

	value := object.FieldByName(name)

	if value.CanSet() {
		value.Set(reflect.ValueOf(val))
	}

	return nil
}

func (self Object[T]) getStructFieldByName(name string, object reflect.Value) (string, bool) {
	if !object.IsValid() {
		return "", false
	}

	for i := 0; i < object.NumField(); i++ {
		field := object.Type().Field(i)
		tag := field.Tag.Get("json")

		if tag == "" {
			tag = field.Name
		}

		if tag == "" || tag == "-" {
			continue
		}

		if tag == name {
			return field.Name, true
		}
	}

	return "", false
}
