CREATE TABLE IF NOT EXISTS aipreferences (
    id VARCHAR(26) PRIMARY KEY,
    userid VARCHAR(26) NOT NULL UNIQUE,
    enablesummarization BOOLEAN DEFAULT TRUE,
    enableanalytics BOOLEAN DEFAULT TRUE,
    enableactionitems BOOLEAN DEFAULT TRUE,
    enableformatting BOOLEAN DEFAULT TRUE,
    defaultmodel VARCHAR(64) DEFAULT 'gpt-3.5-turbo',
    formattingprofile VARCHAR(64) DEFAULT 'professional',
    createdat BIGINT NOT NULL,
    updatedat BIGINT NOT NULL
);

CREATE INDEX idx_aipreferences_userid ON aipreferences(userid);

