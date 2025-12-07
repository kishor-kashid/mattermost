-- Create AI link summaries table
CREATE TABLE IF NOT EXISTS AILinkSummaries (
    Id VARCHAR(26) PRIMARY KEY,
    URL TEXT NOT NULL,
    URLHash VARCHAR(64) NOT NULL,
    Title TEXT,
    Description TEXT,
    Summary TEXT NOT NULL,
    KeyPoints TEXT[],
    ContentType VARCHAR(32),
    ReadingTime INTEGER,
    Domain VARCHAR(255),
    FaviconURL TEXT,
    CreateAt BIGINT NOT NULL,
    ExpiresAt BIGINT NOT NULL
);

-- Index for fast lookup by URL hash and expiry
CREATE INDEX IF NOT EXISTS idx_ailinksummaries_urlhash ON AILinkSummaries (URLHash, ExpiresAt DESC);

