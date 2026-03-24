-- +goose Up
CREATE TABLE set_categories (
    set_id UUID NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (set_id, category_id)
);

-- +goose Down
DROP TABLE set_categories;
