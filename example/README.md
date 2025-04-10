# GoTLD Example Application

This example demonstrates how to use the GoTLD package to extract the fully qualified domain name (FQDN) from URLs and validate origins against a list of allowed origins.

## Usage

```
go run main.go [options] [urls...]
```

### Options

- `-private`: Allow private TLDs (default: false)
- `-timeout`: Timeout for HTTP requests (default: 10s)
- `-url`: Custom URL for public suffix list (default: https://publicsuffix.org/list/public_suffix_list.dat)
- `-verbose`: Enable verbose logging (default: false)

### Examples

Process default URLs:

```
go run main.go
```

Process custom URLs:

```
go run main.go example.com subdomain.example.co.uk
```

Allow private TLDs:

```
go run main.go -private https://github.io
```

Set custom timeout:

```
go run main.go -timeout 5s
```

## Output

The program outputs a table showing the original URL, the extracted FQDN, and the status of the operation. It also demonstrates origin validation against a list of allowed origins.

## Implementation Details

The example uses the GoTLD package to:

1. Initialize the package with custom options
2. Extract FQDNs from URLs
3. Validate origins against a list of allowed origins

This demonstrates how the package can be used for security-related tasks such as CORS validation or domain allow-listing.

// file: example/Makefile
// description: makefile for building and running the example

.PHONY: build run clean

build:
	go build -o gotld-example main.go

run: build
	./gotld-example

clean:
	rm -f gotld-example