package main

import(
	"fmt"
	"sort"
	"os"
	"time"
	"text/tabwriter"
	"math/rand"
)

type Row struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}
type Table struct {
	Rows []*Row
	less func(x, y *Row) bool
}

func (t *Table) Swap(i,j int)  {
	t.Rows[i], t.Rows[j] = t.Rows[j], t.Rows[i]
}

func (t *Table) Len() int {
	return len(t.Rows)
}

func (t *Table) Less(i, j int) bool  {
	return t.less(t.Rows[i], t.Rows[j])
}
type Func func(x, y *Row) bool
var cmpFuncs = map[string]Func {
	"Artist" : func(x, y *Row) bool { return x.Artist < y.Artist },
	"Title"  : func(x, y *Row) bool { return x.Title <y.Title },
	"Album"  : func(x, y *Row) bool { return x.Album < y.Album },
	"Year"   : func(x, y *Row) bool { return x.Year < y.Year },
	"Length" : func(x, y *Row) bool { return x.Length < y.Length },
}



func TableSort (t *Table, columns []string) {
	less := func(x, y *Row) bool{
		for _, s := range columns {
			f := cmpFuncs[s]
			if f(x, y) || f(y, x) {
				return f(x, y)
			}
	 }
	 return false
	}
	t.less = less
	sort.Stable(t)
}

func TableRand (t *Table) {
	rand.Seed(time.Now().UnixNano())
	// for i := 0; i < 100; i++ {
	// 	fmt.Printf("%v ", rand.Int31() > int32(1 << 30))
	// }
	less := func(x, y *Row) bool {
	 return rand.Int31() >= int32(1 << 30)
	}
	t.less = less
	sort.Stable(t)
}

var tracks = []*Row{
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
	{"A", "A", "A", 1, length("3m38s")},
	{"A", "A", "B", 2, length("3m38s")},
	{"A", "B", "C", 3, length("3m38s")},
	{"B", "A", "B", 1, length("3m38s")},
	{"B", "A", "B", 2, length("3m38s")},
	{"B", "B", "B", 2, length("3m38s")},
	{"C", "A", "B", 2012, length("3m38s")},
	{"C", "B", "B", 2012, length("3m38s")},
	{"C", "C", "B", 2012, length("3m38s")},

}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

//!+printTracks
func printTable(table *Table) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tw, "")
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range table.Rows {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

func main()  {
	var t Table
	t.Rows = tracks
	TableRand(&t);TableRand(&t);TableRand(&t)
	printTable(&t)
	cs := []string{"Title", "Artist", "Album", "Year", "Length"}
	// cs = []string{"Title"}
	TableSort(&t, cs);TableSort(&t, cs)
	fmt.Printf("Table is sorted? %v\n", sort.IsSorted(&t))
	fmt.Printf("Table is sorted by %q\n", cs)
	printTable(&t)
}
