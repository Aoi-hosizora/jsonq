package jsonq

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"log"
	"testing"
	"unsafe"
)

var typeDoc = `
{
	"foo": 1,
	"bar": 2,
	"test": "Hello, world!",
	"baz": 123.1,
	"numstring": "42",
	"floatstring": "42.1",
	"array": [
		{"foo": 1},
		{"bar": 2},
		{"baz": 3}
	],
	"subobj": {
		"foo": 1,
		"subarray": [1,2,3],
		"subsubobj": {
			"bar": 2,
			"baz": 3,
			"array": ["hello", "world"]
		}
	},
	"collections": {
		"bools": [false, true, false],
		"strings": ["hello", "strings"],
		"numbers": [1,2,3,4],
		"arrays": [[1.0,2.0],[2.0,3.0],[4.0,3.0]],
		"objects": [
			{"obj1": 1},
			{"obj2": 2}
		]
	},
	"bool": true
}
`

func TestTypes(t *testing.T) {
	bytes := *(*[]byte)(unsafe.Pointer(&typeDoc))
	doc, err := NewJsonDocument(bytes)
	if err != nil {
		log.Fatalln(err)
	}

	jq := NewJsonQuery(doc)

	val1, _ := jq.Int64("foo")
	xtesting.Equal(t, val1, int64(1))

	val2, _ := jq.String("test")
	xtesting.Equal(t, val2, "Hello, world!")

	val3, _ := jq.Float64("baz")
	xtesting.Equal(t, val3, 123.1)

	val4, _ := jq.String("numstring")
	xtesting.Equal(t, val4, "42")

	val5, _ := jq.String("floatstring")
	xtesting.Equal(t, val5, "42.1")

	val6, _ := jq.Array("array")
	xtesting.Equal(t, val6, []interface{}{
		map[string]interface{}{"foo": 1.},
		map[string]interface{}{"bar": 2.},
		map[string]interface{}{"baz": 3.},
	})

	val7, _ := jq.Object("subobj")
	xtesting.Equal(t, val7["foo"], 1.)

	val8, _ := jq.Bool("bool")
	xtesting.Equal(t, val8, true)

	val9, _ := jq.Objects("array")
	xtesting.Equal(t, val9, []map[string]interface{}{{"foo": 1.}, {"bar": 2.}, {"baz": 3.}})

	val10, _ := jq.Int64s("subobj", "subarray")
	xtesting.Equal(t, val10, []int64{1, 2, 3})

	val11, _ := jq.Strings("subobj", "subsubobj", "array")
	xtesting.Equal(t, val11, []string{"hello", "world"})

	val12, _ := jq.Bools("collections", "bools")
	xtesting.Equal(t, val12, []bool{false, true, false})

	val13, _ := jq.Arrays("collections", "arrays")
	xtesting.Equal(t, val13, [][]interface{}{{1., 2.}, {2., 3.}, {4., 3.}})

	val14, _ := jq.ObjectsBySelector("collections objects *")
	xtesting.Equal(t, val14, []map[string]interface{}{{"obj1": 1.}, {"obj2": 2.}})

	val15, _ := jq.Int64BySelector("collections objects #0 obj1")
	xtesting.Equal(t, val15, int64(1))

	val16, _ := jq.BoolBySelector("bool")
	xtesting.Equal(t, val16, true)
}
