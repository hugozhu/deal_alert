package parse

import (
	"testing"
)

func TestLex(t *testing.T) {
	l := lex("name", "Test 123 \"My LA\" +Apple -水果")
	var items []item
	for {
		item := l.nextItem()
		t.Log(item)
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
}
