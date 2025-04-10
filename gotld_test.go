// file: gotld_test.go
// description: tests for the main package functionality

package gotld

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGetFQDN(t *testing.T) {
	// Reset the global manager for testing
	manager = nil
	managerOnce = sync.Once{}

	err := Init(DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to initialize GoTLD: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
		errType  error
	}{
		{
			name:     "Simple domain",
			input:    "example.com",
			expected: "example.com",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "Domain with subdomain",
			input:    "www.example.com",
			expected: "example.com",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "Domain with multiple subdomains",
			input:    "blog.www.example.com",
			expected: "example.com",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "Domain with HTTPS",
			input:    "https://example.com",
			expected: "example.com",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "Domain with path",
			input:    "https://example.com/path/to/resource",
			expected: "example.com",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "Domain with query params",
			input:    "https://example.com?foo=bar",
			expected: "example.com",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "Domain with port",
			input:    "https://example.com:8080",
			expected: "example.com",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "UK domain",
			input:    "example.co.uk",
			expected: "example.co.uk",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:     "UK domain with subdomain",
			input:    "www.example.co.uk",
			expected: "example.co.uk",
			wantErr:  false,
			errType:  nil,
		},
		{
			name:    "Invalid domain - no TLD",
			input:   "invalid",
			wantErr: true,
			errType: ErrInvalidURL,
		},
		{
			name:    "Invalid domain - no dots",
			input:   "invalid",
			wantErr: true,
			errType: ErrInvalidURL,
		},
		{
			name:    "Invalid domain - empty",
			input:   "",
			wantErr: true,
			errType: ErrInvalidURL,
		},
		{
			name:    "Invalid domain - just a dot",
			input:   ".",
			wantErr: true,
			errType: ErrInvalidURL,
		},
		{
			name:    "Invalid domain - starts with dot",
			input:   ".com",
			wantErr: true,
			errType: ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFQDN(tt.input)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFQDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, check it's the right type
			if tt.wantErr && tt.errType != nil && err != nil {
				if !strings.Contains(err.Error(), tt.errType.Error()) {
					t.Errorf("GetFQDN() expected error containing %v, got %v", tt.errType, err)
				}
				return
			}

			// Check the result
			if !tt.wantErr && got != tt.expected {
				t.Errorf("GetFQDN() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInitWithCustomOptions(t *testing.T) {
	// Reset the global manager for testing
	manager = nil
	managerOnce = sync.Once{}

	// Set up a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if ua := r.Header.Get("User-Agent"); !strings.Contains(ua, "GoTLD") {
			t.Errorf("Expected User-Agent header to contain 'GoTLD', got %q", ua)
		}

		// Send a valid response
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// Write the public suffix list header and some test data
		w.Write([]byte("// The Public Suffix List\n// https://publicsuffix.org/list/public_suffix_list.dat\n\n// ===BEGIN ICANN DOMAINS===\ncom\nco.uk\n// ===END ICANN DOMAINS===\n"))
	}))
	defer ts.Close()

	// Create custom options
	opts := &Options{
		AllowPrivateTLDs: true,
		Timeout:          5 * time.Second,
		PublicSuffixURL:  ts.URL,
		Context:          context.Background(),
	}

	// Initialize with custom options
	err := Init(opts)
	if err != nil {
		t.Fatalf("Failed to initialize with custom options: %v", err)
	}

	// Test with a domain that should be recognized
	domain, err := GetFQDN("example.com")
	if err != nil {
		t.Errorf("GetFQDN() error = %v", err)
		return
	}

	if domain != "example.com" {
		t.Errorf("GetFQDN() = %v, want %v", domain, "example.com")
	}
}

func TestPublicSuffixDownloadFailure(t *testing.T) {
	// Reset the global manager for testing
	manager = nil
	managerOnce = sync.Once{}

	// Set up a test server that returns an error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Create options with invalid URL
	opts := &Options{
		PublicSuffixURL: ts.URL,
		Timeout:         2 * time.Second,
	}

	// Initialize should fail
	err := Init(opts)
	if err == nil {
		t.Error("Expected Init() to fail with invalid URL, but it succeeded")
	}
}

func TestContextCancellation(t *testing.T) {
	// Reset the global manager for testing
	manager = nil
	managerOnce = sync.Once{}

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Set up a test server that hangs
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wait for a while to simulate a slow server
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create options with the context
	opts := &Options{
		PublicSuffixURL: ts.URL,
		Context:         ctx,
		Timeout:         5 * time.Second,
	}

	// Cancel the context before initialization completes
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// Initialize should fail with context cancellation
	err := Init(opts)
	if err == nil {
		t.Error("Expected Init() to fail with cancelled context, but it succeeded")
	}
}

func TestValidateOrigin(t *testing.T) {
	// Reset the global manager for testing
	manager = nil
	managerOnce = sync.Once{}

	err := Init(DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to initialize GoTLD: %v", err)
	}

	allowedOrigins := []string{
		"example.com",
		"trusted.org",
		"api.service.com",
	}

	tests := []struct {
		name    string
		origin  string
		allowed bool
	}{
		{
			name:    "Simple allowed origin",
			origin:  "example.com",
			allowed: true,
		},
		{
			name:    "Allowed origin with scheme",
			origin:  "https://example.com",
			allowed: true,
		},
		{
			name:    "Allowed origin with subdomain",
			origin:  "www.example.com",
			allowed: true,
		},
		{
			name:    "Allowed origin with path",
			origin:  "example.com/path",
			allowed: true,
		},
		{
			name:    "Allowed origin with port",
			origin:  "example.com:8080",
			allowed: true,
		},
		{
			name:    "Disallowed origin",
			origin:  "malicious.com",
			allowed: false,
		},
		{
			name:    "Invalid origin",
			origin:  "invalid",
			allowed: false,
		},
		{
			name:    "Multi-part allowed origin",
			origin:  "api.service.com",
			allowed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateOrigin(tt.origin, allowedOrigins)
			if got != tt.allowed {
				t.Errorf("ValidateOrigin() = %v, want %v", got, tt.allowed)
			}
		})
	}
}

// TestETLD tests the ETLD type functionality
func TestETLD(t *testing.T) {
	etld := emptyETLD(1)

	// Test adding
	added := etld.Add("example.com", true)
	if !added {
		t.Error("Failed to add item to ETLD")
	}

	// Test adding duplicate
	added = etld.Add("example.com", true)
	if added {
		t.Error("Added duplicate item to ETLD")
	}

	// Test adding without sorting
	added = etld.Add("aardvark.com", false)
	if !added {
		t.Error("Failed to add second item to ETLD")
	}

	// Test sorting
	etld.Sort()
	if etld.List[0] != "aardvark.com" {
		t.Errorf("Sort failed, expected 'aardvark.com' as first item, got %s", etld.List[0])
	}

	// Test searching for existing item
	found, exists := etld.Search("example.com")
	if !exists {
		t.Error("Failed to find existing item in ETLD")
	}
	if found != "example.com" {
		t.Errorf("Search returned wrong item, expected 'example.com', got %s", found)
	}

	// Test searching for non-existent item
	found, exists = etld.Search("nonexistent.com")
	if exists {
		t.Error("Found non-existent item in ETLD")
	}
	if found != "" {
		t.Errorf("Search for non-existent item returned %s, expected empty string", found)
	}
}

// TestFQDNGuess tests the guess functionality
func TestFQDNGuess(t *testing.T) {
	// Create a new FQDN manager
	fqdn, err := newFQDN(DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to create FQDN manager: %v", err)
	}

	tests := []struct {
		name     string
		domain   string
		count    int
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple domain, level 1",
			domain:   "example.com",
			count:    1,
			expected: "com",
			wantErr:  false,
		},
		{
			name:     "Two-level domain, level 2",
			domain:   "example.co.uk",
			count:    2,
			expected: "co.uk",
			wantErr:  false,
		},
		{
			name:     "Multi-part domain, level 3",
			domain:   "sub.example.co.uk",
			count:    3,
			expected: "example.co.uk",
			wantErr:  false,
		},
		{
			name:     "Deep domain, level 4",
			domain:   "deep.sub.example.co.uk",
			count:    4,
			expected: "sub.example.co.uk",
			wantErr:  false,
		},
		{
			name:     "Very deep domain, level 5",
			domain:   "very.deep.sub.example.co.uk",
			count:    5,
			expected: "deep.sub.example.co.uk",
			wantErr:  false,
		},
		{
			name:    "Invalid domain - no dots",
			domain:  "invalid",
			count:   1,
			wantErr: true,
		},
		{
			name:    "Invalid domain - empty",
			domain:  "",
			count:   1,
			wantErr: true,
		},
		{
			name:    "Count too high for domain",
			domain:  "example.com",
			count:   3,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fqdn.guess(tt.domain, tt.count)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("guess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check the result
			if !tt.wantErr && got != tt.expected {
				t.Errorf("guess() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestHasScheme tests the hasScheme functionality
func TestHasScheme(t *testing.T) {
	// Create a new FQDN manager
	fqdn, err := newFQDN(DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to create FQDN manager: %v", err)
	}

	tests := []struct {
		name      string
		url       string
		remove    bool
		expected  string
		hasScheme bool
	}{
		{
			name:      "HTTP scheme",
			url:       "http://example.com",
			remove:    true,
			expected:  "example.com",
			hasScheme: true,
		},
		{
			name:      "HTTPS scheme",
			url:       "https://example.com",
			remove:    true,
			expected:  "example.com",
			hasScheme: true,
		},
		{
			name:      "FTP scheme",
			url:       "ftp://example.com",
			remove:    true,
			expected:  "example.com",
			hasScheme: true,
		},
		{
			name:      "Fake scheme",
			url:       "fake://example.com",
			remove:    true,
			expected:  "example.com",
			hasScheme: true,
		},
		{
			name:      "No scheme",
			url:       "example.com",
			remove:    true,
			expected:  "example.com",
			hasScheme: false,
		},
		{
			name:      "HTTP scheme, don't remove",
			url:       "http://example.com",
			remove:    false,
			expected:  "http://example.com",
			hasScheme: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, hasScheme := fqdn.hasScheme(tt.url, tt.remove)

			if hasScheme != tt.hasScheme {
				t.Errorf("hasScheme() hasScheme = %v, want %v", hasScheme, tt.hasScheme)
			}

			if got != tt.expected {
				t.Errorf("hasScheme() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestParsePublicSuffixData tests the parsePublicSuffixData functionality
func TestParsePublicSuffixData(t *testing.T) {
	// Create a new FQDN manager
	fqdn, err := newFQDN(DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to create FQDN manager: %v", err)
	}

	// Test with valid data
	validData := []byte(`// The Public Suffix List
// https://publicsuffix.org/list/public_suffix_list.dat

// ===BEGIN ICANN DOMAINS===
com
co.uk
org
// ===END ICANN DOMAINS===

// ===BEGIN PRIVATE DOMAINS===
github.io
amazonaws.com
// ===END PRIVATE DOMAINS===
`)

	// Parse the valid data
	err = fqdn.parsePublicSuffixData(validData)
	if err != nil {
		t.Errorf("parsePublicSuffixData() error = %v", err)
	}

	// Test with invalid data
	invalidData := []byte(`This is not the public suffix list`)
	err = fqdn.parsePublicSuffixData(invalidData)
	if err == nil {
		t.Error("parsePublicSuffixData() expected error with invalid data, got nil")
	}

	// Test with empty data
	emptyData := []byte(``)
	err = fqdn.parsePublicSuffixData(emptyData)
	if err == nil {
		t.Error("parsePublicSuffixData() expected error with empty data, got nil")
	}
}

// TestConcurrentAccess tests concurrent access to ETLD and FQDN
func TestConcurrentAccess(t *testing.T) {
	// Create a new ETLD
	etld := emptyETLD(1)

	// Add some initial data
	etld.Add("example.com", true)
	etld.Add("example.org", true)

	// Test concurrent access to ETLD
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Some goroutines add
			if i%3 == 0 {
				etld.Add(fmt.Sprintf("domain%d.com", i), false)
			}

			// Some goroutines search
			if i%3 == 1 {
				_, _ = etld.Search("example.com")
			}

			// Some goroutines sort
			if i%3 == 2 {
				etld.Sort()
			}
		}(i)
	}

	wg.Wait()

	// Reset the global manager for testing
	manager = nil
	managerOnce = sync.Once{}

	err := Init(DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to initialize GoTLD: %v", err)
	}

	// Test concurrent access to GetFQDN
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = GetFQDN("example.com")
		}()
	}

	wg.Wait()
}

// Benchmark GetFQDN
func BenchmarkGetFQDN(b *testing.B) {
	// Initialize once
	_ = Init(DefaultOptions())

	domains := []string{
		"example.com",
		"www.example.com",
		"blog.example.com",
		"example.co.uk",
		"www.example.co.uk",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain := domains[i%len(domains)]
		_, _ = GetFQDN(domain)
	}
}
