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
				t.Error(err)
			}

			value, ok := res.(map[string]string)

			if !ok {
				t.FailNow()
			}

			if value["id"] != "1" {
				t.Errorf("expected `%s`, received `%s`", "1", value["id"])
			}

			if value["name"] != "dev" {
				t.Errorf("expected `%s`, received `%s`", "dev", value["name"])
			}

			if value["email"] != "dev@gmail.com" {
				t.Errorf("expected `%s`, received `%s`", "dev@gmail.com", value["email"])
			}
		})
	})
}
