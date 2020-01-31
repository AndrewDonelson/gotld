package gotld

import (
	"sort"
)

// ETLD manages all etlds in lists
type ETLD struct {
	List  []string
	Count int
	Dots  int
}

// Add appends a new etld to the list
func (e *ETLD) Add(s string, sort bool) bool {
	// TODO: Add check for duplicates and skip
	oldCount := e.Count
	e.List = append(e.List, s)
	e.Count = len(e.List)
	if sort {
		e.Sort()
	}
	return e.Count > oldCount
}

// Sort will sorth the list of strings
func (e *ETLD) Sort() {
	sort.Strings(e.List)
}

// Search will return true if found as well as the etld from the list
func (e *ETLD) Search(str string) (string, bool) {
	idx := sort.Search(e.Count, func(i int) bool { return e.List[i] >= str })
	if idx < e.Count && e.List[idx] == str {
		return e.List[idx], true // Found (TRUE)
	}
	return "", false // NO FOUND (FALSE)
}

func emptyETLD(dots int) *ETLD {
	return &ETLD{Count: 0, Dots: dots}
} 
