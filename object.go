package gq

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/aacebo/gq/query"
)

type Object[T any] struct {
	Name        string       `json:"name"` // must be unique
	Description string       `json:"description,omitempty"`
	Use         []Middleware `json:"-"`
	Fields      Fields       `json:"fields,omitempty"`
}

func (self Object[T]) Key() string {
	return self.Name
}

func (self Object[T]) Do(params *DoParams) Result {
	parser := query.Parser([]byte(params.Query))
	query, err := parser.Parse()

	if err != nil {
		return Result{Error: err}
	}

	return self.Resolve(&ResolveParams{
		Query:   query,
		Key:     self.Name,
		Value:   params.Value,
		Context: params.Context,
	})
}

func (self Object[T]) Resolve(params *ResolveParams) Result {
	res := Result{Meta: Meta{}}
	routes := []Middleware{}

	if self.Use != nil {
		for _, route := range self.Use {
			routes = append(routes, route)
		}
	}

	routes = append(routes, self.resolve)

	var next Resolver

	i := -1
	next = func(params *ResolveParams) Result {
		i++

		if i > (len(routes) - 1) {
			return Result{}
		}

		return routes[i](params, next)
	}

	result := next(&ResolveParams{
		Query:   params.Query,
		Parent:  params.Parent,
		Key:     "fields",
		Value:   params.Value,
		Context: params.Context,
	})

	if result.Error != nil {
		res.Error = NewEmptyError(params.Key).Add(result.Error)
		return res
	}

	res.Meta = result.Meta
	res.Data = result.Data
	return res
}

func (self Object[T]) Extend(schema Object[T]) Object[T] {
	fields := Fields{}

	for key, value := range self.Fields {
		fields[key] = value
	}

	if schema.Fields != nil {
		for key, value := range schema.Fields {
			fields[key] = value
		}
	}

	middleware := []Middleware{}

	if self.Use != nil {
		for _, use := range self.Use {
			middleware = append(middleware, use)
		}
	}

	if schema.Use != nil {
		for _, use := range schema.Use {
			middleware = append(middleware, use)
		}
	}

	return Object[T]{
		Name:        schema.Name,
		Description: schema.Description,
		Use:         middleware,
		Fields:      fields,
	}
}

func (self Object[T]) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}

func (self Object[T]) resolve(params *ResolveParams, _ Resolver) Result {
	res := Result{}

	if params.Value == nil || self.Fields == nil {
		return res
	}

	err := NewEmptyError(params.Key)
	object := reflect.Indirect(reflect.New(reflect.TypeFor[T]()))

	if object.Kind() == reflect.Pointer {
		object = reflect.New(object.Type().Elem())
	}

	visited := map[string]bool{}

	var resolve func(key string) error
	resolve = func(key string) error {
		if visited[key] {
			return nil
		}

		query := params.Query.Fields[key]
		field, exists := self.Fields[key]

		if !exists {
			return NewError(key, "field not found")
		}

		if field.DependsOn != nil {
			for _, dependsOn := range field.DependsOn {
				e := resolve(dependsOn)

				if e != nil {
					return e
				}
			}
		}

		result := field.Resolve(&ResolveParams{
			Query:   query,
			Parent:  object.Interface(),
			Key:     key,
			Value:   self.getKey(key, reflect.ValueOf(params.Value)),
			Context: params.Context,
		})

		if result.Error != nil {
			return result.Error
		}

		if e := self.setKey(key, result.Data, object); e != nil {
			return e
		}

		if result.Meta != nil && !result.Meta.Empty() {
			res.SetMeta(key, result.Meta)
		}

		visited[key] = true
		return nil
	}

	for key := range params.Query.Fields {
		e := resolve(key)

		if e != nil {
			err = err.Add(e)
		}
	}

	if len(err.Errors) > 0 {
		res.Error = err
		return res
	}

	if res.Meta != nil && res.Meta.Empty() {
		res.Meta = nil
	}

	res.Data = object.Interface()
	return res
}

func (self Object[T]) getKey(key string, value reflect.Value) any {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Map {
		return self.getMapKey(key, value)
	}

	return self.getStructKey(key, value)
}

func (self Object[T]) setKey(key string, val any, value reflect.Value) error {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Interface {
		value = value.Elem()
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
	value := reflect.ValueOf(val)

	if object.CanSet() && object.IsNil() {
		object.Set(reflect.MakeMapWithSize(reflect.TypeFor[T](), 0))
	}

	if object.Type().Elem() != value.Type() && value.CanConvert(object.Type().Elem()) {
		value = value.Convert(object.Type().Elem())
	}

	if object.Type().Elem() != value.Type() {
		return NewError(
			key,
			fmt.Sprintf(
				"expected type '%s', received '%s'",
				object.Type().Elem().String(),
				value.Type().String(),
			),
		)
	}

	object.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	return nil
}

func (self Object[T]) getStructKey(key string, object reflect.Value) any {
	name, exists := self.getStructFieldByName(key, object)

	if !exists {
		return nil
	}

	value := object.FieldByName(name)

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
		return NewError(key, "struct field not found")
	}

	value := object.FieldByName(name)

	if !reflect.ValueOf(val).IsValid() {
		if value.CanSet() {
			value.Set(reflect.New(value.Type()).Elem())
		}

		return nil
	}

	if value.Type() != reflect.ValueOf(val).Type() {
		return NewError(
			key,
			fmt.Sprintf(
				"expected type '%s', received '%s'",
				value.Type().String(),
				reflect.ValueOf(val).Type().String(),
			),
		)
	}

	if value.CanSet() {
		value.Set(reflect.ValueOf(val))
	}

	return nil
}

func (self Object[T]) getStructFieldByName(name string, object reflect.Value) (string, bool) {
	if !object.IsValid() {
		return "", false
	}

	for i := 0; i < object.Type().NumField(); i++ {
		field := object.Type().Field(i)
		tag := field.Tag.Get("json")

		if tag == "" {
			tag = field.Name
		}

		if i := strings.Index(tag, ","); i > -1 {
			tag = tag[:i]
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
