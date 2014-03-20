# jgtr - json go template renderer

jgtr is a renderer for templates written in Go's [text/template](http://golang.org/pkg/template) language. It consumes data written in a standard, human-readable format and uses that data to populate the template, then writes out the rendered template to a file or stdout.

All the building blocks to make jgtr are already a part of Go, so it's mostly a matter of [snapping blocks together](http://blogs.msdn.com/b/oldnewthing/archive/2009/08/04/9856634.aspx). I just wanted a simple way to render templates with data, and I was already comfortable with JSON, and I was already somewhat comfortable with Go templates, and it looked like a good opportunity to learn more Go, so hey, why not.

Despite its name, jgtr supports JSON, YAML 1.1 and TOML v0.2.0. I prefer JSON, so I wrote that first, but I wound up adding other formats. I didn't want to change the name to `gtr` so I kept the `j`.

## Usage

jgtr consumes two input files, the data and the template, and creates an output file by rendering the template with the data. By default, all the input is from stdin, and all the output is to stdout. The `-d`, `-t` and `-o` options can be used to indicate specific files for the data, template and output (respectively). You can't read both the template and the data from stdin at the same time, so at least one of `-d` and `-t` must be given.

To determine the type of the data file, jgtr provides the options `-j`, `-y` and `-T` for JSON, YAML and TOML respectively. If no such option is given, jgtr attempts to guess from the name of the file. Files ending in `.json` are assumed to be JSON, files ending in `.toml` are assumed to be TOML, and `.yml` or `.yaml` are assumed to be YAML. If it's still ambiguous, jgtr defaults to JSON.

Consider the following files. We'll call them `test.json` and and `test.template`.

```JSON
{
    "sweater": "blue",
    "pants": "red",
    "list": [1, 2, 3]
}
```

```
My sweater is {{ .sweater }} and my pants are {{ .pants }}
This list contains {{ range .list }}{{ . }}{{ end }}
```

We could invoke jgtr in any of the following ways:

```
jgtr -d test.json < test.template > test.txt
jgtr -t test.template < test.json > test.txt
jgtr -d test.json -t test.template > test.txt
jgtr -d test.json -t test.template -o test.txt
```

And our output would be stored in `test.txt`, which would look like this:

```
My sweater is blue and my pants are red
This list contains 123
```

If we had `test.yaml` and `test.toml` like this:

```YAML
sweater: "blue"
pants: "red"
list:
 - 1
 - 2
 - 3
```

```TOML
sweater = "blue"
pants = "red"
list = [1, 2, 3]
```

We could get the same output by substituting them for `test.json`. jgtr would guess the type from their extensions and parse accordingly. If their extensions were not so descriptive, we could use the `-y`/`-T` flags to force YAML/TOML format.


## Templates

jgtr uses Go's `text/template` as its template language. If you know that, then you can probably use jgtr. In brief, the top-level value is exposed as `{{ . }}`. If that value is a map, you can access its keys by their names, as shown in the above example. (Note: *Everything* in TOML is a map; there are no top-level values that have no keys.) If that value is an array, you can loop over its contents using the `range` action. Refer to the documentation for more details on how to use the templating language.

You can use certain Go functions in your templates. More details on jgtr's use of this feature can be found in [`funcs.md`](funcs.md).

## Todo

 - add more functions from the Go standard packages. I'm mostly adding as I encounter use cases for them, so feel free to submit issues/pull requests for any that you need.
 - Can't think of any more data formats to add, but open an issue/pull request if you got one. To make life easier for me, it should have a parsing package written in Go (provide a link), which can unmarshal the data into an `interface{}` using the analogous Go types. Take a look at [`data.go`](data.go) to see how the existing formats are handled. If your format can be added in the same way, it's easy enough that I'll actually do it.
 - add a flag to switch to `html/template` for security and proper escaping. I don't personally care much about this use case, but it should be a straightforward addition if I feel like it.

## License

MIT/expat, see [LICENSE.md](LICENSE.md).
