package jsonq

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"testing"
)

func TestParser(t *testing.T) {
	ret1, _ := _NewParser("token1 #2 token3+token4 #5+#6 token7\\ \\+\\ 8 9 #10 \\#0# #-1").Parse()
	xtesting.Equal(t, ret1, []interface{}{"token1", 2, Multi("token3", "token4"), Multi(5, 6), "token7 + 8", "9", 10, "#0#", -1})

	ret2, _ := _NewParser("123123 #000 \\\\456 \\789 \\##### * \\* \\\\**\\\\*+\\**\\\\*+\\*+\\#+\\##\\#").Parse()
	xtesting.Equal(t, ret2, []interface{}{"123123", 0, "\\456", "789", "#####", All(), "*", Multi("\\**\\*", "**\\*", "*", "#", "###")})

	ret3, _ := _NewParser("").Parse()
	xtesting.Equal(t, ret3, []interface{}{})

	ret4, _ := _NewParser(" 0 #1 #2+#3 #4+#5+#6 \\#+#0+\\###\\++\\#+\\+## ").Parse()
	xtesting.Equal(t, ret4, []interface{}{"0", 1, Multi(2, 3), Multi(4, 5, 6), Multi("#", 0, "###+", "#", "+##")})

	ret5, _ := _NewParser("+++ + + +").Parse()
	xtesting.Equal(t, ret5, []interface{}{})

	ret6, _ := _NewParser("\\# \\## \\### #0 0 \\#0 \\##0\"").Parse()
	xtesting.Equal(t, ret6, []interface{}{"#", "##", "###", 0, "0", "#0", "##0\""})

	ret7, _ := _NewParser("\\\\+\\\\\\#+\\\\\\##+\\\\+\\\\\\++\\\\\\\\#0").Parse()
	xtesting.Equal(t, ret7, []interface{}{Multi("\\", "\\#", "\\##", "\\", "\\+", "\\\\#0")})
}
