package parse

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

func (t NodeType) String() string {
	switch t {
	case NodeInvalid:
		return "<INVALID>"
	case NodeAnd:
		return "<AND>"
	case NodeOr:
		return "<OR>"
	case NodeNot:
		return "<NOT>"
	case NodeText:
		return "<TEXT>"
	default:
		return "<UNKNOWN>"
	}
}

const (
	NodeInvalid NodeType = iota
	NodeText
	NodeAnd
	NodeOr
	NodeNot
)

type Node interface {
	Type() NodeType
}

type TextNode struct {
	NodeType

	Field string
	Text  string
}

func newTextNode(field, txt string) *TextNode {
	return &TextNode{Field: field, Text: txt, NodeType: NodeText}
}

type AndNode struct {
	NodeType

	Left  Node
	Right Node
}

func newAndNode(left, right Node) Node {
	return &AndNode{Left: left, Right: right, NodeType: NodeAnd}
}

type OrNode struct {
	NodeType

	Left  Node
	Right Node
}

func newOrNode(left, right Node) Node {
	return &OrNode{Left: left, Right: right, NodeType: NodeOr}
}

type NotNode struct {
	NodeType

	Node Node
}

func newNotNode(n Node) Node {
	return &NotNode{Node: n, NodeType: NodeNot}
}
