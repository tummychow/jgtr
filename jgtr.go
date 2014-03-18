package main

import (
	"encoding/json"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"os"
	"text/template"
)

func main() {
	dataPath := flag.StringP("json", "j", "-", "JSON data file")
	tmplPath := flag.StringP("template", "t", "-", "Go template file")
	outPath := flag.StringP("output", "o", "-", "output file")

	flag.Parse()

	data, err := loadJSONData(*dataPath)
	if err != nil {
		panic(err)
	}
	tmpl, err := loadGoTemplate(*tmplPath)
	if err != nil {
		panic(err)
	}

	var outFile *os.File
	if *outPath == "-" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(*outPath)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()
	}

	err = tmpl.Execute(outFile, data)
	if err != nil {
		panic(err)
	}
}

// loadJSONData unmarshals JSON-encoded data from the file specified by path,
// and returns the result as an interface{}. If the path is "-", then data will
// be acquired from os.Stdin.
func loadJSONData(path string) (ret interface{}, err error) {
	var file *os.File
	if path == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		if err != nil {
			return
		}
		defer file.Close()
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ret)
	return // whether err==nil or not, our work is done
}

// loadGoTemplate parses a Go text template from the file specified by path.
// The file contents are parsed into a top-level template with the name "root".
// If the path is "-", then the template will be parsed from os.Stdin.
func loadGoTemplate(path string) (tmpl *template.Template, err error) {
	var file *os.File
	if path == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		if err != nil {
			return
		}
		defer file.Close()
	}

	rawTmpl, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	// explicitly parse the raw string rather than using ParseFiles
	// ParseFiles creates templates whose names are those of the files
	// and associates them with the parent template
	// this creates some confusing behavior in Template.Parse
	// see http://stackoverflow.com/questions/11805356/text-template-issue-parse-vs-parsefiles
	tmpl, err = template.New("root").Parse(string(rawTmpl))
	return // again, whether err==nil or not, this is finished
}
