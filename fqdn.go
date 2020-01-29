package tld

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// FQDN main object structure
type FQDN struct {
	Options		*Options
	etldList    [eTLDGroupMax]*ETLD
	// etldList[1] 	*ETLD	// eTLD Groups - ie. [com]
	// etldList[2] 	*ETLD	// eTLD Groups - ie. [eu.com]
	// etldList[3] 	*ETLD	// eTLD Groups - ie. [api.stdlib.com]
	// etldList[4] 	*ETLD	// eTLD Groups - ie. [usr.cloud.muni.cz]
	// etldList[5] 	*ETLD	// eTLD Groups - ie. [app.os.stg.fedoraproject.org]
	total 		int		// total combined amount of eTLDs
}

// newFQDN create a new default FQDN Manager
func newFQDN() (*FQDN, error) {
	fqdn := &FQDN{}
	fqdn.Options = &Options{}
	for i := 0; i < eTLDGroupMax; i++ {
		fqdn.etldList[i] = emptyETLD(i)
	}	

	// First step: Get the latest list - dont continue without it
	err := fqdn.downloadPublicSuffixFile(publicSufficFileURL)

	return fqdn, err
}

// Tidy will tally the total numbe rof loaded eTLDs and sort each list
func (f *FQDN) Tidy() {
	f.total = 0
	for i := 0; i < eTLDGroupMax; i++ {
		f.etldList[i].Sort()
		f.total += f.etldList[i].Count
	}	
}

func (f *FQDN) hasScheme(s string, remove bool) (string, bool) {
	var result bool

	if strings.HasPrefix(s,"http://") {
		result = true
		if remove {
			s = strings.Replace(s,"http://","",-1)
		}
	} else if strings.HasPrefix(s,"https://") {
		result = true
		if remove {
			s = strings.Replace(s,"https://","",-1)
		}
	} else if strings.HasPrefix(s,"fake://") {
		result = true
		if remove {
			s = strings.Replace(s,"fake://","",-1)
		}
	}

	return s, result
}

func (f *FQDN) guess(url string, count int) (string,error) {
	dots := strings.Count(url,".")
	if dots < 1 || len(url) < 3 {
		return "", fmt.Errorf("not a valid url")
	}
	groups := strings.Split(url,".")
	grpCnt := len(groups)

	if grpCnt >= count {
		if count == eTLDGroupMax {
			return fmt.Sprintf("%s.%s.%s.%s.%s",groups[grpCnt-5],groups[grpCnt-4],groups[grpCnt-3],groups[grpCnt-2],groups[grpCnt-1]), nil
		} else if count == (eTLDGroupMax - 1) {
			return fmt.Sprintf("%s.%s.%s.%s",groups[grpCnt-4],groups[grpCnt-3],groups[grpCnt-2],groups[grpCnt-1]), nil
		} else if count == (eTLDGroupMax - 2) {
			return fmt.Sprintf("%s.%s.%s",groups[grpCnt-3],groups[grpCnt-2],groups[grpCnt-1]), nil
		} else if count == (eTLDGroupMax - 3) {
			return fmt.Sprintf("%s.%s",groups[grpCnt-2],groups[grpCnt-1]), nil
		} else if count == (eTLDGroupMax - 4) {
			return groups[grpCnt-1], nil
		}		
	}

	return "", fmt.Errorf("unable to make a guess")
}

func  (f *FQDN) findTLD(s string) string {
	var (
		tld, guess string
		found bool
		err error
	)

	dots := strings.Count(s,".")
	if dots >= 1 {
		for i := dots; i > 0; i-- {
			guess,err = f.guess(s,i)
			if err == nil {
				tld, found = f.etldList[i-1].Search(guess)
				if found {
					break
				}
			}
		}		
	}

	return tld
}

// GetFQDN given a valid url will attempt to detect & return the eTLD
// using the loaded public suffix list
func (f *FQDN) GetFQDN(srcurl string) (str string, err error) {

	// shortest domain ex. a.io (4), and must have at least 1 DOT
	if len(srcurl) < 4 || strings.Count(srcurl, ".") < 1 {
		return "", fmt.Errorf("Not a valid URL")
	}

	//if no prefix, add a fake one )net/url.Parser() issue (work around)
	srcurl, yes := f.hasScheme(srcurl,false)
	if !yes {
		srcurl = "fake://" + srcurl
	}

	url, _ := url.Parse(srcurl)

	// We dont need scheme anymore - get rid of it
	srcurl, _ = f.hasScheme(srcurl,true)

	if url.Port() != "" {
		srcurl = strings.Replace(srcurl, ":"+url.Port(), "", -1)
	}

	if url.RawQuery != "" {
		srcurl = strings.Replace(srcurl, "?" + url.RawQuery, "", -1)
	}

	if url.Path != "" {
		srcurl = strings.Replace(srcurl, url.Path, "", -1)
	}

	eTLD := f.findTLD(srcurl)
	if eTLD == "" {
		return "", fmt.Errorf("Not a valid eTLD")	
	}

	srcurl = strings.Replace(srcurl, "."+eTLD, "", -1)

	if srcurl == "" {
		return "", fmt.Errorf("Not a valid URL")
	}

	dots := strings.Count(srcurl,".")
	if dots == 0 {
		return fmt.Sprintf("%s.%s",srcurl,eTLD), nil	
	}

	sub := strings.Split(srcurl,".")
	return  fmt.Sprintf("%s.%s",sub[len(sub)-1],eTLD), nil
}

// DownloadPublicSuffixFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func (f *FQDN) downloadPublicSuffixFile(file string) error {
	var icann bool

	if len(file) <= 0 {
		file = publicSufficFileURL
	}

	// Get the data
	resp, err := http.Get(file)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Public Suffix file was not downloaded")
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(respData) < 32768 {
		return fmt.Errorf("Response data size to small for Public Suffix file")
	}
	sliceData := strings.Split(string(respData), "\n")

	if len(sliceData) > 0 {
		if !strings.Contains(sliceData[4],publicSufficFileURL) {
			return fmt.Errorf("File is not the Public Suffix Data File")
		}

		for _, tld := range sliceData {
			// Skip blank lines and comments
			if len(tld) > 0 {
		
				// detect and toggle icann eTLD state
				if strings.Contains(tld,"===BEGIN ICANN DOMAINS===") {
					icann = true
				} else if strings.Contains(tld,"===END ICANN DOMAINS===") {
					icann = false
				}
	
				// if private tlds not allowed and this is not icann tld
				// skip it
				if !f.Options.AllowPrivateTLDs && !icann {
					continue
				}
				
				// If this is not a comment - continue processing
				if !strings.HasPrefix(tld, "//") {
					if !strings.HasPrefix(tld, "*") {
						if !strings.HasPrefix(tld, "!") {
							// Add eTLD to list
							dots := strings.Count(tld,".")
							tld = strings.ToLower(strings.TrimSpace(tld))
							f.etldList[dots].Add(tld,false)			
						}	
					}
				}
			}
		}
	
		f.Tidy()	
	}

	return nil
}
