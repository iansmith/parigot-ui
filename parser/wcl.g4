parser grammar wcl;
options {
	tokenVocab = wcllex;
}

program
	returns[*Program p]
	@init {
		$p=NewProgram()	
	}:
	i = import_section? { $p.ImportSection = localctx.GetI().GetSection()
	} t = text_section? { $p.TextSection = localctx.GetT().GetSection()
		} css_section? EOF;

import_section
	returns[*ImportSectionNode section]
	@init {
		$section = NewImportSectionNode()
	}:
	Import u = uninterp {
		$section.Text = localctx.GetU().GetItem()[0];
	};

text_section
	returns[*TextSectionNode section]
	@init {
		$section = NewTextSectionNode()
	}:
	Text (
		d = text_decl {$section.Func=append($section.Func,localctx.GetD().GetF())}
	)*;

text_decl
	returns[*TextFuncNode f]:
	i = Id param_spec? t = text_top { 
		$f=NewTextFuncNode(localctx.GetI().GetText(),localctx.GetT().GetItem())};

text_top
	returns[[]TextItem item]:
	DoubleLeftCurly (
		content = text_content {$item = localctx.GetContent().GetItem()}
		|
	) DoubleRightCurly;

text_content
	returns[[]TextItem item]
	@init {
		$item = []TextItem{}
	}:
	(
		c = ContentRawText { $item = append($item, NewTextConstant(localctx.GetC().GetText()))}
		| v = var_subs { }
		| u = uninterp { $item = append($item, localctx.GetU().GetItem()...)}
	)*;

var_subs: ContentDollar sub;

sub: VarLCurly VarId VarRCurly;

uninterp
	returns[[]TextItem item]
	@init {
		$item=[]TextItem{}
	}:
	// jump from content mode to Uninterp mode
	ContentLCurly (
		c = UninterpRawText { $item = append($item, NewTextConstant(localctx.GetC().GetText()))}
		| u = nestedUninterp { $item = append($item, localctx.GetU().GetItem()...)}
		| UninterpDollar sub
	)+ UninterpRCurly;

nestedUninterp
	returns[[]TextItem item]
	@init {
		$item=[]TextItem{}
	}:
	UninterpLCurly {} (
		c = UninterpRawText { $item = append($item, NewTextConstant(localctx.GetC().GetText()))}
		| UninterpDollar sub
		| u = nestedUninterp { $item = append($item, localctx.GetU().GetItem()...)}
	)+ UninterpRCurly;

uninterpVar: UninterpDollar Id;

param_spec: LParen param_seq RParen;

param_seq: (Id Comma)+? Id |;

param_rest: Comma param_seq;

css_section: CSS css_file*;

css_file: Id LCurly css_decl* RCurly;

css_decl: Id;
