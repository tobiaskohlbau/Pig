package module

import (
	"fmt"
	"bytes"
)

type Node interface {
	Type() NodeType
	String() string
	Position() Pos
	unexported()
}

type NodeType int

type Pos int

func (p Pos) Position() Pos {
	return p
}

func (Pos) unexported() {
} 

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeText	NodeType = iota
	NodeModule
	NodeName
	NodePath
	NodeRemote
	NodeBranch
	NodeList
	NodeEnd
	NodeNil
)

type ListNode struct {
	NodeType
	Pos
	Nodes []Node
}

func newList(pos Pos) *ListNode {
	return &ListNode{NodeType: NodeList, Pos: pos}
}

func (l *ListNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *ListNode) String() string {
	b := new(bytes.Buffer)
	for _, n := range l.Nodes {
		fmt.Fprint(b, n)
	}
	return b.String()
}

type ModuleNode struct {
	NodeType
	Pos
	Name Node
	Path Node
	Remote Node
	Branch Node
}

func newModule(pos Pos, name Node, path Node, remote Node, branch Node) *ModuleNode {
	return &ModuleNode{NodeType: NodeModule, Pos: pos, Name: name, Path: path, Remote: remote, Branch: branch}
}

func (t *ModuleNode) String() string {
	return fmt.Sprintf("%s %s %s %s", t.Name.String(), t.Path.String(), t.Remote.String(), t.Branch.String())
}

type TextNode struct {
	NodeType
	Pos
	Text []byte
}

func newText(pos Pos, text string) *TextNode {
	return &TextNode{NodeType: NodeText, Pos: pos, Text: []byte(text)}
}

func (t *TextNode) String() string {
	return fmt.Sprintf("%s", t.Text)
}

type NameNode struct {
	NodeType
	Pos
	Name []byte
}

func newName(pos Pos, name string) *NameNode {
	return &NameNode{NodeType: NodeName, Pos: pos, Name: []byte(name[1:len(name)-1])}
}

func (t *NameNode) String() string {
	return fmt.Sprintf("%s", t.Name)
}

type PathNode struct {
	NodeType
	Pos
	Path []byte
}

func newPath(pos Pos, path string) *PathNode {
	return &PathNode{NodeType: NodePath, Pos: pos, Path: []byte(path[1:len(path)-1])}
}

func (t *PathNode) String() string {
	return fmt.Sprintf("%s", t.Path)
}

type RemoteNode struct {
	NodeType
	Pos
	Remote []byte
}

func newRemote(pos Pos, remote string) *RemoteNode {
	return &RemoteNode{NodeType: NodeRemote, Pos: pos, Remote: []byte(remote[1:len(remote)-1])}
}

func (t *RemoteNode) String() string {
	return fmt.Sprintf("%s", t.Remote)
}

type BranchNode struct {
	NodeType
	Pos
	Branch []byte
}

func newBranch(pos Pos, branch string) *BranchNode {
	return &BranchNode{NodeType: NodeBranch, Pos: pos, Branch: []byte(branch[1:len(branch)-1])}
}

func (t *BranchNode) String() string {
	return fmt.Sprintf("%s", t.Branch)
}

type NilNode struct {
	Pos
}

func newNil(pos Pos) *NilNode {
	return &NilNode{Pos: pos}
}

func (n *NilNode) Type() NodeType {
	return NodeNil
}

func (n *NilNode) String() string {
	return "nil"
}