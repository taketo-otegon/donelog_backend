-- +goose Up
-- 基本スキーマ: categories -> tracks -> donelogs
CREATE TABLE categories (
    id              VARCHAR(64) PRIMARY KEY,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT categories_id_slug CHECK (id ~ '^[a-z][a-z0-9_-]*$')
);

CREATE TABLE tracks (
    id                   VARCHAR(64) PRIMARY KEY,
    default_category_id  VARCHAR(64),
    active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tracks_id_slug CHECK (id ~ '^[a-z][a-z0-9_-]*$'),
    CONSTRAINT tracks_default_category_fk FOREIGN KEY (default_category_id) REFERENCES categories(id)
);

CREATE TABLE donelogs (
    id           VARCHAR(26) PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    track_id     VARCHAR(64) NOT NULL,
    category_id  VARCHAR(64) NOT NULL,
    count        INTEGER NOT NULL,
    occurred_on  DATE NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT donelogs_id_ulid CHECK (id ~ '^[0-9A-HJKMNP-TV-Z]{26}$'),
    CONSTRAINT donelogs_title_length CHECK (char_length(title) > 0 AND char_length(title) <= 120),
    CONSTRAINT donelogs_count_positive CHECK (count > 0),
    CONSTRAINT donelogs_track_fk FOREIGN KEY (track_id) REFERENCES tracks(id),
    CONSTRAINT donelogs_category_fk FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE INDEX idx_donelogs_track_id ON donelogs (track_id);
CREATE INDEX idx_donelogs_category_id ON donelogs (category_id);
CREATE INDEX idx_donelogs_occurred_on ON donelogs (occurred_on);

-- +goose Down
DROP TABLE IF EXISTS donelogs;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS categories;
