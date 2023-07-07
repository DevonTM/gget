package download

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DevonTM/gget/progress"
)

// download file with concurrent chunk and resume capability
func (d *Info) WithRange(s *Setting) error {
	// check if server support specified range
	resp, err := d.getResponse(fmt.Sprintf("bytes=0-%d", d.FileSize))
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 206 {
		return fmt.Errorf("server does not support range request")
	}

	pd := &progress.Data{
		FileSize:  d.FileSize,
		ChunkSize: s.ChunkSize,
	}

	pr := &progress.Report{
		Total: d.FileSize,
	}

	done := make(chan struct{})
	go pr.ShowReport(done)

	// create cancel channel to capture interrupt signal
	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)

	// cancel download if interrupt signal received
	go func() {
		select {
		case <-cancel:
			close(done)
			time.Sleep(1 * time.Second)
			err = pd.SaveProgressFile(d.Path + ".gget")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Download canceled")
			os.Exit(0)
		case <-done:
			return
		}
	}()

	// load download progress
	err = pd.LoadProgress(d.Path + ".gget")
	if err != nil {
		return err
	}

	// save download progress every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err = pd.SaveProgressFile(d.Path + ".gget")
				if err != nil {
					fmt.Println(err)
				}
			case <-done:
				return
			}
		}
	}()

	// if file not exists, create new file
	if _, err = os.Stat(d.Path); os.IsNotExist(err) {
		_, err = os.Create(d.Path)
		if err != nil {
			return err
		}
	}

	// create wait group based on number of chunk
	var wg sync.WaitGroup
	wg.Add(len(pd.Chunk))

	// create semaphore to limit concurrent chunk download
	semaphore := make(chan struct{}, s.Thread)

	// create fail chunk channel to limit max fail chunk
	failChunk := make(chan struct{}, pd.ChunkSize)

	for i := range pd.Chunk {
		// skip chunk if already finished / max fail reached
		if pd.Chunk[i] {
			pr.Current += pd.ChunkSize
			wg.Done()
			continue
		} else if len(failChunk) > s.Thread*5 {
			wg.Done()
			continue
		}

		semaphore <- struct{}{}

		go func(index int) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()

			start := int64(index) * pd.ChunkSize
			end := start + pd.ChunkSize - 1
			if index == len(pd.Chunk)-1 {
				end = d.FileSize
			}

			resp, errr := d.getResponse(fmt.Sprintf("bytes=%d-%d", start, end))
			if errr != nil {
				failChunk <- struct{}{}
				return
			}
			defer resp.Body.Close()

			file, errr := os.OpenFile(d.Path, os.O_WRONLY, 0o644)
			if errr != nil {
				failChunk <- struct{}{}
				return
			}
			defer file.Close()

			_, errr = file.Seek(start, io.SeekStart)
			if errr != nil {
				failChunk <- struct{}{}
				return
			}

			_, errr = io.Copy(file, io.TeeReader(resp.Body, pr))
			if errr != nil {
				failChunk <- struct{}{}
				return
			}

			pd.Chunk[index] = true
		}(i)
	}

	// wait until all chunk finished
	wg.Wait()
	close(done)
	time.Sleep(1 * time.Second)

	// restart download if chunk uncompleted
	if !pd.IsFinished() {
		err = pd.SaveProgressFile(d.Path + ".gget")
		if err != nil {
			return err
		}

		if s.retry < s.MaxRetry {
			s.retry++
			fmt.Printf("Retry %d\n", s.retry)
			return d.WithRange(s)
		} else if s.MaxRetry == -1 {
			fmt.Println("Retry ∞")
			return d.WithRange(s)
		}

		return fmt.Errorf("max retry reached\ntry again later or refresh download url")
	}

	// remove progress file
	err = os.Remove(d.Path + ".gget")
	if err != nil {
		return err
	}

	return nil
}
