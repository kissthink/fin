// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package extra

import (
	"unicode/utf8"

	. "github.com/gizak/termui"
	rw "github.com/mattn/go-runewidth"
)

type Tab struct {
	Label       string
	RuneLen     int
	StringWidth int
	Blocks      []Bufferer
}

func NewTab(label string) *Tab {
	return &Tab{
		Label:       label,
		RuneLen:     utf8.RuneCount([]byte(label)),
		StringWidth: rw.StringWidth(label) + 2,
	}
}

func (tab *Tab) SetLabel(text string) {
	tab.Label = text
	tab.RuneLen = utf8.RuneCount([]byte(text))
	tab.StringWidth = rw.StringWidth(text) + 2
}

func (tab *Tab) AddBlocks(rs ...Bufferer) {
	for _, r := range rs {
		tab.Blocks = append(tab.Blocks, r)
	}
}

func (tab *Tab) Buffer() Buffer {
	buf := NewBuffer()
	for blockNum := 0; blockNum < len(tab.Blocks); blockNum++ {
		b := tab.Blocks[blockNum]
		buf.Merge(b.Buffer())
	}
	return buf
}

type Tabpane struct {
	Block
	Tabs           []*Tab
	activeTabIndex int
	TabpaneFg      Attribute
	TabpaneBg      Attribute
	TabFg          Attribute
	TabBg          Attribute
	ActiveTabFg    Attribute
	ActiveTabBg    Attribute
	posTabText     []int
	offTabText     int
	IsHideMenu     bool
}

func NewTabpane() *Tabpane {
	tp := Tabpane{
		Block:          *NewBlock(),
		activeTabIndex: 0,
		offTabText:     0,
		TabFg:          ThemeAttr("fg.tab"),
		TabBg:          ThemeAttr("bg.tab"),
		ActiveTabFg:    ThemeAttr("fg.tab.active"),
		ActiveTabBg:    ThemeAttr("bg.tab.active")}
	return &tp
}

func (tp *Tabpane) SetTabs(tabs ...*Tab) {
	tp.Tabs = make([]*Tab, len(tabs))
	tp.posTabText = make([]int, len(tabs)+1)
	off := 0
	for i := 0; i < len(tp.Tabs); i++ {
		tp.Tabs[i] = tabs[i]
		tp.posTabText[i] = off
		if false == tp.Block.Border {
			off += tp.Tabs[i].StringWidth
		} else {
			off += tp.Tabs[i].StringWidth + 1 //+1 for space between tabs
		}
	}
	tp.posTabText[len(tabs)] = off - 1 //total length of Tab's text
}

// 向左移动 ActiveTab
// return:
//			bool  是否还能继续右移
func (tp *Tabpane) SetActiveLeft() bool {
	if tp.activeTabIndex == 0 {
		return false
	}
	tp.activeTabIndex -= 1
	if tp.posTabText[tp.activeTabIndex] < tp.offTabText {
		tp.offTabText = tp.posTabText[tp.activeTabIndex]
	}
	return true
}

// 向右移动 ActiveTab
// return:
//			bool  是否还能继续右移
func (tp *Tabpane) SetActiveRight() bool {
	if tp.activeTabIndex == len(tp.Tabs)-1 {
		return false
	}
	tp.activeTabIndex += 1
	endOffset := tp.posTabText[tp.activeTabIndex] + tp.Tabs[tp.activeTabIndex].StringWidth
	if endOffset+tp.offTabText > tp.InnerWidth() {
		tp.offTabText = endOffset - tp.InnerWidth()
	}
	return true
}

func (tp *Tabpane) GetActiveIndex() int {
	return tp.activeTabIndex
}

func (tp *Tabpane) SetActiveTab(index int) bool {
	if index > len(tp.Tabs)-1 || index < 0 {
		return false
	}
	if index == tp.activeTabIndex {
		return false
	}
	tp.activeTabIndex = index
	if tp.posTabText[tp.activeTabIndex] < tp.offTabText {
		tp.offTabText = tp.posTabText[tp.activeTabIndex]
	}
	return true
}

// Checks if left and right tabs are fully visible
// if only left tabs are not visible return -1
// if only right tabs are not visible return 1
// if both return 0
// use only if fitsWidth() returns false
func (tp *Tabpane) checkAlignment() int {
	ret := 0
	if tp.offTabText > 0 {
		ret = -1
	}
	if tp.offTabText+tp.InnerWidth() < tp.posTabText[len(tp.Tabs)] {
		ret += 1
	}
	return ret
}

// Checks if all tabs fits innerWidth of Tabpane
func (tp *Tabpane) fitsWidth() bool {
	return tp.InnerWidth() >= tp.posTabText[len(tp.Tabs)]
}

func (tp *Tabpane) Align() {
	if !tp.fitsWidth() && !tp.Border {
		tp.PaddingLeft += 1
		tp.PaddingRight += 1
		tp.Block.Align()
	}
}

// bridge the old Point stuct
type point struct {
	X  int
	Y  int
	Ch rune
	Fg Attribute
	Bg Attribute
}

func buf2pt(b Buffer) []point {
	ps := make([]point, 0, len(b.CellMap))
	for k, c := range b.CellMap {
		ps = append(ps, point{X: k.X, Y: k.Y, Ch: c.Ch, Fg: c.Fg, Bg: c.Bg})
	}

	return ps
}

// Adds the point only if it is visible in Tabpane.
// Point can be invisible if concatenation of Tab's texts is widther then
// innerWidth of Tabpane
func (tp *Tabpane) addPoint(chWidth int, ptab []point, charOffset *int, oftX *int, points ...point) []point {
	if *charOffset < tp.offTabText || tp.offTabText+tp.InnerWidth() < *charOffset {
		*charOffset += chWidth
		return ptab
	}
	for _, p := range points {
		p.X = *oftX
		ptab = append(ptab, p)
	}
	*oftX += chWidth
	*charOffset += chWidth
	return ptab
}

func (tp *Tabpane) addPoints(ptab []point, points ...point) []point {
	for _, p := range points {
		ptab = append(ptab, p)
	}
	return ptab
}

// Draws the point and redraws upper and lower border points (if it has one)
func (tp *Tabpane) drawPointWithBorder(p point, ch rune, chbord rune, chdown rune, chup rune) []point {
	var addp []point
	p.Ch = ch
	if tp.Border {
		p.Ch = chdown
		p.Y = tp.InnerY() - 1
		addp = append(addp, p)
		p.Ch = chup
		p.Y = tp.InnerY() + 1
		addp = append(addp, p)
		p.Ch = chbord
	}
	p.Y = tp.InnerY()
	return append(addp, p)
}

func (tp *Tabpane) Buffer() Buffer {
	if true == tp.IsHideMenu {
		tp.Height = 0
	} else {
		if tp.Border {
			tp.Height = 3
		} else {
			tp.Height = 1
		}
	}

	if tp.Width > tp.posTabText[len(tp.Tabs)]+2 {
		tp.Width = tp.posTabText[len(tp.Tabs)] + 2
	}

	buf := tp.Block.Buffer()
	ps := []point{}

	tp.Align()
	if false == tp.IsHideMenu {
		if tp.InnerHeight() <= 0 || tp.InnerWidth() <= 0 {
			return NewBuffer()
		}
	}

	if false == tp.IsHideMenu {
		// 画背景
		if false == tp.Block.Border {
			_max := TermWidth()
			var addp []point
			for _oftX := tp.posTabText[len(tp.Tabs)-1]; _oftX < _max; _oftX++ {
				addp = append(addp, point{X: _oftX, Y: tp.InnerY(), Ch: ' ', Fg: tp.TabpaneFg, Bg: tp.TabpaneBg})
			}
			ps = tp.addPoints(ps, addp...)
		}
	}

	oftX := tp.InnerX()
	charOffset := 0
	pt := point{Bg: tp.BorderBg, Fg: tp.BorderFg}
	var chLen int
	for i, tab := range tp.Tabs {
		if false == tp.IsHideMenu {
			if i != 0 && true == tp.Block.Border {
				pt.X = oftX
				pt.Y = tp.InnerY()
				addp := tp.drawPointWithBorder(pt, ' ', VERTICAL_LINE, HORIZONTAL_DOWN, HORIZONTAL_UP)
				ps = tp.addPoint(1, ps, &charOffset, &oftX, addp...)
			}

			if i == tp.activeTabIndex {
				if true == tp.Block.Border {
					ps = tp.addPoint(1, ps, &charOffset, &oftX, []point{
						point{Y: tp.InnerY(), Ch: ' ', Fg: tp.ActiveTabFg, Bg: tp.ActiveTabBg},
						point{Y: tp.InnerY() + 1, Ch: ' ', Fg: tp.BorderFg, Bg: tp.BorderBg},
					}...)
				} else {
					ps = tp.addPoint(1, ps, &charOffset, &oftX,
						point{Y: tp.InnerY(), Ch: ' ', Fg: tp.ActiveTabFg, Bg: tp.ActiveTabBg})
				}
			} else {
				ps = tp.addPoint(1, ps, &charOffset, &oftX,
					point{Y: tp.InnerY(), Ch: ' ', Fg: tp.TabFg, Bg: tp.TabBg})
			}

			pt.Fg = tp.TabFg
			pt.Bg = tp.TabBg

			if i == tp.activeTabIndex {
				pt.Fg = tp.ActiveTabFg
				pt.Bg = tp.ActiveTabBg
			}

			rs := []rune(tab.Label)
			for k := 0; k < len(rs); k++ {

				chLen = rw.RuneWidth(rs[k])

				addp := make([]point, 0, 2)
				if i == tp.activeTabIndex && tp.Border {
					if 2 == chLen {
						pt.Ch = '　'
					} else {
						pt.Ch = ' '
					}

					pt.X = oftX
					pt.Y = tp.InnerY() + 1
					pt.Fg = tp.BorderFg
					pt.Bg = tp.BorderBg
					addp = append(addp, pt)
					pt.Fg = tp.ActiveTabFg
					pt.Bg = tp.ActiveTabBg
				}

				pt.Y = tp.InnerY()
				pt.Ch = rs[k]

				addp = append(addp, pt)

				ps = tp.addPoint(chLen, ps, &charOffset, &oftX, addp...)
			}

			if i == tp.activeTabIndex {
				if true == tp.Block.Border {
					ps = tp.addPoint(1, ps, &charOffset, &oftX, []point{
						point{Y: tp.InnerY(), Ch: ' ', Fg: tp.ActiveTabFg, Bg: tp.ActiveTabBg},
						point{Y: tp.InnerY() + 1, Ch: ' ', Fg: tp.BorderFg, Bg: tp.BorderBg},
					}...)
				} else {
					ps = tp.addPoint(1, ps, &charOffset, &oftX,
						point{Y: tp.InnerY(), Ch: ' ', Fg: tp.ActiveTabFg, Bg: tp.ActiveTabBg})
				}
			} else {
				ps = tp.addPoint(1, ps, &charOffset, &oftX,
					point{Y: tp.InnerY(), Ch: ' ', Fg: tp.TabFg, Bg: tp.TabBg})
			}

			pt.Fg = tp.BorderFg
			pt.Bg = tp.BorderBg

			if !tp.fitsWidth() {
				all := tp.checkAlignment()
				pt.X = tp.InnerX() - 1

				pt.Ch = '*'
				if tp.Border {
					pt.Ch = VERTICAL_LINE
				}
				ps = append(ps, pt)

				if all <= 0 {
					addp := tp.drawPointWithBorder(pt, '<', '«', HORIZONTAL_LINE, HORIZONTAL_LINE)
					ps = append(ps, addp...)
				}

				pt.X = tp.InnerX() + tp.InnerWidth()
				pt.Ch = '*'
				if tp.Border {
					pt.Ch = VERTICAL_LINE
				}
				ps = append(ps, pt)
				if all >= 0 {
					addp := tp.drawPointWithBorder(pt, '>', '»', HORIZONTAL_LINE, HORIZONTAL_LINE)
					ps = append(ps, addp...)
				}
			}
		}
	}

	for _, v := range ps {
		buf.Set(v.X, v.Y, NewCell(v.Ch, v.Fg, v.Bg))
	}
	buf.Sync()
	return buf
}
