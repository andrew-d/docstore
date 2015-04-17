package parse

import (
	"github.com/andrew-d/docstore/query/tokenize"
)

type Node interface {
	Type() NodeType
	Position() Pos
}

type Pos tokenize.Pos

func (p Pos) Position() Pos {
	return p
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeText    NodeType = iota // A single chunk of text
	NodeSubExpr                 // A parenthesized subexpression
	NodeAnd                     // An AND node
	NodeOr                      // An OR node
)

type TextNode struct {
	NodeType
	Pos

	Text string
}

func (p *Parser) newText(pos Pos, text string) *TextNode {
	return &TextNode{Text: text, NodeType: NodeText, Pos: pos}
}

type AndNode struct {
	NodeType
	Pos

	Left  Node
	Right Node
}

func (p *Parser) newAnd(pos Pos) *AndNode {
	return &AndNode{NodeType: NodeAnd, Pos: pos}
}

type OrNode struct {
	NodeType
	Pos

	Left  Node
	Right Node
}

func (p *Parser) newOr(pos Pos) *OrNode {
	return &OrNode{NodeType: NodeOr, Pos: pos}
}
