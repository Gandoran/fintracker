package scraper

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func (f *Fetcher) fetchHTMLSlow(ctx context.Context, targetURL string) string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.ExecPath(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`), //TODO insert on config.yaml
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()
	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()
	taskCtx, cancelTimeout := context.WithTimeout(taskCtx, 30*time.Second)
	defer cancelTimeout()
	var htmlContent string
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML(`html`, &htmlContent),
	)
	if err != nil {
		log.Printf("⚠️ [SLOW SCRAPER] Errore navigazione su %s: %v", targetURL, err)
		return ""
	}
	return htmlContent
}

//For test: https://www.theverge.com
