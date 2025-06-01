// Command visible is a chromedp example demonstrating how to wait until an
// element is visible.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := initChromeContext()
	defer cancel()

	err := chromedp.Run(ctx, findPlaces("Барбершоп")...)
	if err != nil {
		log.Fatal(err)
	}

	var htmlContent string
	err = saveHtml(ctx, htmlContent)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("output.html", []byte(htmlContent), 0644)
}

func initChromeContext() (context.Context, context.CancelFunc) {
	parentContext, parentContextCancel := chromedp.NewExecAllocator(context.Background(), chromedp.ExecPath("/var/lib/flatpak/exports/bin/com.google.Chrome"))
	defer parentContextCancel()
	return chromedp.NewContext(parentContext)
}

func findPlaces(title string) []chromedp.Action {
	return chromedp.Tasks{
		chromedp.Navigate("https://2gis.kz/astana"),
		chromedp.WaitVisible("input._cu5ae4"),
		chromedp.SendKeys("input._cu5ae4", title),
		chromedp.Sleep(5 * time.Second),
		chromedp.KeyEvent("\r"),
		chromedp.Sleep(5 * time.Second),
	}
}

func saveHtml(ctx context.Context, html string) error {
	return chromedp.Run(ctx,
		chromedp.OuterHTML("html", &html),
	)
}
