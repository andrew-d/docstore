package tokenize

import (
	"fmt"
)

// Pos represents the position within the input string
type Pos int

func (p Pos) Position() Pos {
	return p
}

// Token represents a single token or text string returned from the scanner
type Token struct {
	Type TokenType // Type of this token
	Pos  Pos       // Starting position, in bytes, of this token in the input
	Val  string    // Value of this token
}

func (t Token) String() string {
	switch {
	case t.Type == TokenEOF:
		return "EOF"
	case t.Type == TokenError:
		return t.Val
	case t.Type > TokenKeyword:
		return fmt.Sprintf("<%s>", t.Val)
	case len(t.Val) > 10:
		return fmt.Sprintf("%.10q...", t.Val)
	}
	return fmt.Sprintf("%q", t.Val)
}

// TokenType identifies the type of a single Token
type TokenType int

const (
	TokenError      TokenType = iota // An error occurred; value is text of error
	TokenEOF                         // EOF
	TokenLeftParen                   // '('
	TokenRightParen                  // ')'
	TokenNumber                      // A number
	TokenSpace                       // Run of whitespace
	TokenString                      // Quoted string (including quotes)
	TokenText                        // Plain text
	// Keywords appear after the rest
	TokenKeyword // Delimiter only
	TokenAnd     // 'and' keyword
	TokenOr      // 'or' keyword
)
