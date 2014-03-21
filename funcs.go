package main

import (
	"strings"
	"text/template"
	"time"
)

var tmplFuncs = template.FuncMap{
	"notLast": func(slice []interface{}, i int) bool {
		return (i != len(slice)-1)
	},

	"timeParse": time.Parse,
	"timeNow":   time.Now,

	"stringSplit":  strings.Split,
	"stringFields": strings.Fields,
	"stringJoin":   strings.Join,
	"stringUpper":  strings.ToUpper,
	"stringLower":  strings.ToLower,
	"stringTitle":  strings.ToTitle,
}
