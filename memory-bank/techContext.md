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

## AI Feature Integration Technical Context

### Native Integration Technologies

#### Backend (Mattermost Core)
- **Language**: Go 1.24+
- **API Layer**: `server/channels/api4/` (REST handlers)
- **Business Logic**: `server/channels/app/` (services)
- **Data Access**: `server/channels/store/sqlstore/` (PostgreSQL)
- **Jobs**: `server/channels/jobs/` (background workers)
- **Storage**: PostgreSQL database with new AI tables
- **AI/LLM**: OpenAI Go SDK (`sashabaranov/go-openai`)
- **HTTP**: Built-in `net/http` + `gorilla/mux`

#### Frontend (Mattermost Webapp)
- **Framework**: React 18 (integrated into Mattermost channels)
- **State Management**: Redux 5 with thunks
- **Component Location**: `webapp/channels/src/components/ai/`
- **Actions**: `webapp/channels/src/actions/ai_*.ts`
- **Reducers**: `webapp/channels/src/reducers/ai/`
- **Styling**: SCSS (follows Mattermost patterns)
- **Charts**: Recharts (for analytics)
- **Build**: Webpack 5 (standard Mattermost build)

#### External Dependencies
- **OpenAI API**: GPT-4 / GPT-3.5-turbo
- **API Rate Limits**: 60 requests/minute (configurable)
- **Cost Management**: Caching, message limits

### Native Build Integration

```makefile
# Standard Mattermost build commands
make run-server          # Build and run Go server
make run-client          # Build and run React webapp
make test-server         # Run backend tests
make test-client         # Run frontend tests
make build-server        # Build server binary only
make build-client        # Build webapp bundle only
```

### Configuration Schema

**config.json:**
```json
{
  "AISettings": {
    "Enable": true,
    "OpenAIAPIKey": "encrypted_key",
    "OpenAIModel": "gpt-4",
    "MaxMessageLimit": 500,
    "APIRateLimit": 60,
    "EnableSummarization": true,
    "EnableAnalytics": true,
    "EnableActionItems": true,
    "EnableFormatting": true
  }
}
```

**Go Struct:**
```go
type AISettings struct {
    Enable             bool
    OpenAIAPIKey       string
    OpenAIModel        string
    MaxMessageLimit    int
    APIRateLimit       int
    EnableSummarization bool
    EnableAnalytics    bool
    EnableActionItems  bool
    EnableFormatting   bool
}
```

### Data Storage Strategy

**Database Tables:**
- **AIActionItems**: Action items with assignees, deadlines, status
- **AISummaries**: Cached summaries with expiration (24h TTL)
- **AIAnalytics**: Daily aggregated channel metrics
- **AIPreferences**: Per-user AI feature preferences

**Retention Policies:**
- Summaries: 24 hours (cache with ExpiresAt)
- Analytics: 90 days rolling (cleanup job)
- Action Items: Until completed/dismissed (soft delete)
- User Preferences: Indefinite

### API Integration Patterns

#### OpenAI API Usage
- **Endpoint**: `https://api.openai.com/v1/chat/completions`
- **Authentication**: Bearer token (from config.AISettings.OpenAIAPIKey)
- **Error Handling**: Retry with exponential backoff
- **Rate Limiting**: Client-side throttling (configurable)
- **Timeout**: 30 seconds per request

#### Native Mattermost API Usage
- **Get Posts**: Direct store access via `app.GetPostsPage()`
- **Get Channel**: `app.GetChannel()`
- **Get User**: `app.GetUser()`
- **Permissions**: `app.HasPermissionToChannel()`
- **Create Post**: `app.CreatePost()` for notifications

### Performance Considerations

**Optimization Strategies:**
1. **Database Caching**: 24-hour summary cache in AISummaries table
2. **Query Optimization**: Proper indexes on AI tables
3. **Background Processing**: Jobs don't block message flow
4. **Pre-aggregation**: Daily analytics aggregation job
5. **Lazy Loading**: Frontend components fetch on demand
6. **Redux Memoization**: Reselect for computed state

**Resource Limits:**
- Max messages for summarization: 500 (configurable)
- Max concurrent OpenAI requests: 5
- Analytics query limit: 90 days
- Database storage: ~10MB estimated for typical usage

### Security Implementation

**Data Protection:**
- API keys encrypted in config (Mattermost encryption)
- Summaries cached 24h only (auto-expired)
- All API calls respect native Mattermost permissions
- No external data transmission (except OpenAI)

**Permission Checks:**
- Use `app.HasPermissionToChannel()` for channel access
- Respect team membership for analytics
- Action item visibility follows channel permissions
- System admin required for AISettings configuration

### Development Dependencies

**Backend (Go modules):**
```go
require (
    github.com/mattermost/mattermost/server/public v0.0.0
    github.com/sashabaranov/go-openai v1.41.2  // Added in PR #1
    github.com/pkg/errors v0.9.x
)
```

**Frontend (NPM - already in Mattermost):**
```json
{
  "react": "^18.2.0",
  "react-redux": "^9.x",
  "redux": "^5.x",
  "recharts": "^2.x",
  "@mattermost/types": "latest"
}
```

### Testing Strategy

**Backend Unit Tests:**
- Go: Standard `*_test.go` files
- Test app layer services with mocked store
- Test store layer with test database
- Coverage target: 80%+

**Frontend Unit Tests:**
- Jest + React Testing Library
- Test components, actions, reducers, selectors
- Mock Client4 API calls
- Coverage target: 80%+

**Integration Tests:**
- API endpoint tests (`api4/*_test.go`)
- Test complete request/response flow
- Mock OpenAI API responses
- Test permissions enforcement

**E2E Tests (Optional):**
- Cypress/Playwright tests
- Test full user workflows
- Verify UI interactions

