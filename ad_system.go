package go_adstxt

import "fmt"

const (
	ADS_DIRECT   = "DIRECT"
	ADS_RESELLER = "RESELLER"
	ADS_SELLER   = "SELLER"
)

type AdSystem struct {

	/* Domain name of the advertising system
	(Required) The canonical domain name of the SSP, Exchange, Header Wrapper, etc system tha bidders connect to.
	This may be the operational domain of the system, if that is different than the parent corporate domain, to
	facilitate WHOIS and reverse IP lookups to establish clear ownership of the delegate system. Ideally the SSP or
	Exchange publishes a document detailing what domain name to use.
	*/
	Domain string `json:"domain"`

	/* Publisher’s Account ID
	(Required) The identifier associated with the seller or reseller account within the advertising system in
	field #1. This must contain the same value used in transactions (i.e. OpenRTB bid requests) in the
	field specified by the SSP/exchange. Typically, in OpenRTB, this is publisher.id. For OpenDirect it is
	typically the publisher’s organization ID.
	*/
	ID string `json:"id"`

	/* Type of Account/Relationship
	(Required) An enumeration of the type of account. A value of ‘DIRECT’ indicates that the Publisher (content owner)
	directly controls the account indicated in field #2 on the system in field #1. This tends to mean a direct business
	contract between the Publisher and the advertising system. A value of ‘RESELLER’ indicates that the Publisher has
	authorized another entity to control the account indicated in field #2 and resell their ad space via the system in
	field #1. Other types may be added in the future. Note that this field should be treated as case insensitive when
	interpreting the data.
	*/
	Type string `json:"type"`

	/* Certification Authority ID
	(Optional) An ID that uniquely identifies the advertising system within a certification authority
	(this ID maps to the entity listed in field #1). A current certification authority is the Trustworthy
	Accountability Group (aka TAG), and the TAGID would be included here [11].
	*/
	AuthorityID string `json:"authority_id"`
}

func (ads *AdSystem) Validate() (err []error) {

	if len(ads.Domain) < 3 {
		err = append(err, ErrAdsWrongDomain)
	}
	if len(ads.ID) < 1 {
		err = append(err, ErrAdsWrongId)
	}
	if len(ads.Type) < 6 || (ads.Type != ADS_DIRECT && ads.Type != ADS_RESELLER && ads.Type != fmt.Sprintf("%s/%s", ADS_SELLER, ADS_RESELLER)) {
		err = append(err, ErrAdsWrongType)
	}

	return
}
