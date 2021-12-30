CREATE DATABASE ledger ENCODING = UTF8;
USE ledger;
CREATE TABLE IF NOT EXISTS ledger (
    t TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP PRIMARY KEY,
    prev UUID NOT NULL,
    id UUID NOT NULL,
    idsubject UUID NOT NULL,
    subject_type string NOT NULL,
    content JSONB NOT NULL
);
CREATE TABLE IF NOT EXISTS snapshot_users (
    id UUID NOT NULL,
    email STRING NOT NULL,
    preferred_name STRING NOT NULL,
    public_key STRING NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS snapshot_posts (
    id UUID NOT NULL,
    idowner UUID NOT NULL,
    post STRING NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)