package main

import (
	"fmt"
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
	"sliceSortKey": func(s []interface{}, k string) []interface{} {
		ret := make([]interface{}, len(s))
		copy(ret, s)

		sort.Stable(MapSlice{ret, k})
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
	return valueLt(leftV, rightV)
}

// GenericSlice.Len and GenericSlice.Swap were stolen from standard implementations
// in pkg sort
func (p GenericSlice) Len() int      { return len(p) }
func (p GenericSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// another implementation of sort.Interface for slices of map[string]interface{}
// Key is the key into the maps, against which they will be sorted
type MapSlice struct {
	Slice []interface{}
	Key   string
}

func (p MapSlice) Less(i, j int) bool {
	// several layers of indirection in this code
	// you have to reflect.ValueOf on the []interface{}, which reveals that the
	// interface{} is a map[string]interface{}
	// then you have to MapIndex on that map, which gives you another Value,
	// whose Kind is interface (because the map contains interfaces)
	k := reflect.ValueOf(p.Key)
	leftV := reflect.ValueOf(p.Slice[i]).MapIndex(k)
	rightV := reflect.ValueOf(p.Slice[j]).MapIndex(k)

	// if the key was not in the map, the Value will be the zero Value
	// we consider this to be "less than everything"
	// ordering of these comparisons is important:
	// if they're both the zero Value, we return false, not true
	if !rightV.IsValid() {
		return false
	}
	if !leftV.IsValid() {
		return true
	}

	// if the key was in the map, then we have to indirect through the
	// interface using Elem, to get the underlying value
	return valueLt(leftV.Elem(), rightV.Elem())
}
func (p MapSlice) Len() int      { return len(p.Slice) }
func (p MapSlice) Swap(i, j int) { p.Slice[i], p.Slice[j] = p.Slice[j], p.Slice[i] }

// generic compare for pairs of ints, floats or strings
func valueLt(leftV, rightV reflect.Value) bool {
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
	panic(fmt.Errorf("Attempting to compare items of noncomparable types: %v %v", leftV.Kind(), rightV.Kind()))
}
