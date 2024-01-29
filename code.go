package yutex

func parseMixCode(text string) []map[string]string {
	var ok bool
	var codes []map[string]string
	node := NewNode("tmp", nil, nil)
	for text != "" {
		if text, ok = processCodeBlock(node, text); ok {
			continue
		}
		text = processUnknown(text)
	}
	for _, child := range node.Children {
		codes = append(codes, child.Value.(map[string]string))
	}
	return codes
}
