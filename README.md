# jsonq

+ A json query library written by golang

### Function

+ Select object and array by tokens
+ Select by selector (jsonq version)
+ ~~Return a multi-layers object~~ (only support to return an array now)

### Usage

+ see [jsonq_test.go](jsonq_test.go) for details

```go
doc, err := jsonq.NewJsonDocument(objDoc)
if err != nil {
    log.Fatalln(err)
}
jq := jsonq.NewQuery(doc)

// m[1]
val, err := jq.Select(1)
// m[:]
val, err := jq.Select(jsonq.NewStarToken()) // *
// m[len(m)-2][0]
val, err := jq.Select(-2, 0)
// m["a"]["0"]["b"]
val, err := jq.Select("a", "0", "b")
// m["a"][0]["b"][0:2]
val, err := jq.Select("a", 0, "b", jsonq.NewMultiToken(0, 1)) // #0+#1
// m[1]["*"]["a"]["2"][0/2][:]
val, err := jq.SelectBySelector("#1 ** a 2 #0+#2 *")
```

### Selector

+ A convenient language to select json
+ Grammar definition (temporary, will be updated)

```
selector := selector selector // the next layer
selector := selector+selector // multiple field in the current layer
selector := token

token    := #numbers  // array index
token    := *         // all fields
token    := #strings  // map key (ignore #)
token    := *strings  // map key (ignore *)
token    := strings   // map key
numbers  := (0..9)*
```

+ Escape rule:
    + use `\` to escape all tokens first (after ` ` and `+`)
    + use raw number as a field name
    + use `#numbers` to escape as a number, `##` to escape `#` (only at the start of the token)
    + use `*` to represent as a star token, `**` to escape `*` (only at the start of the token)
    + use `#xxx` (`xxx` is a non-numbers) and `*xxx` to represent `xxx`, the first char `#` and `*` will be ignored
+ Example

```
token1 #2 token3+token4 #5+#6 token7\ \+\ 8 9 #10 ##0# #-1
-> "token1", 2, {"token3", "token4"}, {5, 6}, "token7 + 8", "9", 10, "#0#", -1

123123 #000 \\456 \789 ###### * *\* **\\*+***\\*+**+##+###\#
-> "123123", 0, "\456", "789", "#####", *, "*", {"*\*", "**\*", "*", "#", "###"}

0 #1 #2+#3 #4+#5+#6 ##+#0+####\++##+\+##
-> "0", 1, {2, 3}, {4, 5, 6}, {"#", 0, "###+", "#", "+##"}
```

### References

+ [jmoiron/jsonq](https://github.com/jmoiron/jsonq)
