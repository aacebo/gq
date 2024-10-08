package gq_test

import (
	"errors"
	"testing"

	"github.com/aacebo/gq"
)

func Test_Object(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		t.Run("should resolve", func(t *testing.T) {
			schema := gq.Object[map[string]string]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							return "dev@gmail.com", nil
						},
					},
				},
			}

			res, err := schema.Do(nil, "{id,name,email}", map[string]string{
				"id":   "1",
				"name": "dev",
			})

			if err != nil {
				t.Fatal(err)
			}

			value, ok := res.(map[string]string)

			if !ok {
				t.FailNow()
			}

			if value["id"] != "1" {
				t.Fatalf("expected `%s`, received `%s`", "1", value["id"])
			}

			if value["name"] != "dev" {
				t.Fatalf("expected `%s`, received `%s`", "dev", value["name"])
			}

			if value["email"] != "dev@gmail.com" {
				t.Fatalf("expected `%s`, received `%s`", "dev@gmail.com", value["email"])
			}
		})

		t.Run("should fail when wrong field type", func(t *testing.T) {
			schema := gq.Object[map[string]string]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							email := "dev@gmail.com"
							return &email, nil
						},
					},
				},
			}

			_, err := schema.Do(nil, "{id,name,email}", map[string]string{
				"id":   "1",
				"name": "dev",
			})

			if err == nil {
				t.FailNow()
			}

			if err.Error() != `{"key":"User","errors":[{"key":"email","message":"expected type 'string', received '*string'"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"email","message":"expected type 'string', received '*string'"}]}`,
					err.Error(),
				)
			}
		})

		t.Run("should fail when query field not found", func(t *testing.T) {
			schema := gq.Object[map[string]string]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							email := "dev@gmail.com"
							return &email, nil
						},
					},
				},
			}

			_, err := schema.Do(nil, "{id,name,test}", map[string]string{
				"id":   "1",
				"name": "dev",
			})

			if err == nil {
				t.FailNow()
			}

			if err.Error() != `{"key":"User","errors":[{"key":"test","message":"field not found"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"test","message":"field not found"}]}`,
					err.Error(),
				)
			}
		})
	})

	t.Run("struct", func(t *testing.T) {
		type User struct {
			ID        string  `json:"id"`
			Name      string  `json:"name"`
			Email     *string `json:"email,omitempty"`
			CreatedBy *User   `json:"created_by,omitempty"`
		}

		t.Run("should resolve", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							email := "dev@gmail.com"
							return &email, nil
						},
					},
				},
			}

			res, err := schema.Do(nil, "{id,name,email}", User{
				ID:   "1",
				Name: "dev",
			})

			if err != nil {
				t.Fatal(err)
			}

			value, ok := res.(User)

			if !ok {
				t.FailNow()
			}

			if value.ID != "1" {
				t.Fatalf("expected `%s`, received `%s`", "1", value.ID)
			}

			if value.Name != "dev" {
				t.Fatalf("expected `%s`, received `%s`", "dev", value.Name)
			}

			if value.Email == nil {
				t.Fatalf("expected `%s`, received nil", "dev@gmail.com")
			}

			if *value.Email != "dev@gmail.com" {
				t.Fatalf("expected `%s`, received `%s`", "dev@gmail.com", *value.Email)
			}
		})

		t.Run("should fail when wrong field type", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							return "dev@gmail.com", nil
						},
					},
				},
			}

			_, err := schema.Do(nil, "{id,name,email}", User{
				ID:   "1",
				Name: "dev",
			})

			if err == nil {
				t.FailNow()
			}

			if err.Error() != `{"key":"User","errors":[{"key":"email","message":"expected type '*string', received 'string'"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"email","message":"expected type '*string', received 'string'"}]}`,
					err.Error(),
				)
			}
		})

		t.Run("should fail when query field not found", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							return "dev@gmail.com", nil
						},
					},
				},
			}

			_, err := schema.Do(nil, "{id,name,test}", User{
				ID:   "1",
				Name: "dev",
			})

			if err == nil {
				t.FailNow()
			}

			if err.Error() != `{"key":"User","errors":[{"key":"test","message":"field not found"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"test","message":"field not found"}]}`,
					err.Error(),
				)
			}
		})

		t.Run("should fail when field errors", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							return "dev@gmail.com", errors.New("a test error")
						},
					},
				},
			}

			_, err := schema.Do(nil, "{id,name,email}", User{
				ID:   "1",
				Name: "dev",
			})

			if err == nil {
				t.FailNow()
			}

			if err.Error() != `{"key":"User","errors":[{"key":"email","message":"a test error"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"email","message":"a test error"}]}`,
					err.Error(),
				)
			}
		})

		t.Run("should resolve with nested object schema", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"created_by": gq.Field{
						Type: gq.Object[*User]{
							Name: "CreatedBy",
							Fields: gq.Fields{
								"id":   gq.Field{},
								"name": gq.Field{},
							},
						},
						Resolver: func(params gq.Params) (any, error) {
							parent := params.Parent.(User)
							return &parent, nil
						},
					},
				},
			}

			res, err := schema.Do(nil, "{id,name,created_by{id,name}}", User{
				ID:   "1",
				Name: "dev",
			})

			if err != nil {
				t.Fatal(err)
			}

			value, ok := res.(User)

			if !ok {
				t.FailNow()
			}

			if value.CreatedBy == nil {
				t.Fatalf("'created_by' should not be nil")
			}
		})
	})
}
