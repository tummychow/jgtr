package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const helpStr = `    jgtr - JSON Go Template Renderer

USAGE:
    jgtr [OPTIONS]

    jgtr consumes a data file and a template file written in Go's text/template
    language. The values available in the data file are then used to render the
    template file and generate output.

    By default, jgtr reads the data and template from stdin, and writes output
    to stdout. Note that data and template cannot both come from stdin - at
    least one of the two must be specified via an option.

    jgtr can consume data from JSON, YAML 1.1 or TOML v0.2.0. You can specify
    which type to use via an option. If no such option is given, jgtr attempts
    to guess from the extension of the data file (if any). If the format is
    still ambiguous, jgtr uses JSON as the default.

OPTIONS:
    -d FILE, --data=FILE
        Read data data from FILE. Specify "-" (the default) to use stdin.

    -t FILE, --template=FILE
        Read template from FILE. Specify "-" (the default) to use stdin.

    -o FILE, --output=FILE
        Write rendered template to FILE. Specify "-" (the default) to use
        stdout.

    -j, --json
    	Specify the data format as JSON (default).

    -y, --yaml
    	Specify the data format as YAML.

    -T, --toml
    	Specify the data format as TOML.

    -h, --help
        Display this help.

    -V, --version
        Display jgtr version.`

const versionStr = `0.6.0`

func main() {
	help := flag.BoolP("help", "h", false, "show help")
	version := flag.BoolP("version", "V", false, "show version")

	dataPath := flag.StringP("data", "d", "-", "data file (JSON by default)")
	tmplPath := flag.StringP("template", "t", "-", "Go template file")
	outPath := flag.StringP("output", "o", "-", "output file")

	jsonFlag := flag.BoolP("json", "j", false, "interpret data as JSON")
	yamlFlag := flag.BoolP("yaml", "y", false, "interpret data as YAML")
	tomlFlag := flag.BoolP("toml", "T", false, "interpret data as TOML")

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

	var data interface{} = nil
	var err error = nil
	if *yamlFlag { // check the flags first
		data, err = loadYAMLData(*dataPath)
	} else if *tomlFlag {
		data, err = loadTOMLData(*dataPath)
	} else if *jsonFlag {
		data, err = loadJSONData(*dataPath)
	} else if strings.HasSuffix(*dataPath, ".yaml") || strings.HasSuffix(*dataPath, ".yml") { // no flag? check the extension
		data, err = loadYAMLData(*dataPath)
	} else if strings.HasSuffix(*dataPath, ".toml") {
		data, err = loadTOMLData(*dataPath)
	} else { // default case: json
		data, err = loadJSONData(*dataPath)
	}
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
