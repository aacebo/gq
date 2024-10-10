<p align="center">
	<img src="./assets/icon.png" width="120px" style="border-radius:20%" />
</p>
 
<p align="center">
	a zero dependency performant graph query resolver
</p>

<p align="center">
	<a href="https://opensource.org/licenses/MIT" target="_blank" alt="License">
		<img src="https://img.shields.io/badge/License-MIT-blue.svg" />
	</a>
	<a href="https://pkg.go.dev/github.com/aacebo/gq" target="_blank" alt="Go Reference">
		<img src="https://pkg.go.dev/badge/github.com/aacebo/gq.svg" />
	</a>
	<a href="https://goreportcard.com/report/github.com/aacebo/gq" target="_blank" alt="Go Report Card">
		<img src="https://goreportcard.com/badge/github.com/aacebo/gq" />
	</a>
	<a href="https://github.com/aacebo/owl/actions/workflows/ci.yml" target="_blank" alt="Build">
		<img src="https://github.com/aacebo/owl/actions/workflows/ci.yml/badge.svg?branch=main" />
	</a>
	<a href="https://codecov.io/gh/aacebo/gq" > 
		<img src="https://codecov.io/gh/aacebo/gq/graph/badge.svg?token=9XETRUUQUY" /> 
	</a>
</p>

# Install

```bash
go get github.com/aacebo/gq
```

# Usage

```go
schema := gq.Object[User]{
	Name:        "User",
	Description: "...",
	Fields: gq.Fields{
		"id":           gq.Field{},
		"name":         gq.Field{},
		"email":        gq.Field{},
		"created_at":   gq.Field{},
		"updated_at":   gq.Field{},
		"posts": gq.Field{
			Type: gq.List{
				Type: gq.Object[Post]{
					Name: "Post",
					Fields: gq.Fields{
						"id": 			gq.Field{},
						"body":			gq.Field{},
						"created_at":   gq.Field{},
						"updated_at":   gq.Field{},
					}
				}
			},
			Resolver: func(params gq.ResolveParams) (any, error) {
				user := params.Parent.(User)
				posts := // ... get some posts ...
				return posts, nil
			},
		},
	},
}

q := `{
	id,
	name,
	email,
	created_at,
	updated_at,
	posts {id,body}
}`

res := schema.Do(&gq.DoParams{
	Query: q,
	Value: User{
		ID: "1",
		Name: "test",
		Email: "test@test.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
})

if res.Error != nil {
	panic(res.Error)
}
```

# Features

| Name			             	| Status			   |
|-------------------------------|----------------------|
| Object Schema				 	| ✅				  	  |
| List Schema				 	| ✅				  	  |
| Field Schema				 	| ✅				  	  |
| Field Arguments + Validation	| ✅				  	  |
| String						| ⌛					  |
| Date							| ⌛					  |
| Int							| ⌛					  |
| Float							| ⌛					  |
| Bool							| ⌛					  |
| Middleware				 	| ✅				  	  |
| MetaData					 	| ✅				  	  |

# Related

- [Schema Validation](https://github.com/aacebo/owl)
