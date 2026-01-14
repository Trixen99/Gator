-- name: CreatePost :one
INSERT into posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;




-- name: GetPosts :many
select * from posts
where feed_id in (
    select feed_id from feed_follows
    where user_id = $1
    )
order by published_at asc
limit $2;