package parser

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// TextConstant is a simple string.
type TextConstant struct {
	_VarCtx *VarCtx
	Value   string
}

func (t *TextConstant) String() string {
	return t.Value
}

func (t *TextConstant) Generate(_ *VarCtx) string {
	value := t.Value
	if utf8.RuneCountInString(value) == -1 {
		panic("not a utf8 string:" + value)
	}
	if utf8.RuneCountInString(value) > 2 {
		first, fCount := utf8.DecodeRuneInString(value)
		last, lCount := utf8.DecodeLastRuneInString(value)
		if first == utf8.RuneError || last == utf8.RuneError {
			panic("not a utf8 string:" + value)
		}
		if unicode.IsSpace(first) && unicode.IsSpace(last) {
			value = t.Value[fCount : len(t.Value)-lCount]
		}
	}
	return fmt.Sprintf("buf WriteString(%q)\n", value)
}

func (t *TextConstant) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextConstant(s string) *TextConstant {
	return &TextConstant{_VarCtx: nil, Value: s}
}

// TextVar is a text variable that in source form is ${foo}
type TextVar struct {
	_VarCtx *VarCtx
	Name    string
}

func (t *TextVar) String() string {
	return fmt.Sprintf("${%s}", t.Name)
}

func (t *TextVar) Generate(_ *VarCtx) string {
	return fmt.Sprintf("buf.WriteString(%s)\n", t.Name)
}

func (t *TextVar) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextVar() *TextVar {
	return &TextVar{}
}

// TextItem is something that knows how to print itself.
type TextItem interface {
	String() string
	Generate(*VarCtx) string
	VarCtx() *VarCtx
}

// TextExpander is something that can have variables uses in it.
type TextExpander interface {
	Item() []TextItem
}

// TextFuncNode is the that alls the information about a declared
// text function.
type TextFuncNode struct {
	Name  string
	_Item []TextItem
}

func (t *TextFuncNode) Item() []TextItem {
	return t._Item
}

func (t *TextFuncNode) SetItem(item []TextItem) {
	t._Item = item
}

func (t *TextFuncNode) String() string {
	var buf bytes.Buffer
	for _, t := range t._Item {
		buf.WriteString(t.String())
	}
	return buf.String()
}

func NewTextFuncNode() *TextFuncNode {
	return &TextFuncNode{}
}

// TestSection is the collection of text functions.
type TextSectionNode struct {
	Func []*TextFuncNode
}

func NewTextSectionNode() *TextSectionNode {
	return &TextSectionNode{}
}
