-- name: get-books
SELECT
    *,
    COALESCE(finished_at, started_at + INTERVAL '1 day' * ((word_count / ((word_count / page_count) * current_page / GREATEST (EXTRACT(DAY FROM COALESCE(finished_at, CURRENT_DATE) - started_at)::int, 1))::integer))) - started_at AS "days",
    (word_count / page_count) * current_page / GREATEST (EXTRACT(DAY FROM COALESCE(finished_at, CURRENT_DATE) - started_at)::int, 1) AS "pace"
FROM
    public.book
WHERE
    finished_at IS NOT NULL
ORDER BY
    started_at DESC;

-- name: currently-reading
SELECT
    *,
    started_at + INTERVAL '1 day' * EXTRACT(DAY FROM COALESCE(finished_at, started_at + INTERVAL '1 day' * ((word_count / ((word_count / page_count) * current_page / GREATEST (EXTRACT(DAY FROM COALESCE(finished_at, CURRENT_DATE) - started_at)::int, 1))::integer))) - started_at) AS "finished_at",
    COALESCE(finished_at, started_at + INTERVAL '1 day' * ((word_count / ((word_count / page_count) * current_page / GREATEST (EXTRACT(DAY FROM COALESCE(finished_at, CURRENT_DATE) - started_at)::int, 1))::integer))) - started_at AS "days",
    (word_count / page_count) * current_page / GREATEST (EXTRACT(DAY FROM COALESCE(finished_at, CURRENT_DATE) - started_at)::int, 1) AS "pace"
FROM
    public.book
WHERE
    finished_at IS NULL;

