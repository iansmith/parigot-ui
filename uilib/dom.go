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

	d.elemCache[result.Tag.Id] = result
	d.elemIdToValue[result.Tag.Id] = elem
	return result, elem, nil
}

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
