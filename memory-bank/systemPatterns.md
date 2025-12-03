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

## Our Plugin Architecture Patterns

### Plugin Structure Pattern
```
mattermost-plugin-ai-suite/
├── server/                    # Go backend
│   ├── plugin.go              # Entry point (OnActivate, ServeHTTP)
│   ├── configuration.go       # Config handling
│   ├── api.go                 # REST router
│   ├── hooks.go               # Mattermost hooks
│   │
│   ├── openai/                # Service: OpenAI integration
│   ├── summarizer/            # Service: Summarization
│   ├── analytics/             # Service: Analytics
│   ├── actionitems/           # Service: Action items
│   ├── formatter/             # Service: Formatting
│   └── store/                 # Data access (KV store)
│
└── webapp/                    # React frontend
    └── src/
        ├── components/        # UI components
        ├── hooks/             # React hooks
        └── services/          # API clients
```

### Plugin Design Patterns Used

1. **Service Layer Pattern**
   - Each feature has its own service package
   - Services are injected into plugin struct
   - Clear separation of concerns

2. **Key-Value Store Pattern**
   - All data stored in plugin KV store
   - No custom database tables
   - Namespaced keys (e.g., `summary:channel_id:timestamp`)

3. **Hook Pattern**
   - `MessageHasBeenPosted`: Analytics collection, action item detection
   - `MessageWillBePosted`: Optional formatting suggestions
   - Minimal performance impact

4. **API Gateway Pattern**
   - Single `ServeHTTP` entry point
   - Routes to feature handlers
   - Centralized auth/permission checks

5. **Caching Strategy**
   - Summary cache: 24 hours (KV store)
   - Analytics: 1-hour aggregation cache
   - Cache invalidation on relevant events

### Plugin Integration Points

- **REST API**: Custom endpoints at `/plugins/ai-suite/api/v1/*`
- **Slash Commands**: 4 commands registered with Mattermost
- **RHS Panel**: React component injection for summaries
- **Main Menu**: Analytics dashboard entry
- **Channel Header**: Quick action dropdowns
- **Message Composer**: Formatting assistant integration

### Data Flow Patterns

**Summarization Flow:**
1. User triggers `/summarize` or clicks button
2. Webapp → REST API → Summarizer service
3. Fetch messages from Mattermost API
4. Format → Send to OpenAI → Parse response
5. Cache result → Return to user
6. Display in RHS panel

**Action Item Detection Flow:**
1. User posts message → `MessageHasBeenPosted` hook
2. Background check for commitment patterns
3. If detected → OpenAI extraction
4. Create action item → Notify assignee
5. Store in KV → Update dashboards

**Analytics Collection Flow:**
1. Message posted → `MessageHasBeenPosted` hook
2. Extract metrics (non-blocking)
3. Aggregate to hourly buckets
4. Store in KV (rolling 90 days)
5. Dashboard queries aggregated data

