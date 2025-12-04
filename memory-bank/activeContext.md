# Active Context

## Current Work Focus
**Native feature integration of the Mattermost AI Productivity Suite**

We have shifted from a plugin-based approach to **native feature integration** to better demonstrate brownfield development skills. This involves integrating AI features directly into the Mattermost core codebase (api4, app, store layers for backend; components, actions, reducers for frontend).

## Recent Changes
- ✅ Created Memory Bank documentation structure
- ✅ Indexed and analyzed the Mattermost codebase
- ✅ Identified key components and architecture patterns
- ✅ User successfully set up local Mattermost development environment
- ✅ Created comprehensive PRD (Product Requirements Document)
- ✅ Created detailed task list with 7 PRs mapped out
- ✅ **Scope Refinement**: Removed Smart Notifications (too complex)
- ✅ **Scope Refinement**: Removed Scheduled Messages (already exists in Mattermost)
- ✅ **Added**: Action Item Extractor feature
- ✅ **Added**: Message Formatting Assistant feature
- ✅ **Configured**: Message limits (500 max, configurable)
- ✅ **ARCHITECTURE SHIFT**: Changed from plugin to native feature integration
- ✅ **Updated PRD**: Reflects native integration (database tables, api4 endpoints, Redux integration)
- ✅ **Updated Task List**: All 7 PRs rewritten for core Mattermost integration
- ✅ **Database Design**: Created schema for AIActionItems, AISummaries, AIAnalytics, AIPreferences tables
- ✅ **API Design**: Changed endpoints from `/plugins/ai-suite/*` to `/api/v4/ai/*`
- ✅ **Memory Bank Updated**: All documentation reflects brownfield development approach

## Next Steps
1. ✅ Complete codebase analysis
2. ✅ Document architecture and setup
3. ✅ User successfully ran local development setup
4. ✅ Finalize AI feature specifications
5. ✅ Update architecture from plugin to native integration
6. ⏳ **Begin PR #1**: Core infrastructure (database migrations, store layer, OpenAI client)
7. ⏳ **PR #2**: AI API foundation and prompt system
8. ⏳ **PR #3**: AI Message Summarization feature
9. ⏳ **PR #4**: Action Item Extractor feature
10. ⏳ **PR #5**: Message Formatting Assistant feature
11. ⏳ **PR #6**: Channel Analytics Dashboard feature
12. ⏳ **PR #7**: Testing, documentation, and polish

## Active Decisions and Considerations

### Native Feature Integration Approach
- **Development Type**: Brownfield Development (extending existing codebase)
- **Architecture**: Native Mattermost core integration (api4, app, store, jobs layers)
- **Storage**: PostgreSQL database tables (AIActionItems, AISummaries, AIAnalytics, AIPreferences)
- **Frontend**: Redux integration with Mattermost channels webapp
- **AI Provider**: OpenAI GPT-4 / GPT-3.5-turbo
- **Build System**: Standard Mattermost build system (Make + Webpack)

### Feature Decisions
1. **Message Limits** (Summarization)
   - Default: 500 messages max per summary
   - Configurable by system admin (range: 100-1000)
   - Prevents excessive API costs

2. **Action Item Detection**
   - AI-powered extraction from natural language
   - No external integrations (v1.0)
   - Stores in AIActionItems database table
   - Background job for reminders (native jobs framework)

3. **Message Formatting**
   - Real-time preview before applying
   - Multiple formatting profiles (Professional, Casual, Technical, Concise)
   - Preserves user's original meaning
   - Compatible with Mattermost markdown

4. **Analytics**
   - 90-day data retention
   - Stored in AIAnalytics table with daily aggregations
   - Aggregate-only (no individual message storage)
   - Privacy-conscious design
   - Background job for daily aggregation

### Scope Refinements Made
- ❌ **Removed**: Smart Notifications (too complex, high API costs)
- ❌ **Removed**: Scheduled Messages (already exists in Mattermost)
- ✅ **Added**: Action Item Extractor (high value, manageable complexity)
- ✅ **Added**: Message Formatting Assistant (unique value, low complexity)

## Current State
- ✅ Mattermost running locally on user's machine
- ✅ PRD finalized and updated for native integration
- ✅ Task list created (87 tasks, 7 PRs) - updated for core integration
- ✅ Architecture redesigned for native feature approach
- ✅ API specifications defined (`/api/v4/ai/*` endpoints)
- ✅ Database schema designed (4 new tables with migrations)
- ✅ Memory bank updated to reflect brownfield development

## Native Integration Architecture
- **Backend Layers**:
  - `app/ai_*.go` - Business logic services (Summarizer, Analytics, ActionItems, Formatter)
  - `api4/ai_*.go` - REST API handlers
  - `store/sqlstore/ai_*.go` - Database operations
  - `jobs/ai_*.go` - Background jobs (reminders, aggregation)
  - `app/openai/` - OpenAI client package
- **Frontend Layers**:
  - `components/ai/` - React UI components
  - `actions/ai_*.ts` - Redux actions
  - `reducers/ai/` - Redux state management
  - `selectors/ai_*.ts` - Data selectors
- **Integration Points**:
  - Post hooks in `app/post.go`
  - AI routes in `api4/ai.go`
  - Redux store integration
  - Slash commands in `app/slashcommands/`
- **API Endpoints**: 8+ REST endpoints under `/api/v4/ai/*`
- **Slash Commands**: 4 commands (/summarize, /actionitems, /analytics, /format)

## Known Considerations
1. **OpenAI API Costs**: GPT-3.5-turbo recommended as default (10x cheaper than GPT-4)
2. **Rate Limiting**: 60 calls/minute default, configurable in AISettings
3. **Caching Strategy**: 24-hour summary cache in AISummaries table to reduce API costs
4. **Performance**: Target <5 seconds for summarization, <1 second for analytics
5. **Permissions**: All features respect Mattermost's native channel membership permissions
6. **Database Migrations**: Must create and test migrations for 4 new AI tables
7. **Build System**: Standard Mattermost build process (no special flags required)
8. **Configuration**: New AISettings section in config.json for feature toggles and API key
9. **Background Jobs**: Scheduler for reminders and analytics aggregation
10. **Brownfield Development**: Working within existing Mattermost patterns and conventions

## Development Workflow
```
┌─────────────────────────────────────┐
│  1. Start Docker Dependencies       │
│     (make start-docker)              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  2. Build & Run Server               │
│     (make run-server)                │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  3. Build & Run Webapp               │
│     (make run-client)                │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  4. Access at http://localhost:8065 │
└─────────────────────────────────────┘
```

## Quick Reference Commands
- `make run`: Run both server and webapp
- `make stop`: Stop all services
- `make clean`: Clean build artifacts
- `make test-server`: Run backend tests
- `make test-client`: Run frontend tests
- `make help`: Show all available Make targets

