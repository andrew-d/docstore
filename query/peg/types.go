package peg

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

	Text string
}

func newTextNode(t string) *TextNode {
	return &TextNode{Text: t, NodeType: NodeText}
}

type AndNode struct {
	NodeType

	Left  Node
	Right Node
}

func newAndNode(left, right Node) *AndNode {
	return &AndNode{Left: left, Right: right, NodeType: NodeAnd}
}

type OrNode struct {
	NodeType

	Left  Node
	Right Node
}

func newOrNode(left, right Node) *OrNode {
	return &OrNode{Left: left, Right: right, NodeType: NodeOr}
}

type NotNode struct {
	NodeType

	Node Node
}

func newNotNode(n Node) *NotNode {
	return &NotNode{Node: n, NodeType: NodeNot}
}
