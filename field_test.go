package gq_test

import (
	"reflect"
	"testing"

	"github.com/aacebo/gq"
)

func Test_Field(t *testing.T) {
	t.Run("should resolve", func(t *testing.T) {
		schema := gq.Field{
			Resolver: func(params gq.Params) (any, error) {
				return 1, nil
			},
		}

		res, err := schema.Do(nil, "{}", 2)

		if err != nil {
			t.Error(err)
		}

		value := reflect.ValueOf(res)

		if value.Kind() != reflect.Int {
			t.FailNow()
		}

		if value.Int() != 1 {
			t.Errorf("expected %d, received %d", 1, value.Int())
		}
	})

	t.Run("should resolve using default value", func(t *testing.T) {
		schema := gq.Field{}
		res, err := schema.Do(nil, "{}", 1)

		if err != nil {
			t.Error(err)
		}

		value := reflect.ValueOf(res)

		if value.Kind() != reflect.Int {
			t.FailNow()
		}

		if value.Int() != 1 {
			t.Errorf("expected %d, received %d", 1, value.Int())
		}
	})
}