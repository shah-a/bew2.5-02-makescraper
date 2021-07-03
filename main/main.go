package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

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
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://gopherize.me/"),
		chromedp.Click("#shuffle-button", chromedp.ByQuery),
		chromedp.Click("#next-button", chromedp.ByQuery),
		chromedp.WaitNotPresent(".big-gopher", chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}
}
