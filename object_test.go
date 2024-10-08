package gq_test

import (
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
						Resolver: func(params gq.Params) (any, error) {
							parent := params.Parent.(map[string]string)
							return fmt.Sprintf("%s@gmail.com", parent["name"]), nil
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
							parent := params.Parent.(map[string]string)
							email := fmt.Sprintf("%s@gmail.com", parent["name"])
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
		})

		t.Run("should fail when query field not found", func(t *testing.T) {
			schema := gq.Object[map[string]string]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							parent := params.Parent.(map[string]string)
							email := fmt.Sprintf("%s@gmail.com", parent["name"])
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
		})
	})

	t.Run("struct", func(t *testing.T) {
		type User struct {
			ID    string  `json:"id"`
			Name  string  `json:"name"`
			Email *string `json:"email,omitempty"`
		}

		t.Run("should resolve", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							parent := params.Parent.(User)
							email := fmt.Sprintf("%s@gmail.com", parent.Name)
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
							parent := params.Parent.(User)
							return fmt.Sprintf("%s@gmail.com", parent.Name), nil
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
		})

		t.Run("should fail when query field not found", func(t *testing.T) {
			schema := gq.Object[User]{
				Name: "User",
				Fields: gq.Fields{
					"id":   gq.Field{},
					"name": gq.Field{},
					"email": gq.Field{
						Resolver: func(params gq.Params) (any, error) {
							parent := params.Parent.(User)
							return fmt.Sprintf("%s@gmail.com", parent.Name), nil
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
		})
	})
}
