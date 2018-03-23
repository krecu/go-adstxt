package go_adstxt

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AdsTxt struct {
	client *http.Client
	opt    Options
}

type Options struct {
	HttpTimeout             time.Duration
	HttpMaxIdleConns        int
	HttpMaxIdleConnsPerHost int
}

func New(opt Options) (proto *AdsTxt, err error) {

	if opt.HttpTimeout == 0 {
		opt.HttpTimeout = time.Duration(1) * time.Second
	}

	if opt.HttpMaxIdleConns == 0 {
		opt.HttpMaxIdleConns = 100
	}

	if opt.HttpMaxIdleConnsPerHost == 0 {
		opt.HttpMaxIdleConnsPerHost = 100
	}

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

		proto.client = &http.Client{Transport: &defaultTransport, Timeout: opt.HttpTimeout}
	}

	return
}

func (ads *AdsTxt) CheckMulti(hosts ...string) (sites []*AdSite) {

	var (
		requestPool = make(chan *AdSite)
	)

	// формируем пулл запросов
	for _, host := range hosts {
		go func(host string) {
			requestPool <- ads.Check(host)
		}(host)
	}

	for {

		select {
		case ads := <-requestPool:
			sites = append(sites, ads)
			if len(sites) == len(hosts) {
				return
			}
		}
	}

	return
}

// Validation ads.txt lines
func (ads *AdsTxt) Check(host string) (site *AdSite) {

	var (
		req    *http.Request
		res    *http.Response
		reader *csv.Reader
		pos    = 0
		err    error
	)

	site = &AdSite{
		Host: host,
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/ads.txt", host), nil)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("User-Agent", "AdsTxtCrawler/1.0; +https://github.com/krecu/go-adstxt")

	res, err = ads.client.Do(req)
	if err != nil {
		site.Error = append(site.Error, ErrAdsReqFail)
		return
	}

	// check status
	switch res.StatusCode {
	case http.StatusOK:
		reader = csv.NewReader(bufio.NewReader(res.Body))
		// disable check count field
		reader.FieldsPerRecord = -1
	case http.StatusNoContent:
		site.Error = append(site.Error, ErrAdsReqNotFound)
	default:
		site.Error = append(site.Error, ErrAdsReqStatus)
	}

	if len(site.Error) != 0 {
		return
	} else {

		for {

			pos++

			// parse csv string
			if line, err := reader.Read(); err != nil {
				if err == io.EOF {
					break
				} else if err != nil {
					site.Error = append(site.Error, fmt.Errorf("LINE: %d, err: %s", pos, err))
					continue
				}
			} else {
				// clear space
				for i, l := range line {
					line[i] = strings.Trim(l, " ")
				}

				// check require field
				switch len(line) {
				case 0:
					site.Error = append(site.Error, fmt.Errorf("LINE: %d, err: %s", pos, ErrAdsWrongLine))
				case 1:
					site.Error = append(site.Error, fmt.Errorf("LINE: %d, err: %s", pos, ErrAdsWrongId))
				case 2:
					site.Error = append(site.Error, fmt.Errorf("LINE: %d, err: %s", pos, ErrAdsWrongType))
				default:
					ad := &AdSystem{
						Domain: line[0],
						ID:     line[1],
						Type:   strings.ToUpper(line[2]),
					}

					if len(line) > 3 {
						ad.AuthorityID = line[3]
					}

					if errs := ad.Validate(); len(errs) != 0 {
						site.Error = append(site.Error, fmt.Errorf("LINE: %d, err ads: %s", pos, errs))
					}

					site.Ads = append(site.Ads, ad)
				}

			}

		}

	}

	return
}
