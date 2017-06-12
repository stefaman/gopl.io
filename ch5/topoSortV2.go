// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
package main

import (
	"fmt"
	"sort"
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
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"linear algebra":        {"calculus"},
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
	visit := make(map[string]bool)

	seen = make(map[string]bool)
	var visitW func(items []string)
	visitW = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				order = append(order, item)
				seen[item] = true
				fmt.Printf("%s seen\n", item)
			}
		}
		for _, item := range items {
			if !visit[item] {
				visitW(m[item])
				visit[item] = true
				fmt.Printf("%s visited\n", item)
			}
		}
	}

	seen = make(map[string]bool)
	var visitD func(items []string)
	visitD = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				order = append(order, item)
				seen[item] = true
				visitD(m[item])
			}
		}
	}

	// seen := make(map[string]bool)
	// order := []string
	var visitNode func(node string)
	visitNode = func(node string){
		// order = append(order, item)
		for _, item := range m[node] {
			if !seen[item]{
				order = append(order, item)
				seen[item] = true
			}
		}
		for _, item := range m[node] {
			if !visit[item]{
				visitNode(item)
				visit[item] = true
				// fmt.Printf("%s visited\n", item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	m["root"] = keys

	sort.Strings(keys)
	// visitNode("root")
	// visitW(keys)
	// visitD(keys)

	visitFunc := func(item string) []string  {
		order = append(order, item)
		return m[item]
	}

	BreadthFirst(visitFunc, keys)
	return order
}

//!-main

func BreadthFirst(f func(string) []string, worklist [] string){
	seen := make(map[string]bool)
	for len(worklist) >0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item]{
				worklist = append(worklist, f(item)...)
				seen[item] = true
			}
		}
	}
}
