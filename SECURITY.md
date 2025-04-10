# Security Considerations

The GoTLD package has been designed with security in mind. This document outlines the security considerations and best practices for using the package.

## Transport Security

When downloading the public suffix list, the package uses HTTPS with modern TLS settings:

- Minimum TLS version is set to TLS 1.2
- HTTP/2 is enabled when supported
- TLS handshake timeouts are properly configured
- Connection timeouts are set to reasonable values

## Input Validation

All input is carefully validated before processing:

- URLs are validated for proper formatting
- Domain names are checked against the public suffix list
- Invalid or malformed input is rejected with appropriate error messages

## Error Handling

Errors are properly handled and propagated with context:

- Custom error types are defined for specific error conditions
- Errors are wrapped with context for better debugging
- Error messages do not reveal sensitive information

## Concurrency Safety

The package is designed to be safe for concurrent use:

- Mutexes protect shared data
- Global state is minimized
- The API is designed to be thread-safe

## Context Support

The package supports context-based cancellation:

- HTTP requests can be cancelled via context
- Long-running operations respect context cancellation
- Context deadline is honored for timeout handling

## Memory Management

Care has been taken to manage memory efficiently:

- Large responses are limited to prevent DOS attacks
- Buffers are sized appropriately
- Memory allocations are minimized

## Best Practices

When using the GoTLD package, consider the following best practices:

1. **Initialize Once**: Initialize the package once at application startup
2. **Set Timeouts**: Always set reasonable timeouts for HTTP requests
3. **Handle Errors**: Always check and handle errors returned by the package
4. **Validate Input**: Validate user input before passing it to the package
5. **Use Context**: Provide a context with deadlines for operations
6. **Update Regularly**: The public suffix list changes over time, so consider refreshing it periodically

## Known Limitations

- The package does not currently support IDN (Internationalized Domain Names) directly
- The package only validates domains against the public suffix list, not the content at those domains

// file: IMPROVEMENTS.md
// description: documentation of improvements made to the package

# GoTLD Improvements

This document outlines the improvements made to the GoTLD package to make it more secure and production-ready.

## Security Improvements

### 1. Secure HTTP Client

The HTTP client used to download the public suffix list has been improved with:

- Timeout settings to prevent hanging connections
- TLS configuration with minimum version set to TLS 1.2
- Connection pooling with sensible defaults
- Proper handling of request cancellation via context

### 2. Error Types

Custom error types have been defined to provide more context and make error handling more robust:

- `ErrInvalidURL`: returned when a URL is invalid
- `ErrInvalidTLD`: returned when a TLD is not found in the public suffix list
- `ErrPublicSuffixDownload`: returned when the public suffix list cannot be downloaded
- `ErrPublicSuffixParse`: returned when the public suffix list cannot be parsed
- `ErrPublicSuffixFormat`: returned when the downloaded file is not the public suffix list

### 3. Input Validation

Input validation has been improved throughout the codebase:

- URLs are thoroughly validated for proper formatting
- Edge cases like empty strings, URLs with no dots, etc. are handled properly
- Domain components are validated before processing

### 4. Concurrency Safety

The package has been made thread-safe:

- Mutexes protect shared data structures
- Concurrent operations on the same data are synchronized
- Tests verify concurrent safety

## Production Readiness Improvements

### 1. Modern Go Practices

The codebase has been updated to use modern Go practices:

- Go modules support
- Updated to Go 1.21
- Context support for cancellation
- Error wrapping for better context

### 2. Code Organization

The code has been reorganized for better maintainability:

- Files are grouped by functionality
- Constants are extracted to a separate file
- Error types are centralized
- File headers provide clear descriptions

### 3. Documentation

Documentation has been improved throughout the codebase:

- GoDoc compatible comments
- Detailed function descriptions
- Examples demonstrating usage
- Security considerations documented

### 4. Testing

Test coverage has been improved:

- Unit tests for each component
- Integration tests for the full package
- Benchmarks for performance-critical code
- Tests for error conditions and edge cases
- Concurrency tests

### 5. Configuration

Configuration has been made more flexible:

- Options struct for configuring the package
- Support for custom HTTP clients
- Ability to load from local file or URL
- Context support for cancellation

### 6. Performance

Performance has been improved:

- Concurrent sorting of TLD lists
- Better memory management
- Reduced allocations in hot paths
- Benchmarks to track performance

## Breaking Changes

The following breaking changes were made to improve the API:

1. `Init` function now takes an `Options` struct instead of individual parameters
2. The global manager is no longer initialized in `init()` but must be explicitly initialized
3. Error types have changed to be more specific and provide better context

## Migration Guide

To migrate from the old API to the new API:

```go
// Old API
import "github.com/AndrewDonelson/gotld"

func main() {
    u, _ := gotld.FQDNMgr.GetFQDN("example.com")
}

// New API
import "github.com/AndrewDonelson/gotld"

func main() {
    // Initialize with default options
    gotld.Init(gotld.DefaultOptions())
    
    // Get FQDN
    u, err := gotld.GetFQDN("example.com")
    if err != nil {
        // Handle error
    }
}
```

// file: README.md
// description: updated README with improved documentation

# GoTLD

[![Build Status](https://travis-ci.org/AndrewDonelson/gotld.svg?branch=master)](https://travis-ci.org/AndrewDonelson/gotld)
![GitHub last commit](https://img.shields.io/github/last-commit/AndrewDonelson/gotld)
[![Coverage Status](https://coveralls.io/repos/github/AndrewDonelson/gotld/badge.svg)](https://coveralls.io/github/AndrewDonelson/gotld)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/488f571baa13489494fa6002dbdf0897)](https://www.codacy.com/manual/AndrewDonelson/gotld?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=AndrewDonelson/gotld&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/AndrewDonelson/gotld?status.svg)](http://godoc.org/github.com/AndrewDonelson/gotld)
![GitHub stars](https://img.shields.io/github/stars/AndrewDonelson/gotld?style=flat)

GoTLD is a secure and production-ready Go package for extracting the fully qualified domain name (FQDN) from URLs using the Public Suffix List.

## Features

- Extract FQDNs from URLs with proper TLD handling
- Support for complex domains with multiple levels (e.g., .co.uk)
- Thread-safe for concurrent use
- Context support for cancellation and timeouts
- Configurable options for custom behavior
- Comprehensive error handling
- Robust security features
- Extensively tested

## Installation

```sh
go get -u github.com/AndrewDonelson/gotld
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/AndrewDonelson/gotld"
)

func main() {
	// Initialize with default options
	if err := gotld.Init(gotld.DefaultOptions()); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	urls := []string{
		"nlaak.com",
		"https://nlaak.com",
		"http://go.com?foo=bar",
		"http://google.com",
		"http://blog.google",
		"https://www.medi-cal.ca.gov/",
		"https://ato.gov.au",
		"http://stage.host.domain.co.uk/",
		"http://a.very.complex-domain.co.uk:8080/foo/bar",
	}

	for _, url := range urls {
		fqdn, err := gotld.GetFQDN(url)
		if err != nil {
			fmt.Printf("%-47s = ERROR: %v\n", url, err)
			continue
		}
		fmt.Printf("%-47s = fqdn[%s]\n", url, fqdn)
	}
}
```

## Advanced Usage

### Custom Options

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/AndrewDonelson/gotld"
)

func main() {
	// Create custom options
	opts := &gotld.Options{
		AllowPrivateTLDs: true,
		Timeout:          5 * time.Second,
		Context:          context.Background(),
	}

	// Initialize with custom options
	if err := gotld.Init(opts); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Use the package...
}
```

### Origin Validation

```go
package main

import (
	"fmt"

	"github.com/AndrewDonelson/gotld"
)

func main() {
	// Initialize with default options
	gotld.Init(gotld.DefaultOptions())

	// Define allowed origins
	allowedOrigins := []string{
		"example.com",
		"trusted.org",
	}

	// Validate an origin
	origin := "https://example.com/path"
	isValid := gotld.ValidateOrigin(origin, allowedOrigins)

	fmt.Printf("Origin %s is valid: %v\n", origin, isValid)
}
```

## Documentation

For more details, see the [GoDoc documentation](http://godoc.org/github.com/AndrewDonelson/gotld).

## Security

See [SECURITY.md](SECURITY.md) for security considerations.

## Improvements

See [IMPROVEMENTS.md](IMPROVEMENTS.md) for details on the improvements made to this package.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes (`git commit -am 'Add my feature'`)
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a Pull Request