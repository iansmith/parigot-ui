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

build/wcl: driver/*.go driver/template/*.tmpl parser/*.go
	go build -o build/wcl github.com/iansmith/parigot-ui/cmd/wcl

.PHONY: comptest
comptest: build/wcl
	build/wcl testdata/textfunc_test.wcl 