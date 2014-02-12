package main

import (
	"errors"
	"fmt"
	"sort"
	"text/template"
)

type Template struct {
	Text               string
	RequiresNumeric    bool
	RequiresComparable bool
	RequiresOrdered    bool
}

func getTemplate(name string) (result *template.Template, err error) {
	if isProjectionMethod(name) {
		return getProjectionTemplate(name)
	}
	return getStandardTemplate(name)
}

func getHeaderTemplate() *template.Template {
	return template.Must(template.New("header").Parse(header))
}

func getStandardTemplate(name string) (result *template.Template, err error) {
	t, found := StandardTemplates[name]
	if found {
		result = template.Must(template.New(name).Parse(t.Text))
	} else {
		err = fmt.Errorf("%s is not a known method", name)
	}
	return
}

func isStandardMethod(s string) bool {
	_, ok := StandardTemplates[s]
	return ok
}

func getStandardMethodKeys() (result []string) {
	for k := range StandardTemplates {
		result = append(result, k)
	}
	sort.Strings(result)
	return
}

func getProjectionTemplate(name string) (result *template.Template, err error) {
	t, found := ProjectionTemplates[name]
	if found {
		result = template.Must(template.New(name).Parse(t.Text))
	} else {
		err = errors.New(fmt.Sprintf("%s is not a known projection method", name))
	}
	return
}

func isProjectionMethod(s string) bool {
	_, ok := ProjectionTemplates[s]
	return ok
}

func getProjectionMethodKeys() (result []string) {
	for k := range ProjectionTemplates {
		result = append(result, k)
	}
	sort.Strings(result)
	return
}

func getSortSupportTemplate() *template.Template {
	return template.Must(template.New("sortSupport").Parse(sortSupport))
}

func getSortInterfaceTemplate() *template.Template {
	return template.Must(template.New("sortInterface").Parse(sortInterface))
}

func getContainerTemplate(name string) (result *template.Template, err error) {
	t, found := ContainerTemplates[name]
	if found {
		result = template.Must(template.New(name).Parse(t.Text))
	} else {
		err = errors.New(fmt.Sprintf("%s is not a known container", name))
	}
	return
}

const header = `// This file was auto-generated using github.com/clipperhouse/gen
// Modifying this file is not recommended as it will likely be overwritten in the future

// Sort (if included below) is a modification of http://golang.org/pkg/sort/#Sort
// List (if included below) is a modification of http://golang.org/pkg/container/list/
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package {{.Package.Name}}
{{if gt (len .Imports) 0}}
import ({{range .Imports}}
	"{{.}}"{{end}}
)
{{end}}
// {{.Plural}} is a slice of type {{.Pointer}}{{.Name}}, for use with gen methods below. Use this type where you would use []{{.Pointer}}{{.Name}}. (This is required because slices cannot be method receivers.)
type {{.Plural}} []{{.Pointer}}{{.Name}}
`

var StandardTemplates = map[string]*Template{

	"All": &Template{
		Text: `
// All verifies that all elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv {{.Plural}}) All(fn func({{.Pointer}}{{.Name}}) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}
`},

	"Any": &Template{
		Text: `
// Any verifies that one or more elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv {{.Plural}}) Any(fn func({{.Pointer}}{{.Name}}) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}
`},

	"Count": &Template{
		Text: `
// Count gives the number elements of {{.Plural}} that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func (rcv {{.Plural}}) Count(fn func({{.Pointer}}{{.Name}}) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}
`},

	"Distinct": &Template{
		Text: `
// Distinct returns a new {{.Plural}} slice whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
func (rcv {{.Plural}}) Distinct() (result {{.Plural}}) {
	appended := make(map[{{.Pointer}}{{.Name}}]bool)
	for _, v := range rcv {
		if !appended[v] {
			result = append(result, v)
			appended[v] = true
		}
	}
	return result
}
`,
		RequiresComparable: true,
	},

	"DistinctBy": &Template{
		Text: `
// DistinctBy returns a new {{.Plural}} slice whose elements are unique, where equality is defined by a passed func. See: http://clipperhouse.github.io/gen/#DistinctBy
func (rcv {{.Plural}}) DistinctBy(equal func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Plural}}) {
	for _, v := range rcv {
		eq := func(_app {{.Pointer}}{{.Name}}) bool {
			return equal(v, _app)
		}
		if !result.Any(eq) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"Each": &Template{
		Text: `
// Each iterates over {{.Plural}} and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func (rcv {{.Plural}}) Each(fn func({{.Pointer}}{{.Name}})) {
	for _, v := range rcv {
		fn(v)
	}
}
`},

	"First": &Template{
		Text: `
// First returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
func (rcv {{.Plural}}) First(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	for _, v := range rcv {
		if fn(v) {
			result = v
			return
		}
	}
	err = errors.New("no {{.Plural}} elements return true for passed func")
	return
}
`},

	"Max": &Template{
		Text: `
// Max returns the maximum value of {{.Plural}}. In the case of multiple items being equally maximal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
func (rcv {{.Plural}}) Max() (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Max of an empty slice")
		return
	}
	result = rcv[0]
	for _, v := range rcv {
		if v > result {
			result = v
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"Min": &Template{
		Text: `
// Min returns the minimum value of {{.Plural}}. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Min
func (rcv {{.Plural}}) Min() (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Min of an empty slice")
		return
	}
	result = rcv[0]
	for _, v := range rcv {
		if v < result {
			result = v
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"MaxBy": &Template{
		Text: `
// MaxBy returns an element of {{.Plural}} containing the maximum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally maximal, the last such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MaxBy
func (rcv {{.Plural}}) MaxBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the MaxBy of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if rcv[i] != rcv[m] && !less(rcv[i], rcv[m]) {
			m = i
		}
	}
	result = rcv[m]
	return
}
`},

	"MinBy": &Template{
		Text: `
// MinBy returns an element of {{.Plural}} containing the minimum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MinBy
func (rcv {{.Plural}}) MinBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Min of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if less(rcv[i], rcv[m]) {
			m = i
		}
	}
	result = rcv[m]
	return
}
`},

	"Single": &Template{
		Text: `
// Single returns exactly one element of {{.Plural}} that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
func (rcv {{.Plural}}) Single(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	var candidate {{.Pointer}}{{.Name}}
	found := false
	for _, v := range rcv {
		if fn(v) {
			if found {
				err = errors.New("multiple {{.Plural}} elements return true for passed func")
				return
			}
			candidate = v
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("no {{.Plural}} elements return true for passed func")
	}
	return
}
`},

	"Where": &Template{
		Text: `
// Where returns a new {{.Plural}} slice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv {{.Plural}}) Where(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Plural}}) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"Sort": &Template{
		Text: `
// Sort returns a new ordered {{.Plural}} slice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) Sort() {{.Plural}} {
	result := make({{.Plural}}, len(rcv))
	copy(result, rcv)
	sort.Sort(result)
	return result
}
`,
		RequiresOrdered: true,
	},
	"IsSorted": &Template{
		Text: `
// IsSorted reports whether {{.Plural}} is sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) IsSorted() bool {
	return sort.IsSorted(rcv)
}
`,
		RequiresOrdered: true,
	},
	"SortDesc": &Template{
		Text: `
// SortDesc returns a new reverse-ordered {{.Plural}} slice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) SortDesc() {{.Plural}} {
	result := make({{.Plural}}, len(rcv))
	copy(result, rcv)
	sort.Sort(sort.Reverse(result))
	return result
}
`,
		RequiresOrdered: true,
	},
	"IsSortedDesc": &Template{
		Text: `
// IsSortedDesc reports whether {{.Plural}} is reverse-sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) IsSortedDesc() bool {
	return sort.IsSorted(sort.Reverse(rcv))
}
`,
		RequiresOrdered: true,
	},

	"SortBy": &Template{
		Text: `
// SortBy returns a new ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) SortBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) {{.Plural}} {
	result := make({{.Plural}}, len(rcv))
	copy(result, rcv)
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort{{.Plural}}(result, less, 0, n, maxDepth)
	return result
}
`},

	"IsSortedBy": &Template{
		Text: `
// IsSortedBy reports whether an instance of {{.Plural}} is sorted, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) IsSortedBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv[i], rcv[i-1]) {
			return false
		}
	}
	return true
}
`},

	"SortByDesc": &Template{
		Text: `
// SortByDesc returns a new, descending-ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) SortByDesc(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) {{.Plural}} {
	greater := func(a, b {{.Pointer}}{{.Name}}) bool {
		return a != b && !less(a, b)
	}
	return rcv.SortBy(greater)
}
`},

	"IsSortedByDesc": &Template{
		Text: `
// IsSortedDesc reports whether an instance of {{.Plural}} is sorted in descending order, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) IsSortedByDesc(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) bool {
	greater := func(a, b {{.Pointer}}{{.Name}}) bool {
		return a != b && !less(a, b)
	}
	return rcv.IsSortedBy(greater)
}
`},
}

const sortInterface = `
func (rcv {{.Plural}}) Len() int {
	return len(rcv)
}
func (rcv {{.Plural}}) Less(i, j int) bool {
	return rcv[i] < rcv[j]
}
func (rcv {{.Plural}}) Swap(i, j int) {
	rcv[i], rcv[j] = rcv[j], rcv[i]
}
`

const sortSupport = `
// Sort support methods

func swap{{.Plural}}(rcv {{.Plural}}, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

// Insertion sort
func insertionSort{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv[j], rcv[j-1]); j-- {
			swap{{.Plural}}(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(rcv[first+child], rcv[first+child+1]) {
			child++
		}
		if !less(rcv[first+root], rcv[first+child]) {
			return
		}
		swap{{.Plural}}(rcv, first+root, first+child)
		root = child
	}
}

func heapSort{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown{{.Plural}}(rcv, less, i, hi, first)
	}

	// Pop elements, largest first, into end of rcv.
	for i := hi - 1; i >= 0; i-- {
		swap{{.Plural}}(rcv, first, first+i)
		siftDown{{.Plural}}(rcv, less, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values rcv[a], rcv[b], rcv[c] into rcv[a].
func medianOfThree{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv[m1], rcv[m0]) {
		swap{{.Plural}}(rcv, m1, m0)
	}
	if less(rcv[m2], rcv[m1]) {
		swap{{.Plural}}(rcv, m2, m1)
	}
	if less(rcv[m1], rcv[m0]) {
		swap{{.Plural}}(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRange{{.Plural}}(rcv {{.Plural}}, a, b, n int) {
	for i := 0; i < n; i++ {
		swap{{.Plural}}(rcv, a+i, b+i)
	}
}

func doPivot{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree{{.Plural}}(rcv, less, lo, lo+s, lo+2*s)
		medianOfThree{{.Plural}}(rcv, less, m, m-s, m+s)
		medianOfThree{{.Plural}}(rcv, less, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree{{.Plural}}(rcv, less, lo, m, hi-1)

	// Invariants are:
	//	rcv[lo] = pivot (set up by ChoosePivot)
	//	rcv[lo <= i < a] = pivot
	//	rcv[a <= i < b] < pivot
	//	rcv[b <= i < c] is unexamined
	//	rcv[c <= i < d] > pivot
	//	rcv[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if less(rcv[b], rcv[pivot]) { // rcv[b] < pivot
				b++
			} else if !less(rcv[pivot], rcv[b]) { // rcv[b] = pivot
				swap{{.Plural}}(rcv, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less(rcv[pivot], rcv[c-1]) { // rcv[c-1] > pivot
				c--
			} else if !less(rcv[c-1], rcv[pivot]) { // rcv[c-1] = pivot
				swap{{.Plural}}(rcv, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// rcv[b] > pivot; rcv[c-1] < pivot
		swap{{.Plural}}(rcv, b, c-1)
		b++
		c--
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	n := min(b-a, a-lo)
	swapRange{{.Plural}}(rcv, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange{{.Plural}}(rcv, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort{{.Plural}}(rcv {{.Plural}}, less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort{{.Plural}}(rcv, less, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot{{.Plural}}(rcv, less, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort{{.Plural}}(rcv, less, a, mlo, maxDepth)
			a = mhi // i.e., quickSort{{.Plural}}(rcv, mhi, b)
		} else {
			quickSort{{.Plural}}(rcv, less, mhi, b, maxDepth)
			b = mlo // i.e., quickSort{{.Plural}}(rcv, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort{{.Plural}}(rcv, less, a, b)
	}
}
`

var ProjectionTemplates = map[string]*Template{
	"Aggregate": &Template{
		Text: `
// {{.MethodName}} iterates over {{.Parent.Plural}}, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Type}}, {{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}) {
	for _, v := range rcv {
		result = fn(result, v)
	}
	return
}
`},

	"Average": &Template{
		Text: `
// {{.MethodName}} sums {{.Type}} over all elements and divides by len({{.Parent.Plural}}). See: http://clipperhouse.github.io/gen/#Average
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine {{.MethodName}} of zero-length {{.Parent.Plural}}")
		return
	}
	for _, v := range rcv {
		result += fn(v)
	}
	result = result / {{.Type}}(l)
	return
}
`,
		RequiresNumeric: true,
	},

	"GroupBy": &Template{
		Text: `
// {{.MethodName}} groups elements into a map keyed by {{.Type}}. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) map[{{.Type}}]{{.Parent.Plural}} {
	result := make(map[{{.Type}}]{{.Parent.Plural}})
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}
`,
		RequiresComparable: true,
	},

	"Max": &Template{
		Text: `
// {{.MethodName}} selects the largest value of {{.Type}} in {{.Parent.Plural}}. Returns error on {{.Parent.Plural}} with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine {{.MethodName}} of zero-length {{.Parent.Plural}}")
		return
	}
	result = fn(rcv[0])
	if l > 1 {
		for _, v := range rcv[1:] {
			f := fn(v)
			if f > result {
				result = f
			}
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"Min": &Template{
		Text: `
// {{.MethodName}} selects the least value of {{.Type}} in {{.Parent.Plural}}. Returns error on {{.Parent.Plural}} with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine {{.MethodName}} of zero-length {{.Parent.Plural}}")
		return
	}
	result = fn(rcv[0])
	if l > 1 {
		for _, v := range rcv[1:] {
			f := fn(v)
			if f < result {
				result = f
			}
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"Select": &Template{
		Text: `
// {{.MethodName}} returns a slice of {{.Type}} in {{.Parent.Plural}}, projected by passed func. See: http://clipperhouse.github.io/gen/#Select
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result []{{.Type}}) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}
`,
	},

	"Sum": &Template{
		Text: `
// {{.MethodName}} sums {{.Type}} over elements in {{.Parent.Plural}}. See: http://clipperhouse.github.io/gen/#Sum
func (rcv {{.Parent.Plural}}) {{.MethodName}}(fn func({{.Parent.Pointer}}{{.Parent.Name}}) {{.Type}}) (result {{.Type}}) {
	for _, v := range rcv {
		result += fn(v)
	}
	return
}
`,
		RequiresNumeric: true,
	},
}

var ContainerTemplates = map[string]*Template{

	"List": &Template{
		Text: `
// {{.Name}}Element is an element of a linked list.
type {{.Name}}Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *{{.Name}}Element

	// The list to which this element belongs.
	list *{{.Name}}List

	// The value stored with this element.
	Value {{.Pointer}}{{.Name}}
}

// Next returns the next list element or nil.
func (e *{{.Name}}Element) Next() *{{.Name}}Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *{{.Name}}Element) Prev() *{{.Name}}Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// {{.Name}}List represents a doubly linked list.
// The zero value for {{.Name}}List is an empty list ready to use.
type {{.Name}}List struct {
	root {{.Name}}Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *{{.Name}}List) Init() *{{.Name}}List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New() *{{.Name}}List { return new({{.Name}}List).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *{{.Name}}List) Len() int { return l.len }

// Front returns the first element of list l or nil
func (l *{{.Name}}List) Front() *{{.Name}}Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil.
func (l *{{.Name}}List) Back() *{{.Name}}Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero {{.Name}}List value.
func (l *{{.Name}}List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *{{.Name}}List) insert(e, at *{{.Name}}Element) *{{.Name}}Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&{{.Name}}Element{Value: v}, at).
func (l *{{.Name}}List) insertValue(v {{.Pointer}}{{.Name}}, at *{{.Name}}Element) *{{.Name}}Element {
	return l.insert(&{{.Name}}Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *{{.Name}}List) remove(e *{{.Name}}Element) *{{.Name}}Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
func (l *{{.Name}}List) Remove(e *{{.Name}}Element) {{.Pointer}}{{.Name}} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero {{.Name}}Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *{{.Name}}List) PushFront(v {{.Pointer}}{{.Name}}) *{{.Name}}Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *{{.Name}}List) PushBack(v {{.Pointer}}{{.Name}}) *{{.Name}}Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *{{.Name}}List) InsertBefore(v {{.Pointer}}{{.Name}}, mark *{{.Name}}Element) *{{.Name}}Element {
	if mark.list != l {
		return nil
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *{{.Name}}List) InsertAfter(v {{.Pointer}}{{.Name}}, mark *{{.Name}}Element) *{{.Name}}Element {
	if mark.list != l {
		return nil
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
func (l *{{.Name}}List) MoveToFront(e *{{.Name}}Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	l.insert(l.remove(e), &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
func (l *{{.Name}}List) MoveToBack(e *{{.Name}}Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	l.insert(l.remove(e), l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e is not an element of l, or e == mark, the list is not modified.
func (l *{{.Name}}List) MoveBefore(e, mark *{{.Name}}Element) {
	if e.list != l || e == mark {
		return
	}
	l.insert(l.remove(e), mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e is not an element of l, or e == mark, the list is not modified.
func (l *{{.Name}}List) MoveAfter(e, mark *{{.Name}}Element) {
	if e.list != l || e == mark {
		return
	}
	l.insert(l.remove(e), mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.
func (l *{{.Name}}List) PushBackList(other *{{.Name}}List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.
func (l *{{.Name}}List) PushFrontList(other *{{.Name}}List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
`},
}
