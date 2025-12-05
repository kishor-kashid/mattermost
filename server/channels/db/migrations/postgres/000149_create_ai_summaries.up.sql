CREATE TABLE IF NOT EXISTS aisummaries (
    id VARCHAR(26) PRIMARY KEY,
    channelid VARCHAR(26) NOT NULL,
    postid VARCHAR(26),
    summarytype VARCHAR(64) NOT NULL,
    summary TEXT NOT NULL,
    messagecount INT NOT NULL,
    starttime BIGINT NOT NULL,
    endtime BIGINT NOT NULL,
    userid VARCHAR(26) NOT NULL,
    participants TEXT,
    cachekey VARCHAR(255),
    channelname VARCHAR(255),
    createdat BIGINT NOT NULL,
    expiresat BIGINT NOT NULL
);

CREATE INDEX idx_aisummaries_channelid ON aisummaries(channelid);
CREATE INDEX idx_aisummaries_postid ON aisummaries(postid) WHERE postid IS NOT NULL;
CREATE INDEX idx_aisummaries_expiresat ON aisummaries(expiresat);
CREATE INDEX idx_aisummaries_createdat ON aisummaries(createdat);
CREATE INDEX idx_aisummaries_cachekey ON aisummaries(cachekey) WHERE cachekey IS NOT NULL;

