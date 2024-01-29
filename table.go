package yutex

import "strings"

type TableItem []*Node

/*
The test like this:
| Column 1 | Column 2 |	Column 3 |
| :------- | :------: | --------:|
| centered 文本居左 | right-aligned 文本居中 |right-aligned 文本居右|
*/
func parseTable(text string) map[string]interface{} {
	rows := strings.Split(text, "\n")
	if len(rows) < 3 {
		return nil
	}
	// 提取表头
	header := extractTableRow(rows[0])
	column := len(header)
	aligns := parseTableAlign(rows[1])
	if len(aligns) > column {
		aligns = aligns[:column]
	} else if len(aligns) < column {
		for i := len(aligns); i < column; i++ {
			aligns = append(aligns, "left")
		}
	}
	// 提取表格内容
	data := make([][]TableItem, len(rows)-2)
	for i := 2; i < len(rows); i++ {
		data[i-2] = extractTableRow(rows[i])
		if len(data[i-2]) > column {
			data[i-2] = data[i-2][:column]
		} else if len(data[i-2]) < column {
			for j := len(data[i-2]); j < column; j++ {
				data[i-2] = append(data[i-2], nil)
			}
		}
	}
	return map[string]interface{}{
		"header": header,
		"aligns": aligns,
		"data":   data,
	}
}

func extractTableRow(row string) []TableItem {
	row = strings.Trim(row, "|")
	cols := strings.Split(row, " | ")
	items := make([]TableItem, len(cols))
	for i, col := range cols {
		col = strings.TrimSpace(col)
		node := NewNode("tmp", col, nil)
		processInline(node)
		items[i] = node.Children
	}
	return items
}

func parseTableAlign(row string) []string {
	row = strings.Trim(row, "|")
	cols := strings.Split(row, " | ")
	aligns := make([]string, len(cols))
	for i, col := range cols {
		col = strings.TrimSpace(col)
		if strings.HasPrefix(col, ":") && strings.HasSuffix(col, ":") {
			aligns[i] = "center"
		} else if strings.HasPrefix(col, ":") {
			aligns[i] = "left"
		} else if strings.HasSuffix(col, ":") {
			aligns[i] = "right"
		} else {
			aligns[i] = "left"
		}
	}
	return aligns
}
