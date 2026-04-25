package scraper

import "strings"

func (f *Fetcher) checkKnownHosts(targetURL string) []FeedLink {
	lowerURL := strings.ToLower(targetURL)
	if strings.Contains(lowerURL, "ilsole24ore.com") || strings.Contains(lowerURL, "sole24") {
		return []FeedLink{
			{Title: "Il Sole 24 Ore - Finanza", URL: "https://www.ilsole24ore.com/rss/finanza.xml"},
			{Title: "Il Sole 24 Ore - Economia", URL: "https://www.ilsole24ore.com/rss/economia.xml"},
			{Title: "Il Sole 24 Ore - Tecnologia", URL: "https://www.ilsole24ore.com/rss/tecnologia.xml"},
		}
	}
	if strings.Contains(lowerURL, "wallstreetitalia.com") {
		return []FeedLink{
			{Title: "Wall Street Italia - Tutte le News", URL: "https://www.wallstreetitalia.com/feed/"},
		}
	}
	if strings.Contains(lowerURL, "finanzaonline.com") {
		return []FeedLink{
			{Title: "FinanzaOnline - Mercati", URL: "https://www.finanzaonline.com/feed"},
		}
	}
	if strings.Contains(lowerURL, "yahoo.com") || strings.Contains(lowerURL, "finance.yahoo") {
		return []FeedLink{
			{Title: "Yahoo Finance (Big Tech)", URL: "https://query2.finance.yahoo.com/v1/base/rss?s=AAPL,MSFT,TSLA,NVDA,AMZN,META,GOOG"},
			{Title: "Yahoo Finance (Indici Mercato)", URL: "https://query2.finance.yahoo.com/v1/base/rss?s=^GSPC,^DJI,^IXIC"},
			{Title: "Yahoo Finance (Crypto Top 3)", URL: "https://query2.finance.yahoo.com/v1/base/rss?s=BTC-USD,ETH-USD,SOL-USD"},
		}
	}
	if strings.Contains(lowerURL, "cnbc.com") {
		return []FeedLink{
			{Title: "CNBC - Top News", URL: "https://search.cnbc.com/rs/search/combinedcms/view.xml?partnerId=wrss01&id=100003114"},
			{Title: "CNBC - Finance", URL: "https://search.cnbc.com/rs/search/combinedcms/view.xml?partnerId=wrss01&id=10000664"},
			{Title: "CNBC - Tech", URL: "https://search.cnbc.com/rs/search/combinedcms/view.xml?partnerId=wrss01&id=19854910"},
		}
	}
	if strings.Contains(lowerURL, "wsj.com") {
		return []FeedLink{
			{Title: "Wall Street Journal - Business", URL: "https://feeds.a.dj.com/rss/WSJcomUSBusiness.xml"},
			{Title: "Wall Street Journal - Markets", URL: "https://feeds.a.dj.com/rss/RSSMarketsMain.xml"},
		}
	}
	if strings.Contains(lowerURL, "ansa.it") {
		return []FeedLink{
			{Title: "ANSA - Economia", URL: "https://www.ansa.it/sito/notizie/economia/economia_rss.xml"},
			{Title: "ANSA - Ultima Ora", URL: "https://www.ansa.it/sito/ansait_rss.xml"},
			{Title: "ANSA - Tecnologia", URL: "https://www.ansa.it/sito/notizie/tecnologia/tecnologia_rss.xml"},
		}
	}
	if strings.Contains(lowerURL, "corriere.it") {
		return []FeedLink{
			{Title: "Corriere della Sera - Economia", URL: "https://xml2.corriereobjects.it/rss/economia.xml"},
			{Title: "Corriere della Sera - Prima Pagina", URL: "https://xml2.corriereobjects.it/rss/homepage.xml"},
		}
	}
	if strings.Contains(lowerURL, "repubblica.it") {
		return []FeedLink{
			{Title: "La Repubblica - Economia", URL: "https://www.repubblica.it/rss/economia/rss2.0.xml"},
			{Title: "La Repubblica - Prima Pagina", URL: "https://www.repubblica.it/rss/homepage/rss2.0.xml"},
		}
	}
	if strings.Contains(lowerURL, "coindesk.com") {
		return []FeedLink{
			{Title: "CoinDesk - Tutte le News (Crypto)", URL: "https://www.coindesk.com/arc/outboundfeeds/rss/"},
		}
	}
	if strings.Contains(lowerURL, "cointelegraph.com") {
		return []FeedLink{
			{Title: "Cointelegraph - Ultima Ora (Crypto)", URL: "https://cointelegraph.com/rss"},
		}
	}
	if strings.Contains(lowerURL, "techcrunch.com") {
		return []FeedLink{
			{Title: "TechCrunch - Startups & VC", URL: "https://techcrunch.com/feed/"},
		}
	}
	if strings.Contains(lowerURL, "theverge.com") {
		return []FeedLink{
			{Title: "The Verge - Tech News", URL: "https://www.theverge.com/rss/index.xml"},
		}
	}
	if strings.Contains(lowerURL, "wired.it") {
		return []FeedLink{
			{Title: "Wired Italia", URL: "https://www.wired.it/feed/rss"},
		}
	}
	return nil
}
