package gotld

import "testing"

//strArray := [5]string{"India", "Canada", "Japan", "Germany", "Italy"}
var CXHOSTS = []string{
	"abril.ntwk50.com",
	"acmeair.incentivenetworks.com",
	"acmegrocer.incentivenetworks.com",
	"acuonline.bankmyrewards.com",
	"acuonline.eversave.com",
	"affinion.ntwk50.com",
	"ba.shopping.base.com",
	"base.treatstreet.stage.com",
	"base1.incentivenetworks.com",
	"base2.incentivenetworks.com",
	"base3.incentivenetworks.com",
	"base4.incentivenetworks.com",
	"base5.incentivenetworks.com",
	"bauk.incentivenetworks.com",
	"baus.incentivenetworks.com",
	"bestbank.bankmyrewards.com",
	"bestbank.eversave.com",
	"blackbird.eversave.com",
	"bonuspointsmall.anzrewards.com",
	"cadencebank.bankmyrewards.com",
	"cadencebank.eversave.com",
	"candfonline.bankmyrewards.com",
	"careOne.mymembersmarketplace.com",
	"cashback.achatvipservices.com",
	"cashback.acquistierisparmi.it",
	"cashback.alitalia.com",
	"cashback.beneficiosonline.com.mx",
	"cashback.bilevip.com",
	"cashback.butwaitsavemore.com",
	"cashback.buyerassurance.com",
	"cashback.cafepressclub.com",
	"cashback.completehome.com",
	"cashback.completesavings.co.uk",
	"cashback.completesavings.ie",
	"cashback.compraevolta.com.br",
	"cashback.designerdealsclub.com",	
}

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

	src = "http://go.com"
	u, err = FQDNMgr.GetFQDN(src)
	if err!=nil || u != "go.com" {
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

	for _, src = range CXHOSTS {
		u, err = FQDNMgr.GetFQDN(src)
		if err != nil {
			t.Fail()
		}			
	}


}

func TestELD(t *testing.T) {
	etLD := emptyETLD(1)
	etLD.Add("nlaak.com", true)
}

func TestFQDNGuess(t *testing.T) {
	_, err := FQDNMgr.guess("abc",1)
	if err == nil {
		t.Fail()
	}	

	_, err = FQDNMgr.guess("test.etld.with.alot.of.groups",5)
	if err != nil {
		t.Fail()
	}	

	_, err = FQDNMgr.guess("test.etld.with.alot.of.groups",4)
	if err != nil {
		t.Fail()
	}	

	_, err = FQDNMgr.guess("test.to.short",5)
	if err == nil {
		t.Fail()
	}	

}

func TestFileFail(t *testing.T) {
	err := FQDNMgr.downloadPublicSuffixFile("")
	if err != nil {
		t.Fail()
	}	

	err = FQDNMgr.downloadPublicSuffixFile("https://publicsuffix.org/list/this-file-does-not-exist.bad")
	if err == nil {	//Show be 403 Forbidden Error
		t.Fail()
	}	

	err = FQDNMgr.downloadPublicSuffixFile("https://nlaak.com")
	if err == nil {	//200 OK, but Resp Content will be less that 200k
		t.Fail()
	}	

	err = FQDNMgr.downloadPublicSuffixFile("https://en.wikipedia.org/wiki/Main_Page")
	if err == nil {	//200 OK,Resp Content big enough, but is not actual data file
		t.Fail()
	}	
	
}

