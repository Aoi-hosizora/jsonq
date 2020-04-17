package jsonq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	ret1 := NewParser("token1 #2 token3+token4 #5+#6 token7\\ \\+\\ 8 9 #10 \\#0# #-1").Parse()
	assert.Equal(t, ret1, []interface{}{"token1", 2, NewMultiToken("token3", "token4"), NewMultiToken(5, 6), "token7 + 8", "9", 10, "#0#", -1})

	ret2 := NewParser("123123 #000 \\\\456 \\789 \\##### * \\* \\\\**\\\\*+\\**\\\\*+\\*+\\#+\\##\\#").Parse()
	assert.Equal(t, ret2, []interface{}{"123123", 0, "\\456", "789", "#####", NewStarToken(), "*", NewMultiToken("\\**\\*", "**\\*", "*", "#", "###")})

	ret3 := NewParser("").Parse()
	assert.Equal(t, ret3, []interface{}{})

	ret4 := NewParser(" 0 #1 #2+#3 #4+#5+#6 \\#+#0+\\###\\++\\#+\\+## ").Parse()
	assert.Equal(t, ret4, []interface{}{"0", 1, NewMultiToken(2, 3), NewMultiToken(4, 5, 6), NewMultiToken("#", 0, "###+", "#", "+##")})

	ret5 := NewParser("+++ + + +").Parse()
	assert.Equal(t, ret5, []interface{}{})

	ret6 := NewParser("\\# \\## \\### #0 0 \\#0 \\##0\"").Parse()
	assert.Equal(t, ret6, []interface{}{"#", "##", "###", 0, "0", "#0", "##0\""})

	ret7 := NewParser("\\\\+\\\\\\#+\\\\\\##+\\\\+\\\\\\++\\\\\\\\#0").Parse()
	assert.Equal(t, ret7, []interface{}{NewMultiToken("\\", "\\#", "\\##", "\\", "\\+", "\\\\#0")})
}
