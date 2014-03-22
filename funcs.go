package main

import (
	"reflect"
	"sort"
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

	"sliceSort": func(s []interface{}) []interface{} {
		ret := make([]interface{}, len(s))
		copy(ret, s)

		sort.Sort(GenericSlice(ret))
		return ret
	},
	"sliceReverse": func(s []interface{}) []interface{} {
		ret := make([]interface{}, len(s))
		copy(ret, s)

		for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
			ret[i], ret[j] = ret[j], ret[i]
		}
		return ret
	},
}

// to implement functions of pkg sort on the generic slice,
// we need to create a GenericSlice that implements sort.Interface
type GenericSlice []interface{}

func (p GenericSlice) Less(i, j int) bool {
	leftV := reflect.ValueOf(p[i])
	rightV := reflect.ValueOf(p[j])

	switch leftV.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return leftV.Int() < rightV.Int()
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return leftV.Uint() < rightV.Uint()
	case reflect.Float32, reflect.Float64:
		return leftV.Float() < rightV.Float()
	case reflect.String:
		return leftV.String() < rightV.String()
	}
	panic("Attempting to compare two items of noncomparable type")
}

// GenericSlice.Len and GenericSlice.Swap were stolen from standard implementations
// in pkg sort
func (p GenericSlice) Len() int      { return len(p) }
func (p GenericSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
