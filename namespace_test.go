package gq_test

import (
	"testing"

	"github.com/aacebo/gq"
)

func Test_Namespace(t *testing.T) {
	t.Run("should resolve ref", func(t *testing.T) {
		ns := gq.New()
		ns.Register(gq.Object[map[string]string]{
			Name: "User",
			Fields: gq.Fields{
				"id":   gq.Field{},
				"name": gq.Field{},
				"email": gq.Field{
					Resolver: func(params *gq.ResolveParams) gq.Result {
						return gq.Result{Data: "dev@gmail.com"}
					},
				},
			},
		})

		ns.Register(gq.Object[map[string]any]{
			Name: "Org",
			Fields: gq.Fields{
				"id":   gq.Field{},
				"name": gq.Field{},
				"created_by": gq.Field{
					Type: ns.Ref("User"),
					Resolver: func(params *gq.ResolveParams) gq.Result {
						return gq.Result{
							Data: map[string]string{
								"id":   "1",
								"name": "dev",
							},
						}
					},
				},
			},
		})

		res := ns.Do("Org", &gq.DoParams{
			Query: "{id,name,created_by{id,name,email}}",
			Value: map[string]any{
				"id":   "2",
				"name": "dev-org",
			},
		})

		if res.Error != nil {
			t.Fatal(res.Error)
		}
	})
}
