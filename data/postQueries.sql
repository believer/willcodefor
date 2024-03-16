-- name: posts-by-created
SELECT
    title,
    til_id,
    slug,
    created_at at time zone 'utc' at time zone 'Europe/Stockholm' AS created_at,
    updated_at at time zone 'utc' at time zone 'Europe/Stockholm' AS updated_at
FROM
    post
WHERE
    published = TRUE
ORDER BY
    created_at DESC;

-- name: posts-by-updated
SELECT
    title,
    til_id,
    slug,
    created_at at time zone 'utc' at time zone 'Europe/Stockholm' AS created_at,
    updated_at at time zone 'utc' at time zone 'Europe/Stockholm' AS updated_at
FROM
    post
WHERE
    published = TRUE
ORDER BY
    updated_at DESC;

-- name: post-search
SELECT
    title,
    til_id,
    slug,
    created_at at time zone 'utc' at time zone 'Europe/Stockholm' AS created_at,
    updated_at at time zone 'utc' at time zone 'Europe/Stockholm' AS updated_at
FROM
    post
WHERE
    title ILIKE '%' || $1 || '%'
    OR body ILIKE '%' || $1 || '%'
    AND published = TRUE
ORDER BY
    created_at DESC;

-- name: posts-views
SELECT
    p.title,
    p.til_id,
    p.slug,
    COUNT(pv.id) AS views
FROM
    post AS p
    INNER JOIN post_view AS pv ON p.id = pv.post_id
WHERE
    p.published = TRUE
    AND pv.is_bot = FALSE
GROUP BY
    p.id
ORDER BY
    views DESC;

-- name: post-by-slug
SELECT
    title,
    til_id,
    slug,
    id,
    body,
    created_at at time zone 'utc' at time zone 'Europe/Stockholm' AS created_at,
    updated_at at time zone 'utc' at time zone 'Europe/Stockholm' AS updated_at,
    COALESCE(series, '') AS series,
    excerpt
FROM
    post
WHERE
    slug = $1
    OR long_slug = $1
    OR til_id = $2;

-- name: next-post
SELECT
    title,
    slug,
    til_id
FROM
    post
WHERE
    id > $1
    AND published = TRUE
ORDER BY
    id ASC
LIMIT 1;

-- name: previous-post
SELECT
    title,
    slug,
    til_id
FROM
    post
WHERE
    id < $1
    AND published = TRUE
ORDER BY
    id DESC
LIMIT 1;

-- name: insert-view
INSERT INTO post_view (user_agent, post_id, is_bot, browser_name, browser_version, device_type, device_model, device_vendor, os_name, os_version, engine_version, engine_name)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: post-views
SELECT
    COUNT(*)
FROM
    post_view
WHERE
    post_id = $1;

-- name: post-series
SELECT
    slug,
    title
FROM
    post
WHERE
    series = $1
    AND published = TRUE
ORDER BY
    id ASC;

-- name: xml-feed
SELECT
    title,
    slug,
    body,
    updated_at at time zone 'utc' at time zone 'Europe/Stockholm' AS updated_at
FROM
    post
WHERE
    published = TRUE
ORDER BY
    created_at DESC;

-- name: xml-sitemap
SELECT
    slug,
    updated_at
FROM
    post
WHERE
    published = TRUE
ORDER BY
    created_at DESC;

