CREATE TABLE IF NOT EXISTS aiactionitems (
    id VARCHAR(26) PRIMARY KEY,
    channelid VARCHAR(26) NOT NULL,
    postid VARCHAR(26),
    userid VARCHAR(26) NOT NULL,
    assigneeid VARCHAR(26),
    description TEXT NOT NULL,
    deadline BIGINT,
    status VARCHAR(64) DEFAULT 'pending',
    remindedsent BOOLEAN DEFAULT FALSE,
    createdat BIGINT NOT NULL,
    updatedat BIGINT NOT NULL,
    deletedat BIGINT DEFAULT 0
);

CREATE INDEX idx_aiactionitems_channelid ON aiactionitems(channelid);
CREATE INDEX idx_aiactionitems_userid ON aiactionitems(userid);
CREATE INDEX idx_aiactionitems_assigneeid ON aiactionitems(assigneeid);
CREATE INDEX idx_aiactionitems_status ON aiactionitems(status);
CREATE INDEX idx_aiactionitems_deadline ON aiactionitems(deadline) WHERE deadline IS NOT NULL;
CREATE INDEX idx_aiactionitems_deletedat ON aiactionitems(deletedat);

