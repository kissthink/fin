package ui

import (
	"sync"

	"github.com/gizak/termui"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/html"
)

type NodeKeyPress func(e termui.Event)
type NodeFocusMode func()
type NodeUnFocusMode func()
type NodeActiveMode func()
type NodeUnActiveMode func()
type NodeGetValue func() string
type NodeSetText func(content string)
type NodeOnRemove func()

type Node struct {
	Id   string
	page *Page

	ChildrenCount int

	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	// 是否要渲染子节点
	// 子节点将根据其父节点
	// Node.isShouldTermuiRenderChild 来决定是否 将 node.uiBuffer append 进 page.Bufferers
	// 例: TableTrTd 将用到该字段
	isShouldTermuiRenderChild bool

	tmpBorder   bool
	tmpBorderFg termui.Attribute
	ColorFg     string
	ColorBg     string

	uiBuffer interface{}
	uiBlock  *termui.Block

	HtmlData string
	Data     interface{}

	KeyPress     NodeKeyPress
	FocusMode    NodeFocusMode
	UnFocusMode  NodeUnFocusMode
	ActiveMode   NodeActiveMode
	UnActiveMode NodeUnActiveMode

	SetText  NodeSetText
	GetValue NodeGetValue
	OnRemove NodeOnRemove

	KeyPressEnterHandlers map[string]NodeJob
	JobHanderLocker       sync.RWMutex
}

type NodeJobHandler func(node *Node, args ...interface{})
type NodeJob struct {
	*Node
	Handler NodeJobHandler
	Args    []interface{}
}

type NodeBody struct{}

func (p *Node) InitNodeBody() *NodeBody {
	nodeBody := new(NodeBody)
	p.Data = nodeBody
	return nodeBody
}

type NodeDiv struct{}

func (p *Node) InitNodeDiv() *NodeDiv {
	nodeDiv := new(NodeDiv)
	p.Data = nodeDiv
	return nodeDiv
}

func (p *Page) newNode(htmlNode *html.Node) *Node {
	ret := new(Node)
	ret.page = p
	ret.HtmlData = htmlNode.Data
	ret.KeyPressEnterHandlers = make(map[string]NodeJob, 0)
	return ret
}

func (p *Node) addChild(child *Node) {
	if nil == p {
		return
	}

	child.Parent = p
	child.Parent.ChildrenCount += 1

	child.FirstChild = nil
	child.LastChild = nil
	child.PrevSibling = nil
	child.NextSibling = nil

	if nil == p.FirstChild {
		p.FirstChild = child
	}

	if nil != p.LastChild {
		p.LastChild.NextSibling = child
		child.PrevSibling = p.LastChild
	}

	p.LastChild = child
}

func (p *Node) uiRender() {
	if nil == p.uiBuffer {
		return
	}
	uiRender(p.uiBuffer.(termui.Bufferer))
}

func (p *Node) RegisterKeyPressEnterHandler(handler NodeJobHandler, args ...interface{}) string {
	p.JobHanderLocker.Lock()
	defer p.JobHanderLocker.Unlock()

	key := uuid.NewV4().String()
	p.KeyPressEnterHandlers[key] = NodeJob{
		Node:    p,
		Handler: handler,
		Args:    args,
	}
	return key
}

func (p *Node) RemoveKeyPressEnterHandler(key string) {
	p.JobHanderLocker.Lock()
	defer p.JobHanderLocker.Unlock()
	delete(p.KeyPressEnterHandlers, key)
}

func (p *Node) WaitKeyPressEnter() {
	c := make(chan bool, 0)
	key := p.RegisterKeyPressEnterHandler(func(_node *Node, args ...interface{}) {
		c := args[0].(chan bool)
		c <- true
	}, c)
	<-c
	p.RemoveKeyPressEnterHandler(key)
}
