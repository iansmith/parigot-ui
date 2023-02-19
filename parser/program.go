package parser

type ProgramNode struct {
	ImportSection *ImportSectionNode
	CSSSection    *CSSSectionNode
	TextSection   *TextSectionNode
	DocSection    *DocSection
}

func NewProgramNode() *ProgramNode {
	return &ProgramNode{}
}
