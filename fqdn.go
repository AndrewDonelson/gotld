// file: fqdn.go
// description: manages fully qualified domain names

package gotld

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// FQDN main object structure with concurrency support
type FQDN struct {
	Options  *Options
	etldList [eTLDGroupMax]*ETLD
	total    int
	mu       sync.RWMutex
}

// newFQDN creates a new FQDN manager with the specified options
func newFQDN(opts *Options) (*FQDN, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	fqdn := &FQDN{
		Options: opts,
		mu:      sync.RWMutex{},
	}

	for i := 0; i < eTLDGroupMax; i++ {
		fqdn.etldList[i] = emptyETLD(i)
	}

	// Get the public suffix list
	var err error
	if opts.PublicSuffixFile != "" {
		err = fqdn.loadPublicSuffixFromFile(opts.PublicSuffixFile)
	} else {
		err = fqdn.downloadPublicSuffixFile(opts.PublicSuffixURL)
	}

	if err != nil {
		return nil, wrapError(err, "failed to initialize FQDN manager")
	}

	return fqdn, nil
}

// Tidy will tally the total number of loaded eTLDs and sort each list
func (f *FQDN) Tidy() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.total = 0
	var wg errgroup.Group

	for i := 0; i < eTLDGroupMax; i++ {
		i := i // Capture for goroutine
		wg.Go(func() error {
			f.etldList[i].Sort()
			return nil
		})
	}

	_ = wg.Wait() // Ignore error as Sort doesn't return an error

	for i := 0; i < eTLDGroupMax; i++ {
		f.total += f.etldList[i].Count
	}
}

// hasScheme checks if a URL has a scheme and optionally removes it
func (f *FQDN) hasScheme(s string, remove bool) (string, bool) {
	schemes := []string{"http://", "https://", "ftp://", "ws://", "wss://"}

	for _, scheme := range schemes {
		if strings.HasPrefix(s, scheme) {
			if remove {
				return strings.Replace(s, scheme, "", 1), true
			}
			return s, true
		}
	}

	// Handle the fake:// scheme separately for backward compatibility
	if strings.HasPrefix(s, "fake://") {
		if remove {
			return strings.Replace(s, "fake://", "", 1), true
		}
		return s, true
	}

	return s, false
}

// guess attempts to extract a potential eTLD from a URL
func (f *FQDN) guess(domain string, count int) (string, error) {
	if domain == "" {
		return "", ErrInvalidURL
	}

	dots := strings.Count(domain, ".")
	if dots < 1 || len(domain) < 3 {
		return "", ErrInvalidURL
	}

	groups := strings.Split(domain, ".")
	grpCnt := len(groups)

	if grpCnt >= count {
		switch {
		case count == eTLDGroupMax:
			return strings.Join(groups[grpCnt-5:], "."), nil
		case count == (eTLDGroupMax - 1):
			return strings.Join(groups[grpCnt-4:], "."), nil
		case count == (eTLDGroupMax - 2):
			return strings.Join(groups[grpCnt-3:], "."), nil
		case count == (eTLDGroupMax - 3):
			return strings.Join(groups[grpCnt-2:], "."), nil
		case count == (eTLDGroupMax - 4):
			return groups[grpCnt-1], nil
		}
	}

	return "", wrapError(ErrInvalidURL, "unable to make a guess")
}

// findTLD attempts to find the TLD of a domain
func (f *FQDN) findTLD(s string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var (
		tld, guess string
		found      bool
		err        error
	)

	dots := strings.Count(s, ".")
	if dots >= 1 {
		for i := dots; i > 0; i-- {
			guess, err = f.guess(s, i)
			if err == nil {
				tld, found = f.etldList[i-1].Search(guess)
				if found {
					break
				}
			}
		}
	}

	return tld
}

// GetFQDN extracts the FQDN from a URL
func (f *FQDN) GetFQDN(srcURL string) (string, error) {
	if srcURL == "" {
		return "", ErrInvalidURL
	}

	// Shortest domain ex. a.io (4), and must have at least 1 DOT
	if len(srcURL) < 4 || strings.Count(srcURL, ".") < 1 {
		return "", ErrInvalidURL
	}

	// If no prefix, add a fake one for net/url.Parse() (workaround)
	hadScheme := false
	srcURL, hadScheme = f.hasScheme(srcURL, false)
	if !hadScheme {
		srcURL = "fake://" + srcURL
	}

	parsedURL, err := url.Parse(srcURL)
	if err != nil {
		return "", wrapError(ErrInvalidURL, err.Error())
	}

	// We don't need scheme anymore - get rid of it
	srcURL, _ = f.hasScheme(srcURL, true)

	// Remove port if present
	if parsedURL.Port() != "" {
		srcURL = strings.Replace(srcURL, ":"+parsedURL.Port(), "", 1)
	}

	// Remove query parameters
	if parsedURL.RawQuery != "" {
		srcURL = strings.Replace(srcURL, "?"+parsedURL.RawQuery, "", 1)
	}

	// Remove path
	if parsedURL.Path != "" && parsedURL.Path != "/" {
		srcURL = strings.Replace(srcURL, parsedURL.Path, "", 1)
	}

	// Find the TLD
	eTLD := f.findTLD(srcURL)
	if eTLD == "" {
		return "", ErrInvalidTLD
	}

	// Extract the domain from the URL
	domainPart := strings.Replace(srcURL, "."+eTLD, "", 1)

	if domainPart == "" {
		return "", ErrInvalidURL
	}

	// Handle subdomains
	dots := strings.Count(domainPart, ".")
	if dots == 0 {
		return domainPart + "." + eTLD, nil
	}

	parts := strings.Split(domainPart, ".")
	return parts[len(parts)-1] + "." + eTLD, nil
}

// loadPublicSuffixFromFile loads the public suffix list from a local file
func (f *FQDN) loadPublicSuffixFromFile(filePath string) error {
	if filePath == "" {
		return wrapError(ErrPublicSuffixDownload, "no file path provided")
	}

	ctx := f.Options.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// Load file implementation would go here
	// For now, fall back to downloading
	return f.downloadPublicSuffixFile(f.Options.PublicSuffixURL)
}

// downloadPublicSuffixFile downloads and parses the public suffix list
func (f *FQDN) downloadPublicSuffixFile(fileURL string) error {
	if fileURL == "" {
		fileURL = publicSuffixFileURL
	}

	ctx := f.Options.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// Create HTTP client with proper security settings
	var client *http.Client
	if f.Options.CustomHTTPClient != nil {
		client = f.Options.CustomHTTPClient
	} else {
		timeout := f.Options.Timeout
		if timeout == 0 {
			timeout = 10 * time.Second
		}

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		client = &http.Client{
			Transport: transport,
			Timeout:   timeout,
		}
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return wrapError(ErrPublicSuffixDownload, err.Error())
	}

	// Set appropriate headers
	req.Header.Set("User-Agent", "GoTLD/1.0")

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		return wrapError(ErrPublicSuffixDownload, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return wrapError(ErrPublicSuffixDownload, "unexpected status code: "+resp.Status)
	}

	// Read the response body
	respData, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024)) // Limit to 10MB
	if err != nil {
		return wrapError(ErrPublicSuffixParse, err.Error())
	}

	if len(respData) < minDataSize {
		return wrapError(ErrPublicSuffixParse, "response data size too small for public suffix file")
	}

	// Parse the response
	return f.parsePublicSuffixData(respData)
}

// parsePublicSuffixData parses the public suffix list data
func (f *FQDN) parsePublicSuffixData(data []byte) error {
	sliceData := strings.Split(string(data), "\n")

	if len(sliceData) == 0 {
		return ErrPublicSuffixParse
	}

	// Verify that this is the public suffix list
	found := false
	for i := 0; i < 10 && i < len(sliceData); i++ {
		if strings.Contains(sliceData[i], publicSuffixFileURL) {
			found = true
			break
		}
	}

	if !found {
		return ErrPublicSuffixFormat
	}

	var icann bool

	f.mu.Lock()
	defer f.mu.Unlock()

	// Reset the current lists
	for i := 0; i < eTLDGroupMax; i++ {
		f.etldList[i] = emptyETLD(i)
	}

	for _, tld := range sliceData {
		// Skip blank lines
		if len(tld) == 0 {
			continue
		}

		// Detect and toggle ICANN eTLD state
		if strings.Contains(tld, "===BEGIN ICANN DOMAINS===") {
			icann = true
			continue
		} else if strings.Contains(tld, "===END ICANN DOMAINS===") {
			icann = false
			continue
		}

		// If private TLDs not allowed and this is not an ICANN TLD, skip it
		if !f.Options.AllowPrivateTLDs && !icann {
			continue
		}

		// Skip comments, wildcards, and exceptions
		if strings.HasPrefix(tld, "//") || strings.HasPrefix(tld, "*") || strings.HasPrefix(tld, "!") {
			continue
		}

		// Add eTLD to list
		tld = strings.ToLower(strings.TrimSpace(tld))
		if tld == "" {
			continue
		}

		dots := strings.Count(tld, ".")
		if dots < eTLDGroupMax {
			f.etldList[dots].Add(tld, false)
		}
	}

	f.Tidy()
	return nil
}
