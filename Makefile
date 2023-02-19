all: build/wcl

.PHONY: syntaxtest
test:syntaxtest comptest

.PHONY: syntaxtest
syntaxtest: parser/wcl_parser.go parser/wcllex_lexer.go syntaxtest/main.go
	go run cmd/syntaxtest/main.go

parser/wcl_parser.go: parser/wcl.g4 
	cd parser;./generate.sh

parser/wcllex_lexer.go: parser/wcllex.g4 
	cd parser;./generate.sh

build/wcl: driver/*.go driver/template/*.tmpl parser/*.go parser/wcl_parser.go parser/wcllex_lexer.go
	go build -o build/wcl github.com/iansmith/parigot-ui/cmd/wcl

.PHONY: comptest
comptest: build/wcl
	build/wcl testdata/textfunc_test.wcl 

.PHONY: parserclean
parserclean:
	rm -f parser/wcl.interp parser/wcl.tokens \
	parser/wcl_base_listener.go parser/wcl_listener.go \
	parser/wcl_parser.go parser/wcllex.interp \
	parser/wcllex.tokens parser/wcllex_lexer.go \
	parser/wcl_visitor.go parser/wcl_base_visitor.go

semfailtest: build/wcl
	build/wcl -invert testdata/fail_dupparam.wcl 
	build/wcl -invert testdata/fail_duptextname.wcl 
	build/wcl -invert testdata/fail_duptextnameparam.wcl 

semtest: build/wcl
	build/wcl -o /dev/null testdata/textfunc_test.wcl
