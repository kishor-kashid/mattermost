-- Fix AI Action Items Table Schema
-- Run this in your PostgreSQL database

-- Drop the old table and indexes
DROP INDEX IF EXISTS idx_aiactionitems_deletedat;
DROP INDEX IF EXISTS idx_aiactionitems_duedate;
DROP INDEX IF EXISTS idx_aiactionitems_deadline;
DROP INDEX IF EXISTS idx_aiactionitems_priority;
DROP INDEX IF EXISTS idx_aiactionitems_status;
DROP INDEX IF EXISTS idx_aiactionitems_assigneeid;
DROP INDEX IF EXISTS idx_aiactionitems_createdby;
DROP INDEX IF EXISTS idx_aiactionitems_userid;
DROP INDEX IF EXISTS idx_aiactionitems_channelid;
DROP TABLE IF EXISTS aiactionitems;

-- Recreate with correct schema
CREATE TABLE aiactionitems (
    id VARCHAR(26) PRIMARY KEY,
    channelid VARCHAR(26) NOT NULL,
    postid VARCHAR(26),
    createdby VARCHAR(26) NOT NULL,
    assigneeid VARCHAR(26),
    description TEXT NOT NULL,
    duedate BIGINT,
    priority VARCHAR(64) DEFAULT 'medium',
    status VARCHAR(64) DEFAULT 'open',
    completedat BIGINT DEFAULT 0,
    createdat BIGINT NOT NULL,
    updatedat BIGINT NOT NULL,
    deletedat BIGINT DEFAULT 0
);

-- Create indexes
CREATE INDEX idx_aiactionitems_channelid ON aiactionitems(channelid);
CREATE INDEX idx_aiactionitems_createdby ON aiactionitems(createdby);
CREATE INDEX idx_aiactionitems_assigneeid ON aiactionitems(assigneeid);
CREATE INDEX idx_aiactionitems_status ON aiactionitems(status);
CREATE INDEX idx_aiactionitems_priority ON aiactionitems(priority);
CREATE INDEX idx_aiactionitems_duedate ON aiactionitems(duedate) WHERE duedate IS NOT NULL;
CREATE INDEX idx_aiactionitems_deletedat ON aiactionitems(deletedat);

