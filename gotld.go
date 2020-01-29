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