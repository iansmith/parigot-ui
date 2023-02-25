package parser

import (
	"fmt"
	"log"
	"strings"

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

// Check for FN types at this level
func (n *NameCheck) VisitProgram(ctx *ProgramContext) interface{} {
	if ctx.GetP().TextSection != nil && len(ctx.GetP().TextSection.Func) > 0 {
		ctx.p.NeedBytes = true
	}
	if ctx.GetP().DocSection != nil && len(ctx.GetP().DocSection.DocFunc) > 0 {
		ctx.p.NeedElement = true
	}
	return n.VisitChildren(ctx)
}

// VisitText_decl checks that the parameters are not the same name as
// the function name and that the parameter names are distinct.  It also
// checks the Ids for dots.
func (n *NameCheck) VisitText_func(ctx *Text_funcContext) interface{} {
	node := ctx.GetF()
	if !checkFuncForCollisions(node.Name, node.Param, node.Local, true) {
		n.Passed = false
	}
	return nil
}

// VisitDoc_func checks that the parameters are not the same name as
// the function name and that the parameter names are distinct.  It also
// checks the Ids for dots.
func (n *NameCheck) VisitDoc_func(ctx *Doc_funcContext) interface{} {
	dfunc := ctx.GetFn()
	if !checkFuncForCollisions(dfunc.Name, dfunc.Param, dfunc.Local, false) {
		n.Passed = false
	}
	return nil
}

// VisitText_section checks that all the text functions' names are distinct.
func (n *NameCheck) VisitText_section(ctx *Text_sectionContext) interface{} {
	for _, node := range ctx.GetSection().Func {
		if !n.checkForDupNames(ctx.GetParser(), ctx.BaseParserRuleContext, node.Name, true) {
			n.Passed = false
			return nil
		}
	}
	for _, decl := range ctx.AllText_func() {
		n.Visit(decl)
	}
	return nil
}

// VisitDoc_section checks that all the doc functions' names are distinct.
func (n *NameCheck) VisitDoc_section(ctx *Doc_sectionContext) interface{} {
	for _, node := range ctx.GetSection().DocFunc {
		if !n.checkForDupNames(ctx.GetParser(), ctx.BaseParserRuleContext, node.Name, true) {
			n.Passed = false
			return nil
		}
	}

	for _, elem := range ctx.AllDoc_func() {
		n.Visit(elem)
	}
	// mark objects for code generation
	ctx.GetSection().SetNumber()
	return nil

}

// /////////////////////////////// CHECKING FUNCTIONS

func (n *NameCheck) checkForDupNames(parser antlr.Parser, ctx antlr.RuleContext, funcName string, isText bool) bool {
	ok := n.checkFuncName(funcName, isText)
	okDot := strings.Contains(funcName, ".")
	var msg string
	if !ok {
		msg = fmt.Sprintf("text function name '%s' used more than once", funcName)
	}
	if okDot {
		msg = fmt.Sprintf("text function name '%s' contains a dot, which is not allowed in names", funcName)
	}
	if !ok || okDot {
		log.Print(msg)
		ex := antlr.NewBaseRecognitionException(msg, parser, parser.GetInputStream(), ctx)
		ctx.(*antlr.BaseParserRuleContext).SetException(ex)
		return false
	}
	return true
}
func checkFuncForCollisions(name string, p []*PFormal, l []*PFormal, isText bool) bool {
	fnType := "text"
	if !isText {
		fnType = "doc"
	}
	used := make(map[string]struct{})
	for _, param := range p {
		paramName := param.Name
		if paramName == name {
			log.Printf("cannot use '%s' as the name of a %s function and a parameter to that text function", fnType, name)
			return false
		}
		_, ok := used[paramName]
		if ok {
			log.Printf("in %s function '%s': duplicate parameter name '%s' ", fnType, name, paramName)
			return false
		}
		used[paramName] = struct{}{}
	}
	msg := checkParamsAndLocalsForDot(name, p, l)
	if msg != "" {
		log.Print(msg)
		return false
	}
	return true
}

func checkParamsAndLocalsForDot(funcName string, param []*PFormal, local []*PFormal) string {
	for _, p := range param {
		if strings.Contains(p.Name, ".") {
			return fmt.Sprintf("In function '%s', parameter name '%s' may not contain a dot", funcName, p.Name)
		}
	}
	for _, p := range local {
		if strings.Contains(p.Name, ".") {
			return fmt.Sprintf("In function '%s', local name '%s' may not contain a dot", funcName, p.Name)
		}
	}
	return ""
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
