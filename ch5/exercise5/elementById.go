package exercise5

import (
	// "fmt"
	//
	"golang.org/x/net/html"
	// "bytes"
	// "net/http"
	// "time"
	//
	// "log"
	// // "gopl.io/ch5/exercise5"
	// "strings"
)

//bad coding
func forEachNodeV2(n *html.Node, pre, post func(*html.Node) *html.Node) *html.Node {
	if pre != nil {
		if ret := pre(n); ret != nil {
			return ret
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ret := forEachNodeV2(c, pre, post); ret != nil {
			return ret
		}
	}

	if post != nil {
		if ret := post(n); ret != nil {
			return ret
		}
	}

	return nil
}

func forEachNodeAbort(n *html.Node, pre, post func(*html.Node) bool) {
	if pre != nil && pre(n) == true {
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNodeAbort(c, pre, post)
	}

	if post != nil && post(n) == true {
		return
	}
}

//e5.7
func ElementById(doc *html.Node, id string) (nodeRet *html.Node) {
	findId := func(n *html.Node) bool {
		for _, v := range n.Attr {
			if v.Key == "id" && v.Val == id {
				nodeRet = n
				return true
			}
		}
		return false
	}
	forEachNodeAbort(doc, findId, nil)
	return
}
