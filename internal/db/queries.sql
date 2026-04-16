-- name: CreateArticle :one
INSERT INTO articles (title, link, content, source, published_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: CreateAnalysis :one
INSERT INTO analyses (article_id, summary, sentiment, impact, tickers)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRecentAnalyses :many
SELECT 
    sqlc.embed(analyses), 
    sqlc.embed(articles)
FROM analyses
JOIN articles ON analyses.article_id = articles.id
ORDER BY analyses.analyzed_at DESC
LIMIT ?;