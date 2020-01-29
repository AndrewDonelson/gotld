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

	println("Parse()")
	for _, url := range urls {
		u, _ := tld.Parse(url)
		fmt.Printf("%47s = sch[%s] sub[%s] dom[%s] tld[%s] prt[%s] pth[%s] qry[%s] fqdn[%s]\n",
			u, u.Scheme, u.Subdomain, u.Domain, u.TLD, u.Port, u.Path, u.RawQuery,u.FQDN)
	}

	println("\nGetFQDN()")
	for _, url := range urls {
		u, _ := tld.GetFQDN(url)
		fmt.Printf("%47s = %s\n", url, u)
	}	
}