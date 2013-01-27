package parse

import (
	"testing"
)

func TestLex(t *testing.T) {
	query := "   Test 123 \"My LA\"      +  Apple -  \"水果\"  + "
	t.Log("orgi:", "\"", query, "\"")
	l := lex("name", query)
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
