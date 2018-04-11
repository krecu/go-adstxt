package go_adstxt

import "fmt"

var (
	ErrAdsReqFail     = fmt.Errorf("ads.txt request fail")
	ErrAdsReqHeader   = fmt.Errorf("ads.txt request bad content type")
	ErrAdsReqStatus   = fmt.Errorf("ads.txt bad response status")
	ErrAdsReqNotFound = fmt.Errorf("ads.txt not found")

	ErrAdsWrongLine   = fmt.Errorf("ads.txt bad line")
	ErrAdsWrongDomain = fmt.Errorf("ads.txt bad domain")
	ErrAdsWrongId     = fmt.Errorf("ads.txt bad id")
	ErrAdsWrongType   = fmt.Errorf("ads.txt bad type")
)
