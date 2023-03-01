all: build/wcl

.PHONY: syntaxtest
test: semfailtest semtest

.PHONY: syntaxtest
syntaxtest: parser/wcl_parser.go parser/wcllex_lexer.go syntaxtest/main.go
	go run cmd/syntaxtest/main.go

parser/wcl_parser.go: parser/wcl.g4 
	cd parser;./generate.sh

parser/wcllex_lexer.go: parser/wcllex.g4 
	cd parser;./generate.sh

css/css3_lexer.go: css/css3.g4
	cd css;./generate.sh

css/css3_parser.go: css/css3.g4
	cd css;./generate.sh

build/wcl: driver/*.go driver/template/*.tmpl parser/*.go parser/wcl_parser.go parser/wcllex_lexer.go \
	css/css3_lexer.go css/css3_parser.go css/*.go
	go build -o build/wcl github.com/iansmith/parigot-ui/cmd/wcl


.PHONY: parserclean
parserclean:
	rm -f parser/wcl.interp parser/wcl.tokens \
	parser/wcl_base_listener.go parser/wcl_listener.go \
	parser/wcl_parser.go parser/wcllex.interp \
	parser/wcllex.tokens parser/wcllex_lexer.go \
	parser/wcl_visitor.go parser/wcl_base_visitor.go

semfailtest: build/wcl
	build/wcl -invert testdata/fail_dupparam.wcl 
	build/wcl -invert testdata/fail_duplocal.wcl 
	build/wcl -invert testdata/fail_duptextname.wcl 
	build/wcl -invert testdata/fail_duptextnameparam.wcl 
	build/wcl -invert testdata/fail_dupdocfunc.wcl 
	build/wcl -invert testdata/fail_dotindocname.wcl 
	build/wcl -invert testdata/fail_dotintextname.wcl 
	build/wcl -invert testdata/fail_dupnamefuncdoctext.wcl
	build/wcl -invert testdata/fail_dotbadformal.wcl
	build/wcl -invert testdata/fail_dotbadlocal.wcl
	build/wcl -invert testdata/fail_badtag.wcl
	build/wcl -invert testdata/fail_baddocvar.wcl
	build/wcl -invert testdata/fail_badtextvarpre.wcl
	build/wcl -invert testdata/fail_badtextvarpost.wcl
	build/wcl -invert testdata/fail_conflictlocalparamtext.wcl
	build/wcl -invert testdata/fail_conflictlocalnametext.wcl
	@echo "PASS"
	

semtest: build/wcl
	build/wcl -o /dev/null testdata/textfunc_test.wcl
	build/wcl -o /dev/null testdata/docfunc_test.wcl
	@echo PASS

static/t1.wasm: cmd/t1/main.go uilib/*.go
	GOOS=js GOARCH=wasm go build -o static/t1.wasm cmd/t1/main.go