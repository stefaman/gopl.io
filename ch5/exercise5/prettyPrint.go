package exercise5

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
	"bytes"
	"strings"
)

func forEachNodeIgnore(n *html.Node, pre, post func(*html.Node) bool)  {
	if pre != nil && pre(n) == true {
		goto cont
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNodeIgnore(c, pre, post)
	}

	cont:
	if post != nil {
		post(n)
	}
}

/*e5.7
	print indent, ignore text of "script" and "sytle" tags, ignore space text
	print "<img />" style for tags which have no children
*/
func Pretty(n *html.Node)  {
	if n == nil {
		fmt.Fprintf(os.Stderr, "pretty use with nil node\n")
		return
	}
	forEachNodeIgnore(n, prettyPre, prettyPost)
}

var prettyeDepth int

//bad coding
func prettyPre(n *html.Node) (b bool) {
	if n.Type == html.DocumentNode ||
			n.Type == html.DoctypeNode ||
			n.Type == html.ErrorNode {
		return
	}
	//print indent
	buf :=new(bytes.Buffer)
	fmt.Fprintf(buf, "%*s", prettyeDepth * 2, "")
	if n.Type == html.TextNode {
		str := strings.Trim(n.Data, " \t\n\r")
		if len(str) == 0 {//aviod print indent and a "\n" for space text
			prettyeDepth++
			return true
		}
		buf.WriteString(str)
	}
	if n.Type == html.ElementNode{
		fmt.Fprintf(buf, "<%s", n.Data)
		for _, v := range n.Attr {
			if true || v.Key == "href"{
				fmt.Fprintf(buf, " %s=%q", v.Key, v.Val)
			}
		}
		if n.FirstChild == nil{
			buf.WriteString("/>")
		}else{
			buf.WriteString(">")
		}
		if n.Data == "style" || n.Data == "script" {
			b = true //ignore children nodes
		}
	}
	if n.Type == html.CommentNode {
		fmt.Fprintf(buf, "<!--%s", n.Data)
	}

	fmt.Printf("%s\n", buf.String())
	prettyeDepth++
	return
}

//bad coding
func prettyPost(n *html.Node) (b bool) {
	if n.Type == html.DocumentNode ||
	n.Type == html.DoctypeNode ||
	n.Type == html.ErrorNode {
	return
	}
	prettyeDepth--
	if n.Type == html.TextNode {
		return
	}
	buf :=new(bytes.Buffer)
	fmt.Fprintf(buf, "%*s", prettyeDepth * 2, "")
	if n.Type == html.ElementNode {
		if n.FirstChild == nil{
			return
		}
		fmt.Fprintf(buf, "</%s>", n.Data)
	}
	if n.Type == html.CommentNode {
		fmt.Fprintf(buf, "-->")
	}
	fmt.Printf("%s\n", buf.String())
	return
}
