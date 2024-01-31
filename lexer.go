package yutex

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	sectionRegexp    = regexp.MustCompile(`^\\section\{([^{}]*)}.*`)
	paragraphRegexp  = regexp.MustCompile(`^\\begin\{paragraph}\{([^{}]*)}.*([\s\S]*?)\n\\end\{paragraph}.*`)
	blockquoteRegexp = regexp.MustCompile(`^\\begin\{blockquote}.*([\s\S]*?)\n\\end\{blockquote}.*`)
	codeBlockRegexp  = regexp.MustCompile(`^\\begin\{code}\{([^{}]*)}.*([\s\S]*?)\n\\end\{code}\s*`)
	mathBlockRegexp  = regexp.MustCompile(`^\\begin\{math}.*([\s\S]*?)\n\\end\{math}.*`)
	htmlBlockRegexp  = regexp.MustCompile(`^\\begin\{html}.*([\s\S]*?)\n\\end\{html}\s*`)
	tableRegexp      = regexp.MustCompile(`^\\begin\{table}.*([\s\S]*?)\n\\end\{table}.*`)
	sampleRegexp     = regexp.MustCompile(`^\\begin\{sample}\{([1-9]\d*)}.*([\s\S]*?)\n\\sample\{input}.*([\s\S]*?)\n\\sample\{output}.*([\s\S]*?)\n\\end\{sample}.*`)
	mixCodeRegexp    = regexp.MustCompile(`^\\begin\{mixcode}.*([\s\S]*?)\n\\end\{mixcode}.*`)
)

var inlineRegexp = regexp.MustCompile(`\\(link|text|newline|space|math|html)(?:\{([^{}]*)}\{([^{}]*)}|\{([^{}]*)}|\[([^\[\]]*)])?`)

func parseInline(text string, indices []int, v1, v2 *string) *Node {
	if v1 == nil {
		return nil
	}
	if indices[4] != -1 && indices[5] != -1 {
		*v1 = strings.TrimSpace(text[indices[4]:indices[5]])
		if v2 != nil {
			*v2 = strings.TrimSpace(text[indices[6]:indices[7]])
		}
	} else if indices[8] != -1 && indices[9] != -1 && v2 == nil {
		*v1 = strings.TrimSpace(text[indices[8]:indices[9]])
	} else if indices[10] != -1 && indices[11] != -1 && v2 == nil {
		*v1 = strings.TrimSpace(text[indices[10]:indices[11]])
	} else {
		return createTextNode(text[indices[0]:indices[1]])
	}
	return nil
}

func filterBlank(arr []string) []string {
	var res []string
	for _, item := range arr {
		item = strings.TrimSpace(item)
		if item != "" {
			res = append(res, item)
		}
	}
	return res
}

func createTextNode(text string, props ...string) *Node {
	props = filterBlank(props)
	if len(props) == 0 {
		return NewNode("Text", text, nil)
	}
	opts := map[string]interface{}{
		"type":      "default",
		"strong":    false,
		"italic":    false,
		"underline": false,
		"delete":    false,
		"code":      false,
	}
	for _, prop := range props {
		prop = strings.TrimSpace(prop)
		switch prop {
		case "default", "success", "info", "warning", "error":
			opts["type"] = prop
		case "s", "strong":
			opts["strong"] = true
		case "i", "italic":
			opts["italic"] = true
		case "u", "underline":
			opts["underline"] = true
		case "d", "delete":
			opts["delete"] = true
		case "c", "code":
			opts["code"] = true
		}
	}
	val := map[string]interface{}{
		"content": text,
		"options": opts,
	}
	return NewNode("Text", val, nil)
}

func processInline(node *Node) {
	if node == nil {
		return
	}
	text, ok := node.Value.(string)
	if !ok {
		return
	}
	var indices []int
	for text != "" {
		if indices = inlineRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 12 {
			if indices[0] != 0 {
				raw := text[:indices[0]]
				if raw != "" {
					node.Children = append(node.Children, createTextNode(raw))
				}
			}
			tag := text[indices[2]:indices[3]]
			switch tag {
			case "link":
				var link, label string
				if _node := parseInline(text, indices, &link, &label); _node != nil {
					node.Children = append(node.Children, _node)
					break
				}
				if link != "" {
					if label == "" {
						label = link
					}
					node.Children = append(node.Children, NewNode("Link", map[string]string{
						"label": label,
						"link":  link,
					}, nil))
				}
			case "newline":
				node.Children = append(node.Children, NewNode("Newline", nil, nil))
			case "space":
				node.Children = append(node.Children, NewNode("Space", nil, nil))
			case "text":
				var content, opts string
				if _node := parseInline(text, indices, &content, &opts); _node != nil {
					node.Children = append(node.Children, _node)
					break
				}
				if content != "" {
					node.Children = append(node.Children, createTextNode(content, strings.Split(opts, "|")...))
				}
			case "math":
				var formula string
				if _node := parseInline(text, indices, &formula, nil); _node != nil {
					node.Children = append(node.Children, _node)
					break
				}
				if formula != "" {
					node.Children = append(node.Children, NewNode("Math", formula, nil))
				}
			case "html":
				var html string
				if _node := parseInline(text, indices, &html, nil); _node != nil {
					node.Children = append(node.Children, _node)
					break
				}
				if html != "" {
					node.Children = append(node.Children, NewNode("Html", html, nil))
				}
			}
			text = text[indices[1]:]
			continue
		}
		node.Children = append(node.Children, createTextNode(text))
		break
	}
	node.Value = nil
}

func processSection(node *Node, text string) (string, bool) {
	if indices := sectionRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 4 {
		title := strings.TrimSpace(text[indices[2]:indices[3]])
		if title != "" {
			node.Children = append(node.Children, NewNode("Section", title, nil))
			return text[indices[1]:], true
		}
	}
	return text, false
}

func processParagraph(node *Node, text string) (string, bool) {
	if indices := paragraphRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 6 {
		var style interface{}
		style = strings.TrimSpace(text[indices[2]:indices[3]])
		if style == "" {
			style = nil
		}
		raw := strings.TrimSpace(text[indices[4]:indices[5]])
		if raw != "" {
			paragraph := NewNode("Paragraph", raw, style)
			processInline(paragraph)
			node.Children = append(node.Children, paragraph)
			return text[indices[1]:], true
		}
	}
	return text, false
}

func processBlockquote(node *Node, text string) (string, bool) {
	if indices := blockquoteRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 4 {
		raw := strings.TrimSpace(text[indices[2]:indices[3]])
		if raw != "" {
			blockquote := NewNode("Blockquote", strings.TrimSpace(text[indices[2]:indices[3]]), nil)
			processInline(blockquote)
			node.Children = append(node.Children, blockquote)
			return text[indices[1]:], true
		}
	}
	return text, false
}

func processCodeBlock(node *Node, text string) (string, bool) {
	if indices := codeBlockRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 6 {
		language := strings.TrimSpace(text[indices[2]:indices[3]])
		if language == "" {
			language = "text"
		}
		node.Children = append(node.Children, NewNode("CodeBlock", map[string]string{
			"language": language,
			"code":     strings.TrimSpace(text[indices[4]:indices[5]]),
		}, nil))
		return text[indices[1]:], true
	}
	return text, false
}

func processMathBlock(node *Node, text string) (string, bool) {
	if indices := mathBlockRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 4 {
		node.Children = append(node.Children, NewNode("MathBlock", strings.TrimSpace(text[indices[2]:indices[3]]), nil))
		return text[indices[1]:], true
	}
	return text, false
}

func processHtmlBlock(node *Node, text string) (string, bool) {
	if indices := htmlBlockRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 4 {
		node.Children = append(node.Children, NewNode("HtmlBlock", strings.TrimSpace(text[indices[2]:indices[3]]), nil))
		return text[indices[1]:], true
	}
	return text, false
}

func processTable(node *Node, text string) (string, bool) {
	if indices := tableRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 4 {
		raw := strings.TrimSpace(text[indices[2]:indices[3]])
		if raw != "" {
			if data := parseTable(raw); data != nil {
				node.Children = append(node.Children, NewNode("Table", data, nil))
			}
			return text[indices[1]:], true
		}
	}
	return text, false
}

func processSample(node *Node, text string) (string, bool) {
	if indices := sampleRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 10 {
		index, _ := strconv.Atoi(text[indices[2]:indices[3]])
		hintRaw := strings.TrimSpace(text[indices[4]:indices[5]])
		var hint interface{} = ""
		if hintRaw != "" {
			tmp := NewNode("tmp", hintRaw, nil)
			processInline(tmp)
			hint = tmp.Children
		}
		node.Children = append(node.Children, NewNode("Sample", map[string]interface{}{
			"index":  index,
			"hint":   hint,
			"input":  strings.TrimSpace(text[indices[6]:indices[7]]),
			"output": strings.TrimSpace(text[indices[8]:indices[9]]),
		}, nil))
		return text[indices[1]:], true
	}
	return text, false
}

func processMixCode(node *Node, text string) (string, bool) {
	if indices := mixCodeRegexp.FindStringSubmatchIndex(text); indices != nil && len(indices) >= 4 {
		raw := strings.TrimSpace(text[indices[2]:indices[3]])
		if raw != "" {
			if data := parseMixCode(raw); data != nil {
				node.Children = append(node.Children, NewNode("MixCode", data, nil))
			}
			return text[indices[1]:], true
		}
	}
	return text, false
}

func processUnknown(text string) string {
	lines := strings.SplitN(text, "\n", 2)
	if len(lines) == 2 {
		return lines[1]
	}
	return ""
}

func Lex(text string) *Node {
	node := NewNode("Root", nil, nil)
	var ok bool
	for text != "" {
		text = strings.TrimSpace(text)
		if text, ok = processSection(node, text); ok {
			continue
		}
		if text, ok = processParagraph(node, text); ok {
			continue
		}
		if text, ok = processBlockquote(node, text); ok {
			continue
		}
		if text, ok = processCodeBlock(node, text); ok {
			continue
		}
		if text, ok = processMathBlock(node, text); ok {
			continue
		}
		if text, ok = processHtmlBlock(node, text); ok {
			continue
		}
		if text, ok = processTable(node, text); ok {
			continue
		}
		if text, ok = processSample(node, text); ok {
			continue
		}
		if text, ok = processMixCode(node, text); ok {
			continue
		}
		text = processUnknown(text)
	}
	return node
}
