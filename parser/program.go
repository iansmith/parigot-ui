package parser

type ProgramNode struct {
	ImportSection *ImportSectionNode
	CSSSection    *CSSSectionNode
	TextSection   *TextSectionNode
}

func NewProgramNode() *ProgramNode {
	return &ProgramNode{}
}
