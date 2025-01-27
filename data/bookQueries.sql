-- name: get-books
SELECT
    b.*,
    array_agg(f.format_name) AS book_format,
    l.name AS "language"
FROM
    public.book AS b
    INNER JOIN book_format AS bf ON bf.book_id = b.id
    INNER JOIN "format" AS f ON f.id = bf.format_id
    INNER JOIN book_language AS bl ON bl.book_id = b.id
    INNER JOIN "language" AS l ON l.id = bl.language_id
WHERE
    finished_at IS NOT NULL
    AND started_at IS NOT NULL
    AND date_part('year', finished_at) = $1
GROUP BY
    b.id,
    l.id
ORDER BY
    finished_at DESC;

-- name: currently-reading
SELECT
    b.*,
    array_agg(f.format_name) AS book_format,
    l.name AS "language"
FROM
    public.book AS b
    INNER JOIN book_format AS bf ON bf.book_id = b.id
    INNER JOIN "format" AS f ON f.id = bf.format_id
    INNER JOIN book_language AS bl ON bl.book_id = b.id
    INNER JOIN "language" AS l ON l.id = bl.language_id
WHERE
    finished_at IS NULL
    AND started_at IS NOT NULL
GROUP BY
    b.id,
    l.id
ORDER BY
    started_at DESC;

-- name: next-books
SELECT
    b.*,
    array_agg(f.format_name) AS book_format,
    l.name AS "language"
FROM
    public.book AS b
    INNER JOIN book_format AS bf ON bf.book_id = b.id
    INNER JOIN "format" AS f ON f.id = bf.format_id
    INNER JOIN book_language AS bl ON bl.book_id = b.id
    INNER JOIN "language" AS l ON l.id = bl.language_id
WHERE
    started_at IS NULL
    AND finished_at IS NULL
GROUP BY
    b.id,
    l.id
ORDER BY
    started_at DESC;

