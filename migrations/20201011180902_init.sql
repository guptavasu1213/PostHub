-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts(
    post_id INTEGER PRIMARY KEY,
    title  text NOT NULL,
    body text NOT NULL,
    scope text NOT NULL,
    epoch INTEGER NOT NULL
);

CREATE TABLE links(
    link_id text primary key,
    access text NOT NULL,
    post_id INTEGER REFERENCES posts(post_id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);

CREATE TABLE report(
    reason text NOT NULL,
    post_id INTEGER REFERENCES posts(post_id)
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
DROP TABLE links;
DROP TABLE report;
-- +goose StatementEnd
