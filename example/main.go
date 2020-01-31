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
		u, _ := gotld.FQDNMgr.GetFQDN(url)
		fmt.Printf("%47s = fqdn[%s]\n", url, u)
	}
}