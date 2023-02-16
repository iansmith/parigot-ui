package t1

import "github.com/iansmith/parigot-ui/static/uilib"

type const_ struct {
	ParagraphLoc uilib.ElementId
}

var con = const_{
	ParagraphLoc: uilib.ElementId("paraLoc"),
}
