package driver

import (
	"embed"
	"flag"
	"io"
	"log"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot-ui/parser"
)

var langToTempl = map[string]string{
	"go": golang,
}

var language = flag.String("l", "go", "pass the name of a known language to get result in that language")

// GlobalStackScope needs to be visible to every place in the parser
// or it becomes a nightmare to push and pop things.
var GlobalScopeStack *parser.ScopeStack

//go:embed template/*
var templateFS embed.FS

func Main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalf("you must provide a filename that contains a web coordination language description")
	}
	fp, err := os.Open(flag.Arg(0))
	if err != nil {
		wd, _ := os.Getwd()
		log.Fatalf("opening %s: %v (wd is %s)", flag.Arg(0), err, wd)
	}
	buffer, err := io.ReadAll(fp)
	if err != nil {
		log.Fatalf("reading %s: %v", flag.Arg(0), err)
	}
	el := errorListener{0}
	input := antlr.NewInputStream(string(buffer))
	lexer := parser.Newwcllex(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&el)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.Newwcl(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(&el)
	listener := parser.WclGenListener{}
	prog := p.Program()
	antlr.ParseTreeWalkerDefault.Walk(&listener, prog)

	// create a context for this generate
	t, ok := langToTempl[*language]
	if !ok {
		log.Fatalf("unable to find a template for language '%s'", *language)
	}
	ctx := newGenerateContext(t)
	GlobalScopeStack = ctx.scope // pointer copy
	ctx.program = prog.GetP()
	ctx.templateName = t
	ctx.global["import"] = prog.GetP().ImportSection
	ctx.global["text"] = prog.GetP().TextSection
	err = runTemplate(ctx)
	if err != nil {
		log.Fatalf("error trying to execute template %s: %v", ctx.templateName, err)
	}
}
