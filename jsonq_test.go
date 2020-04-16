package jsonq

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var objDoc = `
{
	"a": "b",
	"c": {
		"e": 0,
		"f": [
			{"g": 123, "h": 0.3, "i": "abc"},
			{"g": 456, "h": 0.6, "i": "def"},
			{"g": 789, "h": 0.9, "i": "ghi"}
		],
		"j": {
			"k": null,
			"l": [
				[1, 2, 3],
				[4, 5, 6]
			]
		}
	}
}
`

var arrDoc = `
[
	[1, 2, 3],
	{
		"a": 0,
		"b": {
			"c": "d",
			"e": 0.2
		},
		"bb": {
			"c": "dd",
			"e": 0.22
		}
	},
	{
		"a": 1,
		"b": {
			"c": "dd",
			"e": 0.22,
			"e": 0.22,
			"f": [
				{
					"g": "1g",
					"h": "1h"
				},
				{
					"g": "1gg",
					"h": "1hh"
				}
			]
		}
	},
	{
		"a": 2,
		"b": {
			"c": "ddd",
			"e": 0.222,
			"f": [
				{
					"g": "2g",
					"h": "2h"
				},
				{
					"g": "2gg",
					"h": "2hh"
				},
				{
					"g": "2gg",
					"h": "2hh"
				},
				[1, 2, 3],
				[4.1, 5.2, 6.3]
			]
		}
	},
	null
]
`

func handle(obj interface{}, err error) interface{} {
	if err != nil {
		log.Fatalln(err)
	}
	return obj
}

func TestObject(t *testing.T) {
	doc, err := NewJsonDocument(objDoc)
	if err != nil {
		log.Fatalln(err)
	}

	jq := NewQuery(doc)
	val1 := handle(jq.Select("a"))                 // b
	val2 := handle(jq.Select("c", "e"))            // 0
	val3 := handle(jq.Select("c", "f", 0))         // map[g:123 h:0.3 i:abc]
	val4 := handle(jq.Select("c", "f", 0, "h"))    // 0.3
	val5 := handle(jq.Select("c", "j", "k"))       // <nil>
	val6 := handle(jq.Select("c", "j", "l", 0))    // [1 2 3]
	val7 := handle(jq.Select("c", "j", "l", 0, 0)) // 1

	assert.Equal(t, val1, "b")
	assert.Equal(t, val2, 0.)
	assert.Equal(t, val3, map[string]interface{}{"g": 123., "h": 0.3, "i": "abc"})
	assert.Equal(t, val4, 0.3)
	assert.Equal(t, val5, nil)
	assert.Equal(t, val6, []interface{}{1., 2., 3.})
	assert.Equal(t, val7, 1.)
}

func TestArray(t *testing.T) {
	doc, err := NewJsonDocument(arrDoc)
	if err != nil {
		log.Fatalln(err)
	}

	jq := NewQuery(doc)
	val1 := handle(jq.Select(0))                 // [1 2 3]
	val2 := handle(jq.Select(1))                 // map[a:0 b:map[c:d e:0.2]]
	val3 := handle(jq.Select(2, "a"))            // 1
	val4 := handle(jq.Select(2, "b", "c"))       // dd
	val5 := handle(jq.Select(3, "b", "f", 0))    // map[g:g h:h]
	val6 := handle(jq.Select(3, "b", "f", 3))    // [1 2 3]
	val7 := handle(jq.Select(3, "b", "f", 4, 1)) // 5.2
	val8 := handle(jq.Select(4))                 // <nil>

	assert.Equal(t, val1, []interface{}{1., 2., 3.})
	assert.Equal(t, val2, map[string]interface{}{"a": 0., "b": map[string]interface{}{"c": "d", "e": 0.2}, "bb": map[string]interface{}{"c": "dd", "e": 0.22}})
	assert.Equal(t, val3, 1.)
	assert.Equal(t, val4, "dd")
	assert.Equal(t, val5, map[string]interface{}{"g": "2g", "h": "2h"})
	assert.Equal(t, val6, []interface{}{1., 2., 3.})
	assert.Equal(t, val7, 5.2)
	assert.Equal(t, val8, nil)
}

func TestMultiToken(t *testing.T) {
	doc, err := NewJsonDocument(arrDoc)
	if err != nil {
		log.Fatalln(err)
	}

	jq := NewQuery(doc)
	val1 := handle(jq.Select(NewMultiToken(1, 2, 3), "a"))                             // [0 1 2]
	val2 := handle(jq.Select(NewMultiToken(1, 2, 3), "b", "c"))                        // [d dd ddd]
	val3 := handle(jq.Select(NewMultiToken(1, 2, 3), "b", NewMultiToken("c", "e")))    // [d 0.2 dd 0.22 ddd 0.222]
	val4 := handle(jq.Select(1, NewMultiToken("b", "bb"), "c"))                        // [d dd]
	val5 := handle(jq.Select(1, NewMultiToken("b", "bb"), NewMultiToken("c", "e")))    // [d 0.2 dd 0.22]
	val6 := handle(jq.Select(NewMultiToken(2, 3), "b", "f", NewMultiToken(0, 1), "g")) // [1g 1gg 2g 2gg]

	assert.Equal(t, val1, []interface{}{0., 1., 2.})
	assert.Equal(t, val2, []interface{}{"d", "dd", "ddd"})
	assert.Equal(t, val3, []interface{}{"d", 0.2, "dd", 0.22, "ddd", 0.222})
	assert.Equal(t, val4, []interface{}{"d", "dd"})
	assert.Equal(t, val5, []interface{}{"d", 0.2, "dd", 0.22})
	assert.Equal(t, val6, []interface{}{"1g", "1gg", "2g", "2gg"})
}

/*
=== RUN   TestObject
--- PASS: TestObject (0.00s)
=== RUN   TestArray
--- PASS: TestArray (0.00s)
=== RUN   TestMultiToken
--- PASS: TestMultiToken (0.00s)
PASS
 */
