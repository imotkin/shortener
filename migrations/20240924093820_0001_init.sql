-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    original TEXT NOT NULL,
    shortened TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    views INTEGER NOT NULL DEFAULT(0)
);

CREATE TABLE IF NOT EXISTS stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	link_id INTEGER NOT NULL,
	ip TEXT,
	visit_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	country TEXT,
	region TEXT,
	city TEXT,
	FOREIGN KEY (link_id) REFERENCES links(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links, stats;
-- +goose StatementEnd
