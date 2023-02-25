package parser

type ProgramNode struct {
	ImportSection          *ImportSectionNode
	CSSSection             *CSSSectionNode
	TextSection            *TextSectionNode
	DocSection             *DocSectionNode
	Extern                 []string
	Global                 []*PFormal
	NeedBytes, NeedElement bool
}

func NewProgramNode() *ProgramNode {
	return &ProgramNode{}
}
