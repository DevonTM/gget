package download

import (
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
)

// set http client transport proxy
func SetProxy(p *url.URL) {
	transport.Proxy = http.ProxyURL(p)
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

	if d.CookiesFile != "" {
		var cookies []*http.Cookie
		cookies, err = cookie.ParseFile(d.CookiesFile)
		if err != nil {
			return nil, err
		}
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
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	// get unescaped path
	path, err := url.QueryUnescape(parsedURL.EscapedPath())
	if err != nil {
		return "", err
	}

	// return base path
	return filepath.Base(path), nil
}
