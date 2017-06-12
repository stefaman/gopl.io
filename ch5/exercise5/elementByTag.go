package exercise5

import(
	"golang.org/x/net/html"
)



func forEachNode(node *html.Node, pre, post func(*html.Node)){
	if pre != nil{
		pre(node)
	}

	for c:= node.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(node)
	}
}

func ElementByTag(node *html.Node, tag string) []*html.Node {
	elements := []*html.Node{}
	find := func(n *html.Node){
		if n.Type == html.ElementNode && n.Data == tag {
			elements = append(elements, n)
		}
	}
	forEachNode(node, find, nil)
	return elements
}
