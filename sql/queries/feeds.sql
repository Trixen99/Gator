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


-- name: GetFeedByID :one
select * from feeds where feeds.id = $1;