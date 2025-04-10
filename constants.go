// file: constants.go
// description: defines constants for the package

package gotld

const (
	// eTLDGroupMax is the maximum number of groups in a domain
	eTLDGroupMax = 5

	// publicSuffixFileURL is the URL to download the public suffix list from
	publicSuffixFileURL = "https://publicsuffix.org/list/public_suffix_list.dat"

	// minDataSize is the minimum size of the public suffix list file in bytes
	minDataSize = 32768
)
