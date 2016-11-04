package editor

import "github.com/gizak/termui"

type Editor struct {
	FirstLine, LastLine, CurrentLine *Line

	Lines []*Line

	*Block

	TextFgColor termui.Attribute
	TextBgColor termui.Attribute
	WrapLength  int // words wrap limit. Note it may not work properly with multi-width char
}

func NewEditor() *Editor {
	return &Editor{
		Lines:       make([]*Line, 0),
		Block:       NewBlock(),
		TextFgColor: termui.ThemeAttr("par.text.fg"),
		TextBgColor: termui.ThemeAttr("par.text.bg"),
	}
}

func (p *Editor) Text() string {
	var ret string
	for _, line := range p.Lines {
		ret = ret + "\n" + string(line.Data)
	}
	return ret
}

func (p *Editor) Write(keyStr string) {
	if 0 == len(p.Lines) {
		p.CurrentLine = p.InitNewLine()
	}

	if "<space>" == keyStr {
		keyStr = " "
	}

	if "<tab>" == keyStr {
		keyStr = "\t"
	}

	if "<enter>" == keyStr {
		p.CurrentLine = p.InitNewLine()
		return
	}

	if "C-8" == keyStr {
		if len(p.CurrentLine.Data) > 0 {
			p.CurrentLine.Backspace()
		} else {
			p.RemoveLine(p.CurrentLine)
		}
		return
	}

	p.CurrentLine.Write(keyStr)
}

func (p *Editor) Buffer() termui.Buffer {
	buf := p.Block.Buffer()

	fg, bg := p.TextFgColor, p.TextBgColor
	cs := termui.DefaultTxBuilder.Build(p.Text(), fg, bg)

	// wrap if WrapLength set
	if p.WrapLength < 0 {
		cs = termui.WrapTx(cs, p.Width-2)
	} else if p.WrapLength > 0 {
		cs = termui.WrapTx(cs, p.WrapLength)
	}

	y, x, n := 0, 0, 0
	for y < p.innerArea.Dy() && n < len(cs) {
		w := cs[n].Width()
		if cs[n].Ch == '\n' || x+w > p.innerArea.Dx() {
			y++
			x = 0 // set x = 0
			if cs[n].Ch == '\n' {
				n++
			}

			if y >= p.innerArea.Dy() {
				buf.Set(p.innerArea.Min.X+p.innerArea.Dx()-1,
					p.innerArea.Min.Y+p.innerArea.Dy()-1,
					termui.Cell{Ch: '…', Fg: p.TextFgColor, Bg: p.TextBgColor})
				break
			}
			continue
		}

		buf.Set(p.innerArea.Min.X+x, p.innerArea.Min.Y+y, cs[n])

		n++
		x += w
	}

	return buf
}
