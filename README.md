# jsonq

+ A json query library written by golang

### Function

+ Select object and array by tokens
+ Select by selector (jsonq version)

### Usage

+  see [jsonq_test.go](jsonq_test.go) for details

```go
doc, err := jsonq.NewJsonDocument(objDoc)
if err != nil {
    log.Fatalln(err)
}
jq := jsonq.NewQuery(doc)

// m[1]
val, err := jq.Select(1)
// m["a"]["0"]["b"]
val, err := jq.Select("a", "0", "b")
// m["a"][0]["b"][0:2]
val, err := jq.Select("a", 0, "b", jsonq.NewMultiToken(0, 1))
// m[1]["a"]["2"]["b"][0/2]
val, err := jq.SelectBySelector("#1 a 2 b #0+#2")
```

### Selector

+ A convenient language to select json
+ Grammar definition

```
selector := selector selector // the next layer
selector := selector+selector // multiple field in the current layer
selector := token

token    := #number
token    := strings
number   := 0..9
```

+ Escape rule:
    + use `\` to escape all token
    + use raw number as a string
    + use `#` to escape as a number, `##` to escape `#` (only in the beginning of the token)
+ Example

```
token1 #2 token3+token4 #5+#6 token7\ \+\ 8 9 #10 ##0# #-1
-> "token1", 2, {"token3", "token4"}, {5, 6}, "token7 + 8", "9", 10, "#0#", -1

123123 #000 \\456 \789 ######
-> "123123", 0, "\456", "789", "#####"

0 #1 #2+#3 #4+#5+#6 ##+#0+####\++##+\+##
-> "0", 1, {2, 3}, {4, 5, 6}, {"#", 0, "###+", "#", "+##"}
```

### References

+ [jmoiron/jsonq](https://github.com/jmoiron/jsonq)