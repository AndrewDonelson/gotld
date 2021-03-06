# gotld

[![Build Status](https://travis-ci.org/AndrewDonelson/gotld.svg?branch=master)](https://travis-ci.org/AndrewDonelson/gotld)
![GitHub last commit](https://img.shields.io/github/last-commit/AndrewDonelson/gotld)
[![Coverage Status](https://coveralls.io/repos/github/AndrewDonelson/gotld/badge.svg)](https://coveralls.io/github/AndrewDonelson/gotld)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/488f571baa13489494fa6002dbdf0897)](https://www.codacy.com/manual/AndrewDonelson/gotld?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=AndrewDonelson/gotld&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/AndrewDonelson/gotld?status.svg)](http://godoc.org/github.com/AndrewDonelson/gotld)
![GitHub stars](https://img.shields.io/github/stars/AndrewDonelson/gotld?style=flat)

The `tld` package has the same API ([see godoc](http://godoc.org/github.com/AndrewDonelson/gotld)) as `net/url` except `tld.URL` contains extra fields: `Subdomain`, `Domain`, `TLD` and `Port`.

_Note:_ This was been written using the Google [Public Suffix](http://golang.org/x/net/publicsuffix) package

## Install

```sh
go get github.com/AndrewDonelson/gotld
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/AndrewDonelson/gotld"
)

func main() {
	urls := []string{
		"nlaak.com",			//	net/url bug - returns as path. Workaround add scheme
		"https://nlaak.com",	//	net/url this works :(
		"http://go.com?foo=bar",
		"http://google.com",
		"http://blog.google",
		"https://www.medi-cal.ca.gov/",
		"https://ato.gov.au",
		"http://stage.host.domain.co.uk/",
		"http://a.very.complex-domain.co.uk:8080/foo/bar",
	}

	println("Example #1")
	for _, url := range urls {
		u, _ := tld.FQDNMgr.GetFQDN(url)
		fmt.Printf("%47s = fqdn[%s]\n", url, u)
	}
}
```

```sh
$ go run main.go
Example #1
                                      nlaak.com = fqdn[nlaak.com]
                              https://nlaak.com = fqdn[nlaak.com]
                          http://go.com?foo=bar = fqdn[go.com]
                              http://google.com = fqdn[google.com]
                             http://blog.google = fqdn[blog.google]
                   https://www.medi-cal.ca.gov/ = fqdn[ca.gov]
                             https://ato.gov.au = fqdn[ato.gov.au]
                http://stage.host.domain.co.uk/ = fqdn[domain.co.uk]
http://a.very.complex-domain.co.uk:8080/foo/bar = fqdn[complex-domain.co.uk]
```

### MIT License

Copyright © 2020 Andrew Donelson &lt;me@andrewdonelson.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
