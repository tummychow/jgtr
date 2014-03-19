package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"os"
	"text/template"
)

const helpStr = `    jgtr - JSON Go Template Renderer

USAGE:
    jgtr [OPTIONS]

    jgtr consumes a JSON-encoded data file and a template file written in Go's
    text/template language. The values available in the data file are then used
    to render the template file and generate output.

    By default, jgtr reads the data and template from stdin, and writes output
    to stdout. Note that data and template cannot both come from stdin - at
    least one of the two must be specified via an option.

OPTIONS:
    -j FILE, --json=FILE
        Read JSON data from FILE. Specify "-" (the default) to use stdin.

    -t FILE, --template=FILE
        Read template from FILE. Specify "-" (the default) to use stdin.

    -o FILE, --output=FILE
        Write rendered template to FILE. Specify "-" (the default) to use
        stdout.

    -h, --help
        Display this help.

    -V, --version
        Display jgtr version.`

const versionStr = `0.2.0`

func main() {
	help := flag.BoolP("help", "h", false, "show help")
	version := flag.BoolP("version", "V", false, "show version")

	dataPath := flag.StringP("json", "j", "-", "JSON data file")
	tmplPath := flag.StringP("template", "t", "-", "Go template file")
	outPath := flag.StringP("output", "o", "-", "output file")

	flag.Parse()

	if *help {
		fmt.Println(helpStr)
		return
	}
	if *version {
		println(versionStr)
		return
	}

	if *dataPath == "-" && *tmplPath == "-" {
		println("Cannot use stdin for data and template simultaneously")
		os.Exit(1)
	}

	data, err := loadJSONData(*dataPath)
	if err != nil {
		panic(err)
	}
	tmpl, err := loadGoTemplate(*tmplPath)
	if err != nil {
		panic(err)
	}

	outFile, err := createStream(*outPath)
	if err != nil {
		panic(err)
	}
	defer closeStream(outFile)

	err = tmpl.Execute(outFile, data)
	if err != nil {
		panic(err)
	}
}

// loadJSONData unmarshals JSON-encoded data from the file specified by path,
// and returns the result as an interface{}. If the path is "-", then data will
// be acquired from os.Stdin.
func loadJSONData(path string) (ret interface{}, err error) {
	file, err := openStream(path)
	if err != nil {
		return
	}
	defer closeStream(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ret)
	return // whether err==nil or not, our work is done
}

// loadGoTemplate parses a Go text template from the file specified by path.
// The file contents are parsed into a top-level template with the name "root".
// If the path is "-", then the template will be parsed from os.Stdin.
func loadGoTemplate(path string) (tmpl *template.Template, err error) {
	file, err := openStream(path)
	if err != nil {
		return
	}
	defer closeStream(file)

	rawTmpl, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	// explicitly parse the raw string rather than using ParseFiles
	// ParseFiles creates templates whose names are those of the files
	// and associates them with the parent template
	// this creates some confusing behavior in Template.Parse
	// see http://stackoverflow.com/questions/11805356/text-template-issue-parse-vs-parsefiles
	// also note: functions have to be added before the template is parsed
	tmpl, err = template.New("root").Funcs(tmplFuncs).Parse(string(rawTmpl))
	return // again, whether err==nil or not, this is finished
}

// openStream behaves like os.Open, except that if the path is "-", then it
// simply returns os.Stdin.
func openStream(path string) (file *os.File, err error) {
	if path == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
	}
	return
}

// createStream behaves like os.Create, except that if the path is "-", then it
// simply returns os.Stdout.
func createStream(path string) (file *os.File, err error) {
	if path == "-" {
		file = os.Stdout
	} else {
		file, err = os.Create(path)
	}
	return
}

// closeStream behaves like file.Close, except that if the file is os.Stdin or
// os.Stdout, it does nothing.
func closeStream(file *os.File) (err error) {
	if file == os.Stdout || file == os.Stdin {
		return
	}
	return file.Close()
}
