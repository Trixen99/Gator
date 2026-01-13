-- +goose Up
CREATE table posts(
    id UUID PRIMARY key,
    created_at TIMESTAMP not NULL,
    updated_at TIMESTAMP not null,
    title text not null,
    url text not null unique,
    description text not null,
    published_at TIMESTAMP,
    feed_id UUID not null,
    foreign KEY (feed_id) references feeds(id)
);




-- +goose Down
drop table posts;