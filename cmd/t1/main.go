package main

import (
	"strings"

	dommsg "github.com/iansmith/parigot-ui/g/msg/dom/v1"
	"github.com/iansmith/parigot-ui/uilib"
)

const parentId = "paraLoc"

func main() {
	// xxx should be using locate
	svc, err := uilib.NewDOMService()
	if err != nil {
		panic("unable to get DOMService: " + err.Error())
	}
	elem := &dommsg.Element{}
	elem.Tag = &dommsg.Tag{
		Name:     "span",
		Id:       "myid",
		CssClass: strings.Split("foo bar baz", " "),
	}
	elem.Text = "This is a test."
	svc.SetChild(parentId, []*dommsg.Element{elem})
}
