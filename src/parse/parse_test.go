package parse

import (
	"testing"
)

func TestParse(t *testing.T) {
	text := "   Test 123 \"My LA\"      +  Apple -  \"水果\"  + "
	query, _ := Parse(text)
	t.Log(query.AndNodes)
	t.Log(query.NotNodes)
	t.Log(query.OrNodes)
}
