// generated by gostringer -type=Kind; DO NOT EDIT

package syntaxhighlight

import "fmt"

const _Kind_name = "WhitespaceStringKeywordCommentTypeLiteralPunctuationPlaintextTagHTMLTagHTMLAttrNameHTMLAttrValueDecimal"

var _Kind_index = [...]uint8{0, 10, 16, 23, 30, 34, 41, 52, 61, 64, 71, 83, 96, 103}

func (i Kind) GoString() string {
	if i+1 >= Kind(len(_Kind_index)) {
		return fmt.Sprintf("syntaxhighlight.Kind(%d)", i)
	}
	return "syntaxhighlight." + _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
