// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
package main

import (
	"fmt"
	"sort"
)
// "gopl.io/ch5/exercise5"

//!+table
// prereqs maps computer science courses to their prerequisites.

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
	seen := make(map[string]bool)


	//can report circle
	path := make(map[string]bool) //just for simple, stack structure is better
	var visitDepth func(items []string)
	visitDepth = func(items []string) {
		for _, item := range items {
			if path[item] {
				fmt.Printf("circle for %q in", item)
				for s := range path {
					fmt.Printf("%q ", s)
				}
				fmt.Println("")
			}
			if !seen[item] {
				seen[item] = true
				order = append(order, item)
				path[item] = true
				visitDepth(m[item])
				delete(path, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	visitDepth(keys)

	return order
}

//!-main
