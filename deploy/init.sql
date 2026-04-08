CREATE TABLE IF NOT EXISTS urls (
    id          BIGSERIAL PRIMARY KEY,
    full_url    TEXT NOT NULL,
    short_code  VARCHAR(32) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS url_visits (
    id          BIGSERIAL PRIMARY KEY,
    url_id      BIGINT NOT NULL,
    visited_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_agent  TEXT,

    CONSTRAINT fk_url_visits_url_id
        FOREIGN KEY (url_id)
        REFERENCES urls(id)
        ON DELETE CASCADE
);