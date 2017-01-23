package editor

import (
	"unicode/utf8"

	"github.com/gizak/termui"
)

type Line struct {
	ContentStartX, ContentStartY int

	Editor *Editor
	Data   []byte
	Cells  []termui.Cell
	Next   *Line
	Prev   *Line
}

func (p *Editor) InitNewLine() *Line {
	p.LinesLocker.Lock()
	defer p.LinesLocker.Unlock()

	ret := &Line{
		Editor:        p,
		ContentStartX: p.Block.InnerArea.Min.X,
		ContentStartY: p.Block.InnerArea.Min.Y,
		Data:          make([]byte, 0),
	}
	p.Lines = append(p.Lines, ret)

	if nil == p.FirstLine {
		p.FirstLine = ret
	}

	if nil != p.LastLine {
		p.LastLine.Next = ret
		ret.Prev = p.LastLine
	}

	p.LastLine = ret

	return ret
}

func (p *Editor) RemoveLine(line *Line) {
	p.LinesLocker.Lock()
	defer p.LinesLocker.Unlock()

	if nil != line.Prev {
		line.Prev.Next = line.Next
	}
	if nil != line.Next {
		line.Next.Prev = line.Prev
	}

	if p.FirstLine == line {
		p.FirstLine = p.FirstLine.Next
	}

	if p.LastLine == line {
		p.LastLine = p.LastLine.Prev
	}

	for k, v := range p.Lines {
		if line == v {
			p.Lines = append(p.Lines[:k], p.Lines[k+1:]...)
		}
	}

	p.CurrentLine = line.Prev
	p.CursorLocation.OffXCellIndex = len(p.CurrentLine.Cells)
}

func (p *Editor) ClearLines() {
	p.LinesLocker.Lock()
	defer p.LinesLocker.Unlock()

	p.FirstLine = nil
	p.LastLine = nil
	p.CurrentLine = nil
	p.Lines = []*Line{}
	p.CursorLocation.ResetLocation()
}

func (p *Line) Write(ch string) {
	off := p.Editor.CursorLocation.OffXCellIndex

	if off >= len(p.Cells) {
		p.Data = append(p.Data, []byte(ch)...)

	} else if 0 == off {
		p.Data = append([]byte(ch), p.Data...)

	} else {
		newData := make([]byte, len(p.Data)+len(ch))
		_off, i := 0, 0
		for ; i < off; i += 1 {
			_off += utf8.RuneLen(p.Cells[i].Ch)
		}
		copy(newData, p.Data[:_off])
		copy(newData[_off:], []byte(ch))
		copy(newData[_off+len(ch):], p.Data[_off:])
		p.Data = newData
	}

	fg, bg := p.Editor.TextFgColor, p.Editor.TextBgColor
	cells := termui.DefaultTxBuilder.Build(ch, fg, bg)
	p.Editor.CursorLocation.OffXCellIndex += len(cells)
}

func (p *Line) Backspace() {
	if p.Editor.CursorLocation.OffXCellIndex > len(p.Cells) {
		p.Editor.CursorLocation.OffXCellIndex = len(p.Cells)
	}
	off := p.Editor.CursorLocation.OffXCellIndex

	if off == 0 && 1 == len(p.Editor.Lines) {
		return
	}

	if 0 == off {
		p.Editor.RemoveLine(p)

	} else if off == len(p.Cells) {
		p.Data = p.Data[:len(p.Data)-utf8.RuneLen(p.Cells[off-1].Ch)]
		p.Editor.CursorLocation.OffXCellIndex -= 1

	} else {
		_off, i := 0, 0
		for ; i < off-1; i += 1 {
			_off += utf8.RuneLen(p.Cells[i].Ch)
		}
		p.Data = append(p.Data[:_off], p.Data[_off+utf8.RuneLen(p.Cells[off-1].Ch):]...)
		p.Editor.CursorLocation.OffXCellIndex -= 1
	}
}
