package exercise5

import(
	
)



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
