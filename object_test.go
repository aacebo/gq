package gq_test

import (
	"errors"
	"fmt"
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
						Resolver: func(params gq.ResolveParams) (any, error) {
							return "dev@gmail.com", nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,email}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error != nil {
				t.Fatal(res.Error)
			}

			value, ok := res.Data.(map[string]string)

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
						Resolver: func(params gq.ResolveParams) (any, error) {
							email := "dev@gmail.com"
							return &email, nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,email}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error == nil {
				t.FailNow()
			}

			if res.Error.Error() != `{"key":"User","errors":[{"key":"email","message":"expected type 'string', received '*string'"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"email","message":"expected type 'string', received '*string'"}]}`,
					res.Error.Error(),
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
						Resolver: func(params gq.ResolveParams) (any, error) {
							email := "dev@gmail.com"
							return &email, nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,test}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error == nil {
				t.FailNow()
			}

			if res.Error.Error() != `{"key":"User","errors":[{"key":"test","message":"field not found"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"test","message":"field not found"}]}`,
					res.Error.Error(),
				)
			}
		})
	})

	t.Run("struct", func(t *testing.T) {
		type Org struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		type User struct {
			ID        string  `json:"id"`
			Name      string  `json:"name"`
			Email     *string `json:"email,omitempty"`
			Orgs      []Org   `json:"orgs,omitempty"`
			CreatedBy *User   `json:"created_by,omitempty"`
		}

		t.Run("should resolve", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						DependsOn: []string{"name"},
						Resolver: func(params gq.ResolveParams) (any, error) {
							parent := params.Parent.(User)
							email := fmt.Sprintf("%s@gmail.com", parent.Name)
							return &email, nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,email}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error != nil {
				t.Fatal(res.Error)
			}

			value, ok := res.Data.(User)

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

		t.Run("should resolve null field", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":    gq.Field{},
					"name":  gq.Field{},
					"email": gq.Field{},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,email}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error != nil {
				t.Fatal(res.Error)
			}

			value, ok := res.Data.(User)

			if !ok {
				t.FailNow()
			}

			if value.ID != "1" {
				t.Fatalf("expected `%s`, received `%s`", "1", value.ID)
			}

			if value.Name != "dev" {
				t.Fatalf("expected `%s`, received `%s`", "dev", value.Name)
			}

			if value.Email != nil {
				t.Fatalf("expected nil, received `%s`", "dev@gmail.com")
			}
		})

		t.Run("should fail when wrong field type", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.ResolveParams) (any, error) {
							return "dev@gmail.com", nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,email}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error == nil {
				t.FailNow()
			}

			if res.Error.Error() != `{"key":"User","errors":[{"key":"email","message":"expected type '*string', received 'string'"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"email","message":"expected type '*string', received 'string'"}]}`,
					res.Error.Error(),
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
						Resolver: func(params gq.ResolveParams) (any, error) {
							return "dev@gmail.com", nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,test}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error == nil {
				t.FailNow()
			}

			if res.Error.Error() != `{"key":"User","errors":[{"key":"test","message":"field not found"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"test","message":"field not found"}]}`,
					res.Error.Error(),
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
						Resolver: func(params gq.ResolveParams) (any, error) {
							parent := params.Parent.(User)
							return fmt.Sprintf("%s@gmail.com", parent.Name), errors.New("a test error")
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,email}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error == nil {
				t.FailNow()
			}

			if res.Error.Error() != `{"key":"User","errors":[{"key":"email","message":"a test error"}]}` {
				t.Fatalf(
					"expected `%s`, received `%s`",
					`{"key":"User","errors":[{"key":"email","message":"a test error"}]}`,
					res.Error.Error(),
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
						Resolver: func(params gq.ResolveParams) (any, error) {
							parent := params.Parent.(User)
							return &parent, nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,created_by{id,name}}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error != nil {
				t.Fatal(res.Error)
			}

			value, ok := res.Data.(User)

			if !ok {
				t.FailNow()
			}

			if value.CreatedBy == nil {
				t.Fatalf("'created_by' should not be nil")
			}
		})

		t.Run("should resolve with nested list schema", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"orgs": gq.Field{
						Type: gq.List{
							Type: gq.Object[Org]{
								Name: "Org",
								Fields: gq.Fields{
									"id":   gq.Field{},
									"name": gq.Field{},
								},
							},
						},
						Resolver: func(params gq.ResolveParams) (any, error) {
							return []Org{
								{ID: "1", Name: "one"},
								{ID: "2", Name: "two"},
							}, nil
						},
					},
				},
			}

			res := schema.Do(gq.DoParams{
				Query: "{id,name,orgs{id,name}}",
				Value: map[string]string{
					"id":   "1",
					"name": "dev",
				},
			})

			if res.Error != nil {
				t.Fatal(res.Error)
			}

			value, ok := res.Data.(User)

			if !ok {
				t.FailNow()
			}

			if value.Orgs == nil {
				t.Fatalf("'orgs' should not be nil")
			}

			if len(value.Orgs) != 2 {
				t.Fatalf("should have 2 orgs")
			}

			if value.Orgs[0].ID != "1" {
				t.Fatalf("first org should have `id` = `1`")
			}

			if value.Orgs[0].Name != "one" {
				t.Fatalf("first org should have `name` = `one`")
			}

			if value.Orgs[1].ID != "2" {
				t.Fatalf("second org should have `id` = `2`")
			}

			if value.Orgs[1].Name != "two" {
				t.Fatalf("second org should have `name` = `two`")
			}
		})
	})

	t.Run("extend", func(t *testing.T) {
		t.Run("should extend object", func(t *testing.T) {
			schema := gq.Object[map[string]any]{
				Name: "User",
				Fields: gq.Fields{
					"email":    gq.Field{},
					"password": gq.Field{},
				},
			}.Extend(gq.Object[map[string]any]{
				Name: "User",
				Fields: gq.Fields{
					"staySignedIn": gq.Field{},
				},
			})

			if len(schema.Fields) != 3 {
				t.FailNow()
			}
		})
	})
}
