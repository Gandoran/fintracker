package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type FeedLink struct {
	Title string
	URL   string
}

func (f *Fetcher) DiscoverRSS(ctx context.Context, baseURL string) ([]FeedLink, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	addHeader(req)
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to access the page %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return f.findLinksInHTML(string(body), baseURL), nil
}

func (f *Fetcher) findLinksInHTML(html, baseStr string) []FeedLink {
	var found []FeedLink
	re := regexp.MustCompile(`(?i)<link[^>]+(?:application\/rss\+xml|application\/atom\+xml)[^>]*>`)
	matches := re.FindAllString(html, -1)
	base, _ := url.Parse(baseStr)
	for _, tag := range matches {
		href := extractAttribute(tag, "href")
		title := extractAttribute(tag, "title")
		if href == "" {
			continue
		}
		u, err := url.Parse(href)
		if err == nil {
			absoluteURL := base.ResolveReference(u).String()
			found = append(found, FeedLink{
				Title: title,
				URL:   absoluteURL,
			})
		}
	}
	return found
}

func extractAttribute(tag, attr string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i)%s=["']([^"']+)["']`, attr))
	match := re.FindStringSubmatch(tag)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func addHeader(req *http.Request) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
}
