-- name: CreateFeedFollow :one
with insertedFeedFollow as (INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
) select insertedFeedFollow.*, feeds.name, users.name
from insertedFeedFollow
inner join feeds
on feeds.id = insertedFeedFollow.feed_id
inner join users
on users.id = insertedFeedFollow.user_id;




-- name: GetFeedFollowsForUser :many
select feed_follows.*, feeds.name as feed_name, users.name as user_name 
from feed_follows 
inner join feeds on feeds.id = feed_follows.feed_id 
inner join users on users.id = feed_follows.user_id
where feed_follows.user_id = $1;


-- name: DeleteFeedFollow :exec
delete from feed_follows where user_id = $1 and feed_id = $2;