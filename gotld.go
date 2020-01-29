package tld

const (
	eTLDGroupMax		= 5
	publicSufficFileURL = "https://publicsuffix.org/list/public_suffix_list.dat"
)

// TLDS contains the latest public suffixes via download and parsing. This is called handled automatically
// in the package init()
var (
	FQDNMgr	*FQDN
	TLDS 	[]string
	Count 	int
)

func init() {
	FQDNMgr = newFQDN()
}

func domainPort(host string) (string, string) {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i], host[i+1:]
		} else if host[i] < '0' || host[i] > '9' {
			return host, ""
		}
	}
	//will only land here if the string is all digits,
	//net/url should prevent that from happening
	return host, ""
}
