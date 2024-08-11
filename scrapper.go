// package main

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"golang.org/x/net/html"
// )

// // PageData holds the title and content of the page
// type PageData struct {
// 	URL     string
// 	Title   string
// 	Content string
// }

// // Scrape function retrieves the title and first 100 words of the content
// func Scrape(url string, ch chan PageData, timeout time.Duration) {
// 	client := http.Client{
// 		Timeout: timeout,
// 	}

// 	resp, err := client.Get(url)
// 	if err != nil {
// 		// Check if the error is a timeout error
// 		if errors.Is(err, context.DeadlineExceeded) {
// 			ch <- PageData{URL: url, Title: "Timeout", Content: "The request timed out."}
// 		} else if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
// 			ch <- PageData{URL: url, Title: "Timeout", Content: "The request timed out."}
// 		} else {
// 			ch <- PageData{URL: url, Title: "Error", Content: err.Error()}
// 		}
// 		return
// 	}
// 	defer resp.Body.Close()

// 	doc, err := html.Parse(resp.Body)
// 	if err != nil {
// 		ch <- PageData{URL: url, Title: "Error", Content: err.Error()}
// 		return
// 	}

// 	title, content := extractTitleAndContent(doc)

// 	if len(content) == 0 {
// 		content = "No content found"
// 	}

// 	ch <- PageData{URL: url, Title: title, Content: strings.TrimSpace(content)}
// }

// // extractTitleAndContent extracts the title and the first 100 words of readable text from the HTML document.
// func extractTitleAndContent(n *html.Node) (string, string) {
// 	var title, content string
// 	var wordCount int
// 	var f func(*html.Node)

// 	f = func(n *html.Node) {
// 		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
// 			title = n.FirstChild.Data
// 		}

// 		if n.Type == html.TextNode && n.Parent != nil && n.Parent.Data != "script" && n.Parent.Data != "style" {
// 			words := strings.Fields(n.Data)
// 			for _, word := range words {
// 				if wordCount < 100 {
// 					content += word + " "
// 					wordCount++
// 				} else {
// 					break
// 				}
// 			}
// 		}

// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			f(c)
// 		}
// 	}

// 	f(n)
// 	return title, strings.TrimSpace(content)
// }

// func main() {
// 	urls := []string{
// 		"https://www.wagmitech.co",
// 		"https://www.amazon.com",
// 		"https://go.dev",
// 		"https://httpstat.us/200?sleep=5000", // Slow response
// 		"https://example.com:81",             // Unresponsive server
// 	}
// 	timeout := 3000 * time.Millisecond // Example: 3 seconds
// 	ch := make(chan PageData, len(urls))

// 	for _, url := range urls {
// 		go Scrape(url, ch, timeout)
// 	}

// 	for range urls {
// 		data := <-ch
// 		fmt.Printf("Page: %s\nTitle: %s\nContent: %s\n\n", data.URL, data.Title, data.Content)
// 	}

// 	close(ch)
// }

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// PageData holds the title and content of the page
type PageData struct {
	URL     string
	Title   string
	Content string
}

// Scrape function retrieves the title and first 100 words of the content
func Scrape(url string, ch chan PageData, timeout time.Duration) {
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		// Check if the error is a timeout error
		if errors.Is(err, context.DeadlineExceeded) {
			ch <- PageData{URL: url, Title: "Timeout", Content: "The request timed out."}
		} else if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
			ch <- PageData{URL: url, Title: "Timeout", Content: "The request timed out."}
		} else {
			ch <- PageData{URL: url, Title: "Error", Content: err.Error()}
		}
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		ch <- PageData{URL: url, Title: "Error", Content: err.Error()}
		return
	}

	title, content := extractTitleAndContent(doc)

	if len(strings.TrimSpace(content)) == 0 {
		content = "No content found"
	}

	ch <- PageData{URL: url, Title: title, Content: strings.TrimSpace(content)}
}

// extractTitleAndContent extracts the title and the first 100 words of readable text from the HTML document.
func extractTitleAndContent(n *html.Node) (string, string) {
	var title, content string
	var wordCount int
	var inBody bool
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			if n.Data == "body" {
				inBody = true
			}
		}

		if n.Type == html.TextNode && inBody && n.Parent != nil && n.Parent.Data != "script" && n.Parent.Data != "style" {
			words := strings.Fields(n.Data)
			for _, word := range words {
				if wordCount < 100 {
					content += word + " "
					wordCount++
				} else {
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n.Type == html.ElementNode && n.Data == "body" {
			inBody = false
		}
	}

	f(n)
	return title, strings.TrimSpace(content)
}

func main() {
	urls := []string{
		"https://www.wagmitech.co",
		"https://www.amazon.com",
		"https://go.dev",
		"https://httpstat.us/200?sleep=5000", // Slow response
		"https://example.com:81",             // Unresponsive server
	}
	timeout := 3000 * time.Millisecond // Example: 3 seconds
	ch := make(chan PageData, len(urls))

	for _, url := range urls {
		go Scrape(url, ch, timeout)
	}

	for range urls {
		data := <-ch
		fmt.Printf("Page: %s\nTitle: %s\nContent: %s\n\n", data.URL, data.Title, data.Content)
	}

	close(ch)
}
