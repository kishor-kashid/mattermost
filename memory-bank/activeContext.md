# Active Context

## Current Work Focus
**PR #6 In Progress – Link & Article Summarizer**

Link summarizer scaffolding is partially implemented (backend + frontend). Action Items and Message Formatting remain stable.

### Latest Session (PR6)
Built initial Link Summarizer stack and UI:
1. **Data & Store**: Added `AILinkSummary` model, store interface, SQL store, migrations `000152_create_ai_link_summaries` and `000153_update_ai_preferences_link_settings`.
2. **App & API**: Added link summarizer service, fetcher (timeout, size cap, concurrent fetch limiter, redirect cap), extractor, prompts, and endpoints:  
   - POST `/api/v4/ai/links/summarize`  
   - GET `/api/v4/ai/links/summary?url=`  
   Health reports link feature; post hook triggers async summarize of first URL in posts.
3. **Frontend Data**: Added link summary types, client methods, actions, reducer, selectors.
4. **UI**: Link summary card + key points + styles; post body renders summaries for multiple links per message with refresh.

## Recent Changes (PR6 – Link Summarizer)

- Added link summary model/store, migrations, and API endpoints.
- Added app service with OpenAI prompts, fetcher (timeout, size limit, redirect limit, concurrent fetch limiter), extractor, caching (7-day TTL), and post hook auto-detection of first URL.
- Frontend: client methods, Redux actions/reducer/selectors, multi-link summary rendering, summary card UI with key points and refresh.
- Frontend UX change: link summarization no longer auto-runs. Each link now shows a “Summarize link” button; summaries fetch on click, not on post render. Added inline error display + retry when backend fetch fails (e.g., target URL 403/404).

## Current Feature Status

| Feature | Status | Notes |
|---------|--------|-------|
| AI Message Summarization | ✅ Complete | Fully working |
| Action Item Extractor | ✅ Complete | RHS panel fixed, dashboard working |
| Message Formatting Assistant | ✅ Complete | Stable |
| Link & Article Summarizer | 🚧 In progress | Backend/API/store/UI scaffolding done; prefs/actions/tests pending |

## Next Steps
1. Backend hardening: robots/redirect handling polish, rate limiting for external fetches, scheduled cleanup job for expired link summaries, stronger validation/permissions.
2. Preferences: expose link prefs (auto summarize, default expanded, summary length) via API/UI.
3. UI actions: manual “Summarize Link” action, copy/expand/collapse controls, better loading/error states.
4. Multi-link polish and error handling.
5. Tests: backend (service/API/store), frontend (actions/reducer/components), and E2E happy path.
6. PR #7: Testing, Documentation & Polish after PR6 completes.

## Key Patterns Learned This Session

### RHS Panel Integration Pattern
When adding a new RHS panel state:
1. Add constant to `utils/constants.tsx` under `RHSStates`
2. Add `is[Feature]` prop to `sidebar_right/index.ts`
3. Add exclusion to `searchVisible` in BOTH:
   - `sidebar_right/index.ts`
   - `search/index.tsx` (has its own Redux connection!)
4. Add rendering logic in `sidebar_right.tsx`

### Form Button Pattern
Always use `type='button'` for buttons that shouldn't submit forms:
```tsx
<button
    type='button'  // Prevents form submission
    onClick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        // handler
    }}
>
```

### Redux State Path Pattern
AI reducers are at `state.ai`, NOT `state.entities.ai`:
```typescript
// Correct
const getAIState = (state: GlobalState) => state.ai;

// Wrong (entities is for mattermost-redux data)
const getAIState = (state: GlobalState) => state.entities.ai;
```

## Architecture Reference

### Frontend State Structure
```
state
├── entities/          # mattermost-redux entities (users, channels, posts)
├── views/             # UI state
├── plugins/           # Plugin state
├── storage/           # Persistent storage
└── ai/                # AI feature state (our code)
    ├── summaries
    ├── actionItems
    ├── analytics
    ├── preferences
    ├── formatter
    └── system
```

### RHS Panel State Flow
```
User clicks button
    → dispatch(openRHSForActionItems(channelId))
    → rhsState = 'action-items'
    → sidebar_right checks isActionItems
    → search/index.tsx checks searchVisible (must be false!)
    → Search component renders children (ActionItemsRHS)
    → Dashboard fetches and displays items
```

## Development Workflow
```
┌─────────────────────────────────────┐
│  1. Start Docker Dependencies       │
│     (make start-docker)             │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  2. Build & Run Server              │
│     (make run-server)               │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  3. Build & Run Webapp              │
│     (npm run dev in webapp/)        │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  4. Access at http://localhost:8065 │
└─────────────────────────────────────┘
```

## Files Modified This Session
- `webapp/channels/src/components/search/index.tsx` - Added ACTION_ITEMS exclusion
- `webapp/channels/src/components/sidebar_right/index.ts` - Added ACTION_ITEMS exclusion
- `webapp/channels/src/components/ai/formatter/formatting_menu.tsx` - Fixed button type
- `webapp/channels/src/selectors/ai_formatter.ts` - Fixed state path
- Multiple files - Removed debug logging
- `webapp/channels/src/components/ai/link_summary/post_link_summary.tsx` - Switched to manual summarize trigger, added error display/retry, shallow-equal selector
- `webapp/channels/src/components/advanced_text_editor/formatting_bar/formatting_bar.tsx` - Added stable keys for additional controls

