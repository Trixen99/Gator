-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;



-- name: GetFeeds :many
select * from feeds;



-- name: GetFeedByURL :one
select * from feeds where feeds.url = $1;


-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = $2, updated_at = $2
where feeds.id = $1;


-- name: GetNextFeedToFetch :one
select * from feeds
order by last_fetched_at asc nulls first
limit 1;
