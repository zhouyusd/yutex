package yutex

type Node struct {
	// The token type
	Type string `json:"type"`
	// The token value
	Value interface{} `json:"value,omitempty"`
	// The token style
	Style interface{} `json:"style,omitempty"`
	// The token children
	Children []*Node `json:"children,omitempty"`
}

func NewNode(tp string, val, style interface{}) *Node {
	return &Node{tp, val, style, nil}
}
