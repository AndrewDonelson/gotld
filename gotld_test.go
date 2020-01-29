package tld

import "testing"

func Test1(t *testing.T) {
	var src string
	
	src = ".com"
	u, err := FQDNMgr.GetFQDN(src)
	if err.Error() != "Not a valid URL" || u != "" {
		t.Fail()
	}

	src = "off"
	u, err = FQDNMgr.GetFQDN(src)
	if err.Error() != "Not a valid URL" || u != "" {
		t.Fail()
	}

	src = "."
	u, err = FQDNMgr.GetFQDN(src)
	if err.Error() != "Not a valid URL" || u != "" {
		t.Fail()
	}

	src = "nlaak.com"
	u, err = FQDNMgr.GetFQDN(src)
	if err!=nil || u != src {
		t.Fail()
	}

	src = "https://eu.com"
	u, err = FQDNMgr.GetFQDN(src)
	if err!=nil || u != "eu.com" {
		t.Fail()
	}

	src = "https://staging.www.eu.com"
	u, err = FQDNMgr.GetFQDN(src)
	if err!=nil {
		if FQDNMgr.Options.AllowPrivateTLDs && u != "www.eu.com" {
			t.Fail()		
		} else if u != "eu.com"{
			t.Fail()		
		}
	}
	
	src = "https://sub.domain.dcba:8080/?d=e"
	u, err = FQDNMgr.GetFQDN(src)
	if err.Error() != "Not a valid eTLD" || u != "" {
		t.Fail()
	}

}