-- name: umami-post-views
SELECT
    RTRIM(url_path, '/') AS path,
    COUNT(*) AS views
FROM
    website_event
WHERE
    website_id = $1
    AND event_type = 1 -- Page view
    AND created_at >= $2
    AND (url_path LIKE '/posts/%'
        OR url_path LIKE '/til/%')
GROUP BY
    1;
