package gq_test

import (
	"testing"

	"github.com/aacebo/gq"
)

func Test_List(t *testing.T) {
	type User struct {
		ID    string  `json:"id"`
		Name  string  `json:"name"`
		Email *string `json:"email,omitempty"`
	}

	t.Run("should resolve", func(t *testing.T) {
		schema := gq.List{
			Type: gq.Object[map[string]string]{
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
			},
		}

		res, err := schema.Do(nil, "{id,name,email}", []map[string]string{{
			"id":   "1",
			"name": "dev",
		}})

		if err != nil {
			t.Fatal(err)
		}

		value, ok := res.([]map[string]string)

		if !ok {
			t.FailNow()
		}

		if len(value) != 1 {
			t.FailNow()
		}

		if value[0]["id"] != "1" {
			t.Fatalf("expected `%s`, received `%s`", "1", value[0]["id"])
		}

		if value[0]["name"] != "dev" {
			t.Fatalf("expected `%s`, received `%s`", "dev", value[0]["name"])
		}

		if value[0]["email"] != "dev@gmail.com" {
			t.Fatalf("expected `%s`, received `%s`", "dev@gmail.com", value[0]["email"])
		}
	})

	t.Run("should resolve list of users", func(t *testing.T) {
		schema := gq.List{
			Type: gq.Object[User]{
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
			},
		}

		res, err := schema.Do(nil, "{id,name,email}", []User{{
			ID:   "1",
			Name: "dev",
		}})

		if err != nil {
			t.Fatal(err)
		}

		value, ok := res.([]User)

		if !ok {
			t.FailNow()
		}

		if len(value) != 1 {
			t.Fatalf("should have length of 1")
		}

		if value[0].ID != "1" {
			t.Fatalf("expected `%s`, received `%s`", "1", value[0].ID)
		}

		if value[0].Name != "dev" {
			t.Fatalf("expected `%s`, received `%s`", "dev", value[0].Name)
		}

		if value[0].Email == nil {
			t.Fatalf("expected `%s`, received null", "dev")
		}

		if *value[0].Email != "dev@gmail.com" {
			t.Fatalf("expected `%s`, received `%s`", "dev@gmail.com", *value[0].Email)
		}
	})
}
