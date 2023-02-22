package uilib

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"syscall/js"

	dommsg "github.com/iansmith/parigot-ui/g/msg/dom/v1"
)

type DOMService struct {
	doc js.Value
	//textNameToElementId map[string]string
	elemCache     map[string]*dommsg.Element
	anonCount     int
	elemIdToValue map[string]js.Value
}

func NewDOMService() (*DOMService, error) {
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		return nil, fmt.Errorf("unable to get document object")
	}
	return &DOMService{
		doc: jsDoc,
		//textNameToElementId: make(map[string]string),
		elemCache:     make(map[string]*dommsg.Element),
		elemIdToValue: make(map[string]js.Value),
	}, nil
}

func ErrorToJSMap(err error) map[string]any {
	result := map[string]any{
		"error": "Unable to get document object",
	}
	return result
}

func (d *DOMService) ElemById(id string) (*dommsg.Element, js.Value, error) {
	e, ok := d.elemCache[id]
	if ok {
		return e, d.elemIdToValue[id], nil
	}
	elem := d.doc.Call("getElementById", string(id))
	if !elem.Truthy() {
		return nil, js.Null(), fmt.Errorf("unable to get find element by id '%s'", id)
	}
	result := &dommsg.Element{
		Tag: &dommsg.Tag{
			Id:       elem.Get("id").String(),
			CssClass: strings.Split(elem.Get("className").String(), " "),
			Name:     elem.Get("tagName").String(),
		},
		Text: elem.Get("innerHTML").String(),
	}
	// inner:=value.Get("innerHTML")
	// if inner!="" {
	// 	child:=value.Get("childElementCount")
	// 	children:=value.Get(children)
	// 	for i:=0; i<jsToIntt(child); i++ {
	// 		newElement:=dommsg.Element{
	// 			Tag:
	// 		}
	// 	}
	// }

	d.elemCache[result.Tag.Id] = result
	d.elemIdToValue[result.Tag.Id] = elem
	return result, elem, nil
}

// func (d *DOMService) UpdateText(elementId string, text *dommsg.TextBlob) (*Element, error) {
// 	var elem *dommsg.Element
// 	var err error
// 	n := string(text.Name())
// 	stringVersion, ok := d.textNameToElementId[n]
// 	if !ok {
// 		elem, err = d.ElemById(elementId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		d.textNameToElementId[n] = string(elem.Id())
// 	} else {
// 		// we have the string version and we assume it's in cache, because we
// 		// insert it just above in the !ok case
// 		elem, err = d.ElemById(stringVersion)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	elem.Value().Set("innerHTML", text.Content())
// 	return elem, nil
// }

// SetChild removes all the existing children of the given element (by Id)
// and then sets the children to be the children provided.  An empty set of
// children is ok as the second paramater.
func (d *DOMService) SetChild(elementId string, child []*dommsg.Element) error {
	_, value, err := d.ElemById(elementId)
	if err != nil {
		return err
	}
	for value.Get("firstChild").Truthy() {
		value.Call("RemoveChild", value.Get("firstChild"))
	}
	buf := &bytes.Buffer{}
	for _, e := range child {
		t := e.GetTag()
		if t.GetId() == "" {
			t.Id = fmt.Sprintf("_anon_elem_%08d", d.anonCount)
			d.anonCount++
		}
		buf.WriteString(toHtml(e))
	}
	print("object I've got handle to is ", value.Get("tagName").String(), "\n")
	print("setting inner html of ", elementId, " to ", buf.String(), "\n")
	value.Set("innerHTML", buf.String())
	return nil
}

func toHtml(e *dommsg.Element) string {
	tag := e.GetTag()
	allClass := &bytes.Buffer{}
	for _, clazz := range tag.GetCssClass() {
		allClass.WriteString(clazz + " ")
	}
	t := fmt.Sprintf("<%s id=\"%s\" class=\"%s\">", tag.GetName(), tag.GetId(), allClass)
	print("t is ", t, "\n")
	end := fmt.Sprintf("</%s>", tag.GetName())
	inner := e.GetText()
	if inner == "" {
		child := &bytes.Buffer{}
		for _, c := range e.GetChild() {
			child.WriteString(toHtml(c))
		}
		inner = child.String()
	}

	result := fmt.Sprintf("%s%s%s", t, inner, end)
	print("result of toHtml ", result, "\n")
	return result
}

func jsToInt(value js.Value) int {
	f := value.Float()
	f = math.Floor(f)
	return int(f)
}
