package parser

import "fmt"

type DocSectionNode struct {
	DocFunc      []*DocFuncNode
	AnonymousNum int
}

func NewDocSectionNode(f []*DocFuncNode) *DocSectionNode {
	return &DocSectionNode{DocFunc: f}
}

type DocFuncNode struct {
	Name     string
	DocSexpr *DocSexpr
}

func NewDocFuncNode(n string, s *DocSexpr) *DocFuncNode {
	return &DocFuncNode{Name: n, DocSexpr: s}
}

type DocList struct {
	l []*DocSexpr
}

func NewDocList(content []*DocSexpr) *DocList {
	return &DocList{l: content}
}

type DocSexpr struct {
	Atom *DocAtom
	List *DocList
}

func NewDocSexprFromAtom(a *DocAtom) *DocSexpr {
	return &DocSexpr{Atom: a}
}
func NewDocSexprFromList(l *DocList) *DocSexpr {
	return &DocSexpr{List: l}
}

type DocAtom struct {
	Tag   *DocTag
	Invoc *FuncInvoc
}

func NewDocAtom(t *DocTag, i *FuncInvoc) *DocAtom {
	return &DocAtom{Tag: t, Invoc: i}
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

		"h1", "h2", "h3", "h4", "h5", "title",
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
