lexer grammar wcllex;

// Lexer
Text: 'text';
CSS: 'css';
Id: IdentFirst (IdentAfter)*;

// consistent def of Ident
fragment IdentFirst: ('a' .. 'z' | 'A' .. 'Z' | '.' | '_');

fragment IdentAfter: (
		'a' .. 'z'
		| 'A' .. 'Z'
		| '.'
		| '_'
		| Digit
	);

fragment Digit: '0' ..'9';
DoubleLeftCurly: '{{' {print("push(content)\n")} -> pushMode(CONTENT);
LCurly: '{';
RCurly: '}';
LParen: '(';
RParen: ')';
Dollar: '$';
Comma: ',';
PoundComment: '#' .+? [\n\r] -> skip;
DoubleSlashComment: '//' .+? [\n\r] -> skip;
Whitespace: [ \n\r\t\u000B\u000C\u0000]+ -> skip;

mode CONTENT;
ContentRawText: ~[${}]+;
ContentDollar: '$' -> pushMode(VAR);
ContentLCurly: '{' -> pushMode(UNINTERPRETED);
DoubleRightCurly: '}}' -> popMode;
//ContentRCurly: '}' -> popMode;

mode UNINTERPRETED;
UninterpRawText: ~[{}]+;
UninterpLCurly: '{' -> pushMode(UNINTERPRETED);
UninterpRCurly: '}' -> popMode;
UninterpDollar: '$' -> pushMode(VAR);

mode VAR;
VarLCurly: '{';
VarRCurly: '}' -> popMode;
VarId: IdentFirst (IdentAfter)*;