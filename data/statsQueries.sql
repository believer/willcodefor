-- name: stats-post-views
WITH days AS (
    SELECT
        GENERATE_SERIES((
            SELECT
                created_at
            FROM post
            WHERE
                id = $1), CURRENT_DATE + INTERVAL '1 day', '1 day'::interval)::date AS DAY
)
SELECT
    days.day AS date,
    TO_CHAR(days.day, 'Mon DD YY') AS label,
    COUNT(pv.id)::int AS count
FROM
    days
    LEFT JOIN post_view AS pv ON DATE_TRUNC('day', created_at) = days.day
        AND pv.is_bot = FALSE
        AND post_id = $1
GROUP BY
    1
ORDER BY
    1 ASC;

-- name: stats-post-biggest-day
SELECT
    DATE(created_at) AS DATE,
    COUNT(*) AS COUNT
FROM
    post_view
WHERE
    post_id = $1
    AND is_bot = FALSE
GROUP BY
    DATE(created_at)
ORDER BY
    COUNT DESC
LIMIT 1;

-- name: stats-post-total-views
SELECT
    COUNT(*)
FROM
    post_view
WHERE
    post_id = $1
    AND is_bot = FALSE;

-- name: stats-views-for-period
SELECT
    COUNT(*)
FROM
    post_view
WHERE
    is_bot = FALSE
    AND created_at >= $1;

-- name: stats-bots
SELECT
    COUNT(*)
FROM
    post_view
WHERE
    is_bot = TRUE;

-- name: stats-views-per-week
SELECT
    COUNT(*)
FROM
    post_view
WHERE
    is_bot = FALSE
    AND date_trunc('week', created_at) = date_trunc('week', now());

-- name: stats-most-viewed-posts
SELECT
    COUNT(*) AS views,
    p.title,
    p.slug,
    p.created_at,
    p.id,
    p.updated_at,
    p.til_id
FROM
    post_view AS pv
    INNER JOIN post AS p ON p.id = pv.post_id
WHERE
    pv.is_bot = FALSE
GROUP BY
    p.id
ORDER BY
    views DESC
LIMIT 10;

-- name: stats-most-viewed-posts-today
SELECT
    COUNT(*) AS views,
    p.title,
    p.slug,
    p.created_at,
    p.id,
    p.updated_at,
    p.til_id
FROM
    post_view AS pv
    INNER JOIN post AS p ON p.id = pv.post_id
WHERE
    pv.is_bot = FALSE
    AND pv.created_at >= CURRENT_DATE
GROUP BY
    p.id
ORDER BY
    views DESC;

-- name: stats-os
SELECT
    os_name,
    COUNT(*) AS count,
    TO_CHAR(COUNT(*) / SUM(COUNT(*)) OVER () * 100, 'fm99%') AS percent
FROM
    post_view
WHERE
    is_bot = FALSE
    AND created_at >= $1
GROUP BY
    os_name
ORDER BY
    count DESC
LIMIT 5;

-- name: stats-browsers
SELECT
    browser_name,
    COUNT(*) AS count,
    TO_CHAR(COUNT(*) / SUM(COUNT(*)) OVER () * 100, 'fm99%') AS percent
FROM
    post_view
WHERE
    is_bot = FALSE
    AND created_at >= $1
GROUP BY
    browser_name
ORDER BY
    count DESC
LIMIT 5;

-- name: stats-chart-today
WITH days AS (
    SELECT
        generate_series(CURRENT_DATE, CURRENT_DATE + '1 day'::interval, '1 hour') AS hour
)
SELECT
    days.hour AS date,
    to_char(days.hour, 'HH24:MI') AS label,
    count(pv.id)::int AS count
FROM
    days
    LEFT JOIN post_view AS pv ON DATE_TRUNC('hour', created_at at time zone 'utc' at time zone 'Europe/Stockholm') = days.hour
        AND pv.is_bot = FALSE
    LEFT JOIN post AS p ON p.id = pv.post_id
GROUP BY
    1
ORDER BY
    1 ASC;

-- name: stats-chart-week
WITH days AS (
    SELECT
        generate_series(date_trunc('week', CURRENT_DATE), date_trunc('week', CURRENT_DATE) + '6 days'::interval, '1 day')::date AS day
)
SELECT
    days.day AS date,
    to_char(days.day, 'Mon DD') AS label,
    count(pv.id)::int AS count
FROM
    days
    LEFT JOIN post_view AS pv ON DATE_TRUNC('day', created_at) = days.day
        AND pv.is_bot = FALSE
GROUP BY
    1
ORDER BY
    1 ASC;

-- name: stats-chart-thirty-days
WITH days AS (
    SELECT
        generate_series(CURRENT_DATE - '30 days'::interval, CURRENT_DATE, '1 day')::date AS day
)
SELECT
    days.day AS date,
    to_char(days.day, 'Mon DD') AS label,
    count(pv.id)::int AS count
FROM
    days
    LEFT JOIN post_view AS pv ON DATE_TRUNC('day', created_at) = days.day
        AND pv.is_bot = FALSE
GROUP BY
    1
ORDER BY
    1 ASC;

-- name: stats-chart-this-year
WITH months AS (
    SELECT
        (DATE_TRUNC('year', NOW()) + (INTERVAL '1' MONTH * GENERATE_SERIES(0, 11)))::date AS MONTH
)
SELECT
    months.month AS date,
    to_char(months.month, 'Mon') AS label,
    COUNT(pv.id)::int AS count
FROM
    months
    LEFT JOIN post_view AS pv ON DATE_TRUNC('month', created_at) = months.month
        AND pv.is_bot = FALSE
GROUP BY
    1
ORDER BY
    1 ASC;

-- name: stats-chart-all-time
WITH data AS (
    SELECT
        date_trunc('month', created_at) AS month,
        count(1)::int
    FROM
        post_view
    WHERE
        is_bot = FALSE
    GROUP BY
        1
)
SELECT
    month::date AS date,
    to_char(month, 'Mon YY') AS label,
    sum(count) OVER (ORDER BY month ASC ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)::int AS count
FROM
    data;

-- name: stats-chart-posts-per-month
WITH months AS (
    SELECT
        GENERATE_SERIES('2020-01-01', CURRENT_DATE, '1 month') AS MONTH
)
SELECT
    months.month AS date,
    TO_CHAR(months.month, 'Mon YY') AS label,
    COUNT(p.id) AS count
FROM
    months
    LEFT JOIN post AS p ON DATE_TRUNC('month', p.created_at) = months.month
WHERE
    p.published = TRUE
GROUP BY
    1
ORDER BY
    1;

