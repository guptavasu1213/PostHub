-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts(
    post_id INTEGER PRIMARY KEY,
    title  text,
    body text,
    scope text,
    epoch INTEGER
);

CREATE TABLE links(
    link_id text primary key,
    access text,
    post_id INTEGER REFERENCES posts(post_id)
);

CREATE TABLE report(
    reason text,
    post_id INTEGER REFERENCES posts(post_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
DROP TABLE links;
DROP TABLE report;
-- +goose StatementEnd
