lexer grammar wcllex;

// keywords
Text: '@text';
CSS: '@css';
Import: '@preamble';
Doc: '@doc';
Local: '@local';
Global: '@global';
Extern: '@extern';

//ids
//TypeId: (TypeStarter+)? IdentFirst (IdentAfter)*;
Id: (TypeStarter+)? IdentFirst (IdentAfter)*;

fragment TypeStarter: '[' | ']'|'*';

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

DoubleLeftCurly: '{{' -> pushMode(CONTENT);
LCurly: '{' -> pushMode(UNINTERPRETED);
RCurly: '}';
LParen: '(';
RParen: ')';
Dollar: '$';
Comma: ',';
LessThan: '<';
GreaterThan: '>';
Colon: ':';
Hash: '#';
StringLit: '"' ( Esc | ~[\\"] )* '"';
fragment Esc : '\\"' | '\\\\' ;

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

