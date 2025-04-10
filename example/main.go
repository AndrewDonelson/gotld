// file: example/main.go
// description: example program demonstrating usage of the gotld package

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AndrewDonelson/gotld"
)

func main() {
	// Define and parse command-line flags
	allowPrivate := flag.Bool("private", false, "Allow private TLDs")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for HTTP requests")
	customURL := flag.String("url", "", "Custom URL for public suffix list")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	// Configure logging
	if *verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(os.Stderr)
	}

	// Create custom options
	opts := &gotld.Options{
		AllowPrivateTLDs: *allowPrivate,
		Timeout:          *timeout,
		Context:          context.Background(),
	}

	// Set custom URL if provided
	if *customURL != "" {
		opts.PublicSuffixURL = *customURL
	}

	// Initialize gotld with custom options
	if err := gotld.Init(opts); err != nil {
		log.Fatalf("Failed to initialize gotld: %v", err)
	}

	// Get URLs from command-line arguments or use defaults
	var urls []string
	if flag.NArg() > 0 {
		urls = flag.Args()
	} else {
		// Default examples
		urls = []string{
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
	}

	// Process each URL
	fmt.Println("URL Analysis Results")
	fmt.Println("-------------------")
	fmt.Printf("%-50s | %-30s | %s\n", "Original URL", "FQDN", "Status")
	fmt.Println(strings.Repeat("-", 100))

	for _, url := range urls {
		fqdn, err := gotld.GetFQDN(url)
		
		if err != nil {
			fmt.Printf("%-50s | %-30s | ERROR: %v\n", url, "-", err)
		} else {
			fmt.Printf("%-50s | %-30s | SUCCESS\n", url, fqdn)
		}
	}

	// Demonstrate origin validation
	fmt.Println("\nOrigin Validation")
	fmt.Println("----------------")
	
	allowedOrigins := []string{
		"example.com",
		"trusted.org",
		"api.service.com",
	}

	fmt.Printf("Allowed origins: %s\n\n", strings.Join(allowedOrigins, ", "))
	
	originsToCheck := []string{
		"https://example.com",
		"http://malicious.com",
		"https://trusted.org/path",
		"https://subdomain.example.com",
	}

	for _, origin := range originsToCheck {
		isValid := gotld.ValidateOrigin(origin, allowedOrigins)
		validText := "INVALID"
		if isValid {
			validText = "VALID"
		}
		
		fmt.Printf("%-40s | %s\n", origin, validText)
	}
}