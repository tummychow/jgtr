package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"
)

// tmplFuncs contains the function map for templates that will be rendered by
// jgtr.
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

	// sliceSort sorts a slice of empty interfaces, using the implementation
	// provided by GenericSlice.
	"sliceSort": func(s []interface{}) []interface{} {
		ret := make([]interface{}, len(s))
		copy(ret, s)

		sort.Sort(GenericSlice(ret))
		return ret
	},
	// sliceSortKey sorts a slice of empty interfaces on the given map key,
	// using the implementation provided by MapSlice.
	"sliceSortKey": func(s []interface{}, k string) []interface{} {
		ret := make([]interface{}, len(s))
		copy(ret, s)

		sort.Stable(MapSlice{ret, k})
		return ret
	},
	// sliceReverse reverses a slice of empty interfaces.
	"sliceReverse": func(s []interface{}) []interface{} {
		ret := make([]interface{}, len(s))
		copy(ret, s)

		for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
			ret[i], ret[j] = ret[j], ret[i]
		}
		return ret
	},
}

// GenericSlice implements sort.Interface on a slice of interface{}, with the
// assumption that all of those interface{} have the same underlying type, and
// that type is some kind of int, float or string. The ordering that
// GenericSlice imposes on those elements is the expected ordering (ie sorting
// by ascending magnitude for ints/floats, and lexicographically for strings).
type GenericSlice []interface{}

// Less compares two elements of a GenericSlice. If the two elements are
// comparable (ie they are both ints, uints, floats or strings), then Less will
// return the result of a less-than comparison between those two elements. If
// the two elements are not comparable (ie their types do not match, or they do
// not have a natural ordering), then Less will panic.
func (p GenericSlice) Less(i, j int) bool {
	leftV := reflect.ValueOf(p[i])
	rightV := reflect.ValueOf(p[j])
	return valueLt(leftV, rightV)
}

func (p GenericSlice) Len() int      { return len(p) }
func (p GenericSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// MapSlice implements sort.Interface on a slice of interface{}, with the
// assumption that all of those interface{} have the underlying type
// map[string]interface{}. The ordering that MapSlice imposes on those maps
// is the standard ordering of a certain key, specified by Key. All maps in
// Slice must either omit that key, or contain a value for that key, such
// that all the maps contain the same type of value for that key.
type MapSlice struct {
	Slice []interface{}
	Key   string
}

// Less compares two elements of a MapSlice. If both elements contain the key
// specified by Key, then Less will return the result of comparing the values
// associated with Key. If those values are of noncomparable types, Less will
// panic.
// If an element does not contain the key specified by Key, then it should sort
// before the other element.
// If both elements do not contain the key specified by Key, then they are
// considered to be equal in ordering, so Less will return false.
// If either element is not actually a map[string]interface{}, Less will panic.
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

// valueLT performs a less-than comparison on two reflect.Values. The Values
// must both be valid receivers for one of the following reflection methods:
// Value.Int, Value.Uint, Value.Float, or Value.String. This function will panic
// if the Values do not have comparable underlying values, or if their Kinds
// are not comparable in this way.
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
