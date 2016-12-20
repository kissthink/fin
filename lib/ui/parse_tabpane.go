package ui

import (
	"github.com/gizak/termui/extra"

	"golang.org/x/net/html"
)

func (p *Page) parseBodyTabpane(parentNode *Node, htmlNode *html.Node) (ret *Node, isFallthrough bool) {
	ret = p.newNode(htmlNode)
	parentNode.addChild(ret)
	isFallthrough = true

	ret.InitNodeTabpane()

	return
}

func (p *Page) parseBodyTabpaneTab(parentNode *Node, htmlNode *html.Node) (ret *Node, isFallthrough bool) {
	ret = p.newNode(htmlNode)
	parentNode.addChild(ret)
	isFallthrough = true

	ret.InitNodeTabpaneTab()

	for _, attr := range htmlNode.Attr {
		switch attr.Key {
		case "label":
			uiBuffer := ret.uiBuffer.(*extra.Tab)
			uiBuffer.SetLabel(attr.Val)
		}
	}

	ret.isShouldTermuiRenderChild = true

	return
}
