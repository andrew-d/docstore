package peg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tcases := []struct {
		Name     string
		Input    string
		Expected Node
	}{
		{`number`, `1234`, &TextNode{NodeType: NodeText, Text: `1234`}},
		{`hex number`, `0x123e`, &TextNode{NodeType: NodeText, Text: fmt.Sprint(0x123e)}},
		{`negative number`, `-5678`, &TextNode{NodeType: NodeText, Text: `-5678`}},
		{`negative hex number`, `-0xbcd1`, &TextNode{NodeType: NodeText, Text: fmt.Sprint(-0xbcd1)}},
		{`literal`, `asdf`, &TextNode{NodeType: NodeText, Text: `asdf`}},
		{`quoted string`, `"foo bar baz"`, &TextNode{NodeType: NodeText, Text: `"foo bar baz"`}},

		{`spaces around number`, ` 1234 `, &TextNode{NodeType: NodeText, Text: `1234`}},
		{`spaces around literal`, ` asdf `, &TextNode{NodeType: NodeText, Text: `asdf`}},

		{`simple AND`, `one AND two`, &AndNode{
			NodeType: NodeAnd,
			Left:     &TextNode{NodeType: NodeText, Text: `one`},
			Right:    &TextNode{NodeType: NodeText, Text: `two`},
		}},
		{`simple OR`, `one OR two`, &OrNode{
			NodeType: NodeOr,
			Left:     &TextNode{NodeType: NodeText, Text: `one`},
			Right:    &TextNode{NodeType: NodeText, Text: `two`},
		}},
		{`simple NOT`, `NOT one`, &NotNode{
			NodeType: NodeNot,
			Node:     &TextNode{NodeType: NodeText, Text: `one`},
		}},
		{`double NOT`, `NOT NOT one`, &NotNode{
			NodeType: NodeNot,
			Node: &NotNode{
				NodeType: NodeNot,
				Node:     &TextNode{NodeType: NodeText, Text: `one`},
			},
		}},
		{`AND followed by NOT`, `one AND NOT two`, &AndNode{
			NodeType: NodeAnd,
			Left:     &TextNode{NodeType: NodeText, Text: `one`},
			Right: &NotNode{
				NodeType: NodeNot,
				Node:     &TextNode{NodeType: NodeText, Text: `two`},
			},
		}},
		{`OR followed by NOT`, `one OR NOT two`, &OrNode{
			NodeType: NodeOr,
			Left:     &TextNode{NodeType: NodeText, Text: `one`},
			Right: &NotNode{
				NodeType: NodeNot,
				Node:     &TextNode{NodeType: NodeText, Text: `two`},
			},
		}},
		/*
			{`chained AND`, `one AND two AND three`, &AndNode{
				NodeType: NodeAnd,
				Left: &AndNode{
					NodeType: NodeAnd,
					Left:     &TextNode{NodeType: NodeText, Text: `one`},
					Right:    &TextNode{NodeType: NodeText, Text: `two`},
				},
				Right: &TextNode{NodeType: NodeText, Text: `two`},
			}},
		*/
		{`AND takes precedence over OR`, `one AND two OR three`, &OrNode{
			NodeType: NodeOr,
			Left: &AndNode{
				NodeType: NodeAnd,
				Left:     &TextNode{NodeType: NodeText, Text: `one`},
				Right:    &TextNode{NodeType: NodeText, Text: `two`},
			},
			Right: &TextNode{NodeType: NodeText, Text: `three`},
		}},
		{`braces take precedence`, `one AND (two OR three)`, &AndNode{
			NodeType: NodeAnd,
			Left:     &TextNode{NodeType: NodeText, Text: `one`},
			Right: &OrNode{
				NodeType: NodeOr,
				Left:     &TextNode{NodeType: NodeText, Text: `two`},
				Right:    &TextNode{NodeType: NodeText, Text: `three`},
			},
		}},
	}

	Memoize(true)

	t.Logf("Running %d test cases...", len(tcases))
	for i, tcase := range tcases {
		output, err := Parse(tcase.Name, []byte(tcase.Input))
		assert.NoError(t, err)
		assert.Equal(t, tcase.Expected, output,
			"case '%s' (%d): did not match", tcase.Name, i)
	}
	t.Logf("Ran %d test cases", len(tcases))
}
