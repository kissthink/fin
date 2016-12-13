package ui

import (
	"container/list"
	uiutils "in/ui/utils"
	"log"
	"sync"

	"github.com/gizak/termui"
	"golang.org/x/net/html"
)

type Page struct {
	// 当该 page 为 Modal 类型时，则存在该值
	// MainPage 用于表示 Modal 所属 Page
	MainPage *Page

	Title string

	IdToNodeMap map[string]*Node

	Bufferers []termui.Bufferer

	parseAgentMap  []*ParseAgent
	renderAgentMap []*RenderAgent
	Script         *Script
	Modals         map[string]*NodeModal
	CurrentModal   *NodeModal

	doc                     *html.Node
	parsingNodesStack       *list.List
	FirstChildNode          *Node
	FocusNode               *list.Element
	WorkingNodes            *list.List
	ActiveNode              *Node
	ActiveNodeAfterRerender *Node

	renderingX int
	renderingY int

	recoverVal interface{}

	KeyPressHandleLocker sync.RWMutex
}

func newPage() *Page {
	ret := new(Page)

	ret.IdToNodeMap = make(map[string]*Node, 0)

	ret.parsingNodesStack = list.New()
	ret.WorkingNodes = list.New()

	ret.prepareScript()
	ret.prepareModals()
	ret.prepareParse()
	ret.prepareRender()

	return ret
}

func (p *Page) dumpNodesHtmlData(node *Node) {
	log.Println(node.HtmlData)
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		p.dumpNodesHtmlData(childNode)
	}
}

func (p *Page) DumpNodesHtmlData() {
	p.dumpNodesHtmlData(p.FirstChildNode)
}

func (p *Page) RemoveNode(node *Node) {
	if nodeDataOnRemover, ok := node.Data.(NodeDataOnRemover); true == ok {
		nodeDataOnRemover.NodeDataOnRemove()
	}

	delete(p.IdToNodeMap, node.Id)

	if nil != node.PrevSibling {
		node.PrevSibling.NextSibling = node.NextSibling
	}

	if nil != node.NextSibling {
		node.NextSibling.PrevSibling = node.PrevSibling
	}

	if nil != node.Parent {
		node.Parent.ChildrenCount -= 1
		if node.Parent.FirstChild == node {
			node.Parent.FirstChild = node.NextSibling
		}
		if node.Parent.LastChild == node {
			node.Parent.LastChild = node.PrevSibling
		}
	}

	p.Rerender()
}

func (p *Page) nodeAfterUIRenderHandle(node *Node) {
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		p.nodeAfterUIRenderHandle(childNode)
	}

	if nodeDataAfterUIRenderHandler, ok := node.Data.(NodeDataAfterUIRenderHandler); true == ok {
		nodeDataAfterUIRenderHandler.NodeDataAfterUIRenderHandle()
	}
}

// 将 page 上的内容渲染到屏幕上
// 并更新 FocusNode / ActiveNode / WorkingNodes内元素之间方向关系
func (p *Page) uiRender() {
	GCurrentRenderPage = p
	if 0 == len(p.Bufferers) {
		return
	}

	if nil != p.FocusNode {
		if nodeDataUnFocusModer, ok := p.FocusNode.Value.(*Node).Data.(NodeDataUnFocusModer); true == ok {
			nodeDataUnFocusModer.NodeDataUnFocusMode()
		}
	}

	if nil != p.ActiveNode {
		if nodeDataUnActiveModer, ok := p.ActiveNode.Data.(NodeDataUnActiveModer); true == ok {
			nodeDataUnActiveModer.NodeDataUnActiveMode()
		}
	}

	if nil != p.ActiveNodeAfterRerender {
		p.SetActiveNode(p.ActiveNodeAfterRerender)
		p.FocusNode = nil
	} else if nil != p.FocusNode {
		p.SetActiveNode(p.FocusNode.Value.(*Node))
	}

	var (
		e, e2       *list.Element
		node, node2 *Node
	)
	if p.WorkingNodes.Len() > 0 {
		for e = p.WorkingNodes.Front(); e != nil; e = e.Next() {
			node = e.Value.(*Node)
			node.FocusThisNode = e
			node.FocusTopNode = nil
			node.FocusBottomNode = nil
		}

		for e = p.WorkingNodes.Front(); e != nil; e = e.Next() {
			node = e.Value.(*Node)

			for e2 = e.Next(); e2 != nil; e2 = e2.Next() {
				node2 = e2.Value.(*Node)

				if (node.UIBlock.InnerArea.Min.X <= node2.UIBlock.InnerArea.Min.X &&
					node2.UIBlock.InnerArea.Min.X <= node.UIBlock.InnerArea.Max.X) ||
					(node.UIBlock.InnerArea.Max.X <= node2.UIBlock.InnerArea.Min.X &&
						node2.UIBlock.InnerArea.Max.X <= node.UIBlock.InnerArea.Max.X) ||
					(node.UIBlock.InnerArea.Min.X <= node2.UIBlock.InnerArea.Min.X &&
						node2.UIBlock.InnerArea.Max.X >= node.UIBlock.InnerArea.Max.X) &&
						(node2.UIBlock.InnerArea.Min.Y > node.UIBlock.InnerArea.Max.Y) {

					node.FocusBottomNode = e2
					node2.FocusTopNode = e
					break
				}
			}
		}
	}

	uiutils.UIRender(p.Bufferers...)

	p.nodeAfterUIRenderHandle(p.FirstChildNode)
}

// 刷新 page 内容到屏幕
// 包括刷新 page 上 focus/active 元素
func (p *Page) Refresh() {
	uiClear()
	p.uiRender()
}

// 重新渲染 page 并刷新内容到屏幕
func (p *Page) Rerender() {
	p.Render()
	p.Refresh()
}

func (p *Page) Serve() {
	p.Refresh()

	go func() {
		defer func() {
			if r := recover(); nil != r {
				termui.StopLoop()
				p.recoverVal = r
			}
		}()
		p.Script.Run()
	}()

	termui.Loop()
}
