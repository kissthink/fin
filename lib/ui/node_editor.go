package ui

import (
	"inn/ui/editor"

	"github.com/gizak/termui"
)

type NodeEditor struct {
	*Node
	*editor.Editor
}

func (p *Node) InitNodeEditor() *NodeEditor {
	nodeEditor := new(NodeEditor)
	nodeEditor.Node = p
	nodeEditor.Editor = editor.NewEditor()
	nodeEditor.Border = false
	nodeEditor.BorderFg = COLOR_DEFAULT_BORDERFG
	p.Data = nodeEditor
	p.KeyPress = nodeEditor.KeyPress
	p.FocusMode = nodeEditor.FocusMode
	p.UnFocusMode = nodeEditor.UnFocusMode
	return nodeEditor
}

func (p *NodeEditor) KeyPress(e termui.Event) {
	keyStr := e.Data.(termui.EvtKbd).KeyStr
	if "<escape>" == keyStr {
		p.Node.quitActiveMode()
		return
	}

	p.Editor.Write(keyStr)
	termui.Render(p.Node.uiBuffer.(termui.Bufferer))
}

func (p *NodeEditor) FocusMode() {
	p.Node.uiBuffer.(*editor.Editor).Border = true
	p.Node.uiBuffer.(*editor.Editor).BorderFg = COLOR_FOCUSMODE_BORDERFG
	termui.Render(p.Node.uiBuffer.(termui.Bufferer))
}

func (p *NodeEditor) UnFocusMode() {
	p.Node.uiBuffer.(*editor.Editor).Border = p.Node.Border
	p.Node.uiBuffer.(*editor.Editor).BorderFg = p.Node.BorderFg
	termui.Render(p.Node.uiBuffer.(termui.Bufferer))
}
