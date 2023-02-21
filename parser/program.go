package parser

type ProgramNode struct {
	ImportSection *ImportSectionNode
	CSSSection    *CSSSectionNode
	TextSection   *TextSectionNode
	DocSection    *DocSectionNode
}

func NewProgramNode() *ProgramNode {
	return &ProgramNode{}
}

type FuncInvoc struct {
	Name string
}

func NewFuncInvoc(n string) *FuncInvoc {
	return &FuncInvoc{Name: n}
}
