package parser

import (
	"bytes"
	"fmt"
)

type ImportSectionNode struct {
	Text []TextItem
}

func (i *ImportSectionNode) Dump(indent int) {
	print(fmt.Sprintf("%*s (ImportSectionNode\n%s", indent, "", i.String()))
}

func (i *ImportSectionNode) String() string {
	var buf bytes.Buffer
	for _, t := range i.Text {
		buf.WriteString(t.String())
	}
	return buf.String()
}

func NewImportSectionNode() *ImportSectionNode {
	return &ImportSectionNode{}
}
