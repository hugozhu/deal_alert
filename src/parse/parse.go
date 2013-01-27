package parse

import (
	"bytes"
	"fmt"
)

type Node interface {
	Type() NodeType
	String() string
	Keyword() string
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeTerm NodeType = iota
	NodePhrase
	NodeList
)

type TermNode struct {
	NodeType
	Term string
}

type PhraseNode struct {
	NodeType
	Phrase string
}

func newTerm(text string) *TermNode {
	return &TermNode{NodeType: NodeTerm, Term: text}
}

func (t *TermNode) String() string {
	return fmt.Sprintf("%q ", t.Term)
}

func (t *TermNode) Keyword() string {
	return t.Term
}

func newPhrase(text string) *PhraseNode {
	return &PhraseNode{NodeType: NodePhrase, Phrase: text}
}

func (t *PhraseNode) String() string {
	return fmt.Sprintf("\"%q\"", t.Phrase)
}

func (t *PhraseNode) Keyword() string {
	return t.Phrase[1 : len(t.Phrase)-1]
}

// ListNode holds a sequence of nodes.
type ListNode struct {
	NodeType
	Nodes []Node // The element nodes in lexical order.
}

func newList() *ListNode {
	return &ListNode{NodeType: NodeList}
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

type Query struct {
	lex      *lexer
	OrNodes  *ListNode
	AndNodes *ListNode
	NotNodes *ListNode
}

func (q *Query) AllTerms() (arr []string) {
	arr = []string{}
	for _, n := range q.AndNodes.Nodes {
		arr = append(arr, n.Keyword())
	}
	for _, n := range q.OrNodes.Nodes {
		arr = append(arr, n.Keyword())
	}
	for _, n := range q.NotNodes.Nodes {
		arr = append(arr, n.Keyword())
	}
	return
}

func Parse(input string) (query *Query, err error) {
	query = &Query{
		lex:      lex("query", input),
		OrNodes:  newList(),
		AndNodes: newList(),
		NotNodes: newList(),
	}
	for {
		item := query.lex.nextItem()
		var list *ListNode
		if item.typ == itemOR {
			list = query.OrNodes
			item = query.lex.nextItem()
		} else if item.typ == itemNOT {
			list = query.NotNodes
			item = query.lex.nextItem()
		} else {
			list = query.AndNodes
		}
		if item.typ == itemEOF || item.typ == itemError {
			break
		}

		if item.typ == itemTerm {
			list.append(newTerm(item.val))
		}

		if item.typ == itemPhrase {
			list.append(newPhrase(item.val))
		}
	}
	return
}
