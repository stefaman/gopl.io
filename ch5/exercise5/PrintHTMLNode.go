package exercise5

import (
	"fmt"
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

//!+e5.4 page124
//print indent, "<img />" style for non-children nodes
func PrintHTMLNode(n *html.Node, prefix, indent string) []string {
	depth := 0
	lines := []string{}
	pIndent := func(buf *bytes.Buffer){
		buf.WriteString(prefix)
		for i := 0; i < depth; i++ {
			buf.WriteString(indent)
		}
	}
	pre := func(n *html.Node) {
		//print indent
		buf :=new(bytes.Buffer)
		pIndent(buf)
		switch n.Type {
			case html.TextNode:
				//aviod print indent and a "\n" for space text
				str := strings.Trim(n.Data, " \t\n\r")
				if len(str) == 0 {
					return
				}
				buf.WriteString(str)
			case html.ElementNode:
				fmt.Fprintf(buf, "<%s", n.Data)
				for _, v := range n.Attr {
					if true || v.Key == "href"{
						fmt.Fprintf(buf, " %s=%q", v.Key, v.Val)
					}
				}
				if n.FirstChild == nil{
					buf.WriteString(" />")
				}else{
					buf.WriteString(">")
				}
				depth++
			case html.CommentNode:
				fmt.Fprintf(buf, "<!--%s", n.Data)
				depth++
			default:
				return
		}
		buf.WriteString("\n")
		lines = append(lines, buf.String())
	}
	post := func(n *html.Node) {
		buf :=new(bytes.Buffer)
		depth--
		pIndent(buf)
		switch n.Type {
			case html.ElementNode:
				if n.FirstChild == nil{
					return
				}
				fmt.Fprintf(buf, "</%s>", n.Data)
			case html.CommentNode:
				fmt.Fprintf(buf, "-->")
			default:
				depth++//for blance
				return
		}
		buf.WriteString("\n")
		lines = append(lines, buf.String())
		return
	}

	forEachNode(n, pre, post)
	return lines
}

//bad conding
var level int
func PrintHTML_bad(n *html.Node, pre, indent string) []string {
	return printTree(nil, n, pre, indent)
}

func printTree(tree []string, n *html.Node, pre, indent string) []string {

	tree = printNodePre(tree, pre, indent, n)

	if true && (n.Data != "script" && n.Data != "style") {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			tree = printTree(tree, c, pre, indent)
		}
	}

	tree = printNodePost(tree, pre, indent, n)
	return tree
}

func printNodePre(tree []string, pre, indent string, n *html.Node) []string  {
	// if n.Type == html.CommentNode ||
	// 		n.Type == html.DocumentNode ||
	// 		n.Type == html.DoctypeNode ||
	// 		n.Type == html.ErrorNode {
	// 	return tree
	// }


	//print indent
	buf :=new(bytes.Buffer)
	// fmt.Fprintf(buf, "%s%*s", pre, level*2, "")
	buf.WriteString(pre)
	for i := 0; i < level; i++ {
		buf.WriteString(indent)
	}

	//print text
	if n.Type == html.TextNode {
		fmt.Fprintf(buf, "%s\n", n.Data)
	}

	//print tag
	if n.Type == html.ElementNode || true {
		fmt.Fprintf(buf, "<%s ", n.Data)

		if n.Data != "script" || true {
			for _, v := range n.Attr {
				if true || v.Key == "href"{
					fmt.Fprintf(buf, "%s=%q ", v.Key, v.Val)
				}
			}
		}

		buf.WriteString(">\n")
	}
	tree = append(tree, buf.String()) // add <tag>
	level++
	return tree
}

func printNodePost(tree []string, pre, indent string, n *html.Node) []string  {
	// if n.Type == html.CommentNode ||
	// 		n.Type == html.DocumentNode ||
	// 		n.Type == html.DoctypeNode ||
	// 		n.Type == html.ErrorNode {
	// 	return tree
	// }

	level--
	if n.Type == html.TextNode {
		return tree
	}


	buf :=new(bytes.Buffer)
	// fmt.Fprintf(buf, "%s%*s", pre, level*2, "")
	buf.WriteString(pre)
	for i := 0; i < level; i++ {
		buf.WriteString(indent)
	}

	if n.Type == html.ElementNode || true{
		fmt.Fprintf(buf, "</%s>\n", n.Data)
	}

	tree = append(tree, buf.String()) // add </tag>
	return tree
}
