package parser

import (
	"fmt"
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
	return fmt.Sprintf("buf WriteString(%q)\n", value)
}

func (t *TextConstant) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextConstant(s string) *TextConstant {
	return &TextConstant{_VarCtx: nil, Value: s}
}
func (t *TextConstant) SubTemplate() string {
	return "TextConstant"
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
	return fmt.Sprintf("LookupVar(%s)\n", t.Name)
}

func (t *TextVar) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextVar(n string) *TextVar {
	return &TextVar{Name: n}
}
func (t *TextVar) SubTemplate() string {
	return "TextVar"
}

// TextInline is a blob of code to copied into the output.
type TextInline struct {
	Name      string
	_VarCtx   *VarCtx
	TextItem_ []TextItem
}

func (t *TextInline) String() string {
	return "BLEAH NOT AVAILABLE"
}

func (t *TextInline) Generate(_ *VarCtx) string {
	return "BLEAH NOT AVAILABLE Generate"
}

func (t *TextInline) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextInline() *TextInline {
	return &TextInline{}
}

func (t *TextInline) SubTemplate() string {
	return "TextInline"
}

// TextItem is an interface that represents the things that we
// know how to place inside a text unit.
type TextItem interface {
	//String() string
	VarCtx() *VarCtx
	SubTemplate() string
}

// TextExpander is something that can have variables uses in it.
type TextExpander interface {
	Item() []TextItem
}

// PFormal holds a parameter and type pair.
type PFormal struct {
	Name string
	Type string
}

func NewPFormal(n, t string) *PFormal {
	return &PFormal{Name: n, Type: t}
}

// TextFuncNode is the that alls the information about a declared
// text function.
type TextFuncNode struct {
	Name      string
	NumParams int
	Param     []*PFormal
	_Item     []TextItem
}

func (t *TextFuncNode) Item() []TextItem {
	return t._Item
}

func (t *TextFuncNode) SetItem(item []TextItem) {
	t._Item = item
}

// func (t *TextFuncNode) String() string {
// 	var buf bytes.Buffer
// 	for _, t := range t._Item {
// 		buf.WriteString(t.String())
// 	}
// 	return buf.String()
// }

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
