// file: etld.go
// description: manages effective top-level domains (eTLDs)

package gotld

import (
	"sort"
	"sync"
)

// ETLD manages all eTLDs in lists with thread-safety
type ETLD struct {
	List  []string
	Count int
	Dots  int
	mu    sync.RWMutex
}

// Add appends a new eTLD to the list if it doesn't already exist
func (e *ETLD) Add(s string, sortList bool) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check for duplicates
	for _, item := range e.List {
		if item == s {
			return false
		}
	}

	oldCount := e.Count
	e.List = append(e.List, s)
	e.Count = len(e.List)

	if sortList {
		e.Sort()
	}

	return e.Count > oldCount
}

// Sort will sort the list of strings
func (e *ETLD) Sort() {
	e.mu.Lock()
	defer e.mu.Unlock()

	sort.Strings(e.List)
}

// Search will return true if found as well as the eTLD from the list
func (e *ETLD) Search(str string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.Count == 0 {
		return "", false
	}

	idx := sort.Search(e.Count, func(i int) bool { return e.List[i] >= str })
	if idx < e.Count && e.List[idx] == str {
		return e.List[idx], true // Found (TRUE)
	}

	return "", false // NOT FOUND (FALSE)
}

// emptyETLD creates a new empty ETLD with the specified number of dots
func emptyETLD(dots int) *ETLD {
	return &ETLD{
		List:  make([]string, 0),
		Count: 0,
		Dots:  dots,
		mu:    sync.RWMutex{},
	}
}
