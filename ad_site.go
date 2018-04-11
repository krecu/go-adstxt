package go_adstxt

type AdSite struct {
	Host  string      `json:"host"`
	Ads   []*AdSystem `json:"ads"`
	Error []error     `json:"error"`
}
