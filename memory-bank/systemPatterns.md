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
3. **Selectors**: Memoized data derivation (using reselect)
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
- **Database**: 4 new tables with proper migrations

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

**Analytics Collection Flow:**
1. Message posted → `app.MessageHasBeenPosted()` extended
2. Extract metrics (non-blocking)
3. Update AIAnalytics table (upsert for current day)
4. Background job aggregates hourly
5. Frontend queries `/api/v4/ai/analytics/{channelId}`
6. Redux actions update analytics state
7. Charts re-render with new data

