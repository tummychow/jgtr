# jgtr functions

`text/template` allows the invocation of Go functions in a template. A few simple functions are provided, but if you want more, the renderer has to add them. There are a lot of jobs that you can't do with the global functions alone. If you were writing your own application, you'd add logic to the application to transform the data into the right shape for rendering. Since jgtr is too general to interpret your data, your template *is* the application. Therefore, jgtr has to add functions that let your template do the transforming.

The full list of functions enabled by jgtr is in [`funcs.go`](funcs.go). Here, I will document each of those functions and their potential use cases in your templates. If you need more functions, submit an issue or pull request detailing what you want - try to stick to things in the Go standard packages.

## timeParse

`timeParse` is an alias for Go's standard [`time.Parse`](http://golang.org/pkg/time/#Parse) function. JSON does not consider dates to be first-class values (but TOML does), so manipulating time can be a pain. This function lets you parse a string into a standard `time.Time` object, and then you can invoke standard methods on that object to extract parts of the date or time. You can also use the standard [`time.Format`](http://golang.org/pkg/time/#Time.Format) to print a date in another format of your choice. Here is a basic example.

```JSON
{
    "today": "2013-04-30"
}
```

```
{{ with timeParse "2006-01-02" .today }}The year is {{ .Year }}
The month is {{ .Month }}
The day is {{ .Day }}
Out of 365 days, we're at day number {{ .YearDay }}
All in all, the day is ISO8601 {{ .Format "2006-01-02T15:04:05Z07:00" }}{{ end }}
```

The above JSON and template files will produce this output:

```
The year is 2013
The month is April
The day is 30
Out of 365 days, we're at day number 120
All in all, the day is ISO8601 2013-04-30T00:00:00Z
```

You can parse dates that are written in any format. To specify the format, you write out a certain standard date and time, as if it were written in that format. See [this](http://golang.org/pkg/time/#pkg-constants) for more details.

## timeNow

`timeNow` is an alias for Go's standard [`time.Now`](http://golang.org/pkg/time/#Now) function. You can use this to fetch the time at which the template is being rendered, and then invoke standard methods to read the bits and pieces. If you want to pretty-print the time, use `time.Format` as mentioned above. Here's a very simple example of its usage.

```JSON
{}
```

```
This template was rendered on {{ timeNow.Format "Jan 2 2006" }}
```

```
This template was rendered on Mar 20 2014
```

## notLast
`notLast` takes a slice (ie a list of some kind) and an integer. If the integer is not the last index of the list, then it returns true, and otherwise it returns false. This function helps solve a common problem in templates. Let me describe the problem so you understand why this is useful.

```JSON
[ 41, 42, 43 ]
```

```
My list has the numbers
{{ range . }}{{ . }},{{ end }}
```

That JSON file contains a simple top-level array and the template file iterates over it. Here is the output:

```
My list has the numbers
41,42,43,
```

Note the trailing comma. Obviously it's from the last element of the list. Now suppose we want to omit the comma if the element is the last one. If you are familiar with Liquid templates (I use Jekyll so that's where I learned them), you could do [this](http://github.com/Shopify/liquid/wiki/Liquid-for-Designers#for-loops):

```
My list has the numbers {% for item in list %}{{ item }}{% unless forloop.last %},{% endunless %}{% endfor %}
```

Note the `unless forloop.last` construct. Liquid injects that variable inside every for loop. That lets us skip the comma on the last element, which is what we want. I'm not sure if Liquid's solution to the problem is elegant, but it's pretty good, and it does work, where as `text/template` has no built-in solution at all. This is a pretty common problem when forming text by iterating over lists, so jgtr offers you an alternative. Try this template:

```
My list has the numbers{{ $arr := . }}
{{ range $i, $e := $arr }}{{ $e }}{{ if notLast $arr $i }},{{ end }}{{ end }}
```

With the same JSON file, this gives rendered output of:

```
My list has the numbers
41,42,43
```

Ah. Much better. Let me explain the syntax so you know how to do this on your own lists.

First, we have to assign the list to a variable name. The list is the top-level value, so the dot equals the list. We assign it to the name `$arr`. Then we enter the `range` action. Note the assignment to acquire the index and element of the current iteration. Within the range body, we print out the current element. Then we invoke `notLast` on the list that contains the elements, and the current index. If we aren't on the last element, then we'll also print the comma.

`notLast` is necessary because you can't perform arbitrary arithmetic within a template, so there's no way to compare the current index and the length of the slice minus one. That "minus one" part isn't allowed, so we need a function to do it for us. `notLast` is that function.

The assignment to `$arr` is required as well. Inside the `range` action, the context (ie the value of the dot) is equal to the current element of the iteration. We lose visibility of the list that contains the elements, unless we bind a name to it first.

Overall, I would rate this solution as "ugly" out of 10, but it gets the job done. This is the second idea I've come up with to solve this problem (yes, the first idea was actually worse than this one, if you can believe it). If you have a better idea, I definitely want to hear it, so send me a pull request.
