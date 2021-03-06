package ui

import "golang.org/x/net/html"

func (p *Page) parseHeadTitle(parentNode *Node, htmlNode *html.Node) (ret *Node, isFallthrough bool) {
	ret = nil
	isFallthrough = true

	if nil != htmlNode.FirstChild {
		p.Title = htmlNode.FirstChild.Data
	}
	return
}

func (p *Page) parseBody(parentNode *Node, htmlNode *html.Node) (ret *Node, isFallthrough bool) {
	ret = p.newNode(htmlNode)
	parentNode.addChild(ret)
	isFallthrough = true

	ret.ColorFg = ColorBodyDefaultColorFg

	ret.InitNodeBody()

	return
}

func (p *Page) parseBodyDiv(parentNode *Node, htmlNode *html.Node) (ret *Node, isFallthrough bool) {
	ret = p.newNode(htmlNode)
	parentNode.addChild(ret)
	isFallthrough = true

	ret.InitNodeDiv()

	return
}
