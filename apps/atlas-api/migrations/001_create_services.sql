CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    version TEXT NOT NULL CHECK (length(trim(version)) > 0),
    status TEXT NOT NULL CHECK (length(trim(status)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
