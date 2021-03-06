function UIReRender()
  return base.UIReRender()
end

function WindowWidth()
  return base.WindowWidth()
end

function WindowHeight()
  return base.WindowHeight()
end

function WindowConfirm(title, callback)
    content = string.format([[
    <table top=6>
        <tr>
            <td offset=5 cols=2>
                <select id="SelectConfirm" borderlabel="%s">
                <!--
                paddingtop=1 paddingbottom=1
                paddingleft=1 paddingright=1
                -->
                    <option value="cancel">取消</option>
                    <option value="confirm">确定</option>
                </select>
            </td>
        </tr>
    </table>
    ]], title)
    return base.WindowConfirm(content, callback)
end

local _Node = {}
local _mtNode = {__index = _Node}

function Node(target)
    local nodePointer
    local targetType = type(target)
    if "string" == targetType then
      nodePointer = base.GetNodePointer(target)
      if nil == nodePointer then
          return nil
      end 
    elseif "userdata" == targetType then
      nodePointer = target
    else
      return nil 
    end

    local ret = setmetatable({}, _mtNode)
    ret.nodePointer = nodePointer
    return ret
end

function _Node.Width(self)
    return base.NodeWidth(self.nodePointer)
end

function _Node.Height(self)
    return base.NodeHeight(self.nodePointer)
end

function _Node.InnerAreaWidth(self)
    return base.NodeInnerAreaWidth(self.nodePointer)
end

function _Node.InnerAreaHeight(self)
    return base.NodeInnerAreaHeight(self.nodePointer)
end

function _Node.GetAttribute(self, key, value)
    return base.NodeGetAttribute(self.nodePointer, key, value)
end

function _Node.SetAttribute(self, key, value)
    return base.NodeSetAttribute(self.nodePointer, key, value)
end

function _Node.SetActive(self)
    return base.NodeSetActive(self.nodePointer)
end

function _Node.GetHTMLData(self)
    return base.NodeGetHTMLData(self.nodePointer)
end

function _Node.SetValue(self, text)
    return base.NodeSetValue(self.nodePointer, text)
end

function _Node.GetValue(self)
    return base.NodeGetValue(self.nodePointer)
end

function _Node.SetCursor(self, x, y)
    return base.NodeSetCursor(self.nodePointer, x, y)
end

function _Node.ResumeCursor(self)
    return base.NodeResumeCursor(self.nodePointer)
end

function _Node.HideCursor(self)
    return base.NodeHideCursor(self.nodePointer)
end

function _Node.Trigger(self, eventType, value)
    return base.NodeTrigger(self.nodePointer, eventType, value)
end

function _Node.RegisterLuaActiveModeHandler(self, callback)
    return base.NodeRegisterLuaActiveModeHandler(self.nodePointer, callback)
end

function _Node.RemoveLuaActiveModeHandler(self, key)
    return base.NodeRemoveLuaActiveModeHandler(self.nodePointer, key)
end

function _Node.RegisterKeyPressHandler(self, callback)
    return base.NodeRegisterKeyPressHandler(self.nodePointer, callback)
end

function _Node.RemoveKeyPressHandler(self, key)
    return base.NodeRemoveKeyPressEnterHandler(self.nodePointer, key)
end

function _Node.RegisterKeyPressEnterHandler(self, callback)
    return base.NodeRegisterKeyPressEnterHandler(self.nodePointer, callback)
end

function _Node.RemoveKeyPressEnterHandler(self, key)
    return base.NodeRemoveKeyPressEnterHandler(self.nodePointer, key)
end

function _Node.Remove(self)
    return base.NodeRemove(self.nodePointer)
end

function _Node.Append(self, content)
    return base.NodeAppend(self.nodePointer, content)
end

function _Node.CanvasClean(self)
    return base.NodeCanvasClean(self.nodePointer)
end

function _Node.CanvasUnSet(self, x, y)
    return base.NodeCanvasUnSet(self.nodePointer, x, y)
end

function _Node.CanvasSet(self, x, y, ch, fg, bg)
    return base.NodeCanvasSet(self.nodePointer, x, y, ch, fg, bg)
end

function _Node.CanvasDraw(self)
    return base.NodeCanvasDraw(self.nodePointer)
end

function _Node.SelectAppendOption(self, value, data)
    return base.NodeSelectAppendOption(self.nodePointer, value, data)
end

function _Node.SelectClearOptions(self)
    return base.NodeSelectClearOptions(self.nodePointer)
end

function _Node.SelectSetOptionData(self, value, newData)
    return base.NodeSelectSetOptionData(self.nodePointer, value, newData)
end

function _Node.TerminalSetCommandPrefix(self, commandPrefix)
    return base.NodeTerminalSetCommandPrefix(self.nodePointer, commandPrefix)
end

function _Node.TerminalRegisterCommandHandle(self, callback)
    return base.NodeTerminalRegisterCommandHandle(self.nodePointer, callback)
end

function _Node.TerminalRemoveCommandHandle(self, key)
    return base.NodeTerminalRemoveCommandHandle(self.nodePointer, key)
end

function _Node.TerminalWriteString(self, data)
    return base.NodeTerminalWriteString(self.nodePointer, data)
end

function _Node.TerminalWriteNewLine(self, line)
    return base.NodeTerminalWriteNewLine(self.nodePointer, line)
end

function _Node.TerminalClearLines(self)
    return base.NodeTerminalClearLines(self.nodePointer)
end

function _Node.TerminalClearCommandHistory(self)
    return base.NodeTerminalClearCommandHistory(self.nodePointer)
end

function _Node.EditorLoadFile(self, filePath)
    return base.NodeEditorLoadFile(self.nodePointer, filePath)
end

function _Node.TabpaneSetActiveTab(self, name)
    return base.NodeTabpaneSetActiveTab(self.nodePointer, name)
end

function _Node.ModalDoString(self, callback)
    return base.NodeModalDoString(self.nodePointer, callback)
end

function _Node.ModalShow(self)
    return base.NodeModalShow(self.nodePointer)
end

function ModalClose()
    return base.ModalClose()
end

function MainPageDoString(callback)
    return base.MainPageDoString(callback)
end
