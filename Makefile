all: build/wcl

test:syntaxtest

syntaxtest: parser/wcl_parser.go parser/wcllex_lexer.go syntaxtest/main.go
	go run cmd/syntaxtest/main.go

parser/wcl_parser.go: parser/wcl.g4 
	cd parser;./generate.sh

parser/wcllex_lexer.go: parser/wcllex.g4 
	cd parser;./generate.sh

build/wcl: driver/*.go driver/template/*.tmpl parser/wcl_parser.go parser/wcllex_lexer.go 
	go build -o build/wcl github.com/iansmith/parigot-ui/cmd/wcl
