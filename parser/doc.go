package parser

import (
	"fmt"
)

type DocNode interface {
	SetParentAndNumber(DocNode, int) int
	Parent() DocNode
	Number() int
}

type DocSectionNode struct {
	DocFunc      []*DocFuncNode
	AnonymousNum int
}

func (s *DocSectionNode) SetParentAndNumber() {
	for _, fn := range s.DocFunc {
		fn.SetParentAndNumber()
	}
}

func NewDocSectionNode(f []*DocFuncNode) *DocSectionNode {
	return &DocSectionNode{DocFunc: f}
}

type DocFuncNode struct {
	Name     string
	DocSexpr []*DocSexpr
}

func (f *DocFuncNode) SetParentAndNumber() {
	f.DocSexpr[0].SetParentAndNumber(nil, 0)
}

func NewDocFuncNode(n string, s *DocSexpr) *DocFuncNode {
	return &DocFuncNode{Name: n, DocSexpr: []*DocSexpr{s}}
}

type DocList struct {
	List    []*DocSexpr
	Parent_ DocNode
	Number_ int
}

func (l *DocList) Parent() DocNode {
	return l.Parent_
}
func (l *DocList) Number() int {
	return l.Number_
}

func (l *DocList) SetParentAndNumber(p DocNode, n int) int {
	l.Parent_ = p
	l.Number_ = n
	n++

	for _, sexpr := range l.List {
		n = sexpr.SetParentAndNumber(l, n)
	}
	return n
}

func NewDocList(content []*DocSexpr) *DocList {
	return &DocList{List: content}
}

type DocSexpr struct {
	Atom    *DocAtom
	List    *DocList
	Parent_ DocNode
	Number_ int
}

func (s *DocSexpr) Parent() DocNode {
	return s.Parent_
}

func (s *DocSexpr) SetParentAndNumber(p DocNode, n int) int {
	s.Parent_ = p
	s.Number_ = -32 //marker
	if s.Atom != nil {
		return s.Atom.SetParentAndNumber(s, n)
	} else {
		return s.List.SetParentAndNumber(s, n)
	}
}

func (s *DocSexpr) Number() int {
	if s.Atom != nil {
		return s.Atom.Number()
	}
	return s.List.Number()
}

func NewDocSexprFromAtom(a *DocAtom) *DocSexpr {
	return &DocSexpr{Atom: a}
}
func NewDocSexprFromList(l *DocList) *DocSexpr {
	return &DocSexpr{List: l}
}

type DocAtom struct {
	Tag     *DocTag
	Invoc   *FuncInvoc
	Parent_ DocNode
	Number_ int
}

func (a *DocAtom) Parent() DocNode {
	return a.Parent_
}

func (a *DocAtom) Number() int {
	return a.Number_
}

func (a *DocAtom) SetParentAndNumber(p DocNode, n int) int {
	a.Parent_ = p
	a.Number_ = n
	return n + 1
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
