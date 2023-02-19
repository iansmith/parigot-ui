parser grammar wcl;
options {
	tokenVocab = wcllex;
}

program
	returns[*ProgramNode p]:
	import_section?  
	text_section? 
	css_section?     
	EOF   
	;

import_section
	returns[*ImportSectionNode section]:
	Import LCurly uninterp 
	;

text_section
	returns[*TextSectionNode section]:
	Text (
		text_decl 
	)*;

text_decl
	returns[*TextFuncNode f]:
	i = Id param_spec? text_top 
	;

text_top
	returns[[]TextItem item]:
	DoubleLeftCurly (
		text_content 
		|
	) DoubleRightCurly;

text_content
	returns[[]TextItem item]:
	(
		text_content_inner    
	)*;

text_content_inner
	returns[[]TextItem item]:
		ContentRawText             	#RawText
		| var_subs   				#VarSub
		| ContentLCurly uninterp   	#Unint
;

var_subs
	returns [[]TextItem item]: 
	ContentDollar sub
	;

sub
	returns [TextItem item]: 
	VarLCurly VarId VarRCurly
	;

uninterp
	returns[[]TextItem item]:
	(
		uninterp_inner
	)+ UninterpRCurly;

uninterp_inner 
	returns [[]TextItem Item]:
	UninterpRawText #UninterpRawText
	| UninterpLCurly uninterp  #UninterpNested
	| uninterp_var #UninterpVar
;

uninterp_var
	returns[[]TextItem item]: 
	UninterpDollar Id;

param_spec
	returns[[]*PFormal formal]: 
	LParen (param_pair)* RParen;

param_pair
	returns[[]*PFormal formal]:
	n=Id t=Id Comma    	#Pair
	| n=Id t=Id         #Last
	;

css_section: CSS css_file*;

css_file: Id LCurly css_decl* RCurly;

css_decl: Id;
