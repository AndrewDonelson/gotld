// file: options.go
// description: defines options for the FQDN manager

package gotld

import (
	"context"
	"net/http"
	"time"
)

// Options for the FQDN Manager
type Options struct {
	// AllowPrivateTLDs determines whether private TLDs are allowed
	AllowPrivateTLDs bool

	// Timeout for HTTP requests
	Timeout time.Duration

	// CustomHTTPClient allows setting a custom HTTP client
	CustomHTTPClient *http.Client

	// PublicSuffixURL is the URL to download the public suffix list from
	PublicSuffixURL string

	// PublicSuffixFile is a local file containing the public suffix list
	PublicSuffixFile string

	// Context is used for cancellation
	Context context.Context
}

// DefaultOptions returns default options
func DefaultOptions() *Options {
	return &Options{
		AllowPrivateTLDs: false,
		Timeout:          10 * time.Second,
		CustomHTTPClient: nil,
		PublicSuffixURL:  publicSuffixFileURL,
		PublicSuffixFile: "",
		Context:          context.Background(),
	}
}
