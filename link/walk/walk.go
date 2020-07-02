package walk

import "golang.org/x/net/html"

type (
	Link struct {
		Text string
		HRef string
	}
)

func ExtractText(node *html.Node) string {
	if node == nil {
		return ""
	}

	if node.Type == html.TextNode {
		return node.Data + ExtractText(node.NextSibling)
	}

	return ExtractText(node.FirstChild) + ExtractText(node.NextSibling)
}

func ExtractHRef(node *html.Node) string {
	if node == nil {
		return ""
	}

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func ExtractAnchors(root *html.Node, links chan<- Link) {
	if root == nil {
		return
	}

	if root.Type == html.ElementNode && root.Data == "a" {
		links <- Link{
			Text: ExtractText(root.FirstChild),
			HRef: ExtractHRef(root),
		}
	}

	ExtractAnchors(root.FirstChild, links)
	ExtractAnchors(root.NextSibling, links)
}

func ExtractSync(root *html.Node) []Link {
	ch := make(chan Link)
	go func() {
		defer close(ch)
		ExtractAnchors(root, ch)
	}()
	var res []Link
	for l := range ch {
		res = append(res, l)
	}
	return res
}
