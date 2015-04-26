package parse

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Println

func TestParse(t *testing.T) {
	type testSpec struct {
		Name     string
		Input    string
		Expected Node
	}

	tcases := []testSpec{
		{`number`, `1234`, &TextNode{NodeType: NodeText, Text: `1234`}},
		{`hex number`, `0x123e`, &TextNode{NodeType: NodeText, Text: fmt.Sprint(0x123e)}},
		{`negative number`, `-5678`, &TextNode{NodeType: NodeText, Text: `-5678`}},
		{`negative hex number`, `-0xbcd1`, &TextNode{NodeType: NodeText, Text: fmt.Sprint(-0xbcd1)}},
		{`literal`, `asdf`, &TextNode{NodeType: NodeText, Text: `asdf`}},
		{`quoted string`, `"foo bar baz"`, &TextNode{NodeType: NodeText, Text: `foo bar baz`}},
		{`quoted string with escape`, `"foo\u0020bar"`, &TextNode{NodeType: NodeText, Text: `foo bar`}},

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
		{`braces take precedence`, `one AND (two OR three)`, &AndNode{
			NodeType: NodeAnd,
			Left:     &TextNode{NodeType: NodeText, Text: `one`},
			Right: &OrNode{
				NodeType: NodeOr,
				Left:     &TextNode{NodeType: NodeText, Text: `two`},
				Right:    &TextNode{NodeType: NodeText, Text: `three`},
			},
		}},
		{`chained AND`, `one AND two AND three`, &AndNode{
			NodeType: NodeAnd,
			Left:     &TextNode{NodeType: NodeText, Text: `one`},
			Right: &AndNode{
				NodeType: NodeAnd,
				Left:     &TextNode{NodeType: NodeText, Text: `two`},
				Right:    &TextNode{NodeType: NodeText, Text: `three`},
			},
		}},
	}

	Memoize(true)
	spew.Config.DisableMethods = true

	runTest := func(i int, tcase testSpec) {
		output, err := Parse(tcase.Name, []byte(tcase.Input))
		if assert.NoError(t, err, "%d: case '%s' should have matched", i, tcase.Name) {
			if !assert.Equal(t, tcase.Expected, output,
				"%d: case '%s' did not equal expected output", i, tcase.Name) {
				t.Logf("Dumping actual output")
				spew.Dump(output)
			}
		}
	}

	// Check for 'only' rules
	only := ``
	for i, tcase := range tcases {
		if tcase.Name == only {
			t.Logf("%d: Only running test case '%s'", i, tcase.Name)
			runTest(i, tcase)
			return
		}
	}

	t.Logf("Running %d test cases...", len(tcases))
	for i, tcase := range tcases {
		runTest(i, tcase)
	}
	t.Logf("Ran %d test cases", len(tcases))
}
