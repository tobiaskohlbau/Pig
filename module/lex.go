package module

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType
	pos Pos
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type itemType int

const (
	itemError	itemType = iota
	itemBool
	itemChar
	itemEOF
	itemText
	itemIdentifier
	itemLeftDelim
	itemRightDelim
	itemPipe
	itemSpace
	itemEqual
	itemQuote
	itemString
	itemKeyword
	itemModName
	itemModPath
	itemModRemote
	itemModBranch
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	name	string //the name of the input
	input	string // the string being scanned
	leftDelim string // left delim
	rightDelim	string // right delim
	state 	stateFn
	pos 	Pos // current position
	start 	Pos
	width   Pos
	lastPos Pos
	items	chan item
	parenDepth int
}

var key = map[string]itemType {
	"name": itemModName,
	"path": itemModPath,
	"remote": itemModRemote,
	"branch": itemModBranch,
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:	name,
		input:	input,
		leftDelim: "{{",
		rightDelim: "}}",
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
}

const (
	leftDelim = "{{"
	rightDelim = "}}"
)

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], l.leftDelim) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftDelim
		}
		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexLeftDelim(l *lexer) stateFn {
	l.pos += Pos(len(l.leftDelim))
	l.emit(itemLeftDelim)
	l.parenDepth = 0
	return lexInsideAction
}

func lexRightDelim(l *lexer) stateFn {
	l.pos += Pos(len(l.rightDelim))
	l.emit(itemRightDelim)
	return lexText
}

func lexInsideAction(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], l.rightDelim) {
		if l.parenDepth == 0 {
			return lexRightDelim
		}
	}
		//add error handling
	switch r := l.next(); {
	case r == eof || isEndOfLine(r):
		// add error handling
	case isSpace(r):
		return lexSpace
	case isText(r):
		l.backup()
		return lexIdentifier
	case r == '=':
		l.emit(itemEqual)
	case r == '"':
		return lexQuote
	default:
		//return error handle
	}
	return lexInsideAction
}

func lexSpace(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexInsideAction
}

func lexIdentifier(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isText(r):
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch {
			case key[word] > itemKeyword:
				l.emit(key[word])
			default:
				l.emit(itemIdentifier)
			}
			break Loop
		}
	}
	return lexInsideAction
}

func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
		case eof, '\n':
			//return error
		case '"':
			break Loop
		}
	}
	l.emit(itemString)
	return lexInsideAction
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

func isText(r rune) bool {
	return unicode.IsLetter(r)
}