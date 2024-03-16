-- name: five-latest-posts
SELECT
    title,
    til_id,
    slug,
    created_at at time zone 'utc' at time zone 'Europe/Stockholm' AS created_at
FROM
    post
WHERE
    published = TRUE
ORDER BY
    id DESC
LIMIT 5;

-- name: command-menu-search
SELECT
    title,
    slug
FROM
    post
WHERE
    CASE WHEN $1 <> '"%%"' THEN
        title ILIKE $1
        AND published = TRUE
    ELSE
        published = TRUE
    END
ORDER BY
    id DESC
LIMIT 5;

