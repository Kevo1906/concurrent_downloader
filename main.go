package main

import (
	"concurrent_downloader/downloader"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	fileFlag := flag.String("file", "", "Path to file containing list of URLs (one per line)" )
	urlsFlag := flag.String("urls", "", "Comma-separated list of URLs to download")
	outputFlag := flag.String("output", `./downloaded_files`, "Output path to download files")
	workersFlag := flag.Int("workers", 5, "Number of concurrent downloads (dafault: 5)")

	flag.Parse()

	var urls []string

	// Load URLs from file
	if *fileFlag != ""{
		fileContent, err := os.ReadFile(*fileFlag)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		lines := strings.Split(string(fileContent), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != ""{
				urls = append(urls, line)
			}
		}
	}

	// Load URLs from command line
	if *urlsFlag != ""{
		parts := strings.Split(*urlsFlag,",")
		for _, url:= range parts {
			url = strings.TrimSpace(url)
			if url != ""{
				urls = append(urls, url)
			}
		}
	}

	// Validate input
	if len(urls) == 0 {
		fmt.Println("Error: You must provide URLs using --file or --urls")
		os.Exit(1)
	}

	fmt.Printf("Starting download of %d files using %d workers...\n", len(urls), *workersFlag)

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C or kill signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(){
		<- c
		fmt.Println("\nInterrumpt signal received. Cancelling downloads...")
		cancel()
	}()


	downloader.StartDownloadPool(ctx,urls, *outputFlag, *workersFlag)
}