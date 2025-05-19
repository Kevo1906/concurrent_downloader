package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func DownloadFile(ctx context.Context, url string, outputPath string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	
	select {
	case <-ctx.Done():
		fmt.Printf("Cancelled before starting: %s\n", url)
		return
	default:

	}

	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	fileName := getFileNameFromURL(url)
	fmt.Printf("Downloading: %s -> %s\n", url, fileName)

	// Download logic
	resp, err := DownloadFileWithRetries(ctx, url, 5)
	if err != nil {
		fmt.Printf("Failed to download %s after retries: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Create output file
	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	filePath := path.Join(outputPath, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", fileName, err)
		return
	}
	defer out.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("Saving %s", fileName),
	)

	// Copy the content
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		fmt.Printf("Error saving %s: %v\n", filePath, err)
		return
	}

	fmt.Printf("\nDownloaded: %s\n", filePath)
}

func DownloadFileWithRetries(ctx context.Context, url string, maxRetries int) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 1; i <= maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("download cancelled: %s", url)
		default:

		}

		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		fmt.Printf("Attempt %d failed for %s. Retrying in 3 seconds...\n", i, url)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("download cancelled during wait: %s", url)
		case <-time.After(3 * time.Second):
		}
	}

	return nil, fmt.Errorf("failed to download %s after %d attempts: %v", url, maxRetries, err)

}

func getFileNameFromURL(url string) string {
	base := path.Base(url)

	if base == "." || base == "/" || base == "" {
		return "downloaded_file"
	}

	if strings.Contains(base, "?") {
		base = strings.Split(base, "?")[0]
	}
	fmt.Printf("BASE: %s\n",base)
	return base
}

func StartDownloadPool(ctx context.Context, urls []string, outputPath string, workerCount int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workerCount)

	for _, url := range urls {
		wg.Add(1)
		go DownloadFile(ctx, url, outputPath, &wg, semaphore)
	}

	wg.Wait()
	fmt.Println("All downloads completed")
}
