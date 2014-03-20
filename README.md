# jgtr - json go template renderer
jgtr is a very simple tool that uses a data file, and a template file written in Go's [`text/template`](http://golang.org/pkg/text/template) as input. It outputs rendered templates to stdout or a file. It can be used as a general purpose tool for rendering templates of whatever.

All the building blocks to make jgtr are already a part of Go, so it's mostly a matter of [snapping blocks together](http://blogs.msdn.com/b/oldnewthing/archive/2009/08/04/9856634.aspx). I just wanted a simple way to render templates with data, and I was already comfortable with JSON, and I was already somewhat comfortable with Go templates, and it looked like a good opportunity to learn more Go, so hey, why not.

Although I originally wrote jgtr for JSON because it's my preferred format, I am busy adding other formats, which you can specify via flags. I'm going to leave the j in the name because I like it that way.

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
jgtr -d test.json < test.template > test.txt
jgtr -t test.template < test.json > test.txt
jgtr -d test.json -t test.template > test.txt
jgtr -d test.json -t test.template -o test.txt
```

Namely, `test.txt`, which will contain:
```
Sweater color is blue
Pants color is red
```

If you had a YAML file like this:
```
sweater: blue
pants: red
```
Then you could use the -y flag to specify that the data is YAML, and do this to get the same output.
```
jgtr -d test.yaml -y < test.template > test.txt
```

## Templates

jgtr uses Go's `text/template` language. If you know that, then you can probably use jgtr. In brief, the top-level value is exposed as `{{ . }}`. If that value is a map, you can access its keys by their names, as shown in the above example. If that value is an array, you can loop over its contents using the `range` action. Refer to the documentation for more details on how to use the templating language.

jgtr exposes some standard Go functions that you can use in your templates. Take a look at [`funcs.go`](funcs.go) for the full list. As an example, jgtr exposes Go's `time.Parse` under the name `timeParse`. You can this to generate a Go `time.Time` object and invoke standard functions on it.
```
"2013-04-30"
```
```
{{ with timeParse "2006-01-02" . }}
The year is {{ .Year }}
The month is {{ .Month }}
The day is {{ .Day }}
Out of 365 days, we're at day number {{ .YearDay }}
All in all, the day is ISO8601 {{ .Format "2006-01-02T15:04:05Z07:00" }}
{{ end }}
```
will produce
```

The year is 2013
The month is April
The day is 30
Out of 365 days, we're at day number 120
All in all, the day is ISO8601 2013-04-30T00:00:00Z

```
Note the empty lines in the output. Those are from the `with` and `end` templates. Those templates evaluate to the empty string, so we only use them for their side effects. Since I put them on their own lines (for readability), they leave behind empty lines when rendered.

## Todo
 - add more functions. `text/template` supports introducing more functions into a template via [`Funcs`](http://golang.org/pkg/text/template/#Template.Funcs). Feel free to submit issues or pull requests for more functions.
 - I added YAML, TOML should be coming soon.
 - add a flag to switch to `html/template` for security and proper escaping. I don't personally care much about this use case, but it should be a straightforward addition if I feel like it.

## License
MIT/expat, see [LICENSE.md](LICENSE.md).
