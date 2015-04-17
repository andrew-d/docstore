package tokenize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func allTokens(str string) []Token {
	s := NewScanner(str)
	toks := []Token{}
	for {
		curr := s.NextToken()
		curr.Pos = 0 // TODO: do this?
		toks = append(toks, curr)
		if curr.Type == TokenEOF || curr.Type == TokenError {
			break
		}
	}

	return toks
}

var (
	tEOF    = Token{TokenEOF, 0, ""}
	tSpace  = Token{TokenSpace, 0, " "}
	tLParen = Token{TokenLeftParen, 0, "("}
	tRParen = Token{TokenRightParen, 0, ")"}
	tAND    = Token{TokenAnd, 0, "AND"}
	tOR     = Token{TokenOr, 0, "OR"}
)

func TestTokenize(t *testing.T) {
	tcases := []struct {
		Name     string
		Input    string
		Expected []Token
	}{
		{`spaces`, " \t ", []Token{
			{TokenSpace, 0, " \t "},
			tEOF,
		}},
		{`text`, "foo bar baz", []Token{
			{TokenText, 0, "foo"},
			tSpace,
			{TokenText, 0, "bar"},
			tSpace,
			{TokenText, 0, "baz"},
			tEOF,
		}},
		{`numbers`, "1 02 0x1e -1.23 -0xf1", []Token{
			{TokenNumber, 0, "1"},
			tSpace,
			{TokenNumber, 0, "02"},
			tSpace,
			{TokenNumber, 0, "0x1e"},
			tSpace,
			{TokenNumber, 0, "-1.23"},
			tSpace,
			{TokenNumber, 0, "-0xf1"},
			tEOF,
		}},
		{`quoted string`, `"foo bar baz"`, []Token{
			{TokenString, 0, `"foo bar baz"`},
			tEOF,
		}},
		{`quoted string with escape`, `"foo \" bar"`, []Token{
			{TokenString, 0, `"foo \" bar"`},
			tEOF,
		}},
		{`identifiers`, `foo AND bar OR baz`, []Token{
			{TokenText, 0, "foo"},
			tSpace,
			tAND,
			tSpace,
			{TokenText, 0, "bar"},
			tSpace,
			tOR,
			tSpace,
			{TokenText, 0, "baz"},
			tEOF,
		}},
		{`braces`, `(((123)))`, []Token{
			tLParen,
			tLParen,
			tLParen,
			{TokenNumber, 0, "123"},
			tRParen,
			tRParen,
			tRParen,
			tEOF,
		}},
		{`complex query`, `(one AND two) OR ("three four" five OR six)`, []Token{
			tLParen,
			{TokenText, 0, "one"},
			tSpace,
			tAND,
			tSpace,
			{TokenText, 0, "two"},
			tRParen,
			tSpace,
			tOR,
			tSpace,
			tLParen,
			{TokenString, 0, `"three four"`},
			tSpace,
			{TokenText, 0, "five"},
			tSpace,
			tOR,
			tSpace,
			{TokenText, 0, "six"},
			tRParen,
			tEOF,
		}},

		// errors

		{`invalid identifier terminater`, `foo&`, []Token{
			{TokenError, 0, `bad character U+0026 '&'`},
		}},
		{`unclosed quote`, `"foo`, []Token{
			{TokenError, 0, "unterminated quoted string"},
		}},
		{`unclosed quote with escape`, `"foo \"`, []Token{
			{TokenError, 0, "unterminated quoted string"},
		}},
		{`unclosed quote with escape and eof`, `"foo \`, []Token{
			{TokenError, 0, "unterminated quoted string"},
		}},
		{`unclosed paren`, `(foo`, []Token{
			tLParen,
			{TokenText, 0, "foo"},
			{TokenError, 0, "unclosed left paren"},
		}},
		{`extra right paren`, `(foo))`, []Token{
			tLParen,
			{TokenText, 0, "foo"},
			tRParen,
			tRParen,
			{TokenError, 0, "unexpected right paren U+0029 ')'"},
		}},
		{`bad number`, `9q`, []Token{
			{TokenError, 0, `bad number syntax: "9q"`},
		}},
		{`bad character`, "\x01", []Token{
			{TokenError, 0, `unrecognized character in query: U+0001`},
		}},
	}

	t.Logf("Running %d test cases...", len(tcases))
	for i, tcase := range tcases {
		toks := allTokens(tcase.Input)
		assert.Equal(t, tcase.Expected, toks,
			"case '%s' (%d): did not match", tcase.Name, i)
	}
	t.Logf("Ran %d test cases", len(tcases))
}
