package download

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type Info struct {
	URL         string
	Path        string
	Header      map[string]string
	CookiesFile string
	FileName    string
	FileSize    int64
}

type Setting struct {
	ChunkSize int64
	Thread    int
	MaxRetry  int
	retry     int
	Force     bool
}

// manage and download file
func (d *Info) Manager(s *Setting) error {
	resp, err := d.getResponse()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check http status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	// set download url to the redirected url
	if resp.Request.URL.String() != d.URL {
		d.URL = resp.Request.URL.String()
	}

	// set file size from http response
	d.FileSize = resp.ContentLength

	// set file name if not set
	if d.FileName == "" {
		d.FileName = getFileName(resp)
		d.Path = filepath.Join(d.Path, d.FileName)
	}

	// check if download file and progress file exist
	if _, err = os.Stat(d.Path); err == nil {
		if _, err = os.Stat(d.Path + ".gget"); os.IsNotExist(err) {
			if s.Force {
				err = os.Remove(d.Path)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("file already exist")
			}
		} else {
			if s.Force {
				err = os.Remove(d.Path)
				if err != nil {
					return err
				}
				err = os.Remove(d.Path + ".gget")
				if err != nil {
					return err
				}
			}
		}
	}

	fmt.Println("Downloading:", d.FileName)

	// check if server accept ranges, else download without concurrency and resume
	if resp.ContentLength > 0 && resp.Header.Get("Accept-Ranges") == "bytes" {
		err = d.WithRange(s)
		if err != nil {
			return err
		}
	} else {
		err = d.WithoutRange()
		if err != nil {
			return err
		}
	}

	return nil
}
