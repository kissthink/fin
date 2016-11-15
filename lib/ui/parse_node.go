package ui

import (
	. "in/ui/utils"
	"strconv"

	"golang.org/x/net/html"
)

func (p *Node) ParseAttribute(attr []html.Attribute) (isUIChange, isNeedRerenderPage bool) {
	isUIChange = false
	isNeedRerenderPage = false

	if nil == p.UIBlock {
		return
	}
	p.UIBlock.BorderLabelFg = COLOR_DEFAULT_BORDER_LABEL_FG
	p.UIBlock.BorderFg = COLOR_DEFAULT_BORDER_FG

	for _, v := range attr {
		switch v.Key {
		case "paddingtop":
			isUIChange = true
			p.UIBlock.PaddingTop, _ = strconv.Atoi(v.Val)

		case "paddingbottom":
			isUIChange = true
			p.UIBlock.PaddingBottom, _ = strconv.Atoi(v.Val)

		case "paddingleft":
			isUIChange = true
			p.UIBlock.PaddingLeft, _ = strconv.Atoi(v.Val)

		case "paddingright":
			isUIChange = true
			p.UIBlock.PaddingRight, _ = strconv.Atoi(v.Val)

		case "borderlabelfg":
			isUIChange = true
			p.UIBlock.BorderLabelFg = ColorToTermuiAttribute(v.Val, COLOR_DEFAULT_BORDER_LABEL_FG)

		case "borderlabel":
			isUIChange = true
			p.UIBlock.BorderLabel = v.Val

		case "borderfg":
			isUIChange = true
			p.UIBlock.BorderFg = ColorToTermuiAttribute(v.Val, COLOR_DEFAULT_BORDER_FG)

		case "border":
			isUIChange = true
			p.UIBlock.Border = StringToBool(v.Val, p.UIBlock.Border)

		case "borderleft":
			isUIChange = true
			p.UIBlock.BorderLeft = StringToBool(v.Val, p.UIBlock.BorderLeft)

		case "borderright":
			isUIChange = true
			p.UIBlock.BorderRight = StringToBool(v.Val, p.UIBlock.BorderRight)

		case "bordertop":
			isUIChange = true
			p.UIBlock.BorderTop = StringToBool(v.Val, p.UIBlock.BorderTop)

		case "borderbottom":
			isUIChange = true
			p.UIBlock.BorderBottom = StringToBool(v.Val, p.UIBlock.BorderBottom)

		case "height":
			isUIChange = true
			isNeedRerenderPage = true
			p.UIBlock.Height, _ = strconv.Atoi(v.Val)
			if p.UIBlock.Height < 0 {
				p.UIBlock.Height = 0
			}
			p.isShouldCalculateHeight = false

		case "width":
			isUIChange = true
			isNeedRerenderPage = true
			p.UIBlock.Width, _ = strconv.Atoi(v.Val)
			if p.UIBlock.Width < 0 {
				p.UIBlock.Width = 0
			}
			p.isShouldCalculateWidth = false
		}
	}

	if nodeDataParseAttributer, ok := p.Data.(NodeDataParseAttributer); true == ok {
		_isUIChange, _isNeedRerenderPage := nodeDataParseAttributer.NodeDataParseAttribute(attr)
		if true == (isUIChange || _isUIChange) {
			isUIChange = true
		}
		if true == (isNeedRerenderPage || _isNeedRerenderPage) {
			isNeedRerenderPage = true
		}
	}

	return
}