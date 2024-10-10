package gq_test

import (
	"reflect"
	"testing"

	"github.com/aacebo/gq"
)

func Test_Field(t *testing.T) {
	t.Run("should resolve", func(t *testing.T) {
		schema := gq.Field{
			Resolver: func(params *gq.ResolveParams) (any, error) {
				return 1, nil
			},
		}

		res := schema.Do(&gq.DoParams{
			Query: "{}",
			Value: 2,
		})

		if res.Error != nil {
			t.Error(res.Error)
		}

		value := reflect.ValueOf(res.Data)

		if value.Kind() != reflect.Int {
			t.FailNow()
		}

		if value.Int() != 1 {
			t.Errorf("expected %d, received %d", 1, value.Int())
		}
	})

	t.Run("should resolve using default value", func(t *testing.T) {
		schema := gq.Field{}
		res := schema.Do(&gq.DoParams{
			Query: "{}",
			Value: 1,
		})

		if res.Error != nil {
			t.Error(res.Error)
		}

		value := reflect.ValueOf(res.Data)

		if value.Kind() != reflect.Int {
			t.FailNow()
		}

		if value.Int() != 1 {
			t.Errorf("expected %d, received %d", 1, value.Int())
		}
	})
}
