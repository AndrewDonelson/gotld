// file: gotld.go
// description: main package file with exported functions

package gotld

import (
	"sync"
)

// FQDNManager is the main interface for the GoTLD package
type FQDNManager interface {
	GetFQDN(url string) (string, error)
}

var (
	// Global manager instance
	manager     FQDNManager
	managerOnce sync.Once
	managerErr  error
)

// Init initializes the GoTLD package with custom options
func Init(opts *Options) error {
	var err error

	managerOnce.Do(func() {
		manager, managerErr = newFQDN(opts)
	})

	if managerErr != nil {
		return managerErr
	}

	return err
}

// GetFQDN extracts the FQDN from a URL using the global manager
func GetFQDN(url string) (string, error) {
	// Initialize with default options if not already initialized
	if manager == nil {
		err := Init(DefaultOptions())
		if err != nil {
			return "", err
		}
	}

	return manager.GetFQDN(url)
}

// ValidateOrigin checks if a given origin is in the allowed origins list
func ValidateOrigin(origin string, allowedOrigins []string) bool {
	u, err := GetFQDN(origin)
	if err != nil {
		return false
	}

	// Check if the FQDN is in the allowed origins list
	for _, allowed := range allowedOrigins {
		if u == allowed {
			return true
		}
	}

	return false
}
