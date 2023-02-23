package parser

import (
	"fmt"
	"log"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type NameCheck struct {
	*BasewclVisitor
	Passed bool
	Func   map[string]bool
}

var _ wclVisitor = &NameCheck{}

func NewNameCheck() *NameCheck {
	return &NameCheck{
		BasewclVisitor: &BasewclVisitor{
			BaseParseTreeVisitor: &antlr.BaseParseTreeVisitor{},
		},
		Passed: true,
		Func:   make(map[string]bool),
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
		ok := n.checkFuncName(node.Name, true)
		if !ok {
			msg := fmt.Sprintf("text function name '%s' used more than once", node.Name)
			ex := antlr.NewBaseRecognitionException(msg, ctx.GetParser(), ctx.GetParser().GetInputStream(), ctx)
			ctx.SetException(ex)
			n.Passed = false
			return nil
		}
	}
	for _, decl := range ctx.AllText_decl() {
		n.Visit(decl)
	}
	return nil
}

// VisitDoc_section checks that all the doc functions' names are distinct.
func (n *NameCheck) VisitDoc_section(ctx *Doc_sectionContext) interface{} {
	for _, node := range ctx.GetSection().DocFunc {
		ok := n.checkFuncName(node.Name, false)
		if !ok {
			msg := fmt.Sprintf("doc function name '%s' used more than once", node.Name)
			ex := antlr.NewBaseRecognitionException(msg, ctx.GetParser(), ctx.GetParser().GetInputStream(), ctx)
			ctx.SetException(ex)
			n.Passed = false
			return nil
		}
	}
	for _, sexpr := range ctx.AllDoc_sexpr() {
		n.Visit(sexpr)
	}
	// mark objects for code generation
	ctx.GetSection().SetParentAndNumber() //recursive traversal downward, pre-order
	return nil

}

// checkFuncName checks to see if the specified function is already
// present.
func (n *NameCheck) checkFuncName(name string, isText bool) bool {
	typeOfFunc := "doc"
	other := "text"
	if isText {
		typeOfFunc = "text"
		other = "doc"
	}

	if b, ok := n.Func[name]; ok {
		if b == isText {
			if isText {
				log.Printf("two text functions named '%s'", name)
			} else {
				log.Printf("two doc functions named '%s'", name)
			}
		} else {
			log.Printf("%s function '%s' conflicts with %s function", typeOfFunc, name, other)
		}
		return false
	}
	n.Func[name] = isText
	return true
}

// //////////////////////////////// BOILERPLATE /////////////////////////////
func (n *NameCheck) VisitProgram(ctx *ProgramContext) interface{} {
	return n.VisitChildren(ctx)
}

func (n *NameCheck) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(n)
}
func (n *NameCheck) VisitChildren(ctx antlr.RuleNode) interface{} {
	count := ctx.GetChildCount()
	var last interface{}
	for i := 0; i < count; i++ {
		tree := ctx.GetChild(i).(antlr.ParseTree)
		last = n.Visit(tree)
	}
	return last
}
