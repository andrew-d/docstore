package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrew-d/docstore/query/tokenize"
)

func tokensEqual(t tokenize.Token, ty tokenize.TokenType, val string) bool {
	return t.Type == ty && t.Val == val
}

func TestParserFetch(t *testing.T) {
	p := &Parser{
		lex: tokenize.NewScanner("one two three four five"),
	}

	assert.True(t, tokensEqual(p.next(), tokenize.TokenText, "one"))
	assert.True(t, tokensEqual(p.next(), tokenize.TokenSpace, " "))
	assert.True(t, tokensEqual(p.next(), tokenize.TokenText, "two"))

	p.backup()

	assert.True(t, tokensEqual(p.next(), tokenize.TokenText, "two"))
	assert.True(t, tokensEqual(p.next(), tokenize.TokenSpace, " "))
	assert.True(t, tokensEqual(p.next(), tokenize.TokenText, "three"))
}
