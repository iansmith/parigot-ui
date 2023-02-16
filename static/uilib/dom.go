package uilib

import (
	"fmt"
	"syscall/js"
)

type ElementId string

type Element struct {
	value js.Value
	id    ElementId
}

func (e *Element) Value() js.Value {
	return e.value
}
func (e *Element) Id() ElementId {
	return e.id
}

type DOMService struct {
	doc                 js.Value
	textNameToElementId map[string]string
	elemCache           map[string]*Element
}

type TextName string

type Text struct {
	content string
	name    TextName
}

func (t *Text) Name() TextName {
	return t.name
}
func (t *Text) Content() string {
	return t.content
}
func NewText(name, content string) *Text {
	return &Text{content: content, name: TextName(name)}
}

func NewDOMService() (*DOMService, error) {
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		return nil, fmt.Errorf("unable to get document object")
	}
	return &DOMService{
		doc:                 jsDoc,
		textNameToElementId: make(map[string]string),
		elemCache:           make(map[string]*Element),
	}, nil
}

func ErrorToJSMap(err error) map[string]any {
	result := map[string]any{
		"error": "Unable to get document object",
	}
	return result
}

func (d *DOMService) ElemById(id ElementId) (*Element, error) {
	e, ok := d.elemCache[string(id)]
	if ok {
		return e, nil
	}
	elem := d.doc.Call("getElementById", string(id))
	if !elem.Truthy() {
		return nil, fmt.Errorf("unable to get find element by id '%s'", id)
	}
	result := &Element{
		value: elem,
		id:    ElementId(id),
	}
	d.elemCache[string(id)] = result
	return result, nil
}

func (d *DOMService) UpdateText(elementId ElementId, text *Text) (*Element, error) {
	var elem *Element
	var err error
	n := string(text.Name())
	stringVersion, ok := d.textNameToElementId[n]
	if !ok {
		elem, err = d.ElemById(elementId)
		if err != nil {
			return nil, err
		}
		d.textNameToElementId[n] = string(elem.Id())
	} else {
		// we have the string version and we assume it's in cache, because we
		// insert it just above in the !ok case
		elem, err = d.ElemById(ElementId(stringVersion))
		if err != nil {
			return nil, err
		}
	}
	elem.Value().Set("innerHTML", text.Content())
	return elem, nil
}
