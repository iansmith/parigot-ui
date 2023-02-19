package driver

import (
	"embed"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot-ui/parser"
)

var langToTempl = map[string]string{
	"go": golang,
}

var language = flag.String("l", "go", "pass the name of a known language to get result in that language")
var outputFile = flag.String("o", "", "output file (default is stdout)")
var gopkg = flag.String("gopkg", "main", "golang package code should be generated for")

var buildSuccess = true

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
	fp.Close()
	el := errorListener{0}
	input := antlr.NewInputStream(string(buffer))
	lexer := parser.Newwcllex(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&el)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.Newwcl(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(&el)
	builder := parser.WclBuildListener{}
	prog := p.Program()
	antlr.ParseTreeWalkerDefault.Walk(&builder, prog)
	if !buildSuccess {
		log.Fatalf("failed due to syntax errors")
	}
	// create a context for this generate
	t, ok := langToTempl[*language]
	if !ok {
		log.Fatalf("unable to find a template for language '%s'", *language)
	}
	ctx := newGenerateContext(t)
	if prog.GetP() == nil {
		log.Fatalf("no program")
	}
	ctx.program = prog.GetP()
	ctx.templateName = t
	ctx.global["import"] = prog.GetP().ImportSection
	ctx.global["text"] = prog.GetP().TextSection
	golang := make(map[string]string)
	ctx.global["golang"] = golang
	golang["package"] = *gopkg
	// deal with output file
	dir, err := os.MkdirTemp(os.TempDir(), "wcl*")
	if err != nil {
		log.Fatalf("unable to create temp dir: %v", err)
	}
	defer func() {
		log.Printf("cleaning up temp dir %s", dir)
		//os.RemoveAll(dir) // clean up
	}()
	file := filepath.Join(dir, "output_program.go")
	fp, err = os.Create(file)
	if err != nil {
		log.Fatalf("unable to create output file: %v", err)
	}
	err = runTemplate(ctx, fp)
	if err != nil {
		log.Fatalf("error trying to execute template %s: %v", ctx.templateName, err)
	}

	cmd := exec.Command("gofmt", file)
	var outFp io.Writer
	if *outputFile != "" {
		outFp, err = os.Create(*outputFile)
		if err != nil {
			log.Fatalf("unable to create output file %s: %v", *outputFile, err)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if outFp != nil {
		cmd.Stdout = outFp
	}
	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to run gofmt: %v, %v", err, outFp)
	}
}
