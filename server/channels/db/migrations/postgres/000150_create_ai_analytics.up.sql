CREATE TABLE IF NOT EXISTS aianalytics (
    id VARCHAR(26) PRIMARY KEY,
    channelid VARCHAR(26) NOT NULL,
    date VARCHAR(10) NOT NULL,
    messagecount INT DEFAULT 0,
    usercount INT DEFAULT 0,
    avgresponsetime BIGINT DEFAULT 0,
    topcontributors JSONB,
    hourlydistribution JSONB,
    createdat BIGINT NOT NULL,
    updatedat BIGINT NOT NULL,
    UNIQUE(channelid, date)
);

CREATE INDEX idx_aianalytics_channelid ON aianalytics(channelid);
CREATE INDEX idx_aianalytics_date ON aianalytics(date);
CREATE INDEX idx_aianalytics_channelid_date ON aianalytics(channelid, date);

