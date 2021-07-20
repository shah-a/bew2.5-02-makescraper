package main

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func BenchmarkScrape(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		Scrape(ctx)
	}
}
