package go_adstxt

import (
	"testing"
	"time"
)

func TestAdxTxt_Check(t *testing.T) {

	adx, err := New(Options{
		HttpMaxIdleConnsPerHost: 100,
		HttpMaxIdleConns:        100,
	})

	if err != nil {
		t.Fail()
	}

	adsGood := adx.Check("http://www.shkolazhizni.ru")
	if len(adsGood.Ads) == 0 {
		t.Fail()
	} else {
		t.Logf("Good %+v", adsGood)
	}

	adsBad := adx.Check("htt://www.sports.ru")
	if len(adsBad.Error) == 0 {
		t.Fail()
	} else {
		t.Logf("Bad %+v", adsBad)
	}

}

func TestAdsTxt_CheckMulti(t *testing.T) {
	adx, err := New(Options{
		HttpTimeout:             time.Duration(100) * time.Millisecond,
		HttpMaxIdleConnsPerHost: 100,
		HttpMaxIdleConns:        100,
	})

	if err != nil {
		t.Fail()
	}

	ads := adx.CheckMulti("https://www.sports.ru", "https://www.sports.ru", "https://www.sports.ru")

	if len(ads) != 3 {
		t.Fail()
	} else {
		t.Logf("Sites %+v", ads)
	}

}
