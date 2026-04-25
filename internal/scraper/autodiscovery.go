package scraper

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
)

type FeedLink struct {
	Title string
	URL   string
}

func (f *Fetcher) DiscoverRSS(ctx context.Context, baseURL string) ([]FeedLink, error) {
	if knownFeeds := f.checkKnownHosts(baseURL); len(knownFeeds) > 0 {
		return knownFeeds, nil
	}
	html, status := f.fetchHTMLFast(ctx, baseURL)
	if status == 200 && html != "" {
		links := f.findLinksInHTML(html, baseURL)
		if len(links) > 0 {
			return links, nil
		}
	}
	if status == 403 || status == 429 || status == 401 || status == 503 || len(html) < 1000 {
		html = f.fetchHTMLSlow(ctx, baseURL)
		if html != "" {
			links := f.findLinksInHTML(html, baseURL)
			if len(links) > 0 {
				return links, nil
			}
		}
	}
	return []FeedLink{}, nil
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
			found = append(found, FeedLink{Title: title, URL: absoluteURL})
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
