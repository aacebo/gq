package gq_test

import (
	"testing"

	"github.com/aacebo/gq"
)

func Test_List(t *testing.T) {
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
}
