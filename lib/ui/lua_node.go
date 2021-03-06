package ui

import (
	"golang.org/x/net/html"

	lua "github.com/yuin/gopher-lua"
)

func (p *Script) _getNodePointerFromUserData(L *lua.LState, lu *lua.LUserData) *Node {
	if nil == lu || nil == lu.Value {
		return nil
	}

	node, ok := lu.Value.(*Node)
	if false == ok || nil == node {
		return nil
	}

	return node
}

func (p *Script) luaFuncGetNodePointer(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	nodeID := L.ToString(1)

	var (
		node *Node
		ok   bool
	)

	node, ok = p.page.IDToNodeMap[nodeID]

	if true == ok {
		luaNode := L.NewUserData()
		luaNode.Value = node
		L.Push(luaNode)
	} else {
		L.Push(lua.LNil)
	}

	return 1
}

func (p *Script) luaFuncNodeWidth(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	L.Push(lua.LNumber(node.UIBlock.Width))
	return 1
}

func (p *Script) luaFuncNodeHeight(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	L.Push(lua.LNumber(node.UIBlock.Height))
	return 1
}

func (p *Script) luaFuncNodeInnerAreaWidth(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	L.Push(lua.LNumber(node.UIBlock.InnerArea.Dx()))
	return 1
}

func (p *Script) luaFuncNodeInnerAreaHeight(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	L.Push(lua.LNumber(node.UIBlock.InnerArea.Dy()))
	return 1
}

func (p *Script) luaFuncNodeGetAttribute(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}
	if attr, ok := node.HTMLAttribute[L.ToString(2)]; true == ok {
		L.Push(lua.LString(attr.Val))
		return 1
	}
	return 0
}

func (p *Script) luaFuncNodeSetAttribute(L *lua.LState) int {
	if L.GetTop() < 3 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}
	isUIChange, isNeedReRenderPage := node.ParseAttribute([]html.Attribute{
		html.Attribute{Key: L.ToString(2), Val: L.ToString(3)},
	})
	if false == isUIChange {
		return 0
	}
	if true == isNeedReRenderPage {
		p.page.ReRender()
	} else {
		node.UIRender()
	}
	return 0
}

func (p *Script) luaFuncNodeSetActive(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}
	p.page.SetActiveNode(node)
	node.UIRender()
	return 0
}

func (p *Script) luaFuncNodeGetHTMLData(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	L.Push(lua.LString(node.HTMLData))
	return 1
}

func (p *Script) luaFuncNodeSetValue(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	text := L.ToString(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	if nodeDataSetValueer, ok := node.Data.(NodeDataSetValueer); true == ok {
		nodeDataSetValueer.NodeDataSetValue(text)
	}

	return 0
}

func (p *Script) luaFuncNodeGetValue(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	nodeDataGetValuer, ok := node.Data.(NodeDataGetValuer)
	if false == ok {
		return 0
	}

	ret, ok := nodeDataGetValuer.NodeDataGetValue()
	if false == ok {
		return 0
	}

	L.Push(lua.LString(ret))
	return 1
}

func (p *Script) luaFuncNodeSetCursor(L *lua.LState) int {
	if L.GetTop() < 3 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	relativeX := L.ToInt(2)
	relativeY := L.ToInt(3)

	relativeX, relativeY = node.SetRelativeCursor(relativeX, relativeY)

	L.Push(lua.LNumber(relativeX))
	L.Push(lua.LNumber(relativeY))

	return 2
}

func (p *Script) luaFuncNodeResumeCursor(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	node.ResumeCursor()

	return 0
}

func (p *Script) luaFuncNodeHideCursor(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	node.HideCursor()

	return 0
}

// node 模拟输入事件
func (p *Script) luaFuncNodeTrigger(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	eventType := L.ToString(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	switch eventType {
	case "keypress":
		keyStr := L.ToString(3)
		if nil != node.KeyPress {
			node.KeyPress(keyStr)
		}
	}
	return 0
}

func (p *Script) luaFuncNodeRegisterLuaActiveModeHandler(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	callback := L.ToFunction(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	key := node.RegisterLuaActiveModeHandler(func(_node *Node, args ...interface{}) {
		_L := args[0].(*lua.LState)
		_callback := args[1].(*lua.LFunction)
		luaNode := _L.NewUserData()
		luaNode.Value = _node
		if err := p.Script.LuaCallByParam(_L, lua.P{
			Fn:      _callback,
			NRet:    0,
			Protect: true,
		}, luaNode); err != nil {
			panic(err)
		}
	}, L, callback)

	L.Push(lua.LString(key))
	L.Push(lua.LString(key))
	return 1
}

func (p *Script) luaFuncNodeRemoveLuaActiveModeHandler(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	key := L.ToString(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	node.RemoveLuaActiveModeHandler(key)
	return 0
}

func (p *Script) luaFuncNodeRegisterKeyPressHandler(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	callback := L.ToFunction(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	key := node.RegisterKeyPressHandler(func(_node *Node, args ...interface{}) {
		_L := args[0].(*lua.LState)
		_callback := args[1].(*lua.LFunction)
		luaNode := _L.NewUserData()
		luaNode.Value = _node
		_keyStr := args[2].(string)
		if err := p.Script.LuaCallByParam(_L, lua.P{
			Fn:      _callback,
			NRet:    0,
			Protect: true,
		}, luaNode, lua.LString(_keyStr)); err != nil {
			panic(err)
		}
	}, L, callback)

	L.Push(lua.LString(key))
	return 1
}

func (p *Script) luaFuncNodeRemoveKeyPressHandler(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	key := L.ToString(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	node.RemoveKeyPressHandler(key)
	return 0
}

func (p *Script) luaFuncNodeRegisterKeyPressEnterHandler(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	callback := L.ToFunction(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	key := node.RegisterKeyPressEnterHandler(func(_node *Node, args ...interface{}) {
		_L := args[0].(*lua.LState)
		_callback := args[1].(*lua.LFunction)
		luaNode := _L.NewUserData()
		luaNode.Value = _node
		if err := p.Script.LuaCallByParam(_L, lua.P{
			Fn:      _callback,
			NRet:    0,
			Protect: true,
		}, luaNode); err != nil {
			panic(err)
		}
	}, L, callback)

	L.Push(lua.LString(key))
	return 1
}

func (p *Script) luaFuncNodeRemoveKeyPressEnterHandler(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	key := L.ToString(2)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	node.RemoveKeyPressEnterHandler(key)
	return 0
}

func (p *Script) luaFuncNodeRemove(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	if nil == node {
		return 0
	}

	p.page.RemoveNode(node)

	return 0
}

func (p *Script) luaFuncNodeAppend(L *lua.LState) int {
	if L.GetTop() < 2 {
		return 0
	}

	lu := L.ToUserData(1)
	node := p._getNodePointerFromUserData(L, lu)
	content := L.ToString(2)
	if nil == node {
		return 0
	}

	err := p.page.AppendNode(node, content)
	if nil != err {
		L.Push(lua.LString(err.Error()))
	}

	return 0
}
