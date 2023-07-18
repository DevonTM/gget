package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DevonTM/gget/download"
)

const (
	VERSION          = "0.1.0"
	DefaultUserAgent = "Gget/" + VERSION
)

type headerFlags []string

func (h *headerFlags) String() string {
	return ""
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

var (
	urlPath  = flag.String("url", "", "Download URL")
	path     = flag.String("o", "", "Output Path (default current directory)")
	fileName = flag.String("O", "", "Set File Name (default from server)")
	force    = flag.Bool("f", false, "Overwrite Existing File")

	headers   = headerFlags{}
	cookies   = flag.String("C", "", "Set Cookies File")
	referer   = flag.String("r", "", "Set HTTP Referer")
	userAgent = flag.String("ua", DefaultUserAgent, "HTTP User-Agent")
	proxy     = flag.String("p", "", "Set Proxy, support http, https, and socks5")

	chunkSize = flag.String("c", "1M", "Chunk Size")
	thread    = flag.Int("j", 4, "Number of Threads (1-16)")
	maxRetry  = flag.Int("t", 3, "Maximum Retry, set -1 for infinite retry")
	noH2      = flag.Bool("no-h2", false, "Disable HTTP/2")

	help    = flag.Bool("help", false, "Print Help")
	version = flag.Bool("version", false, "Print Version")
)

func init() {
	flag.Var(&headers, "H", "Set Custom HTTP Header")
	flag.Parse()
}

func main() {
	if *help {
		flag.Usage()
		return
	}
	if *version {
		fmt.Println("Gget Version: " + VERSION)
		return
	}
	if *urlPath == "" {
		if len(flag.Args()) > 0 {
			*urlPath = flag.Args()[0]
		} else {
			fmt.Println("please input URL")
			return
		}
	}
	if *noH2 {
		os.Setenv("GODEBUG", "http2client=0")
	}

	parsedURL, err := url.Parse(*urlPath)
	if err != nil {
		fmt.Println("invalid URL")
		return
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		fmt.Printf("%s scheme not supported", parsedURL.Scheme)
		return
	}

	if *cookies != "" {
		err = download.SetCookies(*cookies)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if *proxy != "" {
		err = download.SetProxy(*proxy)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if *thread < 1 || *thread > 16 {
		fmt.Println("thread must be between 1 and 16")
		return
	}
	chunk, err := parseSize(*chunkSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	if chunk < 1024*1024 {
		fmt.Println("chunk size must be at least 1 megabytes")
		return
	}

	s := &download.Setting{
		ChunkSize: chunk,
		Thread:    *thread,
		MaxRetry:  *maxRetry,
		Force:     *force,
	}

	header, err := parseHeader(headers)
	if err != nil {
		fmt.Println(err)
		return
	}

	paths := "./"
	if *path != "" {
		paths = *path
	}
	if *fileName != "" {
		paths = filepath.Join(paths, *fileName)
	}
	dir, fileName := filepath.Split(paths)
	if dir == "" {
		dir = "."
	}
	if _, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0o755)
			if err != nil {
				fmt.Println("cannot create directory")
				return
			}
		}
	}

	d := &download.Info{
		URL:      *urlPath,
		Path:     paths,
		Header:   header,
		FileName: fileName,
	}

	err = d.Manager(s)
	if err != nil {
		fmt.Println("Download failed:", err)
		return
	}
	fmt.Println("Download completed")
}

func parseSize(s string) (size int64, err error) {
	s = strings.ToUpper(s)
	suffixes := map[string]int64{
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
	}
	for suffix, factor := range suffixes {
		if strings.HasSuffix(s, suffix) {
			value := strings.TrimSuffix(s, suffix)
			size, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return 0, err
			}
			return size * factor, nil
		}
	}
	size, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func parseHeader(h []string) (header map[string]string, err error) {
	header = make(map[string]string)
	for _, hh := range h {
		hh = strings.TrimSpace(hh)
		if hh == "" {
			continue
		}
		keyVal := strings.SplitN(hh, ":", 2)
		if len(keyVal) != 2 {
			return nil, fmt.Errorf("invalid header format")
		}
		header[strings.TrimSpace(keyVal[0])] = strings.TrimSpace(keyVal[1])
	}
	if header["Accept"] == "" {
		header["Accept"] = "*/*"
	}
	if *referer != "" {
		header["Referer"] = *referer
	}
	if *userAgent != "" {
		header["User-Agent"] = *userAgent
	} else if header["User-Agent"] == "" {
		header["User-Agent"] = DefaultUserAgent
	}
	return header, nil
}
