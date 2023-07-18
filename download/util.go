package download

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"

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

// get file name from http response
func getFileName(resp *http.Response) string {
	// get file name from content-disposition
	disposition := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(disposition)
	if err == nil {
		return params["filename"]
	}

	// get file name from url
	fileName, err := getFileNameFromURL(resp.Request.URL.String())
	if err == nil && fileName != "." && filepath.Ext(fileName) != "" {
		return fileName
	}

	if fileName == "" || fileName == "\\" || fileName == "/" || fileName == "." {
		fileName = "unknown"
	}

	// append extension from content-type
	extension, err := mime.ExtensionsByType(resp.Header.Get("Content-Type"))
	if err == nil {
		return fileName + extension[0]
	}

	return fileName
}

// get filename from url
func getFileNameFromURL(u string) (string, error) {
	// parse url
	URL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	// get unescaped path
	path, err := url.QueryUnescape(URL.EscapedPath())
	if err != nil {
		return "", err
	}

	// return base path
	return filepath.Base(path), nil
}
