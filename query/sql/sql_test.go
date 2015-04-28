package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrew-d/docstore/query/parse"
)

func TestNodesToSQL(t *testing.T) {
	testCases := []struct {
		Desc string
		SQL  string
		Args []interface{}
		Node parse.Node
	}{
		{`simple value`, "default = ?", []interface{}{"foo"}, &parse.TextNode{
			Text: "foo",
		}},
		{`value with field`, "myfield = ?", []interface{}{"foo"}, &parse.TextNode{
			Field: "myfield",
			Text:  "foo",
		}},
		{`simple OR`, "(default = ?) OR (default = ?)", []interface{}{"foo", "bar"}, &parse.OrNode{
			Left:  &parse.TextNode{Text: "foo"},
			Right: &parse.TextNode{Text: "bar"},
		}},
		{`simple AND`, "(default = ?) AND (default = ?)", []interface{}{"foo", "bar"}, &parse.AndNode{
			Left:  &parse.TextNode{Text: "foo"},
			Right: &parse.TextNode{Text: "bar"},
		}},
		{`simple NOT`, "NOT (default = ?)", []interface{}{"foo"}, &parse.NotNode{
			Node: &parse.TextNode{Text: "foo"},
		}},
	}

	for _, tcase := range testCases {
		sql, args := NodesToSQL(tcase.Node, `default`)
		assert.Equal(t, tcase.SQL, sql)
		assert.Equal(t, tcase.Args, args)
	}
}
