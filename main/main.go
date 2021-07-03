package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

/*
 * TODO:
 * - Figure out how to grab the download URL from <img class=".big-gopher"> ✅
 * - Figure out how to actually save the image to current directory ❌
 */

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		// chromedp.Flag("disable-gpu", false),
		// chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	// ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var downloadURL string
	var ok bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://gopherize.me/"),
		chromedp.Click("#shuffle-button", chromedp.ByQuery),
		chromedp.Click("#next-button", chromedp.ByQuery),
		chromedp.AttributeValue(".big-gopher", "src", &downloadURL, &ok, chromedp.ByQuery),
		// chromedp.WaitNotPresent(".big-gopher", chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal("Could not scrape img src")
	}

	fmt.Printf("Download URL: %s\n", downloadURL)
}
