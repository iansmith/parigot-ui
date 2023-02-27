parser grammar wcl;
options {
	tokenVocab = wcllex;
}

program
	returns[*ProgramNode p]:
	import_section?  
	extern?
	global?
	text_section? 
	css_section?     
	doc_section?
	EOF   
	;

global
	returns [[]*PFormal g]:
	Global param_spec
	;

extern
	returns [[]string e]:
	Extern LParen Id* RParen
	;

import_section
	returns[*ImportSectionNode section]:
	Import LCurly uninterp 
	;

text_section
	returns[*TextSectionNode section]:
	Text (
		text_func 
	)*;

text_func
	returns[*TextFuncNode f]:
	i = Id param_spec? text_func_local? pre_code? text_top post_code?
	;

pre_code 
	returns [[]TextItem item]:
	Pre LCurly uninterp
	;

post_code 
	returns [[]TextItem item]:
	Post LCurly uninterp
	;

text_func_local
	returns [[]*PFormal formal]:
	Local param_spec
	;

text_top
	returns[[]TextItem item]:
	BackTick 	(
		text_content
		|
	) ContentBackTick ;

text_content
	returns[[]TextItem item]:
	(
		text_content_inner    
	)*;

text_content_inner
	returns[[]TextItem item]:
		ContentRawText             	#RawText
		| var_subs   				#VarSub
;

var_subs
	returns [[]TextItem item]: 
	ContentDollar sub
	;

sub
	returns [TextItem item]: 
	VarId VarRCurly
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
	UninterpDollar VarId VarRCurly;

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

doc_section
	returns [*DocSectionNode section]: 
	Doc (doc_func)*;

doc_func
	returns [*DocFuncNode fn]:
	Id doc_func_formal pre_code? doc_func_local? doc_elem post_code?
	;

doc_func_local
	returns [[]*PFormal formal]:
	Local param_spec
	;

doc_func_formal
	returns [[]*PFormal formal]:
	param_spec
	|
	;

doc_tag
	returns [*DocTag tag]:
	LessThan Id 
	doc_id?
	doc_class?
	GreaterThan 
	;

doc_id
	returns [string s]:
	Hash Id
	;

doc_class
	returns [[]string clazz]:
	(Colon Id)+
	;

doc_elem
	returns [*DocElement elem]:
	doc_tag doc_elem_content?  # haveTag
	| doc_elem_child           # haveList
	;

doc_elem_content
	returns [*DocElement element]:
	doc_elem_text 
	| doc_elem_child 
	;

doc_elem_text
	returns [*FuncInvoc invoc]:
	func_invoc       #doc_elem_text_func_call
	| text_top       #doc_elem_text_anon
	;

doc_elem_child
	returns [*DocElement elem]:
	LParen (doc_elem)* RParen
	;

func_invoc
	returns [*FuncInvoc invoc]:
	Id LParen func_actual_seq RParen
	;

func_actual_seq
	returns [[]*FuncActual actual]:
	( a=func_actual (Comma b=func_actual)* )?
	;

func_actual 
	returns [*FuncActual actual]:
	Id
	| StringLit
	;