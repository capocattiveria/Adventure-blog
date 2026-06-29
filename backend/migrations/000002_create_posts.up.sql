CREATE TABLE posts (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id     UUID        NOT NULL REFERENCES users(id),
    title         TEXT        NOT NULL,
    slug          TEXT        NOT NULL UNIQUE,
    description   TEXT,
    content       TEXT        NOT NULL DEFAULT '',
    thumbnail_url TEXT,
    published_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- velocizza la ricerca per slug (usato nelle URL) e il filtro published
CREATE INDEX posts_slug_idx        ON posts (slug);
CREATE INDEX posts_published_at_idx ON posts (published_at) WHERE published_at IS NOT NULL;

-- aggiorna updated_at automaticamente ad ogni UPDATE
CREATE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_set_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
