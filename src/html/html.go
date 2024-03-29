// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/html/html.go
// Go representations of HTML objects

package html

import "fmt"

// Element is the basic interface implemented by all HTML elements.
type Element interface {
	Render() string
}

// HTML is a <html> tag.
type HTML struct {
	*Head
	*Body
}

var _ Element = &HTML{}

func (h *HTML) Render() string {
	str := "<html>"
	if h.Head != nil {
		str += h.Head.Render()
	}
	if h.Body != nil {
		str += h.Body.Render()
	}
	str += "</html>"
	return str
}

// Head is a <head> tag.
type Head struct {
	Title string
}

var _ Element = &HTML{}

func (h *Head) Render() string {
	return fmt.Sprintf("<head><title>%s</title></head>", h.Title)
}

// Body is a <body> tag.
type Body struct {
	children []Element
}

var _ Element = &Body{}

func (b *Body) Insert(e Element) {
	b.children = append(b.children, e)
}

func (b *Body) Render() string {
	str := "<body>"
	for _, child := range b.children {
		str += child.Render()
	}
	str += "</body>"
	return str
}

// Heading1 is a <h1> tag.
type Heading1 struct {
	heading string
}

var _ Element = &Heading1{}

func NewHeading1(heading string) *Heading1 {
	return &Heading1{heading}
}

func (h1 *Heading1) Render() string {
	return fmt.Sprintf("<h1>%s</h1>", h1.heading)
}

// Heading2 is a <h2> tag.
type Heading2 struct {
	heading string
}

var _ Element = &Heading2{}

func NewHeading2(heading string) *Heading2 {
	return &Heading2{heading}
}

func (h2 *Heading2) Render() string {
	return fmt.Sprintf("<h2>%s</h2>", h2.heading)
}

// Paragraph is a <p> tag.
type Paragraph struct {
	children []Element
}

var _ Element = &Paragraph{}

func NewParagraph() *Paragraph {
	return &Paragraph{}
}

func (b *Paragraph) Insert(e Element) {
	b.children = append(b.children, e)
}

func (b *Paragraph) Render() string {
	str := "<p>"
	for _, child := range b.children {
		str += child.Render()
	}
	str += "</p>"
	return str
}

// Anchor is an <a> tag.
type Anchor struct {
	href string
	text string
}

var _ Element = &Anchor{}

func NewAnchor(href, text string) *Anchor {
	return &Anchor{href, text}
}

func (a *Anchor) Render() string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, a.href, a.text)
}

// Lists

// List is implemented by OrderedList and UnorderedList.
type List interface {
	Element
	Insert(li *ListItem)
}

// ListBase provides a partial implementation of List, to embed in List types.
type ListBase struct {
	items []ListItem
}

func (l *ListBase) Insert(li *ListItem) {
	l.items = append(l.items, *li)
}

// OrderedList is a <ol> tag.
type OrderedList struct {
	ListBase
}

var _ List = &OrderedList{}

func NewOrderedList() *OrderedList {
	return &OrderedList{}
}

func (ol *OrderedList) Render() string {
	str := "<ol>"
	for _, item := range ol.items {
		str += item.Render()
	}
	str += "</ol>"
	return str
}

// UnorderedList is a <ul> tag.
type UnorderedList struct {
	ListBase
}

var _ List = &UnorderedList{}

func NewUnorderedList() *UnorderedList {
	return &UnorderedList{}
}

func (ul *UnorderedList) Render() string {
	str := "<ul>"
	for _, item := range ul.items {
		str += item.Render()
	}
	str += "</ul>"
	return str
}

// ListItem is an <li> tag.
type ListItem struct {
	children []Element
	value    string
}

func NewListItem(children ...Element) *ListItem {
	return &ListItem{children, ""}
}

func (li *ListItem) SetValue(val string) {
	li.value = val
}

func (li *ListItem) Render() string {
	var str string
	if li.value == "" {
		str = "<li>"
	} else {
		str = fmt.Sprintf(`<li value="%s">`, li.value)
	}

	for _, child := range li.children {
		str += child.Render()
	}
	str += "</li>"
	return str
}

// String is a string as an HTML element (e.g. the inner HTML for a tag).
type String string

var _ Element = String("")

func (s String) Render() string {
	return string(s)
}
