package download

import (
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
)

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
