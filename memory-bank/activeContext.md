# Active Context

## Current Work Focus
**Plugin implementation of the Mattermost AI Productivity Suite (PR #2 complete, PR #3 next)**

We now have the shared OpenAI, storage, and API foundations in place; next up is delivering the feature layers (summaries, action items, analytics, formatter).

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
- ✅ **Aligned**: Task list PRs now match scoped features (PR4 Action Items, PR5 Formatting, PR6 Analytics) with updated testing coverage
- ✅ **Synced**: PRD updated to remove legacy scheduled-messages scope and document action items + formatter flows
- ✅ **Completed PR #1**: Added plugin manifest, Go module (with `server/public` replace), Makefile, README/SETUP docs, server entrypoint/config, and placeholder webapp bundle
- ✅ **Completed PR #2**: Built OpenAI GPT-4 client wrapper + prompt templates, KV store abstraction, REST API/router middleware, plugin configuration settings, and base webapp API client/common UI components
- ✅ **Build Guidance**: Documented that `GOWORK=off` is required when running `go build`/`make bundle` inside the plugin and noted npm’s `--legacy-peer-deps` workaround for TypeScript peer conflicts

## Next Steps
1. ✅ Complete codebase analysis
2. ✅ Document architecture and setup
3. ✅ User successfully ran local development setup
4. ✅ Finalize plugin feature specifications
5. ✅ **PR #1 Complete**: Project initialization and plugin scaffold
6. ✅ **PR #2 Complete**: OpenAI integration and core services
7. ⏳ **Begin PR #3**: Summarization feature implementation
8. ⏳ Continue through remaining PRs

## Active Decisions and Considerations

### Plugin Development Approach
- **Plugin Type**: Mattermost Server + Webapp Plugin
- **Architecture**: Follows Mattermost plugin SDK patterns
- **Storage**: Using plugin Key-Value Store (no custom DB tables)
- **AI Provider**: OpenAI GPT-4 / GPT-3.5-turbo
- **Build System**: Makefile-based (standard plugin structure)

### Feature Decisions
1. **Message Limits** (Summarization)
   - Default: 500 messages max per summary
   - Configurable by system admin (range: 100-1000)
   - Prevents excessive API costs

2. **Action Item Detection**
   - AI-powered extraction from natural language
   - No external integrations (v1.0)
   - Stores in plugin KV store
   - Background processing for reminders

3. **Message Formatting**
   - Real-time preview before applying
   - Multiple formatting profiles (Professional, Casual, Technical, Concise)
   - Preserves user's original meaning
   - Compatible with Mattermost markdown

4. **Analytics**
   - 90-day data retention
   - Aggregate-only (no individual message storage)
   - Privacy-conscious design

### Scope Refinements Made
- ❌ **Removed**: Smart Notifications (too complex, high API costs)
- ❌ **Removed**: Scheduled Messages (already exists in Mattermost)
- ✅ **Added**: Action Item Extractor (high value, manageable complexity)
- ✅ **Added**: Message Formatting Assistant (unique value, low complexity)

## Current State
- ✅ Mattermost running locally on user's machine
- ✅ PRD finalized (1,453 lines, comprehensive)
- ✅ Task list created (86 tasks, 7 PRs)
- ✅ Architecture designed
- ✅ API specifications defined
- ✅ Plugin scaffold builds (`make bundle`) when run with `GOWORK=off` and npm legacy peer deps
- ✅ GPT-4 OpenAI client, KV store service, REST API router, and base webapp infrastructure implemented (PR #2)
- ⏳ Ready to implement summarization feature (PR #3)

## Plugin Architecture Decisions
- **Backend Services**: Summarizer, Analytics, ActionItems, Formatter
- **Core Services**: OpenAI Client, Message Processor, Store wrapper
- **Plugin Hooks**: MessageHasBeenPosted (analytics + action items), MessageWillBePosted (formatting)
- **API Endpoints**: 8 REST endpoints defined
- **Slash Commands**: 4 commands (/summarize, /actionitems, /analytics, /format)

## Known Considerations
1. **OpenAI API Costs**: GPT-3.5-turbo recommended as default (10x cheaper than GPT-4)
2. **Rate Limiting**: 60 calls/minute default, configurable
3. **Caching Strategy**: 24-hour summary cache to reduce API costs
4. **Performance**: Target <5 seconds for summarization, <1 second for analytics
5. **Permissions**: All features respect Mattermost's channel membership permissions
6. **Build Tooling**: When compiling from `server/plugins/ai-suite`, set `GOWORK=off` (workspace otherwise points to monorepo root). Webapp npm install currently requires `--legacy-peer-deps` due to `@mattermost/types` optional TypeScript 4.x peer.

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

