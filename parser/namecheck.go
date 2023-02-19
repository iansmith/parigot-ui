package parser

import (
	"log"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type NameCheck struct {
	*BasewclVisitor
	Passed   bool
	TextFunc map[string]struct{}
}

var _ wclVisitor = &NameCheck{}

func NewNameCheck() *NameCheck {
	return &NameCheck{
		BasewclVisitor: &BasewclVisitor{
			BaseParseTreeVisitor: &antlr.BaseParseTreeVisitor{},
		},
		Passed:   true,
		TextFunc: make(map[string]struct{}),
	}
}

// NameCheckVisit returns true if the visiting pass on tree
// completes without error.
func NameCheckVisit(tree antlr.ParseTree) bool {
	n := NewNameCheck()
	n.Visit(tree)
	return n.Passed
}

// VisitText_decl checks that the parameters are not the same name as
// the function name and that the parameter names are distinct.
func (n *NameCheck) VisitText_decl(ctx *Text_declContext) interface{} {
	node := ctx.GetF()
	used := make(map[string]struct{})
	for _, param := range node.Param {
		paramName := param.Name
		if paramName == node.Name {
			log.Printf("cannot use '%s' as the name of a text function and a parameter to that text function", node.Name)
			n.Passed = false
			return nil
		}
		_, ok := used[paramName]
		if ok {
			log.Printf("in text function '%s': duplicate parameter name '%s' ", node.Name, paramName)
			n.Passed = false
			return nil
		}
		used[paramName] = struct{}{}
	}
	return nil
}

// VisitText_section checks that all the text functions' names are distinct.
func (n *NameCheck) VisitText_section(ctx *Text_sectionContext) interface{} {
	for _, node := range ctx.GetSection().Func {
		if _, ok := n.TextFunc[node.Name]; ok {
			log.Printf("text function name '%s' used more than once", node.Name)
			n.Passed = false
			return nil
		}
		n.TextFunc[node.Name] = struct{}{}
	}
	for _, decl := range ctx.AllText_decl() {
		n.Visit(decl)
	}
	return nil
}

// //////////////////////////////// BOILERPLATE /////////////////////////////
func (n *NameCheck) VisitProgram(ctx *ProgramContext) interface{} {
	return n.Visit(ctx.Text_section())
}

func (n *NameCheck) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(n)
}
