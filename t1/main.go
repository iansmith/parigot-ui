package t1

import (
	"fmt"
	"syscall/js"

	"github.com/iansmith/parigot-ui/static/uilib"
)

var textExample1 = uilib.NewText("example1", "This is the time for all good "+
	"men to come to the aid of their country.")

func programWrapper() any {
	fn := func(_ js.Value, _ []js.Value) any /* js error map*/ {
		domSvc, err := uilib.NewDOMService()
		if err != nil {
			return uilib.ErrorToJSMap(err)
		}
		_, err = domSvc.UpdateText(con.ParagraphLoc, textExample1)
		if err != nil {
			return uilib.ErrorToJSMap(err)
		}
		return js.Undefined()
	}
	return js.FuncOf(fn)
}

// This is the programs entry point.
func Main() {
	fmt.Printf("reached Main()")
	js.Global().Set("t1main", programWrapper())
	js.Global().Call("t1main")
	<-make(chan struct{})
}
