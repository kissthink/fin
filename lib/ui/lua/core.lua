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

function _Node.Render(self)
    return base.NodeRender(self.nodePointer)
end

function _Node.SetActive(self)
    return base.NodeSetActive(self.nodePointer)
end

function _Node.GetHtmlData(self)
    return base.NodeGetHtmlData(self.nodePointer)
end

function _Node.SetText(self, text)
    return base.NodeSetText(self.nodePointer, text)
end

function _Node.GetValue(self)
    return base.NodeGetValue(self.nodePointer)
end

function _Node.OnKeyPressEnter(self, callback)
    return base.NodeOnKeyPressEnter(self.nodePointer, callback)
end

function _Node.Remove(self)
    return base.NodeRemove(self.nodePointer)
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

function _Node.TerminalPopNewCommand(self)
    return base.NodeTerminalPopNewCommand(self.nodePointer)
end

function _Node.TerminalWriteNewLine(self, line)
    return base.NodeTerminalWriteNewLine(self.nodePointer, line)
end

function WindowConfirm(title)
    content = string.format([[
    <table>
        <tr>
            <td offset=4 cols=4><par height=6></par></td>
        </tr>
        <tr>
            <td offset=4 cols=4><par>%s</par></td>
        </tr>
        <tr>
            <td offset=5 cols=2>
                <select id="SelectConfirm">
                    <option value="cancel">取消</option>
                    <option value="confirm">确定</option>
                </select>
            </td>
        </tr>
    </table>
    ]], title)
    return base.WindowConfirm(content)
end

function Log(...)
    base.Log(unpack(arg))
end
