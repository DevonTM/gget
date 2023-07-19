package progress

import (
	"fmt"
	"sync"
	"time"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

type Report struct {
	Total   int64 // total file size
	Current int64 // current file size
	written int64 // current written bytes
	mutex   sync.Mutex
}

// monitor written bytes
func (pr *Report) Write(p []byte) (int, error) {
	n := len(p)

	pr.mutex.Lock()
	pr.Current += int64(n)
	pr.written += int64(n)
	pr.mutex.Unlock()

	return n, nil
}

// display progress report
func (pr *Report) ShowReport(done chan struct{}) {
	startTime := time.Now()
	interval := 1 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// calculate progress report
			percent := float64(pr.Current) / float64(pr.Total) * 100
			remaining := pr.Total - pr.Current
			speed := float64(pr.written) / time.Since(startTime).Seconds()
			eta := time.Duration(float64(remaining)/speed) * time.Second

			// convert download speed to human readable format
			var speedString string
			switch {
			case speed > GB:
				speedString = fmt.Sprintf("%.2f GB/s", speed/GB)
			case speed > MB:
				speedString = fmt.Sprintf("%.2f MB/s", speed/MB)
			case speed > KB:
				speedString = fmt.Sprintf("%.f KB/s", speed/KB)
			default:
				speedString = fmt.Sprintf("%.f B/s", speed)
			}

			fmt.Print("\033[2K\r")
			fmt.Printf("%s / %s (%.2f%%) | %s | ETA %s", convertSize(pr.Current), convertSize(pr.Total), percent, speedString, eta)

		case <-done:
			fmt.Print("\033[2K\r")
			elapsed := time.Since(startTime).Round(time.Second)
			fmt.Printf("Elapsed time: %s\n", elapsed)
			return
		}
	}
}

// convert size from bytes to human readable format
func convertSize(size int64) string {
	switch {
	case size > GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size > MB:
		return fmt.Sprintf("%d MB", size/MB)
	case size > KB:
		return fmt.Sprintf("%d KB", size/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
