package peg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func debugf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	// fmt.Printf(format, args...)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 12, col: 1, offset: 176},
			expr: &actionExpr{
				pos: position{line: 12, col: 10, offset: 185},
				run: (*parser).callonInput1,
				expr: &seqExpr{
					pos: position{line: 12, col: 10, offset: 185},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 12, col: 10, offset: 185},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 12, col: 12, offset: 187},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 17, offset: 192},
								name: "Expr",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 12, col: 22, offset: 197},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 17, col: 1, offset: 267},
			expr: &actionExpr{
				pos: position{line: 17, col: 9, offset: 275},
				run: (*parser).callonExpr1,
				expr: &seqExpr{
					pos: position{line: 17, col: 9, offset: 275},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 17, col: 9, offset: 275},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 17, col: 14, offset: 280},
								name: "OrExpr",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 17, col: 21, offset: 287},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "OrExpr",
			pos:  position{line: 21, col: 1, offset: 312},
			expr: &choiceExpr{
				pos: position{line: 21, col: 11, offset: 322},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 21, col: 11, offset: 322},
						run: (*parser).callonOrExpr2,
						expr: &seqExpr{
							pos: position{line: 21, col: 11, offset: 322},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 21, col: 11, offset: 322},
									label: "one",
									expr: &ruleRefExpr{
										pos:  position{line: 21, col: 15, offset: 326},
										name: "AndExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 21, col: 23, offset: 334},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 21, col: 25, offset: 336},
									val:        "OR",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 21, col: 30, offset: 341},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 21, col: 32, offset: 343},
									label: "two",
									expr: &ruleRefExpr{
										pos:  position{line: 21, col: 36, offset: 347},
										name: "AndExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 21, col: 44, offset: 355},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 25, col: 5, offset: 463},
						run: (*parser).callonOrExpr12,
						expr: &seqExpr{
							pos: position{line: 25, col: 5, offset: 463},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 25, col: 5, offset: 463},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 25, col: 10, offset: 468},
										name: "AndExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 25, col: 18, offset: 476},
									name: "_",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "AndExpr",
			pos:  position{line: 30, col: 1, offset: 541},
			expr: &choiceExpr{
				pos: position{line: 30, col: 12, offset: 552},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 30, col: 12, offset: 552},
						run: (*parser).callonAndExpr2,
						expr: &seqExpr{
							pos: position{line: 30, col: 12, offset: 552},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 30, col: 12, offset: 552},
									label: "one",
									expr: &ruleRefExpr{
										pos:  position{line: 30, col: 16, offset: 556},
										name: "NotExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 30, col: 24, offset: 564},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 30, col: 26, offset: 566},
									val:        "AND",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 30, col: 32, offset: 572},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 30, col: 34, offset: 574},
									label: "two",
									expr: &ruleRefExpr{
										pos:  position{line: 30, col: 38, offset: 578},
										name: "NotExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 30, col: 46, offset: 586},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 34, col: 5, offset: 696},
						run: (*parser).callonAndExpr12,
						expr: &seqExpr{
							pos: position{line: 34, col: 5, offset: 696},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 34, col: 5, offset: 696},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 34, col: 10, offset: 701},
										name: "NotExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 34, col: 18, offset: 709},
									name: "_",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NotExpr",
			pos:  position{line: 39, col: 1, offset: 775},
			expr: &choiceExpr{
				pos: position{line: 39, col: 12, offset: 786},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 39, col: 12, offset: 786},
						run: (*parser).callonNotExpr2,
						expr: &seqExpr{
							pos: position{line: 39, col: 12, offset: 786},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 39, col: 12, offset: 786},
									val:        "NOT",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 39, col: 18, offset: 792},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 39, col: 20, offset: 794},
									label: "one",
									expr: &ruleRefExpr{
										pos:  position{line: 39, col: 24, offset: 798},
										name: "SimpleExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 39, col: 35, offset: 809},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 43, col: 5, offset: 907},
						run: (*parser).callonNotExpr9,
						expr: &labeledExpr{
							pos:   position{line: 43, col: 5, offset: 907},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 10, offset: 912},
								name: "SimpleExpr",
							},
						},
					},
				},
			},
		},
		{
			name: "SimpleExpr",
			pos:  position{line: 48, col: 1, offset: 987},
			expr: &choiceExpr{
				pos: position{line: 48, col: 15, offset: 1001},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 48, col: 15, offset: 1001},
						run: (*parser).callonSimpleExpr2,
						expr: &seqExpr{
							pos: position{line: 48, col: 15, offset: 1001},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 48, col: 15, offset: 1001},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 48, col: 19, offset: 1005},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 48, col: 21, offset: 1007},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 48, col: 26, offset: 1012},
										name: "Expr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 48, col: 31, offset: 1017},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 48, col: 33, offset: 1019},
									val:        ")",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 48, col: 37, offset: 1023},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 51, col: 5, offset: 1100},
						run: (*parser).callonSimpleExpr11,
						expr: &seqExpr{
							pos: position{line: 51, col: 5, offset: 1100},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 51, col: 5, offset: 1100},
									label: "lit",
									expr: &ruleRefExpr{
										pos:  position{line: 51, col: 9, offset: 1104},
										name: "Literal",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 51, col: 17, offset: 1112},
									name: "_",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 56, col: 1, offset: 1182},
			expr: &choiceExpr{
				pos: position{line: 56, col: 12, offset: 1193},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 56, col: 12, offset: 1193},
						run: (*parser).callonLiteral2,
						expr: &labeledExpr{
							pos:   position{line: 56, col: 12, offset: 1193},
							label: "text",
							expr: &ruleRefExpr{
								pos:  position{line: 56, col: 17, offset: 1198},
								name: "Text",
							},
						},
					},
					&actionExpr{
						pos: position{line: 59, col: 5, offset: 1304},
						run: (*parser).callonLiteral5,
						expr: &labeledExpr{
							pos:   position{line: 59, col: 5, offset: 1304},
							label: "quot",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 10, offset: 1309},
								name: "Quoted",
							},
						},
					},
					&actionExpr{
						pos: position{line: 63, col: 5, offset: 1444},
						run: (*parser).callonLiteral8,
						expr: &labeledExpr{
							pos:   position{line: 63, col: 5, offset: 1444},
							label: "num",
							expr: &ruleRefExpr{
								pos:  position{line: 63, col: 9, offset: 1448},
								name: "Number",
							},
						},
					},
				},
			},
		},
		{
			name: "Text",
			pos:  position{line: 68, col: 1, offset: 1563},
			expr: &actionExpr{
				pos: position{line: 68, col: 9, offset: 1571},
				run: (*parser).callonText1,
				expr: &seqExpr{
					pos: position{line: 68, col: 9, offset: 1571},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 68, col: 9, offset: 1571},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 68, col: 18, offset: 1580},
							expr: &charClassMatcher{
								pos:        position{line: 68, col: 18, offset: 1580},
								val:        "[a-zA-Z0-9_'-]",
								chars:      []rune{'_', '\'', '-'},
								ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Quoted",
			pos:  position{line: 72, col: 1, offset: 1629},
			expr: &actionExpr{
				pos: position{line: 72, col: 11, offset: 1639},
				run: (*parser).callonQuoted1,
				expr: &seqExpr{
					pos: position{line: 72, col: 11, offset: 1639},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 72, col: 11, offset: 1639},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 72, col: 15, offset: 1643},
							expr: &charClassMatcher{
								pos:        position{line: 72, col: 15, offset: 1643},
								val:        "[^\"]",
								chars:      []rune{'"'},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&litMatcher{
							pos:        position{line: 72, col: 21, offset: 1649},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Number",
			pos:  position{line: 76, col: 1, offset: 1686},
			expr: &choiceExpr{
				pos: position{line: 76, col: 11, offset: 1696},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 76, col: 11, offset: 1696},
						run: (*parser).callonNumber2,
						expr: &labeledExpr{
							pos:   position{line: 76, col: 11, offset: 1696},
							label: "hexint",
							expr: &ruleRefExpr{
								pos:  position{line: 76, col: 18, offset: 1703},
								name: "HexInteger",
							},
						},
					},
					&actionExpr{
						pos: position{line: 78, col: 5, offset: 1740},
						run: (*parser).callonNumber5,
						expr: &labeledExpr{
							pos:   position{line: 78, col: 5, offset: 1740},
							label: "integer",
							expr: &ruleRefExpr{
								pos:  position{line: 78, col: 13, offset: 1748},
								name: "Integer",
							},
						},
					},
				},
			},
		},
		{
			name: "HexInteger",
			pos:  position{line: 82, col: 1, offset: 1782},
			expr: &actionExpr{
				pos: position{line: 82, col: 15, offset: 1796},
				run: (*parser).callonHexInteger1,
				expr: &seqExpr{
					pos: position{line: 82, col: 15, offset: 1796},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 82, col: 15, offset: 1796},
							expr: &litMatcher{
								pos:        position{line: 82, col: 15, offset: 1796},
								val:        "-",
								ignoreCase: false,
							},
						},
						&litMatcher{
							pos:        position{line: 82, col: 20, offset: 1801},
							val:        "0",
							ignoreCase: false,
						},
						&charClassMatcher{
							pos:        position{line: 82, col: 24, offset: 1805},
							val:        "[xX]",
							chars:      []rune{'x', 'X'},
							ignoreCase: false,
							inverted:   false,
						},
						&oneOrMoreExpr{
							pos: position{line: 82, col: 29, offset: 1810},
							expr: &charClassMatcher{
								pos:        position{line: 82, col: 29, offset: 1810},
								val:        "[0-9a-fA-F]",
								ranges:     []rune{'0', '9', 'a', 'f', 'A', 'F'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 86, col: 1, offset: 1875},
			expr: &actionExpr{
				pos: position{line: 86, col: 12, offset: 1886},
				run: (*parser).callonInteger1,
				expr: &seqExpr{
					pos: position{line: 86, col: 12, offset: 1886},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 86, col: 12, offset: 1886},
							expr: &litMatcher{
								pos:        position{line: 86, col: 12, offset: 1886},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 86, col: 17, offset: 1891},
							expr: &charClassMatcher{
								pos:        position{line: 86, col: 17, offset: 1891},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 90, col: 1, offset: 1951},
			expr: &zeroOrMoreExpr{
				pos: position{line: 90, col: 19, offset: 1969},
				expr: &charClassMatcher{
					pos:        position{line: 90, col: 19, offset: 1969},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 92, col: 1, offset: 1977},
			expr: &notExpr{
				pos: position{line: 92, col: 8, offset: 1984},
				expr: &anyMatcher{
					line: 92, col: 9, offset: 1985,
				},
			},
		},
	},
}

func (c *current) onInput1(expr interface{}) (interface{}, error) {
	debugf("Input: returning expr %#v", expr)
	return expr, nil
}

func (p *parser) callonInput1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput1(stack["expr"])
}

func (c *current) onExpr1(expr interface{}) (interface{}, error) {
	return expr, nil
}

func (p *parser) callonExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpr1(stack["expr"])
}

func (c *current) onOrExpr2(one, two interface{}) (interface{}, error) {
	var n Node = newOrNode(one.(Node), two.(Node))
	debugf("OR: returning node %#v", n)
	return n, nil
}

func (p *parser) callonOrExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOrExpr2(stack["one"], stack["two"])
}

func (c *current) onOrExpr12(expr interface{}) (interface{}, error) {
	debugf("OR: returning expr %#v", expr)
	return expr, nil
}

func (p *parser) callonOrExpr12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOrExpr12(stack["expr"])
}

func (c *current) onAndExpr2(one, two interface{}) (interface{}, error) {
	var n Node = newAndNode(one.(Node), two.(Node))
	debugf("AND: returning node %#v", n)
	return n, nil
}

func (p *parser) callonAndExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAndExpr2(stack["one"], stack["two"])
}

func (c *current) onAndExpr12(expr interface{}) (interface{}, error) {
	debugf("AND: returning expr %#v", expr)
	return expr, nil
}

func (p *parser) callonAndExpr12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAndExpr12(stack["expr"])
}

func (c *current) onNotExpr2(one interface{}) (interface{}, error) {
	var n Node = newNotNode(one.(Node))
	debugf("NOT: returning node %#v", n)
	return n, nil
}

func (p *parser) callonNotExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNotExpr2(stack["one"])
}

func (c *current) onNotExpr9(expr interface{}) (interface{}, error) {
	debugf("NOT: returning expr %#v", expr)
	return expr, nil
}

func (p *parser) callonNotExpr9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNotExpr9(stack["expr"])
}

func (c *current) onSimpleExpr2(expr interface{}) (interface{}, error) {
	debugf("Simple: returning braced expr %#v", expr)
	return expr, nil
}

func (p *parser) callonSimpleExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleExpr2(stack["expr"])
}

func (c *current) onSimpleExpr11(lit interface{}) (interface{}, error) {
	debugf("Simple: returning literal %#v", lit)
	return lit, nil
}

func (p *parser) callonSimpleExpr11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleExpr11(stack["lit"])
}

func (c *current) onLiteral2(text interface{}) (interface{}, error) {
	debugf(`Literal: returning text "%s"`, text.(string))
	return newTextNode(text.(string)), nil
}

func (p *parser) callonLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteral2(stack["text"])
}

func (c *current) onLiteral5(quot interface{}) (interface{}, error) {
	// TODO: unquote
	debugf("Literal: returning quoted string `%s`", quot.(string))
	return newTextNode(quot.(string)), nil
}

func (p *parser) callonLiteral5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteral5(stack["quot"])
}

func (c *current) onLiteral8(num interface{}) (interface{}, error) {
	debugf("Literal: returning number %d", num.(int64))
	return newTextNode(fmt.Sprint(num.(int64))), nil
}

func (p *parser) callonLiteral8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteral8(stack["num"])
}

func (c *current) onText1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onText1()
}

func (c *current) onQuoted1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonQuoted1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuoted1()
}

func (c *current) onNumber2(hexint interface{}) (interface{}, error) {
	return hexint, nil
}

func (p *parser) callonNumber2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumber2(stack["hexint"])
}

func (c *current) onNumber5(integer interface{}) (interface{}, error) {
	return integer, nil
}

func (p *parser) callonNumber5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumber5(stack["integer"])
}

func (c *current) onHexInteger1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 0, 0)
}

func (p *parser) callonHexInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHexInteger1()
}

func (c *current) onInteger1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 0)
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
