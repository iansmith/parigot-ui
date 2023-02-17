package parser

import (
	"fmt"
)

type ImportSectionNode struct {
	Text TextItem
}

func (i *ImportSectionNode) Dump(indent int) {
	print(fmt.Sprintf("%*s (ImportSectionNode\n%s", indent, "", i.Text.String()))
}

func NewImportSectionNode() *ImportSectionNode {
	return &ImportSectionNode{}
}
