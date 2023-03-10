package css

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

func ReadCSS(path string) (map[string]struct{}, error) {
	builder, p := readInput(path)
	sheet := p.Stylesheet()
	antlr.ParseTreeWalkerDefault.Walk(builder, sheet)
	return builder.(*CSSBuild).ClassName, nil
}

func readInput(path string) (css3Listener, *css3Parser) {
	fp, err := os.Open(path)
	if err != nil {
		wd, _ := os.Getwd()
		log.Fatalf("%v (wd is %s), %v", flag.Arg(0), wd, err)
	}
	buffer, err := io.ReadAll(fp)
	if err != nil {
		log.Fatalf("reading %s: %v", flag.Arg(0), err)
	}
	fp.Close()
	//el := errorListener{0}
	input := antlr.NewInputStream(string(buffer))
	lexer := Newcss3Lexer(input)
	//lexer.RemoveErrorListeners()
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := Newcss3Parser(stream)
	p.RemoveErrorListeners()
	// the diagnostic listener is good for debugging (displays good error msgs)
	p.AddErrorListener(&antlr.DiagnosticErrorListener{
		DefaultErrorListener: &antlr.DefaultErrorListener{},
	})
	//p.AddErrorListener(&el)
	return NewCSSBuild(), p
}
