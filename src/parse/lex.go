// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ == itemAND:
		return "+"
	case i.typ == itemOR:
		return "|"
	case i.typ == itemNOT:
		return "-"
	case i.typ > itemError:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemTerm
	itemPhrase
	itemAND
	itemOR
	itemNOT
	itemEOF
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:  "error",
	itemTerm:   "term",
	itemPhrase: "phrase",
	itemAND:    "and",
	itemOR:     "or",
	itemNOT:    "not",
	itemEOF:    "eof",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // the name of the input; used only for error reports.
	input string    // the string being scanned.
	state stateFn   // the next lexing function to enter.
	pos   int       // current position in the input.
	start int       // start position of this item.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		state: lexAction,
		items: make(chan item, 2), // Two items of buffering is sufficient for all state functions
	}
	return l
}

// state functions
func lexAction(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof || r == '\n':
		return l.errorf("unclosed action")
	case isSpace(r):
		return lexSpace
	case r == '"':
		return lexQuote
	case r == '|':
		l.emit(itemOR)
	case r == '+':
		l.emit(itemAND)
	case r == '-' || r == '!':
		l.emit(itemNOT)
		return lexAction
	case unicode.IsPrint(r):
		return lexTerm
	default:
		return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexAction
}

func lexSpace(l *lexer) stateFn {
	l.ignore()
	for {
		r := l.next()
		if r == eof {
			l.emit(itemEOF)
			return nil
		}
		if isSpace(r) {
			l.ignore()
		} else {
			l.backup()
			break
		}
	}
	return lexAction
}

func lexTerm(l *lexer) stateFn {
	for {
		if l.peek() == ' ' {
			if l.pos > l.start {
				l.emit(itemTerm)
			}
			return lexAction
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemTerm)
	}
	l.emit(itemEOF)
	return nil
}

// lexQuote scans a quoted string.
func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(itemPhrase)
	return lexAction
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
