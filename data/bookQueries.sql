-- name: get-books
SELECT
    b.*,
    array_agg(f.format_name) AS book_format
FROM
    public.book AS b
    INNER JOIN book_format AS bf ON bf.book_id = b.id
    INNER JOIN FORMAT AS f ON f.id = bf.format_id
WHERE
    finished_at IS NOT NULL
    AND started_at IS NOT NULL
GROUP BY
    b.id
ORDER BY
    finished_at DESC;

-- name: currently-reading
SELECT
    b.*,
    array_agg(f.format_name) AS book_format
FROM
    public.book AS b
    INNER JOIN book_format AS bf ON bf.book_id = b.id
    INNER JOIN FORMAT AS f ON f.id = bf.format_id
WHERE
    finished_at IS NULL
    AND started_at IS NOT NULL
GROUP BY
    b.id
ORDER BY
    started_at DESC;
GROUP BY
    b.id
ORDER BY
    started_at DESC;

