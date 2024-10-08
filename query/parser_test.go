package query_test

import (
	"testing"

	"github.com/aacebo/gq/query"
)

func Test_Parser(t *testing.T) {
	t.Run("should succeed with fields", func(t *testing.T) {
		q, err := query.Parser([]byte(`{
			id,
			name,
			created_at
		}`)).Parse()

		if err != nil {
			t.Fatal(err)
		}

		t.Log(q)

		if len(q.Args) != 0 {
			t.Fatalf("should have 0 args")
		}

		if len(q.Fields) != 3 {
			t.Fatalf("should have 3 fields")
		}
	})

	t.Run("should succeed with args", func(t *testing.T) {
		q, err := query.Parser([]byte(`{
			users(page: 1, pageSize: 10)
		}`)).Parse()

		if err != nil {
			t.Fatal(err)
		}

		t.Log(q)

		if len(q.Args) != 0 {
			t.Fatalf("should have 0 args")
		}

		if len(q.Fields) != 1 {
			t.Fatalf("should have 1 field")
		}

		q, exists := q.Fields["users"]

		if !exists {
			t.Fatalf("should have users field")
		}

		if len(q.Args) != 2 {
			t.Fatalf("should have 2 args for users field")
		}

		if q.Args.Get("page") != 1 {
			t.Fatalf("page arg should be 1")
		}

		if q.Args.Get("pageSize") != 10 {
			t.Fatalf("pageSize arg should be 10")
		}
	})

	t.Run("should succeed with subquery", func(t *testing.T) {
		q, err := query.Parser([]byte(`{
			users {
				id,
				name,
				created_at
			}
		}`)).Parse()

		if err != nil {
			t.Fatal(err)
		}

		t.Log(q)

		if len(q.Args) != 0 {
			t.Fatalf("should have 0 args")
		}

		if len(q.Fields) != 1 {
			t.Fatalf("should have 1 field")
		}

		q, exists := q.Fields["users"]

		if !exists {
			t.Fatalf("should have users field")
		}

		if len(q.Args) != 0 {
			t.Fatalf("should have 0 args for users field")
		}

		if len(q.Fields) != 3 {
			t.Fatalf("should have 3 fields for users field")
		}
	})

	t.Run("should succeed with args and subquery", func(t *testing.T) {
		q, err := query.Parser([]byte(`{
			users(search: "some text", deleted: true) {
				id,
				name,
				created_at
			}
		}`)).Parse()

		if err != nil {
			t.Fatal(err)
		}

		t.Log(q)

		if len(q.Args) != 0 {
			t.Fatalf("should have 0 args")
		}

		if len(q.Fields) != 1 {
			t.Fatalf("should have 1 field")
		}

		q, exists := q.Fields["users"]

		if !exists {
			t.Fatalf("should have users field")
		}

		if len(q.Args) != 2 {
			t.Fatalf("should have 2 args for users field")
		}

		if len(q.Fields) != 3 {
			t.Fatalf("should have 3 fields for users field")
		}

		if q.Args.Get("search") != "some text" {
			t.Fatalf("page arg should be 'some text'")
		}

		if q.Args.Get("deleted") != true {
			t.Fatalf("page arg should be true")
		}
	})
}
