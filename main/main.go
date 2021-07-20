package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// parseFlags checks for `--name` flag
func parseFlags() string {
	var name string
	flag.StringVar(&name, "name", "my", "name to parse for gopher's filename")
	flag.StringVar(&name, "n", "my", "short form of \"name\" flag")

	flag.Parse()

	if name == "" {
		return "my-gopher.png"
	}

	return fmt.Sprintf("%s-gopher.png", name)
}

func Scrape(ctx context.Context) string {
	var downloadUrl string
	var ok bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://gopherize.me/"),
		chromedp.WaitVisible("#shuffle-button", chromedp.ByQuery),
		chromedp.Click("#shuffle-button", chromedp.ByQuery),
		chromedp.Click("#next-button", chromedp.ByQuery),
		chromedp.WaitVisible(".big-gopher", chromedp.ByQuery),
		chromedp.AttributeValue(".big-gopher", "src", &downloadUrl, &ok, chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal("Could not scrape img src")
	}

	return downloadUrl
}

func main() {
	filename := parseFlags()

	opts := append(chromedp.DefaultExecAllocatorOptions[:])
	// Add this to the `opts` append operation to disable headless mode (i.e. to see what the scraper is doing):
	// chromedp.Flag("headless", false)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	// ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Call scraper function and get resultant downloadUrl
	downloadUrl := Scrape(ctx)

	// set up a channel so we can block later while we monitor the download progress
	downloadComplete := make(chan bool)

	// this will be used to capture the request id for matching network events
	var requestId network.RequestID

	// set up a listener to watch the network events and close the channel when complete
	// the request id matching is important both to filter out unwanted network events
	// and to reference the downloaded file later
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			// fmt.Printf("EventRequestWillBeSent: %v: %v\n", ev.RequestID, ev.Request.URL)
			if ev.Request.URL == downloadUrl {
				requestId = ev.RequestID
			}
		case *network.EventLoadingFinished:
			// fmt.Printf("EventLoadingFinished: %v\n", ev.RequestID)
			if ev.RequestID == requestId {
				close(downloadComplete)
			}
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(downloadUrl),
	); err != nil {
		log.Fatal(err)
	}

	// this will block until the chromedp listener closes the channel
	<-downloadComplete

	var downloadBytes []byte
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		downloadBytes, err = network.GetResponseBody(requestId).Do(ctx)
		return err
	})); err != nil {
		log.Fatal(err)
	}

	// write the file to disk - since we hold the bytes we dictate the name and location
	downloadDest := filepath.Join(".", filename)
	if err := ioutil.WriteFile(downloadDest, downloadBytes, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Download Complete! Your gopher was saved as `%s`\n", downloadDest)
}
