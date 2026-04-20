-- name: CreateArticle :one
INSERT INTO articles (title, link, content, source, published_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: CreateAnalysis :one
INSERT INTO analyses (article_id, summary, sentiment, impact, tickers, reference_links, reliability_score)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRecentAnalyses :many
SELECT 
    sqlc.embed(analyses), 
    sqlc.embed(articles)
FROM analyses
JOIN articles ON analyses.article_id = articles.id
ORDER BY analyses.analyzed_at DESC
LIMIT ?;

-- name: SearchAnalyses :many
SELECT 
    sqlc.embed(analyses), 
    sqlc.embed(articles)
FROM analyses
JOIN articles ON analyses.article_id = articles.id
WHERE articles.title LIKE '%' || ? || '%' 
   OR analyses.tickers LIKE '%' || ? || '%'
ORDER BY analyses.analyzed_at DESC
LIMIT 30;

-- name: GetArticleByID :one
SELECT * FROM articles WHERE id = ? LIMIT 1;

-- name: GetAnalysesByDate :many
SELECT 
    sqlc.embed(analyses), 
    sqlc.embed(articles)
FROM analyses
JOIN articles ON analyses.article_id = articles.id
WHERE DATE(analyses.analyzed_at) = DATE(?)
ORDER BY analyses.analyzed_at DESC;

-- name: GetNextPendingArticle :one
SELECT * FROM articles 
WHERE status = 'PENDING' 
ORDER BY published_at ASC 
LIMIT 1;

-- name: UpdateArticleStatus :exec
UPDATE articles 
SET status = ? 
WHERE id = ?;