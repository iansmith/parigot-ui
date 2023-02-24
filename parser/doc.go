package parser

import (
	"fmt"
)

type DocSectionNode struct {
	DocFunc      []*DocFuncNode
	AnonymousNum int
}

func (s *DocSectionNode) SetNumber() {
	for _, fn := range s.DocFunc {
		fn.SetNumber()
	}
}

func NewDocSectionNode(f []*DocFuncNode) *DocSectionNode {
	return &DocSectionNode{DocFunc: f}
}

type DocFuncNode struct {
	Name string
	Elem *DocElement
}

func (f *DocFuncNode) SetNumber() {
	f.Elem.SetNumber(0)
}

func NewDocFuncNode(n string, s *DocElement) *DocFuncNode {
	return &DocFuncNode{Name: n, Elem: s}
}

type DocElement struct {
	Number      int
	Tag         *DocTag
	TextContent *FuncInvoc
	Child       []*DocElement
}

func (e *DocElement) SetNumber(n int) int {
	if e.TextContent == nil && len(e.Child) == 0 {
		e.Number = n
		return n + 1
	}
	if e.TextContent != nil {
		e.Number = n
		return n + 1
	}
	e.Number = n
	n++
	for _, c := range e.Child {
		n = c.SetNumber(n)
	}
	return n
}

type DocTag struct {
	Tag   string
	Id    string
	Class []string
}

func NewDocTag(tag string, id string, class []string) (*DocTag, error) {
	if !validTag(tag) {
		return nil, fmt.Errorf("unknown tag '%s'", tag)
	}
	return &DocTag{Tag: tag, Id: id, Class: class}, nil
}

func validTag(tag string) bool {
	switch tag {
	case
		"article", "aside", "details", "figcaption", "figure", "footer", "header", "legend", "main",
		"mark", "nav", "section", "summary", "time",
		"abbr", "address", "base", "blockquote", "body", "col", "head", "hr", "link", "meta", "noscript",
		"object", "param", "progress", "q", "sub", "sup", "track", "var", "video", "wbr",

		"h1", "h2", "h3", "h4", "h5", "title", "br",
		"strong", "em",
		"a", "p", "span", "div",
		"form", "input", "fieldset", "label", "keygen", "optgroup", "option", "textarea",
		"ul", "ol", "dl", "dd", "dt", "li",
		"img",
		"code", "kbd", "pre", "samp",
		"script",
		"table", "tbody", "td", "tr":
		return true
	}
	return false
}
