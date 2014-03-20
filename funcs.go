package main

import (
	"text/template"
	"time"
)

var tmplFuncs = template.FuncMap{
	"timeParse": time.Parse,
	"timeNow":   time.Now,
	"notLast": func(slice []interface{}, i int) bool {
		return (i != len(slice)-1)
	},
}
