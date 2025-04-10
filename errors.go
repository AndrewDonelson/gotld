// file: errors.go
// description: defines error types for the package

package gotld

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidURL is returned when a URL is invalid
	ErrInvalidURL = errors.New("invalid URL")

	// ErrInvalidTLD is returned when a TLD is not found in the public suffix list
	ErrInvalidTLD = errors.New("invalid TLD")

	// ErrPublicSuffixDownload is returned when the public suffix file cannot be downloaded
	ErrPublicSuffixDownload = errors.New("failed to download public suffix file")

	// ErrPublicSuffixParse is returned when the public suffix file cannot be parsed
	ErrPublicSuffixParse = errors.New("failed to parse public suffix file")

	// ErrPublicSuffixFormat is returned when the downloaded file is not the public suffix file
	ErrPublicSuffixFormat = errors.New("file is not the public suffix file")
)

// wrapError wraps an error with additional context
func wrapError(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
