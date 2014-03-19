package main

import (
	"text/template"
	"time"
)

var tmplFuncs = template.FuncMap{
	"timeParse": time.Parse,
}
