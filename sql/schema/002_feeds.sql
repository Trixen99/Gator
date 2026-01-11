-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT not NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL, 
    foreign key (user_id) references users(id)
    on delete cascade
);





-- +goose Down
DROP TABLE feeds;