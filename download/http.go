package download

import (
	"fmt"
	"net/http"
	"net/url"

	cookie "github.com/MercuryEngineering/CookieMonster"
)

var (
	transport = &http.Transport{
		MaxIdleConnsPerHost: 100,
	}

	client = &http.Client{
		Transport: transport,
	}

	cookies []*http.Cookie
)

// set http client transport proxy
func SetProxy(p string) error {
	URL, err := url.Parse(p)
	if err != nil {
		return err
	}

	if URL.Scheme != "http" && URL.Scheme != "https" && URL.Scheme != "socks5" {
		return fmt.Errorf("only http, https and socks5 proxy are supported")
	}

	transport.Proxy = http.ProxyURL(URL)
	return nil
}

// set cookies from netscape cookies file
func SetCookies(f string) (err error) {
	cookies, err = cookie.ParseFile(f)
	return
}

// send http request and return http response
// with optional range header
func (d *Info) getResponse(ranges ...string) (*http.Response, error) {
	req, err := http.NewRequest("GET", d.URL, http.NoBody)
	if err != nil {
		return nil, err
	}

	for key, value := range d.Header {
		req.Header.Set(key, value)
	}

	if len(cookies) > 0 {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	if len(ranges) > 0 {
		req.Header.Set("Range", ranges[0])
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
