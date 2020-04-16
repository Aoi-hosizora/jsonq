package jsonq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEscapeString(t *testing.T) {
	ret1, _ := renderSelector("token1 #2 token3+token4 #5+#6 token7\\ \\+\\ 8 9 #10 ##0#") // token1 #2 token3+token4 #5+#6 token7\ \+\ 8 9 #10 ##0#
	assert.Equal(t, ret1, []interface{}{"token1", 2, NewMultiToken("token3", "token4"), NewMultiToken(5, 6), "token7 + 8", "9", 10, "#0#"})

	ret2, _ := renderSelector("123123 #000 \\\\456 \\789 ######")
	assert.Equal(t, ret2, []interface{}{"123123", 0, "\\456", "789", "#####"})

	ret3, _ := renderSelector("")
	assert.Equal(t, ret3, []interface{}{})

	ret4, _ := renderSelector(" 0 #1 #2+#3 #4+#5+#6 ##+#0+####\\++##+\\+## ")
	assert.Equal(t, ret4, []interface{}{"0", 1, NewMultiToken(2, 3), NewMultiToken(4, 5, 6), NewMultiToken("#", 0, "###+", "#", "+##")})

	ret5, _ := renderSelector("+++ + + +")
	assert.Equal(t, ret5, []interface{}{})

	ret6, _ := renderSelector("## ### #### #0 0 ##0 ###0\"")
	assert.Equal(t, ret6, []interface{}{"#", "##", "###", 0, "0", "#0", "##0\""})

	ret7, _ := renderSelector("\\\\+\\\\\\#+\\\\\\##+\\\\+\\\\\\++\\\\\\\\#0") // \\+\\\#+\\\##+\\+\\\++\\\\#0 -> \, \#, \##, \, \+, \\#0
	assert.Equal(t, ret7, []interface{}{NewMultiToken("\\", "\\#", "\\##", "\\", "\\+", "\\\\#0")})
}
