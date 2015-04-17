package tokenize

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type stateFn func(*Scanner) stateFn

// Scanner represents a lexical scanner.
type Scanner struct {
	input      string     // The string being scanned
	state      stateFn    // Next lexing function to enter
	pos        Pos        // Current position in the input
	start      Pos        // Start position of this token
	width      Pos        // Width of the last rune read from the input
	lastPos    Pos        // Position of the most recent token returned by NextToken
	tokens     chan Token // Channel of scanned tokens
	parenDepth int        // Nesting depth of '()' expressions
}

// next returns the next rune in the input
func (s *Scanner) next() rune {
	if int(s.pos) >= len(s.input) {
		s.width = 0
		return eof
	}

	r, width := utf8.DecodeRuneInString(s.input[s.pos:])
	s.width = Pos(width)
	s.pos += s.width
	return r
}

// backup steps back one rune.  Can only be safely called once per call of next.
func (s *Scanner) backup() {
	s.pos -= s.width
}

// peek returns but does not consume the next rune in the input.
func (s *Scanner) peek() rune {
	r := s.next()
	s.backup()
	return r
}

// emit passes an token back to the user of this Scanner.
func (s *Scanner) emit(ty TokenType) {
	s.tokens <- Token{
		Type: ty,
		Pos:  s.start,
		Val:  s.input[s.start:s.pos],
	}
	s.start = s.pos
}

// ignore skips over the pending input before this position.
func (s *Scanner) ignore() {
	s.start = s.pos
}

// accept consumes the next rune if it's contained in the given set.
func (s *Scanner) accept(valid string) bool {
	if strings.IndexRune(valid, s.next()) >= 0 {
		return true
	}
	s.backup()
	return false
}

// acceptRun consumes a run (sequence) of runes from the given set.
func (s *Scanner) acceptRun(valid string) {
	for strings.IndexRune(valid, s.next()) >= 0 {
	}
	s.backup()
}

// errorf returns an error token and terminates the scan.
func (s *Scanner) errorf(format string, args ...interface{}) stateFn {
	s.tokens <- Token{
		Type: TokenError,
		Pos:  s.start,
		Val:  fmt.Sprintf(format, args...),
	}
	return nil
}

// NextToken returns the next token from the input
func (s *Scanner) NextToken() Token {
	tok := <-s.tokens
	s.lastPos = tok.Pos
	return tok
}

// NewScanner returns a new instance of Scanner that reads from the given
// string.
func NewScanner(input string) *Scanner {
	s := &Scanner{
		input:  input,
		tokens: make(chan Token),
	}
	go s.run()
	return s
}
