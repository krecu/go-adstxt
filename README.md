## VideoNow ads.txt crawler & validator


## Installation
```
go get https://github.com/krecu/go-adstxt
cd $GOPATH/src/github.com/krecu/go-adstxt && go test -v
```

## Usage
``` go
if adx, err := New(Options{
    HttpTimeout: time.Duration(1) * time.Second,
    HttpMaxIdleConnsPerHost: 100,
    HttpMaxIdleConns:        100,
}); err == nil {
    // use once or mutiple check
    site := adx.Check("https://www.sports.ru")
    fmt.Errorf("ADS: %+v", site.Ads)
    fmt.Errorf("ERR: %+v", site.Error)
} 
```

## License

The open source license used is the 2-clause BSD license
