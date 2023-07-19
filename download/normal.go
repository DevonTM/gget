package download

import (
	"io"
	"os"
	"time"

	"github.com/DevonTM/gget/progress"
)

// download file normally without concurrency and resume capability
func (d *Info) WithoutRange() error {
	file, err := os.Create(d.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := d.getResponse()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	pr := &progress.Report{
		Total: d.fileSize,
	}

	done := make(chan struct{}, 1)
	go pr.ShowReport(done)

	_, err = io.Copy(file, io.TeeReader(resp.Body, pr))
	if err != nil {
		return err
	}

	done <- struct{}{}
	time.Sleep(1 * time.Second)

	return nil
}
