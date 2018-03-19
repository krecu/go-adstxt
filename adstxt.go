package go_adstxt

import (
	"net/http"
	"fmt"
	"github.com/k0kubun/pp"
	"bufio"
	"encoding/csv"
	"io"
	"strings"
)

type Ads struct {

	/* Domain name of the advertising system
	(Required) The canonical domain name of the SSP, Exchange, Header Wrapper, etc system tha bidders connect to.
	This may be the operational domain of the system, if that is different than the parent corporate domain, to
	facilitate WHOIS and reverse IP lookups to establish clear ownership of the delegate system. Ideally the SSP or
	Exchange publishes a document detailing what domain name to use.
	*/
	Domain string

	/* Publisher’s Account ID
	(Required) The identifier associated with the seller or reseller account within the advertising system in
	field #1. This must contain the same value used in transactions (i.e. OpenRTB bid requests) in the
	field specified by the SSP/exchange. Typically, in OpenRTB, this is publisher.id. For OpenDirect it is
	typically the publisher’s organization ID.
	 */
	ID string

	/* Type of Account/Relationship
	(Required) An enumeration of the type of account. A value of ‘DIRECT’ indicates that the Publisher (content owner)
	directly controls the account indicated in field #2 on the system in field #1. This tends to mean a direct business
	contract between the Publisher and the advertising system. A value of ‘RESELLER’ indicates that the Publisher has
	authorized another entity to control the account indicated in field #2 and resell their ad space via the system in
	field #1. Other types may be added in the future. Note that this field should be treated as case insensitive when
	interpreting the data.
	 */
	Type string

	/* Certification Authority ID
	(Optional) An ID that uniquely identifies the advertising system within a certification authority
	(this ID maps to the entity listed in field #1). A current certification authority is the Trustworthy
	Accountability Group (aka TAG), and the TAGID would be included here [11].
	 */
	AuthorityID string
}

func (ads *Ads) Validate() (err []error){

	if len(ads.Domain) < 3 {

	}
	if len(ads.ID) < 1 {

	}
	if len(ads.Type) < 6 {

	}


	return
}

type AdsTxt struct {
	client *http.Client
	opt Options
}

type Options struct {
	HttpMaxIdleConns int
	HttpMaxIdleConnsPerHost int
}

var (
	ErrAdsNotFound = fmt.Errorf("ads.txt not found")
	ErrAdsRequire = fmt.Errorf("ads.txt not found")
	ErrAdsRequireDomain = fmt.Errorf("ads.txt not found")
	ErrAdsRequireId = fmt.Errorf("ads.txt not found")
	ErrAdsRequireType = fmt.Errorf("ads.txt not found")
	ErrAdsWrongDomain = fmt.Errorf("ads.txt not found")
	ErrAdsWrongId = fmt.Errorf("ads.txt not found")
	ErrAdsWrongType = fmt.Errorf("ads.txt not found")
	ErrAdsReqFail = fmt.Errorf("ads.txt request fail")
)

const (
	ADS_DIRECT = "DIRECT"
	ADS_RESELLER = "RESELLER"
)

func New(opt Options) (proto *AdsTxt, err error) {

	proto = &AdsTxt{
		opt: opt,
	}

	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		err = fmt.Errorf("defaultRoundTripper not an *http.Transport")
	} else {
		// dereference it to get a copy of the struct that the pointer points to
		defaultTransport := *defaultTransportPointer
		defaultTransport.MaxIdleConns = proto.opt.HttpMaxIdleConns
		defaultTransport.MaxIdleConnsPerHost = proto.opt.HttpMaxIdleConnsPerHost

		proto.client = &http.Client{Transport: &defaultTransport}
	}

	return
}

func (adx *AdsTxt) Validate(url string) (err error) {

	var (
		req *http.Request
		res *http.Response
		reader *csv.Reader
		ads []*Ads
	)

	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("User-Agent", "Ads.txt Crawler")

	res, err = adx.client.Do(req)

	switch res.StatusCode {
	case http.StatusOK :
		reader = csv.NewReader(bufio.NewReader(res.Body))
		// disable check count field
		reader.FieldsPerRecord = -1
	case http.StatusNoContent :
		err = ErrAdsNotFound
	default:
		err = ErrAdsReqFail
	}

	if err != nil {
		return
	} else {

		for {

			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				pp.Fatal(err)
			}

			for i, l := range line {
				line[i] = strings.Trim(l, " ")
			}

			// check require field
			switch len(line) {
			case 0:
				err = ErrAdsRequire
			case 1:
				err = ErrAdsRequireId
			case 2:
				err = ErrAdsRequireType
			default:
				ad := &Ads{
					Domain:      line[0],
					ID:          line[1],
					Type:        line[2],
				}
				if len(line) > 3 {
					ad.AuthorityID = line[3]
				}

				if errs := ad.Validate(); errs == nil {
					ads = append(ads, ad)
				}
			}

		}
	}

	pp.Println(ads)
	return
}