// Command visible is a chromedp example demonstrating how to wait until an
// element is visible.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	parentContext, parentContextCancel := chromedp.NewExecAllocator(context.Background(), chromedp.ExecPath("/var/lib/flatpak/exports/bin/com.google.Chrome"))
	defer parentContextCancel()

	ctx, cancel := chromedp.NewContext(parentContext)
	defer cancel()

	// TODO: fix this, sometimes pressing enter does not work
	err := chromedp.Run(ctx, findPlaces("Барбершоп")...)
	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(ctx, acceptCookies()...)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: stop going next page once its not there
	// numOfPages := getNumOfPages()
	i := 0

	for {
		err = chromedp.Run(ctx, scrollUntilVisible()...)
		if err != nil {
			log.Fatal(err)
		}

		err = chromedp.Run(ctx, savePageHTML(fmt.Sprintf("html/output_%d.html", i)))
		if err != nil {
			log.Fatal(err)
		}
		i += 1

		err = chromedp.Run(ctx, goNextPage()...)
		if err != nil {
			log.Fatal(err)
		}

		// if numOfPages == getNumOfPages() {
		// 	break
		// }
	}
}

func acceptCookies() []chromedp.Action {
	var nodes []*cdp.Node
	return []chromedp.Action{
		chromedp.Nodes("._n1367pl ._13xlah4 svg", &nodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if len(nodes) > 0 {
				chromedp.Click(`._n1367pl ._13xlah4 svg`, chromedp.ByQuery).Do(ctx)
			}
			return nil
		}),
		chromedp.Sleep(3 * time.Second),
	}
}

func savePageHTML(filename string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var htmlContent string

		err := chromedp.OuterHTML("html", &htmlContent).Do(ctx)
		if err != nil {
			return err
		}

		return os.WriteFile(filename, []byte(htmlContent), 0644)
	})
}

func getNumOfPages() int {
	return 2 // mocked for now
}

func goNextPage() []chromedp.Action {
	return []chromedp.Action{
		chromedp.Click("._5ocwns ._n5hmn94:nth-child(2) svg"),
		chromedp.Sleep(3 * time.Second),
	}
}

func scrollUntilVisible() []chromedp.Action {
	return []chromedp.Action{
		chromedp.ScrollIntoView("._1x4k6z7 div"),
		chromedp.WaitVisible("._1x4k6z7 div"),
		chromedp.Sleep(3 * time.Second),
	}
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
