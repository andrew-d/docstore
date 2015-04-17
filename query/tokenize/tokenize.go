package tokenize

import (
	"unicode"
)

func (s *Scanner) run() {
	for s.state = lexExpr; s.state != nil; {
		s.state = s.state(s)
	}
}

func lexExpr(s *Scanner) stateFn {
	// We either have a number, a quoted string, a bracketed expression,
	// or an identifier.
	switch r := s.next(); {
	case r == eof:
		if s.parenDepth != 0 {
			return s.errorf("unclosed left paren")
		}

		s.emit(TokenEOF)
		return nil

	case unicode.IsSpace(r):
		return lexSpace

	case isIdentChar(r):
		return lexIdentifier

	case r == '"':
		return lexQuote

	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		s.backup()
		return lexNumber

	case r == '(':
		s.emit(TokenLeftParen)
		s.parenDepth++
		return lexExpr

	case r == ')':
		s.emit(TokenRightParen)
		s.parenDepth--
		if s.parenDepth < 0 {
			return s.errorf("unexpected right paren %#U", r)
		}
		return lexExpr

	default:
		return s.errorf("unrecognized character in query: %#U", r)
	}

	panic("unreachable")
}

func lexSpace(s *Scanner) stateFn {
	for unicode.IsSpace(s.peek()) {
		s.next()
	}
	s.emit(TokenSpace)
	return lexExpr
}

func lexIdentifier(s *Scanner) stateFn {
Loop:
	for {
		switch r := s.next(); {
		case isIdentChar(r):
			// Saved into current input

		default:
			s.backup()
			word := s.input[s.start:s.pos]

			// Only certain chars are allowed to terminate an identifier
			if !s.atTerminator() {
				return s.errorf("bad character %#U", r)
			}

			switch {
			case keywords[word] > TokenKeyword:
				s.emit(keywords[word])

			default:
				s.emit(TokenText)
			}

			break Loop
		}
	}

	return lexExpr
}

func (s *Scanner) atTerminator() bool {
	ch := s.peek()

	if unicode.IsSpace(ch) {
		return true
	}

	switch ch {
	case eof, '.', ')', '(':
		return true
	}

	return false
}

func lexNumber(s *Scanner) stateFn {
	if !s.scanNumber() {
		return s.errorf("bad number syntax: %q", s.input[s.start:s.pos])
	}

	s.emit(TokenNumber)
	return lexExpr
}

func (s *Scanner) scanNumber() bool {
	// Optional leading sign
	s.accept("+-")

	// Set proper character set if it's hex
	digits := "0123456789"
	if s.accept("0") && s.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}

	// Accept some digits, optional decimal, then more
	s.acceptRun(digits)
	if s.accept(".") {
		s.acceptRun(digits)
	}

	// The next thing should NOT be non-digit chars in our identifier set.
	ch := s.peek()
	if unicode.IsLetter(ch) || ch == '_' || ch == '\'' {
		s.next()
		return false
	}

	return true
}

func lexQuote(s *Scanner) stateFn {
Loop:
	for {
		switch s.next() {
		case '\\':
			if r := s.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return s.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}

	s.emit(TokenString)
	return lexExpr
}
