package driver

import (
	"fmt"
	"log"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot-ui/parser"
)

const test1 = `text foo  {{ this is a template with embedded { fnCall()}}}`
const test2 = `text foo(EN,v)  #to eol 
{{ this is a template with embedded ${v} }}`
const test3 = `text foo  {{ this is a template with embedded {functions()} and ${vars} }}`
const test4 = `text _foo.bar  {{ this is a text blob that calls a function {f(x,y)} }}`
const test5 = `
#first one
text foo {{ this is a test}} 
      // second one
bar{{ this is aloso a test}}`

const test6 = `css bootstrap {foo bar baz} simple {blech fleazil}`

const mismatch1 = `text foo {{ this is a template with embedded { bad stuff }}`
const badId1 = `text foo@bar`
const badId2 = `text foo:bar`
const twoTextSection1 = `text foo {{ this is a test}} text bar {{ should fail }}`

func Main() {
	expectedErrors := []int{0, 0, 0, 0, 0, 0, 2, 3, 3, 1}
	for i, text := range []string{test1, test2, test3, test4, test5, test6, mismatch1, badId1, badId2, twoTextSection1} {
		fmt.Printf("test is: %s\n", text)
		el := errorListener{0}
		input := antlr.NewInputStream(text)
		lexer := parser.Newwcllex(input)
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(&el)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.Newwcl(stream)
		p.RemoveErrorListeners()
		p.AddErrorListener(&el)
		p.Program()
		if el.count != expectedErrors[i] {
			fmt.Printf("test %d failed (expected %d but got %d)", i, expectedErrors[i], el.count)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

type errorListener struct {
	count int
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	log.Printf("syntax error %d:%d %s %v", line, column, msg, offendingSymbol)
	el.count++
}
func (el *errorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	log.Printf("ambiguous alternatives in DFA")
	el.count++
}
func (el *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	log.Printf("attempting full context in DFA:")
	el.count++
}
func (el *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	log.Printf("context sensitivity in DFA")
	el.count++
}
