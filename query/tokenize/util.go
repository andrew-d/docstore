package tokenize

import (
	"unicode"
)

var eof = rune(0)

// Map of literal strings to keywords
var keywords = map[string]TokenType{
	"AND": TokenAnd,
	"OR":  TokenOr,
}

// Helper function that returns whether something could be a character in
// an identifier.  Currently, this is letters, digits, the single quote,
// underscores and dashes ['_-]
func isIdentChar(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '\'' || ch == '_'
}
