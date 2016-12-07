package ui

import lua "github.com/yuin/gopher-lua"

func (p *Script) _getNodeModalPointerFromUserData(L *lua.LState, lu *lua.LUserData) *NodeModal {
	if nil == lu || nil == lu.Value {
		return nil
	}

	var (
		node      *Node
		nodeModal *NodeModal
		ok        bool
	)

	node, ok = lu.Value.(*Node)
	if false == ok || nil == node {
		return nil
	}

	if nodeModal, ok = node.Data.(*NodeModal); false == ok {
		return nil
	}

	return nodeModal
}

func (p *Script) luaFuncNodeModalShow(L *lua.LState) int {
	if L.GetTop() < 1 {
		return 0
	}

	lu := L.ToUserData(1)
	nodeModal := p._getNodeModalPointerFromUserData(L, lu)
	if nil == nodeModal {
		return 0
	}

	p.page.ClearActiveNode()
	nodeModal.page.uiRender()

	p.page.CurrentModal = nodeModal

	return 0
}
