package parser

type Program struct {
	ImportSection *ImportSectionNode
	CSSSection    *CSSSectionNode
	TextSection   *TextSectionNode
}

func NewProgram() *Program {
	return &Program{}
}
