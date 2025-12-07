# System Patterns

## Architecture Overview
Mattermost follows a **client-server architecture** with clear separation between frontend and backend:

```
┌─────────────────┐
│  Web Clients    │
│  (React/TS)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Go Server      │
│  (REST API +    │
│   WebSockets)   │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌────────┐ ┌──────────┐
│ PostgreSQL│ │  MinIO   │
│           │ │ (S3-like)│
└───────────┘ └──────────┘
```

## Key Technical Decisions

### Backend Architecture
- **Single Binary Deployment**: Entire server compiles to one Go binary
- **Plugin System**: Go and JavaScript plugins via RPC
- **Database**: PostgreSQL with migration system (using `morph`)
- **Object Storage**: MinIO for file storage (S3-compatible)
- **Caching**: Redis for session and performance optimization
- **Search**: Elasticsearch or OpenSearch for full-text search

### Frontend Architecture
- **Framework**: React 18 with TypeScript
- **State Management**: Redux with redux-persist
- **Build Tool**: Webpack 5
- **Styling**: SCSS with styled-components
- **Routing**: React Router v5
- **Workspaces**: NPM workspaces for monorepo management

### Communication Patterns
- **REST API**: Version 4 API (`/api/v4/*`)
- **WebSockets**: Real-time updates and messaging
- **Event System**: Internal event bus for plugin hooks

## Component Relationships

### Server Components
```
channels/
├── api4/        # HTTP API handlers
├── app/         # Business logic layer
├── store/       # Data access layer (PostgreSQL)
├── jobs/        # Background job processing
├── wsapi/       # WebSocket API
└── web/         # Static file serving
```

### Webapp Components
```
webapp/
├── channels/                  # Main web app
│   └── src/
│       ├── components/        # React components
│       ├── actions/           # Redux actions
│       ├── reducers/          # Redux reducers
│       └── selectors/         # Redux selectors
└── platform/
    ├── client/                # API client library
    ├── types/                 # TypeScript types
    └── components/            # Shared components
```

## Design Patterns in Use

### Backend Patterns
1. **Repository Pattern**: Store layer abstracts database operations
2. **Dependency Injection**: Services injected into app server
3. **Plugin Architecture**: Hashicorp plugin system with RPC
4. **Middleware Chain**: HTTP middleware for auth, logging, rate limiting
5. **Job Queue**: Background processing for async tasks

### Frontend Patterns
1. **Container/Presenter**: Separation of logic and UI components
2. **Redux Thunks**: Async action creators
3. **Selectors**: Memoized data derivation (using `mattermost-redux/selectors/create_selector`)
4. **HOCs**: Higher-order components for cross-cutting concerns
5. **Code Splitting**: Webpack dynamic imports for lazy loading

## Clustering & High Availability
- **Leader/Follower**: Multi-node clustering with HAProxy
- **Distributed Events**: Cross-node communication via cluster
- **Session Management**: Redis for shared session state
- **Database Pooling**: Connection pooling for performance

## Security Patterns
- **Session Tokens**: JWT-based authentication
- **RBAC**: Role-based access control with permissions
- **CSRF Protection**: Token-based CSRF prevention
- **Rate Limiting**: Throttling at multiple layers
- **Plugin Sandbox**: Isolated plugin execution environments

---

## Our AI Feature Integration Patterns

### Native Integration Structure
```
mattermost/
├── server/channels/
│   ├── api4/
│   │   └── ai.go              # AI route registration
│   │   └── ai_*.go            # Feature-specific API handlers
│   ├── app/
│   │   ├── ai.go              # AI service initialization
│   │   ├── ai_*.go            # Business logic services
│   │   └── openai/            # OpenAI client package
│   ├── store/
│   │   └── sqlstore/
│   │       └── ai_*.go        # Database operations
│   ├── jobs/
│   │   └── ai_*.go            # Background jobs
│   └── db/migrations/
│       └── postgres/
│           └── 000XXX_create_ai_*.sql  # Table migrations
│
└── webapp/channels/src/
    ├── components/ai/         # AI UI components
    ├── actions/ai_*.ts        # Redux actions
    ├── reducers/ai/           # Redux state
    ├── selectors/ai_*.ts      # Data selectors
    └── client/ai.ts           # API client methods
```

### Native Integration Patterns Used

1. **Layered Architecture Pattern**
   - API Layer (`api4/`) - HTTP handlers, validation, permissions
   - Business Logic (`app/`) - Core AI services
   - Data Layer (`store/`) - Database operations
   - Jobs Layer (`jobs/`) - Background processing
   - Clear separation with well-defined interfaces

2. **Database-First Storage**
   - Dedicated PostgreSQL tables for each feature
   - Proper migrations with up/down scripts
   - Indexed for query performance
   - Transactional integrity

3. **Redux State Management (Frontend)**
   - Actions → Reducers → Selectors pattern
   - Immutable state updates
   - Memoized selectors for performance
   - Integration with existing Mattermost Redux store

4. **Post Hook Integration**
   - Extend existing `app.MessageHasBeenPosted()`
   - Call AI services asynchronously
   - Non-blocking message flow
   - Respects existing Mattermost patterns

5. **Background Jobs Pattern**
   - Use native Mattermost jobs framework
   - Schedulers for reminders and aggregation
   - Persistent across restarts
   - Configurable intervals

6. **Caching Strategy**
   - Summary cache: 24 hours (AISummaries table)
   - Analytics: Pre-aggregated daily (AIAnalytics table)
   - Cache invalidation via TTL (ExpiresAt column)

### Native Integration Points

- **REST API**: Native endpoints at `/api/v4/ai/*`
- **Slash Commands**: Registered via `app/slashcommands/`
- **RHS Panel**: React components in `components/ai/`
- **Channel Header**: Extended with AI menu items
- **Message Composer**: AI formatting integration
- **Redux Store**: AI reducers integrated into root reducer
- **Database**: 4 new tables with proper migrations (AIActionItems, AISummaries, AILinkSummaries, AIPreferences)

### Data Flow Patterns

**Summarization Flow:**
1. User triggers `/summarize` or clicks button
2. Webapp → Client4.summarizeThread() → `/api/v4/ai/summarize`
3. API handler validates permissions
4. App service fetches messages from store
5. Format → Send to OpenAI → Parse response
6. Save to AISummaries table (cache)
7. Return to user → Display in RHS panel via Redux

**Action Item Detection Flow:**
1. User posts message → `app.MessageHasBeenPosted()` extended
2. AI detector checks for commitments (async)
3. If detected → OpenAI extraction
4. Create action item → Store in AIActionItems table
5. Notify assignee via DM
6. Frontend fetches via `/api/v4/ai/actionitems`
7. Update Redux store → Re-render dashboard

**Message Formatting Flow:**
1. User types message → Clicks AI formatting button in composer toolbar
2. FormattingMenu component loads profiles (if not loaded)
3. User selects formatting profile from dropdown
4. Webapp → Client4.formatPreview() → `/api/v4/ai/format/preview`
5. API handler validates and calls formatter service
6. Formatter service sends to OpenAI with profile-specific prompt
7. Response parsed → Diff generated → Returned to frontend
8. Preview modal displays (side-by-side or diff view)
9. User clicks "Apply" → Formatted text replaces message in composer
10. Redux state updated → Preview cleared

**Link Summarization Flow:**
1. User posts message with URL → URL detected via regex
2. System fetches URL content asynchronously
3. Content extracted and cleaned (readability algorithm)
4. Sent to OpenAI for summary generation
5. Summary cached in AILinkSummaries table (7-day TTL)
6. Frontend displays rich preview card below message
7. Subsequent requests served from cache

**Action Items Creation Flow (with Global Modal):**
1. User clicks "Create Action Item" from post menu
2. PostActionItemMenuItem calls global `showCreateModal()` function
3. Global function creates container div, appends to document.body
4. Modal rendered via `ReactDOM.render()` with Provider + IntlProvider wrappers
5. Modal independent of menu component lifecycle (persists after menu closes)
6. User fills form → Submit → `dispatch(createActionItem(request))`
7. API call → `/api/v4/ai/actionitems` → Backend validation and save
8. Success → `dispatch(getActionItems({}))` to refresh list
9. Modal cleanup via global cleanup function

### UI Patterns Implemented

**Global Modal Manager Pattern:**
- **Problem**: Modals triggered from dropdown menus disappear when parent unmounts
- **Solution**: Render modal directly to `document.body` using `ReactDOM.render()`
- **Implementation**:
  ```typescript
  let currentModal: {element: HTMLElement; cleanup: () => void} | null = null;
  
  const showCreateModal = (postId: string, channelId: string) => {
      if (currentModal) {
          currentModal.cleanup();
      }
      
      const container = document.createElement('div');
      document.body.appendChild(container);
      
      const cleanup = () => {
          ReactDOM.unmountComponentAtNode(container);
          document.body.removeChild(container);
          currentModal = null;
      };
      
      ReactDOM.render(
          <Provider store={store}>
              <IntlProvider>
                  <CreateActionItemModal onClose={cleanup} />
              </IntlProvider>
          </Provider>,
          container
      );
      
      currentModal = {element: container, cleanup};
  };
  ```
- **Benefits**: Modal lifecycle independent of triggering component
- **Trade-offs**: Manual context management (Provider, IntlProvider must be wrapped)

**Channel Header Button Pattern:**
- **Pattern**: Add icon button to channel header for quick access to RHS features
- **Example**: Action Items button (✓ icon) next to Files button
- **Implementation**: 
  - Add icon and class in `channel_header.tsx` render method
  - Add click handler that dispatches RHS action (`openRHSForActionItems`)
  - Wire action in Redux connector (`index.ts`)
  - Icon highlights when RHS state matches (`rhsState === RHSStates.ACTION_ITEMS`)
- **Benefits**: Immediate visual access without menu navigation

**RHS Panel State Pattern (Critical!):**
- **Pattern**: When adding new RHS panel types, must update TWO files for `searchVisible`
- **Problem**: The `Search` component wraps RHS content and has its own Redux connection
- **Files to Update**:
  1. `webapp/channels/src/components/sidebar_right/index.ts` - Add to `searchVisible` exclusion
  2. `webapp/channels/src/components/search/index.tsx` - Add to `searchVisible` exclusion
- **Implementation**:
  ```typescript
  // In sidebar_right/index.ts
  searchVisible: Boolean(rhsState) && rhsState !== RHSStates.PLUGIN && rhsState !== RHSStates.ACTION_ITEMS,
  
  // In search/index.tsx
  searchVisible: rhsState !== null && (![
      RHSStates.PLUGIN,
      RHSStates.CHANNEL_INFO,
      RHSStates.CHANNEL_MEMBERS,
      RHSStates.EDIT_HISTORY,
      RHSStates.ACTION_ITEMS,  // Must add new RHS states here!
  ].includes(rhsState)),
  ```
- **Why**: The Search component checks `searchVisible` and shows SearchResults instead of `children` when true
- **Symptom**: New RHS panel shows Search UI instead of your custom content

**Form Button Pattern:**
- **Pattern**: Always use `type='button'` for buttons that shouldn't submit forms
- **Problem**: Buttons without explicit type default to `type='submit'` inside forms
- **Solution**:
  ```tsx
  <button
      type='button'  // Prevents form submission
      onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          // handler logic
      }}
  >
  ```
- **When to Use**: Any button inside a form that should NOT submit (menus, toggles, etc.)

**Redux State Path Pattern:**
- **Pattern**: AI reducers are at `state.ai`, NOT `state.entities.ai`
- **Mattermost Redux Structure**:
  ```
  state
  ├── entities/     # mattermost-redux entities (users, channels, posts, teams)
  ├── views/        # UI state (modals, selections)
  ├── plugins/      # Plugin state
  ├── storage/      # Persistent storage
  └── ai/           # Our AI feature state (custom reducers)
  ```
- **Correct Selector**:
  ```typescript
  const getAIState = (state: GlobalState) => state.ai;
  ```
- **Wrong** (common mistake):
  ```typescript
  const getAIState = (state: GlobalState) => state.entities.ai;  // entities is mattermost-redux
  ```

**Go JSON Serialization Pattern:**
- **Pattern**: Always add JSON tags to Go structs that will be serialized to JSON
- **Problem**: Go uses PascalCase by default, but frontend expects snake_case
- **Solution**:
  ```go
  // Wrong - Go serializes as {"Title": "...", "Summary": "..."}
  type LinkSummary struct {
      Title   string
      Summary string
  }
  
  // Correct - Serializes as {"title": "...", "summary": "..."}
  type LinkSummary struct {
      Title   string   `json:"title,omitempty"`
      Summary string   `json:"summary"`
  }
  ```
- **Symptom**: Frontend receives data but can't access properties (undefined)

**HTTP Content Fetching Pattern:**
- **Pattern**: Use realistic browser headers when fetching external URLs
- **Problem**: Bot User-Agents get blocked (403 Forbidden)
- **Solution**:
  ```go
  // Use realistic Chrome User-Agent
  req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
  req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
  req.Header.Set("Accept-Language", "en-US,en;q=0.9")
  
  // DON'T set Accept-Encoding manually - Go handles gzip automatically
  // req.Header.Set("Accept-Encoding", "gzip")  // BAD - causes garbled content
  ```
- **Why**: Go's default HTTP transport handles gzip compression transparently, but if you set Accept-Encoding manually, Go expects you to decompress yourself

**React/Redux Race Condition Pattern:**
- **Pattern**: Account for async Redux updates when using local state
- **Problem**: Local state updates sync, Redux updates async → UI flickers
- **Solution**:
  ```tsx
  const [requested, setRequested] = useState<Set<string>>(new Set());
  const loading = useSelector(isLinkSummaryLoading);
  const summary = useSelector(getLinkSummaryForUrl);
  const error = useSelector(getLinkSummaryError);
  
  // Consider loading if Redux says loading OR if we just requested but Redux hasn't updated
  const effectiveLoading = loading || (requested.has(url) && !summary && !error);
  ```
- **Symptom**: Clicking button briefly shows empty state before loading skeleton

**Map Destructuring Pattern:**
- **Pattern**: Always destructure ALL needed variables from array.map()
- **Problem**: Variable used in JSX but not destructured → undefined
- **Solution**:
  ```tsx
  // Wrong - error is undefined
  {summaries.map(({url, summary, loading}) => (
      {error && <div>{error}</div>}  // error is not defined!
  ))}
  
  // Correct - destructure error too
  {summaries.map(({url, summary, loading, error}) => (
      {error && <div>{error}</div>}
  ))}
  ```
- **Symptom**: Condition never evaluates to true even when data exists

