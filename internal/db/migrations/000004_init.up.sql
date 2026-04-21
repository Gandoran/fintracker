CREATE TABLE sources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    category TEXT DEFAULT 'Generale',
    is_active BOOLEAN NOT NULL DEFAULT 1,
    error_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO sources (name, url) VALUES 
('Yahoo Finance', 'https://feeds.finance.yahoo.com/rss/2.0/headline?s=AAPL,TSLA,MSFT,NVDA'),
('CNBC Markets', 'https://search.cnbc.com/rs/search/combinedcms/view.xml?partnerId=wrss01&id=10000664');