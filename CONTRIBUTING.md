# Contributing to jgtr

Interested in contributing to jgtr? Trying to do something with jgtr but it's just not coming out right? You have come to the right document. Issues and pull requests on pretty much anything regarding jgtr are welcome. Here are some extra notes so we can be on the same page.

## Adding a data format

Although I originally wrote jgtr for JSON, I have also added TOML and YAML, both of which I have come across in other programming adventures. If you have a data format that you need support for, feel free to open an issue. Here are some guidelines on what I am looking for in a data format:

 - structurally similar to JSON. It should have the same data structures - ordered lists, string-value mappings, numbers, strings, etc. I imagine most data formats that are actually useful will already satisfy this requirement.
 - already has a support package written in Go. It should be `go get`-able, no weird dependencies. Pure Go is ideal but nothing's wrong with cgo either. In addition, this package should be able to unmarshal data structures in the following manner:
    - maps/objects/tables should be unmarshalled to `map[string]interface{}`
    - lists/arrays should be unmarshalled to `[]interface`
    - strings should be unmarshalled to `string` (duh)
    - numbers should be unmarshalled to some sensible combination of standard Go numeric types, and there should be only one size of number (eg for YAML, all floats map to `float64` and all ints map to `int`; you don't have some floats mapping to `float32` and others mapping to `float64`). Preferably no mingling of uints and ints, and minimal mingling of floats and ints, but jgtr should be able to handle it.
 - can unmarshal out of a standard UTF-8 `string` or `[]byte`. Decoding straight from an `io.Reader` is a nice bonus, such as with [`json.NewDecoder`](http://golang.org/pkg/encoding/json/#NewDecoder).

Why all these restrictions? Well, if you read them, I think they're all pretty obvious. But more importantly, it minimizes the amount of thinking and testing needed to handle specific formats. For example, jgtr's YAML decoder unmarshals integers and floats separately, where as the standard `encoding/json` package unmarshals all numbers to `float64`. This leads to some weird edge cases if you try to sort a list of YAML numbers and it turns out that some of them were unmarshalled to int, and others to float, and you can't compare them and oops jgtr panicked. But the same list in JSON will sort fine because all the numbers are floats. Now, maybe the int/float distinction is a good thing, maybe it's not, but either way, the inconsistency is a pain in the neck. If all decoders have consistently shaped output, maintenance is easier.

New data formats are added to [`data.go`](data.go) in a single function of the form `func (string) (interface{}, error)`. The existing formats give a self-explanatory description of how a new format should appear. If your format can't be parsed out of an `io.Reader`, use `ioutil.ReadAll` as with the YAML function.

## Adding a function

As I mentioned in [`funcs.md`](funcs.md), the template is the application. If the functions provided by jgtr are too sparse, your templates become less powerful and jgtr becomes less useful. You'd have to make manual changes to your data to get the output you want, and at some point you might as well just substitute in the data into your template by hand, which defeats the purpose of having templates at all. So having functions that solve problems is pretty important.

jgtr has a pretty good list of functions right now; I'm quite pleased with the sorting functions and they were fun to write. I'll make note of things I want to add as I come up with them.

If there's functionality you need, try hacking it out with `if` and `range` actions first. I bet you could write a sorting function entirely with actions if you really wanted. If it starts to feel really clunky, then you probably need the expressive power of an actual function. Feel free to open an issue describing your use case, or better yet implement it and open a pull request. The function should be added to the `tmplFuncs` map in [`funcs.go`](funcs.go). Be sure to update `funcs.md` as well to document its usage.

Try to stick to the standard packages and keep things simple. If the function you have in mind is *really* specific, it might not make sense to add it to jgtr. Be advised that if your function messes with lists/maps, you will probably need to familiarize yourself with the [`reflect`](http://golang.org/pkg/reflect) package.

## Miscellaneous improvements

Got other ideas or problems? Just open an issue or pull request describing the situation.
