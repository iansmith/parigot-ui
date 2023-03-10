package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type NameCheck struct {
	*BasewclVisitor
	Passed  bool
	Func    map[string]bool
	Program *ProgramNode
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
	n.Program = ctx.GetP()
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
	if msg := dfunc.CheckForBadVariableUse(); msg != "" {
		log.Printf("%s\n", msg)
		n.Passed = false
	}
	if !checkFuncForCollisions(dfunc.Name, dfunc.Param, dfunc.Local, false) {
		n.Passed = false
	}
	return nil
}

// VisitText_section checks that all the text functions' names are distinct.
func (n *NameCheck) VisitText_section(ctx *Text_sectionContext) interface{} {
	for _, node := range ctx.GetSection().Func {
		node.Section = ctx.GetSection()
		if !n.checkForDupNames(ctx.GetParser(), ctx.BaseParserRuleContext, node.Name, true) {
			n.Passed = false
			return nil
		}
	}
	for _, decl := range ctx.AllText_func() {
		if msg := decl.GetF().CheckForBadVariableUse(); msg != "" {
			log.Print(msg)
			n.Passed = false
		}
		n.Visit(decl)

	}
	return nil
}

// VisitDoc_section checks that all the doc functions' names are distinct.
func (n *NameCheck) VisitDoc_section(ctx *Doc_sectionContext) interface{} {
	for _, node := range ctx.GetSection().DocFunc {
		node.Section = ctx.GetSection()
		if !n.checkForDupNames(ctx.GetParser(), ctx.BaseParserRuleContext, node.Name, true) {
			n.Passed = false
			return nil
		}
	}

	for _, f := range ctx.AllDoc_func() {
		n.Visit(f)
		errorMsg := n.checkFuncCallName(f.GetFn())
		if errorMsg != "" {
			log.Print(errorMsg)
			n.Passed = false
			return nil
		}
	}
	// mark objects for code generation
	ctx.GetSection().SetNumber()
	return nil

}

// /////////////////////////////// CHECKING FUNCTIONS

func (n *NameCheck) checkFuncCallName(fn *DocFuncNode) string {
	e := fn.Elem
	if e.TextContent != nil {
		if !e.TextContent.Name.IsVar && strings.HasPrefix(e.TextContent.Name.Name, anonPrefix) {
			return ""
		}
		// is it a variable ref?
		if e.TextContent.Name.IsVar {
			return fmt.Sprintf("You cannot use variable '%s' as the name of a function currently", e.TextContent.Name.Name)
		} else {
			f, ok := n.Func[e.TextContent.Name.Name]
			if !ok {
				found := false
				for _, ext := range n.Program.Extern {
					if ext == e.TextContent.Name.Name {
						found = true
						break
					}
				}
				if !found {
					return fmt.Sprintf("in function '%s', unable to find function '%s'", fn.Name, e.TextContent.Name.Name)
				}
			} else { //found the name
				if !f {
					return fmt.Sprintf("use of doc functions to create content is currently not supported (such as '%s' in function '%s')", e.TextContent.Name.Name, fn.Name)
				}
				return ""
			}
		}

	}
	return n.checkFuncCallParameters(fn, e.TextContent)
}

// checkFuncCallParameters checks that that an invocation only uses variables that are
// known.  This is called only after we have checked that the name of the function being
// invoced is ok.
func (n *NameCheck) checkFuncCallParameters(fn *DocFuncNode, invoc *FuncInvoc) string {
	if invoc == nil || invoc.Actual == nil {
		return ""
	}
	for _, p := range invoc.Actual {
		if p.Literal != "" {
			continue
		}
		if formalContainsActual(fn.Local, p.Var) {
			continue
		}
		if formalContainsActual(fn.Param, p.Var) {
			continue
		}
		if formalContainsActual(n.Program.Global, p.Var) {
			continue
		}
		found := false
		for _, e := range n.Program.Extern {
			if e == p.Var {
				found = true
				break
			}
		}
		if !found {
			return fmt.Sprintf("in function '%s', use of unknown variable '%s'", fn.Name, p.Var)
		}
	}
	return ""
}

func formalContainsActual(formal []*PFormal, actual string) bool {
	for _, f := range formal {
		if f.Name == actual {
			return true
		}
	}
	return false
}

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
		// param is also name?
		paramName := param.Name
		if paramName == name {
			log.Printf("cannot use '%s' as the name of a %s function and a parameter to that text function", fnType, name)
			return false
		}
		// same local and param
		if formalContainsActual(l, param.Name) {
			log.Printf("cannot use '%s' as the name of both a local and a parameter", param.Name)
			return false

		}
		// dup param
		_, ok := used[paramName]
		if ok {
			log.Printf("in %s function '%s': duplicate parameter name '%s' ", fnType, name, paramName)
			return false
		}
		used[paramName] = struct{}{}
	}
	for _, local := range l {
		// local is also name?
		localName := local.Name
		if localName == name {
			log.Printf("cannot use '%s' as the name of a %s function and a local to that text function", fnType, name)
			return false
		}
		// dup param
		_, ok := used[localName]
		if ok {
			log.Printf("in %s function '%s': duplicate local name '%s' ", fnType, name, localName)
			return false
		}
		used[localName] = struct{}{}
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

func IsSelfVar(name string) bool {
	return name == "result"
}
