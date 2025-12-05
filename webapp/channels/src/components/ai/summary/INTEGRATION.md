# AI Summary Integration Guide

This document explains how to integrate AI summarization features into the Mattermost UI.

## Components Created

1. **SummaryPanel** - Main RHS panel for displaying summaries
2. **SummaryContent** - Renders the summary text with Markdown
3. **SummaryMetadata** - Displays metadata (message count, participants, etc.)
4. **DateRangeModal** - Modal for selecting date range for channel summaries
5. **ThreadSummarizeMenuItem** - Menu item for thread summarization
6. **ChannelSummarizeMenuItem** - Menu item for channel summarization

## Integration Points

### 1. Post Dropdown Menu (Thread Summarization)

**File**: `webapp/channels/src/components/dot_menu/dot_menu.tsx`

Add the following import:
```typescript
import ThreadSummarizeMenuItem from 'components/ai/summary/thread_summarize_menu_item';
```

Add the menu item in the dropdown menu render:
```typescript
{/* Add after other menu items */}
<ThreadSummarizeMenuItem
    postId={post.id}
    channelId={post.channel_id}
    onClose={this.handleClose}
/>
```

### 2. Channel Header Menu (Channel Summarization)

**File**: `webapp/channels/src/components/channel_header_menu/channel_header_menu.tsx`

Add the following import:
```typescript
import ChannelSummarizeMenuItem from 'components/ai/summary/channel_summarize_menu_item';
```

Add the menu item in the channel header dropdown:
```typescript
{/* Add after other menu items */}
<ChannelSummarizeMenuItem
    channelId={channel.id}
    onClose={this.handleClose}
/>
```

### 3. Right Hand Sidebar (RHS) Panel

**File**: `webapp/channels/src/components/rhs/rhs.tsx` or similar

Import and use SummaryPanel:
```typescript
import SummaryPanel from 'components/ai/summary/summary_panel';
import {useSelector} from 'react-redux';
import {getSummariesState, getLatestChannelSummary} from 'selectors/ai_summarizer';

// In component:
const summariesState = useSelector(getSummariesState);
const latestSummary = useSelector((state) => getLatestChannelSummary(state, channelId));

// Render:
{summariesState.loading || latestSummary ? (
    <SummaryPanel
        summary={latestSummary}
        loading={summariesState.loading}
        error={summariesState.error}
        fromCache={false}
        onClose={() => dispatch(clearSummary())}
        onRegenerate={() => {
            // Re-trigger summarization
        }}
    />
) : null}
```

## Usage Flow

### Thread Summarization
1. User clicks three-dot menu on a post
2. User clicks "Summarize Thread"
3. Action dispatched to Redux
4. API call made to `/api/v4/ai/summarize/thread/{post_id}`
5. Summary displayed in RHS panel

### Channel Summarization
1. User clicks channel header dropdown
2. User clicks "Summarize Channel"
3. Date range modal appears
4. User selects date range and summary level
5. Action dispatched to Redux
6. API call made to `/api/v4/ai/summarize/channel/{channel_id}`
7. Summary displayed in RHS panel

### Slash Command
Users can also use: `/summarize` command
- `/summarize` - Summarizes last 24 hours of current channel
- `/summarize thread` - Summarizes current thread
- `/summarize brief` - Brief summary
- `/summarize detailed` - Detailed summary

## State Management

### Redux State Structure
```typescript
state.ai.summaries = {
    byId: {
        [summaryId]: AISummary,
    },
    byChannel: {
        [channelId]: [summaryId1, summaryId2, ...],
    },
    loading: boolean,
    error: any,
}
```

### Actions
- `AI_SUMMARIZE_THREAD_REQUEST`
- `AI_SUMMARIZE_THREAD_SUCCESS`
- `AI_SUMMARIZE_THREAD_FAILURE`
- `AI_SUMMARIZE_CHANNEL_REQUEST`
- `AI_SUMMARIZE_CHANNEL_SUCCESS`
- `AI_SUMMARIZE_CHANNEL_FAILURE`
- `AI_CLEAR_SUMMARY`

### Selectors
- `getSummariesState(state)`
- `getSummaryById(state, summaryId)`
- `getSummariesByChannel(state, channelId)`
- `getLatestChannelSummary(state, channelId)`
- `getThreadSummary(state, postId)`
- `isSummariesLoading(state)`
- `getSummariesError(state)`

## API Client Methods

```typescript
import {aiClient} from 'client/ai';

// Summarize thread
const response = await aiClient.summarizeThread(postId, 'standard', true);

// Summarize channel
const response = await aiClient.summarizeChannel(
    channelId,
    startTime,
    endTime,
    'standard',
    true
);
```

## Styling

Add styles to `webapp/channels/src/components/ai/ai.scss`:

```scss
.ai-summary-panel {
    padding: 16px;
    
    .panel-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
    }
    
    .panel-content {
        margin-bottom: 16px;
    }
    
    .panel-actions {
        display: flex;
        gap: 8px;
    }
    
    .cache-indicator {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12px;
        color: var(--sys-text-secondary);
        margin-bottom: 8px;
    }
}

.ai-summary-metadata {
    padding: 12px;
    background: var(--sys-surface-secondary);
    border-radius: 4px;
    font-size: 13px;
    
    .metadata-row {
        margin-bottom: 8px;
        
        &:last-child {
            margin-bottom: 0;
        }
    }
    
    .metadata-label {
        font-weight: 600;
    }
}
```

## Testing

1. **Manual Testing**:
   - Test thread summarization via post menu
   - Test channel summarization via header menu
   - Test slash command `/summarize`
   - Verify caching works (second request faster)
   - Test different summary levels (brief, standard, detailed)
   - Test date range selection

2. **Permissions Testing**:
   - Verify users can only summarize channels they have access to
   - Test with private channels
   - Test with DMs

3. **Error Handling**:
   - Test with API disabled
   - Test with invalid API key
   - Test with no messages in range
   - Test with network errors

## Configuration

System admins must configure in `config.json`:
```json
{
    "AISettings": {
        "Enable": true,
        "OpenAIAPIKey": "sk-...",
        "OpenAIModel": "gpt-3.5-turbo",
        "MaxMessageLimit": 500,
        "EnableSummarization": true
    }
}
```

