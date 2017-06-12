// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
package main

import (
	"fmt"
	"sort"
	"gopl.io/ch5/exercise5"
)

//!+table
// prereqs maps computer science courses to their prerequisites.
var pre = map[string][]string{
	"a": {"b", },
	"b": {"d", "e"},
	"d": {"f","g","j"},
	"c": {"h", "i"},
	"g": {"c"},
}

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},//circle

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"linear algebra":        {"calculus"},//circle
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

//!-table

//!+main
func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

//e5.14
	visitFunc := func(item string) []string  {
		order = append(order, item)
		return m[item]
	}
	exercise5.BreadthFirst(visitFunc, keys)


	return order
}

//!-main
