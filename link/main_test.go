package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestExtractText(t *testing.T) {
	if s := extractText(nil); s != "" {
		t.Error("expected emptry string, got: ", s)
	}

	node := &html.Node{
		Type: html.TextNode,
		Data: "test",
	}
	if s := extractText(node); s != "test" {
		t.Error("expected 'test', got: ", s)
	}

	nested := &html.Node{
		Type: html.TextNode,
		Data: "nes",
		NextSibling: &html.Node{
			Type: html.ElementNode,
			Data: "i",
			FirstChild: &html.Node{
				Type: html.TextNode,
				Data: "ted",
			},
		},
	}

	if s := extractText(nested); s != "nested" {
		t.Error("expected 'nested', got: ", s)
	}
}

func TestExtractHRef(t *testing.T) {
	if h := extractHRef(nil); h != "" {
		t.Error("expected '', got: ", h)
	}

	missing := &html.Node{
		Type: html.ElementNode,
		Data: "a",
	}
	if h := extractHRef(missing); h != "" {
		t.Error("expected '', got: ", h)
	}

	node := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{
				Key: "href",
				Val: "/1234",
			},
		},
	}
	if h := extractHRef(node); h != "/1234" {
		t.Error("expected '/1234', got: ", h)
	}
}

func TestExtractAnchors(t *testing.T) {
	if ls := extractAll(nil); len(ls) > 0 {
		t.Error("expected no links, got ", ls)
	}

	fixtures := []struct {
		name     string
		html     string
		expected []Link
	}{
		{
			name: "empty",
			html: ``,
		},
		{
			name: "simple",
			html: `<a href="/123">test</a>`,
			expected: []Link{
				{Text: "test", HRef: "/123"},
			},
		},
		{
			name: "nested",
			html: `
			<html>
			  <body>
				<a href="/123">test</a>
			  </body>
			</html>
			`,
			expected: []Link{
				{Text: "test", HRef: "/123"},
			},
		},
		{
			name: "siblings",
			html: `
			<html>
			  <body>
			  	<div>
				  <a href="/123">test</a>
				  <a href="/456">other</a>
				</div>
				<a href="/789">outer</a>
			  </body>
			</html>
			`,
			expected: []Link{
				{Text: "test", HRef: "/123"},
				{Text: "other", HRef: "/456"},
				{Text: "outer", HRef: "/789"},
			},
		},
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			assert := assert.New(t)
			r, err := html.Parse(strings.NewReader(f.html))
			if err != nil {
				t.Fatal("unexpected html parse error: ", err)
			}
			res := extractAll(r)
			assert.ElementsMatch(f.expected, res)
		})
	}
}

func extractAll(root *html.Node) []Link {
	ch := make(chan Link)
	go func() {
		defer close(ch)
		extractAnchors(root, ch)
	}()
	var res []Link
	for l := range ch {
		res = append(res, l)
	}
	return res
}
