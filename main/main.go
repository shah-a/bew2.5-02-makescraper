package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

/*
 * TODO:
 * - Figure out how to grab the download URL from <img class=".big-gopher">
 * - Figure out how to actually save the image to current directory
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
	defer cancel()

	downloadURL := "https://storage.googleapis.com/gopherizeme.appspot.com/gophers/3fbbd6ea507241a0e663656182ac58d2b811b5f1.png"

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://gopherize.me/"),
		chromedp.Click("#shuffle-button", chromedp.ByQuery),
		chromedp.Click("#next-button", chromedp.ByQuery),
		chromedp.WaitNotPresent(".big-gopher", chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}
}
