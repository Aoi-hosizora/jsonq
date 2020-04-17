# jsonq

+ A json query library written in golang

### Function

+ Select object and array by tokens
+ Select by selector (jsonq version)
+ ~~Return a multi-layers object~~ (only support to return an array now)

### Install

```bash
# use go get
go get github.com/Aoi-hosizora/jsonq

# use go mod
# import "github.com/Aoi-hosizora/jsonq"
go mod tidy
```

### Usage

+ see [jsonq_test.go](jsonq_test.go) and [types_test.go](types_test.go) for details

```go
doc, err := jsonq.NewJsonDocument(objDoc)
if err != nil {
    log.Fatalln(err)
}
jq := jsonq.NewJsonQuery(doc)

// m[1]
val, err := jq.Select(1)
// m[:]
val, err := jq.Strings(jsonq.All()) // *
// m[len(m)-2][0]
val, err := jq.Int64(-2, 0)
// m["a"]["0"]["b"]
val, err := jq.Select("a", "0", "b")
// m["a"][0]["b"][0:2]
val, err := jq.Select("a", 0, "b", jsonq.Multi(0, 1)) // #0+#1
// m[1]["*"]["a"]["2"][0/2][:]
val, err := jq.SelectBySelector("#1 \\* a 2 #0+#2 *")
```

### Selector

+ A convenient language to select json
+ Grammar definition

```
selector := mtok        // multi token

mtok     := mtok mtok   // the next layer
mtok     := *           // all fields in the current layer
mtok     := stok+stok   // multiple fields in the current layer
mtok     := stok        // single token
stok     := token       // string or number

token    := #numbers    // array index
token    := strings     // map key
```

+ Rules: (`WS` means `whitespace`)
    + use `WS` to split layers
    + use `+` to split fields
    + use `*` to represent all fields (could not use with `+`) 
    + use `\` to escape all tokens (especially for `WS` `+` `#` `*`)
    + use `#numbers` as an array index (token start with `#`)
    + use raw number and other string as a map field name
    + if a field name starts with `#` or `*`, use `\#` and `\*` (if `#` and `*` is inside string, it is not necessary to escape)
    + if a field name includes a `WS` or `+`, use `\WS` and `\+`
+ Example

```
token1 #2 token3+token4 #5+#6 token7\ \+\ 8 9 #10 \#0# #-1
-> "token1", 2, {"token3", "token4"}, {5, 6}, "token7 + 8", "9", 10, "#0#", -1

123123 #000 \\456 \789 \##### * \* \\**\\*+\**\\*+\*+\#+\##\#
-> "123123", 0, "\456", "789", "#####", *, "*", {"*\*", "**\*", "*", "#", "###"}

0 #1 #2+#3 #4+#5+#6 \#+#0+\###\++\#+\+##
-> "0", 1, {2, 3}, {4, 5, 6}, {"#", 0, "###+", "#", "+##"}
```

### References

+ [jmoiron/jsonq](https://github.com/jmoiron/jsonq)
+ [Handwritten Parsers & Lexers in Go](https://blog.gopheracademy.com/advent-2014/parsers-lexers/)
