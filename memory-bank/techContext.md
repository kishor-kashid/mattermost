# Technical Context

## Technologies Used

### Backend Stack
- **Language**: Go 1.24.6
- **Web Framework**: Custom (built on `gorilla/mux`)
- **Database**: PostgreSQL 14+
- **Database Migration**: `golang-migrate/migrate`
- **ORM/Query Builder**: `sqlx` + `squirrel`
- **WebSocket**: `gorilla/websocket`
- **Object Storage**: MinIO (S3-compatible)
- **Caching**: Redis
- **Search**: Elasticsearch 8.x / OpenSearch 4.x
- **Authentication**: JWT (`golang-jwt/jwt`)
- **LDAP**: Custom LDAP client
- **SAML**: `gosaml2`

### Frontend Stack
- **Framework**: React 18.2.0
- **Language**: TypeScript 5.6.3
- **State Management**: Redux 5.0.1 + Redux Thunk
- **Build Tool**: Webpack 5.95.0
- **Package Manager**: NPM (workspaces)
- **Styling**: SCSS + styled-components
- **Testing**: Jest 30.1.3 + React Testing Library
- **Linting**: ESLint 8.57.0
- **Code Editor**: Monaco Editor (in-app)

### Development Tools
- **Build System**: Make (GNU Make)
- **Containerization**: Docker + Docker Compose
- **E2E Testing**: Cypress + Playwright
- **API Testing**: Go testing framework
- **Mocking**: `vektra/mockery` (Go), Jest mocks (JS)
- **Code Coverage**: Go coverage + Jest coverage

### Infrastructure Services (Development)
- **PostgreSQL**: Database
- **MinIO**: Object storage
- **Inbucket**: Email testing
- **OpenLDAP**: LDAP testing
- **Elasticsearch/OpenSearch**: Search indexing
- **Redis**: Caching and session storage
- **Keycloak**: SAML/OAuth testing
- **Prometheus + Grafana**: Metrics and monitoring
- **Loki + Promtail**: Log aggregation

## Development Setup Requirements

### Prerequisites
- **Go**: 1.24.6 or compatible
- **Node.js**: >=18.10.0
- **NPM**: ^9.0.0 or ^10.0.0
- **Docker**: For running development dependencies
- **Make**: Build automation
- **Git**: Version control

### Optional Tools
- **Delve**: Go debugger
- **jq**: JSON processing (for config management)
- **golangci-lint**: Go linting
- **mmctl**: CLI tool for Mattermost administration

## Technical Constraints

### Performance Requirements
- Sub-second message delivery
- Support for thousands of concurrent users
- Efficient database query patterns
- WebSocket connection management at scale

### Compatibility Requirements
- Browser support: Modern browsers (Chrome, Firefox, Safari, Edge)
- Mobile apps: iOS and Android
- Desktop apps: Windows, macOS, Linux
- Database: PostgreSQL 11+ (14+ recommended)

### Security Constraints
- TLS/SSL for production deployments
- Secure password hashing (bcrypt)
- CSRF protection enabled
- Rate limiting on API endpoints
- Session token rotation
- Plugin signature verification

## Build Configuration

### Server Build
- **Entry Point**: `server/cmd/mattermost/main.go`
- **Build Tags**: `enterprise`, `sourceavailable`, `requirefips` (optional)
- **LDFLAGS**: Build metadata injection
- **Output**: Single binary (`mattermost`)

### Webapp Build
- **Entry Point**: `webapp/channels/src/`
- **Webpack Config**: `webapp/channels/webpack.config.js`
- **Output**: `webapp/channels/dist/`
- **Assets**: Bundled JS, CSS, fonts, images

## Dependencies Management
- **Go Modules**: `go.mod`, `go.sum` for backend dependencies
- **NPM Workspaces**: Monorepo with multiple packages
- **Prepackaged Plugins**: Downloaded from official releases
- **Vendor Directory**: Not used (relies on module cache)

## Environment Variables
Key environment variables for development:
- `MM_SQLSETTINGS_DRIVERNAME`: Database driver (postgres)
- `MM_SQLSETTINGS_DATASOURCE`: Database connection string
- `MM_SERVICESETTINGS_SITEURL`: Base URL
- `MM_SERVICESETTINGS_ENABLELOCALMODE`: Local development mode
- `MM_NO_DOCKER`: Disable Docker dependency startup
- `BUILD_ENTERPRISE_READY`: Enable enterprise features

---

## Plugin Development Technical Context

### Plugin Technologies

#### Backend
- **Framework**: Mattermost Plugin SDK (server)
- **Language**: Go 1.24+
- **API Client**: Built-in `plugin.API` interface
- **Storage**: Plugin Key-Value Store
- **AI/LLM**: OpenAI Go SDK (`sashabaranov/go-openai`)
- **JSON**: Standard library `encoding/json`
- **HTTP**: Built-in `net/http`

#### Frontend
- **Framework**: React (same as Mattermost webapp)
- **Plugin SDK**: Mattermost webapp plugin SDK
- **State**: React Hooks + Context API
- **Styling**: SCSS (follows Mattermost patterns)
- **Charts**: Recharts (for analytics)
- **Build**: Webpack (plugin template)

#### External Dependencies
- **OpenAI API**: GPT-4 / GPT-3.5-turbo
- **API Rate Limits**: 60 requests/minute (configurable)
- **Cost Management**: Caching, message limits

### Plugin Build System

```makefile
# Key Make targets for plugin development
make                    # Build both server and webapp
make dist               # Create distribution package
make deploy             # Deploy to local Mattermost
make check-style        # Lint Go and JS code
make test               # Run all tests
```

### Plugin Configuration Schema

```go
type Configuration struct {
    OpenAIAPIKey      string  // Encrypted in DB
    OpenAIModel       string  // "gpt-4" or "gpt-3.5-turbo"
    MaxMessageLimit   int     // Default: 500
    APIRateLimit      int     // Default: 60/min
    EnableSummarization bool  // Default: true
    EnableAnalytics   bool    // Default: true
}
```

### Data Storage Strategy

**Key-Value Store Schema:**
- Summaries: `summary:{channel_id}:{hash}` → JSON
- Analytics: `analytics:{channel_id}:{hour}` → Aggregated metrics
- Action Items: `actionitem:{id}` → Action item object
- User Prefs: `userprefs:{user_id}` → Preferences JSON
- Config: `plugin_configuration` → Plugin settings

**Retention Policies:**
- Summaries: 24 hours (cache only)
- Analytics: 90 days rolling
- Action Items: Until completed/dismissed
- User Preferences: Indefinite

### API Integration Patterns

#### OpenAI API Usage
- **Endpoint**: `https://api.openai.com/v1/chat/completions`
- **Authentication**: Bearer token (API key)
- **Error Handling**: Retry with exponential backoff
- **Rate Limiting**: Client-side throttling
- **Timeout**: 30 seconds per request

#### Mattermost Plugin API Usage
- **Get Messages**: `plugin.API.GetPostsForChannel()`
- **Create Post**: `plugin.API.CreatePost()`
- **Get User**: `plugin.API.GetUser()`
- **KV Operations**: `plugin.API.KVSet()`, `KVGet()`
- **Permissions**: `plugin.API.HasPermissionToChannel()`

### Performance Considerations

**Optimization Strategies:**
1. **Caching**: 24-hour summary cache reduces API calls by 90%+
2. **Pagination**: Fetch messages in batches (100 per request)
3. **Background Processing**: Action item detection doesn't block message posting
4. **Aggregation**: Analytics pre-aggregated hourly
5. **Lazy Loading**: Dashboard components load data on demand

**Resource Limits:**
- Max message length for summarization: 500 messages
- Max concurrent OpenAI requests: 5
- Analytics query limit: 90 days
- KV storage per plugin: ~100MB estimated

### Security Implementation

**Data Protection:**
- API keys encrypted at rest (Mattermost encryption)
- No message content stored (except summaries, 24h cache)
- All API calls respect Mattermost permissions
- No external data transmission (except OpenAI)

**Permission Checks:**
- User must have channel read access
- Analytics requires channel membership
- Action item visibility follows channel permissions
- System admin required for plugin configuration

### Development Dependencies

```json
// Go dependencies (go.mod)
{
  "github.com/mattermost/mattermost-plugin-starter-template": "latest",
  "github.com/mattermost/mattermost/server/v8": "v8.x",
  "github.com/sashabaranov/go-openai": "v1.x",
  "github.com/pkg/errors": "v0.9.x"
}

// NPM dependencies (package.json)
{
  "react": "^18.2.0",
  "react-dom": "^18.2.0",
  "recharts": "^2.x",
  "@mattermost/types": "latest"
}
```

### Testing Strategy

**Unit Tests:**
- Go: `*_test.go` files using standard testing package
- JS: Jest with React Testing Library
- Coverage target: 90%+

**Integration Tests:**
- Mock Mattermost Plugin API
- Mock OpenAI API responses
- Test complete feature workflows

**Manual Testing:**
- Install plugin on local Mattermost
- Test each feature end-to-end
- Verify error handling
- Check performance metrics

