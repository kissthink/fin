package ui

import (
	"fin/ui/canvas"

	"github.com/gizak/termui"
)

type NodeCanvas struct {
	*Node
	*canvas.Canvas
}

func (p *Node) InitNodeCanvas() {
	nodeCanvas := new(NodeCanvas)
	nodeCanvas.Node = p
	nodeCanvas.Canvas = canvas.NewCanvas()

	p.Data = nodeCanvas
	p.KeyPress = nodeCanvas.KeyPress

	p.uiBuffer = nodeCanvas.Canvas
	p.UIBlock = &nodeCanvas.Canvas.Block
	p.Display = &p.UIBlock.Display

	p.isShouldCalculateWidth = true
	p.isShouldCalculateHeight = false
	p.UIBlock.Height = 10
	p.UIBlock.Border = true

	p.isWorkNode = true

	return
}

func (p *NodeCanvas) KeyPress(e termui.Event) {
	if len(p.Node.KeyPressHandlers) > 0 {
		for _, v := range p.Node.KeyPressHandlers {
			v.Args = append(v.Args, e)
			v.Handler(p.Node, v.Args...)
		}
	}
}

func (p *NodeCanvas) NodeDataFocusMode() {
	if false == p.Node.isCalledFocusMode && true == p.Node.UIBlock.Border {
		p.Node.isCalledFocusMode = true
		p.Node.tmpFocusModeBorder = p.Node.UIBlock.Border
		p.Node.tmpFocusModeBorderFg = p.Node.UIBlock.BorderFg
		p.Node.UIBlock.Border = true
		p.Node.UIBlock.BorderFg = COLOR_FOCUS_MODE_BORDERFG
		p.Node.uiRender()
	}
}

func (p *NodeCanvas) NodeDataUnFocusMode() {
	if true == p.Node.isCalledFocusMode && true == p.Node.UIBlock.Border {
		p.Node.isCalledFocusMode = false
		p.Node.UIBlock.Border = p.Node.tmpFocusModeBorder
		p.Node.UIBlock.BorderFg = p.Node.tmpFocusModeBorderFg
		p.Node.uiRender()
	}
}

func (p *NodeCanvas) NodeDataActiveMode() {
	if false == p.Node.isCalledActiveMode && true == p.Node.UIBlock.Border {
		p.Node.isCalledActiveMode = true
		p.Node.tmpActiveModeBorderFg = p.Node.UIBlock.BorderFg
		p.Node.UIBlock.BorderFg = COLOR_ACTIVE_MODE_BORDERFG
		p.Node.ResumeCursor()
	}
}

func (p *NodeCanvas) NodeDataUnActiveMode() {
	if true == p.Node.isCalledActiveMode && true == p.Node.UIBlock.Border {
		p.Node.isCalledActiveMode = false
		p.Node.UIBlock.BorderFg = p.Node.tmpActiveModeBorderFg
		p.Node.HideCursor()
	}
}
