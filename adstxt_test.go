package go_adstxt

import (
	"testing"
)

func TestAdxTxt_Validate(t *testing.T) {

	adx, err := New(Options{
		HttpMaxIdleConnsPerHost: 100,
		HttpMaxIdleConns: 100,
	})

	if err = adx.Validate("https://www.sports.ru/ads.txt"); err != nil {
		t.Fatal(err.Error())
	}

}
