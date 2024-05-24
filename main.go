package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// "fyne.io/fyne/v2/canvas",
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly/v2"
)

func loadImageFromURL(imageURL string) (image.Image, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func performSearch(searchQuery string, linksContainer *fyne.Container) {
	searchQuery = strings.TrimSpace(searchQuery)
	linksContainer.Objects = nil

	var movieURLs []string

	c := colly.NewCollector(
		colly.AllowedDomains("filmweb.pl", "www.filmweb.pl"),
	)

	searchURL := fmt.Sprintf("https://www.filmweb.pl/films/search?q=%s", url.QueryEscape(searchQuery))

	c.OnHTML(".previewFilm", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a[href]", "href")
		if link != "" {
			movieURL := "https://www.filmweb.pl" + link
			movieURLs = append(movieURLs, movieURL)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		for _, urlString := range movieURLs {
			parsedURL, err := url.Parse(urlString)
			if err != nil {
				log.Printf("Error parsing URL '%s': %v", urlString, err)
				continue
			}
			link := widget.NewHyperlink(parsedURL.String(), parsedURL)
			linksContainer.Add(link)
		}
		linksContainer.Refresh()
	})

	err := c.Visit(searchURL)
	if err != nil {
		log.Fatal(err)
	}
}

func makeUI(app fyne.App) *fyne.Container {
	in := widget.NewEntry()
	out := widget.NewLabel("Enter a movie title...")
	button := widget.NewButton("Search", nil)
	linksContainer := container.NewVBox()

	button.OnTapped = func() {
		performSearch(in.Text, linksContainer)
	}

	in.OnSubmitted = func(content string) {
		performSearch(content, linksContainer)
	}

	in.OnChanged = func(content string) {
		out.SetText("Film: " + content)
	}

	return container.NewVBox(out, in, button, linksContainer)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.SetContent(makeUI(myApp))
	myWindow.ShowAndRun()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Exited")
}
