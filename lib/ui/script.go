package ui

import (
	"fin/script"
	"log"
	"path/filepath"
	"sync"

	"github.com/gizak/termui"
	lua "github.com/yuin/gopher-lua"
)

type ScriptDoc struct {
	DataType string
	Data     string
}

type Script struct {
	Script               *script.Script
	page                 *Page
	luaDocs              []ScriptDoc
	luaState             *lua.LState
	LuaCallByParamLocker *sync.RWMutex
}

func (p *Page) prepareScript() {
	var err error
	s := new(Script)
	s.LuaCallByParamLocker = new(sync.RWMutex)
	s.Script = script.NewScript(s.LuaCallByParamLocker)

	s.page = p

	s.luaDocs = make([]ScriptDoc, 0)

	s.luaState = lua.NewState()

	luaBase := s.luaState.NewTable()

	s.Script.RegisterScript(s.luaState)

	s.luaState.SetGlobal("base", luaBase)

	s.Script.RegisterBaseTable(s.luaState, luaBase)

	s.luaState.SetField(luaBase, "UIReRender", s.luaState.NewFunction(s.luaFuncUIReRender))

	s.luaState.SetField(luaBase, "WindowWidth", s.luaState.NewFunction(s.luaFuncWindowWidth))
	s.luaState.SetField(luaBase, "WindowHeight", s.luaState.NewFunction(s.luaFuncWindowHeight))
	s.luaState.SetField(luaBase, "WindowConfirm", s.luaState.NewFunction(s.luaFuncWindowConfirm))

	s.luaState.SetField(luaBase, "GetNodePointer", s.luaState.NewFunction(s.luaFuncGetNodePointer))

	s.luaState.SetField(luaBase, "NodeWidth", s.luaState.NewFunction(s.luaFuncNodeWidth))
	s.luaState.SetField(luaBase, "NodeHeight", s.luaState.NewFunction(s.luaFuncNodeHeight))
	s.luaState.SetField(luaBase, "NodeInnerAreaWidth", s.luaState.NewFunction(s.luaFuncNodeInnerAreaWidth))
	s.luaState.SetField(luaBase, "NodeInnerAreaHeight", s.luaState.NewFunction(s.luaFuncNodeInnerAreaHeight))
	s.luaState.SetField(luaBase, "NodeGetAttribute", s.luaState.NewFunction(s.luaFuncNodeGetAttribute))
	s.luaState.SetField(luaBase, "NodeSetAttribute", s.luaState.NewFunction(s.luaFuncNodeSetAttribute))
	s.luaState.SetField(luaBase, "NodeSetActive", s.luaState.NewFunction(s.luaFuncNodeSetActive))

	s.luaState.SetField(luaBase, "NodeGetHTMLData", s.luaState.NewFunction(s.luaFuncNodeGetHTMLData))
	s.luaState.SetField(luaBase, "NodeSetValue", s.luaState.NewFunction(s.luaFuncNodeSetValue))
	s.luaState.SetField(luaBase, "NodeGetValue", s.luaState.NewFunction(s.luaFuncNodeGetValue))

	s.luaState.SetField(luaBase, "NodeSetCursor", s.luaState.NewFunction(s.luaFuncNodeSetCursor))
	s.luaState.SetField(luaBase, "NodeResumeCursor", s.luaState.NewFunction(s.luaFuncNodeResumeCursor))
	s.luaState.SetField(luaBase, "NodeHideCursor", s.luaState.NewFunction(s.luaFuncNodeHideCursor))

	s.luaState.SetField(luaBase, "NodeTrigger",
		s.luaState.NewFunction(s.luaFuncNodeTrigger))

	s.luaState.SetField(luaBase, "NodeRegisterLuaActiveModeHandler",
		s.luaState.NewFunction(s.luaFuncNodeRegisterLuaActiveModeHandler))
	s.luaState.SetField(luaBase, "NodeRemoveLuaActiveModeHandler",
		s.luaState.NewFunction(s.luaFuncNodeRemoveLuaActiveModeHandler))
	s.luaState.SetField(luaBase, "NodeRegisterKeyPressHandler",
		s.luaState.NewFunction(s.luaFuncNodeRegisterKeyPressHandler))
	s.luaState.SetField(luaBase, "NodeRemoveKeyPressHandler",
		s.luaState.NewFunction(s.luaFuncNodeRemoveKeyPressHandler))
	s.luaState.SetField(luaBase, "NodeRegisterKeyPressEnterHandler",
		s.luaState.NewFunction(s.luaFuncNodeRegisterKeyPressEnterHandler))
	s.luaState.SetField(luaBase, "NodeRemoveKeyPressEnterHandler",
		s.luaState.NewFunction(s.luaFuncNodeRemoveKeyPressEnterHandler))

	s.luaState.SetField(luaBase, "NodeRemove", s.luaState.NewFunction(s.luaFuncNodeRemove))
	s.luaState.SetField(luaBase, "NodeAppend", s.luaState.NewFunction(s.luaFuncNodeAppend))

	s.luaState.SetField(luaBase, "NodeCanvasClean", s.luaState.NewFunction(s.luaFuncNodeCanvasClean))
	s.luaState.SetField(luaBase, "NodeCanvasUnSet", s.luaState.NewFunction(s.luaFuncNodeCanvasUnSet))
	s.luaState.SetField(luaBase, "NodeCanvasSet", s.luaState.NewFunction(s.luaFuncNodeCanvasSet))
	s.luaState.SetField(luaBase, "NodeCanvasDraw", s.luaState.NewFunction(s.luaFuncNodeCanvasDraw))

	s.luaState.SetField(luaBase, "NodeSelectAppendOption",
		s.luaState.NewFunction(s.luaFuncNodeSelectAppendOption))
	s.luaState.SetField(luaBase, "NodeSelectClearOptions",
		s.luaState.NewFunction(s.luaFuncNodeSelectClearOptions))
	s.luaState.SetField(luaBase, "NodeSelectSetOptionData",
		s.luaState.NewFunction(s.luaFuncNodeSelectSetOptionData))

	s.luaState.SetField(luaBase, "NodeTerminalSetCommandPrefix",
		s.luaState.NewFunction(s.luaFuncNodeTerminalSetCommandPrefix))
	s.luaState.SetField(luaBase, "NodeTerminalRegisterCommandHandle",
		s.luaState.NewFunction(s.luaFuncNodeTerminalRegisterCommandHandle))
	s.luaState.SetField(luaBase, "NodeTerminalRemoveCommandHandle",
		s.luaState.NewFunction(s.luaFuncNodeTerminalRemoveCommandHandle))
	s.luaState.SetField(luaBase, "NodeTerminalWriteString",
		s.luaState.NewFunction(s.luaFuncNodeTerminalWriteString))
	s.luaState.SetField(luaBase, "NodeTerminalWriteNewLine",
		s.luaState.NewFunction(s.luaFuncNodeTerminalWriteNewLine))
	s.luaState.SetField(luaBase, "NodeTerminalClearLines",
		s.luaState.NewFunction(s.luaFuncNodeTerminalClearLines))
	s.luaState.SetField(luaBase, "NodeTerminalClearCommandHistory",
		s.luaState.NewFunction(s.luaFuncNodeTerminalClearCommandHistory))

	s.luaState.SetField(luaBase, "NodeEditorLoadFile",
		s.luaState.NewFunction(s.luaFuncNodeEditorLoadFile))

	s.luaState.SetField(luaBase, "NodeTabpaneSetActiveTab",
		s.luaState.NewFunction(s.luaFuncNodeTabpaneSetActiveTab))

	s.luaState.SetField(luaBase, "NodeModalDoString", s.luaState.NewFunction(s.luaFuncNodeModalDoString))
	s.luaState.SetField(luaBase, "NodeModalShow", s.luaState.NewFunction(s.luaFuncNodeModalShow))

	s.luaState.SetField(luaBase, "ModalClose", s.luaState.NewFunction(s.luaFuncModalClose))
	s.luaState.SetField(luaBase, "MainPageDoString", s.luaState.NewFunction(s.luaFuncMainPageDoString))

	err = s.luaState.DoFile(filepath.Join(GlobalOption.ResBaseDir, "lua/script/core.lua"))
	if nil != err {
		panic(err)
	}
	err = s.luaState.DoFile(filepath.Join(GlobalOption.ResBaseDir, "lua/ui/core.lua"))
	if nil != err {
		panic(err)
	}

	p.Script = s
}

func (p *Script) appendDoc(doc ScriptDoc) {
	p.luaDocs = append(p.luaDocs, doc)
}

func (p *Page) GetLuaDocs(index int) ScriptDoc {
	return p.Script.luaDocs[index]
}

func (p *Page) AppendScript(doc ScriptDoc) {
	p.Script.appendDoc(doc)
}

func (p *Script) Run() {
	defer func() {
		if r := recover(); nil != r {
			termui.StopLoop()
			log.Println(r)
		}
	}()

	var err error
	for _, doc := range p.luaDocs {
		switch doc.DataType {
		case "file":
			err = p.luaState.DoFile(
				filepath.Join(GlobalOption.ProjectPath, doc.Data))
		case "string":
			err = p.luaState.DoString(doc.Data)
		}
		if nil != err {
			panic(err)
		}
	}
}
