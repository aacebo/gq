package middleware_test

import (
	"testing"
	"time"

	"github.com/aacebo/gq"
	"github.com/aacebo/gq/middleware"
)

func TestElapse(t *testing.T) {
	t.Run("should have $elapse", func(t *testing.T) {
		schema := gq.Object[map[string]string]{
			Name: "User",
			Use:  []gq.Middleware{middleware.Elapse},
			Fields: gq.Fields{
				"id":   gq.Field{},
				"name": gq.Field{},
				"email": gq.Field{
					Resolver: func(params *gq.ResolveParams) gq.Result {
						time.Sleep(100 * time.Millisecond)
						return gq.Result{Data: "dev@gmail.com"}
					},
				},
			},
		}

		res := schema.Do(&gq.DoParams{
			Query: "{id,name,email}",
			Value: map[string]string{
				"id":   "1",
				"name": "dev",
			},
		})

		if res.Error != nil {
			t.Fatal(res.Error)
		}

		if res.Meta == nil {
			t.Fatalf("should have $meta object")
		}

		elapse, exists := res.Meta["$elapse"].(int64)

		if !exists {
			t.Fatalf("should have $meta.elapse")
		}

		if elapse < 100 {
			t.Fatalf("should have taken at least 100ms to resolve")
		}
	})
}
