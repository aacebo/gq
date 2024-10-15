package middleware_test

import (
	"testing"
	"time"

	"github.com/aacebo/gq"
	"github.com/aacebo/gq/middleware"
)

func Test_Logger(t *testing.T) {
	t.Run("should log resolvers", func(t *testing.T) {
		schema := gq.Object[map[string]string]{
			Name: "User",
			Use:  []gq.Middleware{middleware.Logger(nil)},
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
	})
}
