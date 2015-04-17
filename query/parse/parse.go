package parse

import (
	"fmt"
	"runtime"

	"github.com/andrew-d/docstore/query/tokenize"
)

const kLOOKAHEAD = 1

type Parser struct {
	lex       *tokenize.Scanner
	lahead    [kLOOKAHEAD]tokenize.Token
	peekCount int
}

// recover is the handler that turns panics into returns from the top level
// of Parse.
func (p *Parser) recover(errptr *error) {
	err := recover()
	if err != nil {
		if _, ok := err.(runtime.Error); ok {
			panic(err)
		}

		*errptr = err.(error)
	}
	return
}

// errorf halts parsing with the given error
func (p *Parser) errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

// unexpected complains about the unexpected token and halts parsing
func (p *Parser) unexpected(token tokenize.Token, context string) {
	p.errorf("unexpected %s in %s", token, context)
}

// next returns the next token
func (p *Parser) next() tokenize.Token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.lahead[0] = p.lex.NextToken()
	}
	return p.lahead[p.peekCount]
}

// backup backs the input stream up one token
func (p *Parser) backup() {
	p.peekCount++
}

// nextNonSpace returns the next non-space token
func (p *Parser) nextNonSpace() (token tokenize.Token) {
	for {
		token = p.next()
		if token.Type != tokenize.TokenSpace {
			break
		}
	}

	return token
}

// expect consumes the next token and guarantees that it had the required
// type, halting parsing otherwise
func (p *Parser) expect(expected tokenize.TokenType, context string) tokenize.Token {
	token := p.nextNonSpace()
	if token.Type != expected {
		p.unexpected(token, context)
	}
	return token
}

// peek returns but does not consume the next token
func (p *Parser) peek() tokenize.Token {
	tok := p.next()
	p.backup()
	return tok
}

// peekNonSpace returns but does not consume the next non-space token
func (p *Parser) peekNonSpace() tokenize.Token {
	tok := p.nextNonSpace()
	p.backup()
	return tok
}

// Parse is the function that actually performs parsing.  It takes an input
// string and will either return a set of nodes representing the parsed query,
// or an error.
func Parse(text string) (n Node, err error) {
	p := &Parser{
		lex: tokenize.NewScanner(text),
	}
	defer p.recover(&err)

	// The top-level item is an expression
	n = p.parseExpr()
	return n, nil
}

// parseExpr parses and returns an expression
func (p *Parser) parseExpr() Node {
	for p.peek().Type != tokenize.TokenEOF {
		switch token := p.nextNonSpace(); token.Type {
		case tokenize.TokenLeftParen:
			// Recursively parse this expression
			return p.parseExpr()

		case tokenize.TokenRightParen:
			// We must be nested - we're done with this expression.
			return nil // TODO: real return

		case tokenize.TokenText:
			// Consume this

		case tokenize.TokenError:
			p.errorf("%s", token.Val)

		default:
			panic("Unhandled token")
		}
	}

	return nil // TODO
}

func (p *Parser) parseOperator() Node {
	switch token := p.nextNonSpace(); token.Type {
	case tokenize.TokenAnd:
		return p.newAnd(-1)
	case tokenize.TokenOr:
		return p.newOr(-1)
	default:
		p.errorf("expected operator, found %s", token.Type)
	}
	panic("unreachable")
}
