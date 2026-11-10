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
-- $1 and $2 are Umami paths and their views, $3 is when Umami took over
WITH umami_paths AS (
    SELECT
        REGEXP_REPLACE(path, '^/(posts|til)/', '') AS slug,
        views
    FROM
        UNNEST($1::text[], $2::int[]) AS u (path, views)
),
-- Match paths to posts like post-by-slug, preferring a slug match
umami_views AS (
    SELECT
        p.id AS post_id,
        SUM(up.views) AS views
    FROM
        umami_paths AS up
        CROSS JOIN LATERAL (
            SELECT
                id
            FROM
                post
            WHERE
                slug = up.slug
                OR long_slug = up.slug
                OR til_id = SUBSTRING(up.slug FROM '(?:^|-)([0-9]{1,9})$')::int
            ORDER BY
                (slug = up.slug OR long_slug = up.slug) IS TRUE DESC
            LIMIT 1) AS p
    GROUP BY
        p.id
),
frozen_views AS (
    SELECT
        post_id,
        COUNT(*) AS views
    FROM
        post_view
    WHERE
        is_bot = FALSE
        AND created_at < $3::timestamptz AT TIME ZONE 'UTC'
    GROUP BY
        post_id
)
SELECT
    p.title,
    p.til_id,
    p.slug,
    (COALESCE(fv.views, 0) + COALESCE(uv.views, 0))::int AS views
FROM
    post AS p
    LEFT JOIN frozen_views AS fv ON fv.post_id = p.id
    LEFT JOIN umami_views AS uv ON uv.post_id = p.id
WHERE
    p.published = TRUE
    AND (fv.views IS NOT NULL
        OR uv.views IS NOT NULL)
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

-- name: post-views
-- Same as posts-views for a single post, $4 is the post id
WITH umami_paths AS (
    SELECT
        REGEXP_REPLACE(path, '^/(posts|til)/', '') AS slug,
        views
    FROM
        UNNEST($1::text[], $2::int[]) AS u (path, views)
),
umami_views AS (
    SELECT
        COALESCE(SUM(up.views), 0) AS views
    FROM
        umami_paths AS up
        CROSS JOIN LATERAL (
            SELECT
                id
            FROM
                post
            WHERE
                slug = up.slug
                OR long_slug = up.slug
                OR til_id = SUBSTRING(up.slug FROM '(?:^|-)([0-9]{1,9})$')::int
            ORDER BY
                (slug = up.slug OR long_slug = up.slug) IS TRUE DESC
            LIMIT 1) AS p
    WHERE
        p.id = $4
),
frozen_views AS (
    SELECT
        COUNT(*) AS views
    FROM
        post_view
    WHERE
        post_id = $4
        AND is_bot = FALSE
        AND created_at < $3::timestamptz AT TIME ZONE 'UTC'
)
SELECT
    (frozen_views.views + umami_views.views)::int
FROM
    frozen_views,
    umami_views;

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

