package go_adstxt

type AdSite struct {
	Host  string
	Ads   []*AdSystem
	Error []error
}
