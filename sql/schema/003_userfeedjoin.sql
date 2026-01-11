-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not NULL,
    updated_at TIMESTAMP not NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (feed_id) references feeds(id) on delete cascade,
    unique(user_id, feed_id)
);






-- +goose Down
DROP TABLE feed_follows;