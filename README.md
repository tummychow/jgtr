# jgtr - json go template renderer
jgtr is a very simple tool that uses JSON-encoded data and template files written in Go's [`text/template`](http://golang.org/pkg/text/template) as input. It outputs rendered templates to stdout or a file. It can be used as a general purpose tool for rendering templates of whatever.

All the building blocks to make jgtr are already a part of Go, so it's mostly a matter of [snapping blocks together](http://blogs.msdn.com/b/oldnewthing/archive/2009/08/04/9856634.aspx). I just wanted a simple way to render templates with data, and I was already comfortable with JSON, and I was already somewhat comfortable with Go templates, and it looked like a good opportunity to learn more Go, so hey, why not.

## Usage
jgtr normally reads all its input from stdin and writes all its output to stdout. Standard GNU-style flags allow you to set files for these purposes instead. Reading the template and the data from stdin at the same time is finicky, so at least one of those two has to be specified explicitly.

Consider these example files, `test.json` and `test.template`:

```
{
    "sweater": "blue",
    "pants": "red"
}
```

```
Sweater color is {{ .sweater }}
Pants color is {{ .pants }}
```

Then any of the following invocations will all result in the same output:
```
jgtr -j test.json < test.template > test.txt
jgtr -t test.template < test.json > test.txt
jgtr -j test.json -t test.template > test.txt
jgtr -j test.json -t test.template -o test.txt
```

Namely, `test.txt`, which will contain:
```
Sweater color is blue
Pants color is red
```

If you know Go's `text/template` language, then you can probably use jgtr. In brief, the top-level JSON value is exposed as `{{ . }}`. If that value is a JSON object, you can access its keys by their names, as shown in the above example. If that value is a JSON array, you can loop over its contents using the `range` action. Refer to the documentation for more details on how to use the templating language.

## Todo
 - add more functions. `text/template` supports introducing more functions into a template via [`Funcs`](http://golang.org/pkg/text/template/#Template.Funcs). Some obvious things that come to mind are date and time manipulation functions.
 - add more data file types. While I personally prefer JSON, there are plenty of other encodings that are isomorphic to it and can represent the same data structures. As long as you can unmarshal it into an `interface{}`, it can probably be dropped into this code in place of JSON. [TOML](https://github.com/BurntSushi/toml) and [YAML](https://github.com/go-yaml/yaml) come to mind.
 - add other general-purpose template languages. There aren't a lot of pure Go implementations, but I am looking at [mustache](https://github.com/hoisie/mustache). If you can think of any others, please open an issue or, better yet, add it and submit a pull request.
 - add a flag to switch to `html/template` for security and proper escaping. I don't personally care much about this use case, but it should be a straightforward addition if I feel like it.

## License
MIT/expat, see [LICENSE.md](LICENSE.md).
