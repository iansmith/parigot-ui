package driver

import (
	"fmt"
	"log"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot-ui/parser"
)

const test1 = `text foo  {{ this is a template with embedded {{ stuff }} }}`
const test2 = `text foo  {{ this is a template with embedded {{ bad stuff {{}} }} }}`
const test3 = `text foo  {{ this is a template with embedded {{ {{}}bad {{{{}}}}stuff {{really bad stuff}}}}}}`

const mismatch1 = `text foo  {{ this is a template with embedded {{ bad stuff }} }} }}`

func Main() {
	expectedErrors := []int{0, 0, 0, 1}
	for i, text := range []string{test1, test2, test3, mismatch1} {
		fmt.Printf("test is: %s\n", text)
		el := errorListener{0}
		input := antlr.NewInputStream(text)
		lexer := parser.NewwclLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewwclParser(stream)
		p.AddErrorListener(&el)
		p.Program()
		if el.count != expectedErrors[i] {
			panic(fmt.Sprintf("test %d failed (expected %d but got %d)", i, expectedErrors[i], el.count))
		}
	}
}

type errorListener struct {
	count int
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	log.Printf("xxxline %d:%d %s %v", line, column, msg, offendingSymbol)
	el.count++
}
func (el *errorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	log.Printf("ambiguous alternatives in DFA")
	el.count++
}
func (el *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	log.Printf("attempting full context in DFA")
	el.count++
}
func (el *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	log.Printf("context sensitivity in DFA")
	el.count++
}
