package parser

// TextConstant is a simple string.
type TextConstant struct {
	Value string
}

func (t *TextConstant) String() string {
	return t.Value
}

func NewTextConstant(s string) *TextConstant {
	return &TextConstant{s}
}

// TextVar is a text variable that in source form is ${foo}
type TextVar struct {
	Name string
}

func (t *TextVar) String() string {
	return "TEXT VAR NOT IMPL"
}

// TextItem is something that knows how to print itself.
type TextItem interface {
	String() string
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

func NewTextFuncNode(name string, item []TextItem) *TextFuncNode {
	return &TextFuncNode{Name: name, _Item: item}
}

// TestSection is the collection of text functions.
type TextSectionNode struct {
	Func []*TextFuncNode
}

func NewTextSectionNode() *TextSectionNode {
	return &TextSectionNode{}
}
