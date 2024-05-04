-- name: get-books
SELECT
    *
FROM
    public.book
WHERE
    finished_at IS NOT NULL
ORDER BY
    started_at DESC;

-- name: currently-reading
SELECT
    *
FROM
    public.book
WHERE
    finished_at IS NULL
ORDER BY
    started_at DESC;

