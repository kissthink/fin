package ui

func (p *Page) normalLayoutNodeBlock(node *Node) (isFallthrough bool) {
	isFallthrough = true
	if nil == node.UIBlock {
		return
	}

	if false == node.isSettedPositionX {
		node.UIBlock.X = p.layoutingX
	}

	if false == node.isSettedPositionY {
		node.UIBlock.Y = p.layoutingY
	}

	if nil != node.UIBlock {
		node.UIBlock.Align()
	}

	if "absolute" == node.Position {
	} else {
		p.layoutingY = node.UIBlock.Y + node.UIBlock.Height
	}

	return
}
