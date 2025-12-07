ALTER TABLE AIPreferences
DROP COLUMN IF EXISTS LinkSummaryLength,
DROP COLUMN IF EXISTS DefaultLinkExpanded,
DROP COLUMN IF EXISTS AutoSummarizeLinks,
DROP COLUMN IF EXISTS EnableLinkSummaries;

