package module

import (
	"fmt"
)

type Tree struct {
	Name	string
	text	string
	lex		*lexer
	peekCount int
	token	[3]item
	Root	*ListNode
}

func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

func (t *Tree) pos() Pos {
	return t.token[0].pos
}

func (t *Tree) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}

func (t *Tree) nextNonSpace() (token item) {
	//if token.typ != itemSpace {
	//	return token
	//}
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	return token
}

func (t *Tree) expect(expected itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected {
		t.unexpected(token, context)
	}
	return token
}

func (t *Tree) unexpected(token item, context string) {
	fmt.Printf("%s in %s", token, context)
}

func (t *Tree) parse() {
	t.Root = newList(t.peek().pos)
	for t.peek().typ != itemEOF {
		n := t.textOrModule()
		t.Root.append(n)
	}
}

func (t *Tree) textOrModule() Node {
	switch token := t.nextNonSpace(); token.typ {
	case itemText:
		return newText(token.pos, token.val)
	case itemLeftDelim:
		return t.module()
	}
	return nil
}

func (t *Tree) module() Node {
	var nameNode, pathNode, remoteNode, branchNode Node
	for t.peek().typ != itemRightDelim {
		switch token := t.nextNonSpace(); token.typ {
		case itemModName:
			nameNode = t.modName()
		case itemModPath:
			pathNode = t.modPath()
		case itemModRemote:
			remoteNode = t.modRemote()
		case itemModBranch:
			branchNode = t.modBranch()
		default:
		}
	}
	t.next()
	return newModule(t.pos(), nameNode, pathNode, remoteNode, branchNode)
}

func (t *Tree) modName() Node {
	if t.nextNonSpace().typ == itemEqual {
		if token := t.nextNonSpace(); token.typ == itemString {
			return newName(token.pos, token.val)
		}
	}
	return nil
}

func (t *Tree) modPath() Node {
	if t.nextNonSpace().typ == itemEqual {
		if token := t.nextNonSpace(); token.typ == itemString {
			return newPath(token.pos, token.val)
		}
	}
	return nil
}

func (t *Tree) modRemote() Node {
	if t.nextNonSpace().typ == itemEqual {
		if token := t.nextNonSpace(); token.typ == itemString {
			return newRemote(token.pos, token.val)
		}
	}
	return nil
}

func (t *Tree) modBranch() Node {
	if t.nextNonSpace().typ == itemEqual {
		if token := t.nextNonSpace(); token.typ == itemString {
			return newBranch(token.pos, token.val)
		}
	}
	return nil
}