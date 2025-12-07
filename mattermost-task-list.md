# Mattermost AI Productivity Suite
## Development Task List & PR Tracker

**Project:** Mattermost AI Productivity Suite (Native Feature Integration)  
**Repository:** `mattermost` (forked from mattermost/mattermost)  
**Date:** December 4, 2024

---

## Project File Structure

```
mattermost/
├── server/                      # Go backend, build scripts, configs
│   ├── cmd/                     # CLI entrypoints (mattermost, platform, etc.)
│   ├── channels/                # Core app (api4, app, store, jobs, wsapi)
│   │   ├── api4/                # REST API layer
│   │   │   └── ai_*.go          # NEW: AI API endpoints
│   │   ├── app/                 # Business logic layer
│   │   │   ├── ai_*.go          # NEW: AI services (summarizer, link_summarizer, etc.)
│   │   │   └── openai/          # NEW: OpenAI client package
│   │   ├── store/               # Data access layer
│   │   │   └── sqlstore/        # SQL implementations
│   │   │       └── ai_*.go      # NEW: AI data stores
│   │   └── jobs/                # Background workers
│   │       └── ai_*.go          # NEW: AI jobs (reminders, cache cleanup)
│   ├── enterprise/              # Enterprise-only logic
│   ├── config/, data/, logs/    # Runtime config & local dev assets
│   └── scripts/, tests/, public # Supporting assets and helpers
├── webapp/                      # React/TypeScript clients
│   ├── channels/                # Primary Mattermost web client
│   │   └── src/                 # Frontend source code
│   │       ├── actions/         # Redux action creators
│   │       │   └── ai_*.ts      # NEW: AI actions
│   │       ├── components/      # React components
│   │       │   └── ai/          # NEW: AI UI components
│   │       ├── reducers/        # Redux reducers
│   │       │   └── ai_*.ts      # NEW: AI reducers
│   │       ├── selectors/       # Redux selectors
│   │       │   └── ai_*.ts      # NEW: AI selectors
│   │       └── utils/           # Utility functions
│   │           └── ai_*.ts      # NEW: AI utilities
│   ├── platform/                # Shared platform packages
│   └── scripts/                 # Frontend build helpers
├── api/                         # Swagger + REST/V4 API reference
├── e2e-tests/                   # Cypress/Playwright end-to-end suites
├── memory-bank/                 # Project knowledge base (this exercise)
└── Root docs & configs          # PRD, task list, README, licenses, etc.
```

> **Note:** All AI Productivity Suite features will be integrated directly into the Mattermost core codebase (`server/channels` for backend, `webapp/channels/src` for frontend) as native functionality, demonstrating brownfield development practices.

---

## PR Tracker & Task Checklist

---

## PR #1: Project Initialization & Core Infrastructure

**Branch:** `feature/ai-infrastructure`  
**Description:** Set up the foundation for AI features including configuration, database schema, and basic service structure integrated into Mattermost core.

### Tasks

- [ ] **1.1 Repository Setup**
  - [ ] Fork Mattermost repository from mattermost/mattermost
  - [ ] Create feature branch for AI integration
  - [ ] Update workspace documentation
  - **Files Modified:**
    - `README.md` (add AI features section)
    - `memory-bank/` (project documentation)

- [ ] **1.2 Configuration Schema**
  - [ ] Add AISettings section to config structure
  - [ ] Define OpenAI API configuration fields
  - [ ] Add feature toggle flags
  - [ ] Create default configuration values
  - **Files Modified:**
    - `server/public/model/config.go` (add AISettings struct)
    - `server/config/default.json` (add default AI config)

- [ ] **1.3 Database Schema & Migrations**
  - [ ] Create AIActionItems table migration
  - [ ] Create AISummaries table migration
  - [ ] Create AILinkSummaries table migration
  - [ ] Create AIPreferences table migration
  - [ ] Add indexes for query optimization
  - **Files Created:**
    - `server/channels/db/migrations/postgres/000XXX_create_ai_action_items.up.sql`
    - `server/channels/db/migrations/postgres/000XXX_create_ai_summaries.up.sql`
    - `server/channels/db/migrations/postgres/000XXX_create_ai_link_summaries.up.sql`
    - `server/channels/db/migrations/postgres/000XXX_create_ai_preferences.up.sql`

- [ ] **1.4 Data Store Layer**
  - [ ] Create AI store interface definitions
  - [ ] Implement AIActionItemStore (SQL)
  - [ ] Implement AISummaryStore (SQL)
  - [ ] Implement AILinkSummaryStore (SQL)
  - [ ] Implement AIPreferencesStore (SQL)
  - **Files Created:**
    - `server/channels/store/ai_store.go` (interfaces)
    - `server/channels/store/sqlstore/ai_action_item_store.go`
    - `server/channels/store/sqlstore/ai_summary_store.go`
    - `server/channels/store/sqlstore/ai_link_summary_store.go`
    - `server/channels/store/sqlstore/ai_preferences_store.go`

- [ ] **1.5 OpenAI Client Package**
  - [ ] Create OpenAI client wrapper
  - [ ] Implement chat completion methods
  - [ ] Add error handling and retries
  - [ ] Implement rate limiting
  - [ ] Add request/response logging
  - **Files Created:**
    - `server/channels/app/openai/client.go`
    - `server/channels/app/openai/types.go`
    - `server/channels/app/openai/errors.go`

- [ ] **1.6 Base AI Service Layer**
  - [ ] Create AI service initialization in app layer
  - [ ] Set up OpenAI client in app context
  - [ ] Add AI configuration loading
  - [ ] Create shared AI utilities
  - **Files Created:**
    - `server/channels/app/ai.go` (service initialization)
    - `server/channels/app/ai_utils.go`

- [ ] **1.7 Frontend Infrastructure**
  - [ ] Create AI action types constants
  - [ ] Set up AI client utilities
  - [ ] Create base AI component structure
  - [ ] Add AI-specific TypeScript types
  - **Files Created:**
    - `webapp/channels/src/actions/ai_types.ts`
    - `webapp/channels/src/utils/ai_client.ts`
    - `webapp/channels/src/types/ai.ts`

- [ ] **1.8 Build System Updates**
  - [ ] Verify Go module dependencies (add OpenAI library)
  - [ ] Update webapp build configuration if needed
  - [ ] Add AI-related build targets to Makefile
  - **Files Modified:**
    - `server/go.mod` (add dependencies)
    - `Makefile` (if needed)

- [ ] **1.9 Verification**
  - [ ] Run database migrations successfully
  - [ ] Build server successfully
  - [ ] Build webapp successfully
  - [ ] Verify OpenAI client can connect (with test key)
  - [ ] Verify configuration loads correctly

---

## PR #2: AI Prompt System & API Foundation

**Branch:** `feature/ai-api-foundation`  
**Description:** Implement prompt template system, register AI API routes, and create base Redux infrastructure for AI features.

### Tasks

- [ ] **2.1 Prompt Template System**
  - [ ] Create prompt template structure
  - [ ] Implement summarization prompts (brief, standard, detailed)
  - [ ] Implement action item extraction prompts
  - [ ] Implement message formatting prompts
  - [ ] Add prompt variable substitution helpers
  - **Files Created:**
    - `server/channels/app/openai/prompts.go`
    - `server/channels/app/openai/prompt_templates.go`

- [ ] **2.2 AI API Routes Registration**
  - [ ] Create AI router initialization function
  - [ ] Register AI routes in api4 layer
  - [ ] Add permission checking middleware for AI endpoints
  - [ ] Create API response helpers
  - **Files Created:**
    - `server/channels/api4/ai.go` (route registration)
  - **Files Modified:**
    - `server/channels/api4/api.go` (call InitAI)

- [ ] **2.3 Base API Handlers**
  - [ ] Create health check endpoint for AI services
  - [ ] Implement AI configuration validation endpoint
  - [ ] Add OpenAI connectivity test endpoint
  - **Files Created:**
    - `server/channels/api4/ai_system.go`

- [ ] **2.4 Redux Store Setup (Frontend)**
  - [ ] Create AI root reducer
  - [ ] Register AI reducer in root reducer
  - [ ] Create AI state type definitions
  - [ ] Set up AI action constants
  - **Files Created:**
    - `webapp/channels/src/reducers/ai/index.ts`
    - `webapp/channels/src/types/store/ai.ts`
    - `webapp/channels/src/utils/constants/ai.ts`
  - **Files Modified:**
    - `webapp/channels/src/reducers/index.ts` (add AI reducer)

- [ ] **2.5 AI Client Service (Frontend)**
  - [ ] Create Client4 AI methods
  - [ ] Implement request/response handling
  - [ ] Add error transformation utilities
  - **Files Created:**
    - `webapp/channels/src/client/ai.ts`
  - **Files Modified:**
    - `webapp/channels/src/client/client4.ts` (add AI methods)

- [ ] **2.6 Common AI UI Components**
  - [ ] Create AI loading states component
  - [ ] Create AI error display component
  - [ ] Create AI feature badge component
  - [ ] Add AI-specific styles
  - **Files Created:**
    - `webapp/channels/src/components/ai/common/loading_state.tsx`
    - `webapp/channels/src/components/ai/common/error_display.tsx`
    - `webapp/channels/src/components/ai/common/feature_badge.tsx`
    - `webapp/channels/src/components/ai/ai.scss`

- [ ] **2.7 Shared AI Utilities (Backend)**
  - [ ] Create message formatting utilities
  - [ ] Implement participant extraction helpers
  - [ ] Add time range parsing utilities
  - **Files Created:**
    - `server/channels/app/ai_message_utils.go`
    - `server/channels/app/ai_time_utils.go`

- [ ] **2.8 Verification**
  - [ ] Test OpenAI client with sample prompts
  - [ ] Verify API routes are registered correctly
  - [ ] Test permission middleware
  - [ ] Verify frontend can call AI health check endpoint
  - [ ] Verify Redux state is properly initialized

---

## PR #3: AI Message Summarization Feature

**Branch:** `feature/ai-summarization`  
**Description:** Implement complete AI message summarization feature including backend services, REST API endpoints, slash commands, caching, and UI components.

### Tasks

- [ ] **3.1 Summarizer Service (Backend)**
  - [ ] Create summarizer service in app layer
  - [ ] Implement message fetching and filtering
  - [ ] Implement message formatting for LLM prompts
  - [ ] Add user mention resolution
  - [ ] Implement participant extraction
  - **Files Created:**
    - `server/channels/app/ai_summarizer.go`
    - `server/channels/app/ai_summarizer_types.go`

- [ ] **3.2 Thread Summarization Logic**
  - [ ] Implement thread message retrieval from store
  - [ ] Build thread context (root post + replies)
  - [ ] Generate thread summary via OpenAI client
  - [ ] Parse and structure summary response
  - **Files Modified:**
    - `server/channels/app/ai_summarizer.go`

- [ ] **3.3 Channel Summarization Logic**
  - [ ] Implement channel message retrieval with time range
  - [ ] Handle pagination for large channels
  - [ ] Implement message batching (max 500 messages)
  - [ ] Generate channel summary via OpenAI
  - **Files Modified:**
    - `server/channels/app/ai_summarizer.go`

- [ ] **3.4 Summary Caching (Database)**
  - [ ] Implement cache key generation logic
  - [ ] Add summary storage to AISummaries table
  - [ ] Implement cache retrieval and expiration checks
  - [ ] Add cache cleanup background job
  - **Files Modified:**
    - `server/channels/store/sqlstore/ai_summary_store.go`
    - `server/channels/app/ai_summarizer.go`

- [ ] **3.5 Summarizer API Endpoints**
  - [ ] Create POST /api/v4/ai/summarize endpoint
  - [ ] Implement request validation
  - [ ] Add channel membership permission checking
  - [ ] Return structured summary response
  - **Files Created:**
    - `server/channels/api4/ai_summarizer.go`
  - **Files Modified:**
    - `server/channels/api4/ai.go`

- [ ] **3.6 Slash Command Implementation**
  - [ ] Register `/summarize` command in command provider
  - [ ] Parse command arguments (thread, channel, time range)
  - [ ] Implement command execution logic
  - [ ] Return ephemeral response with summary
  - **Files Created:**
    - `server/channels/app/slashcommands/command_ai_summarize.go`
  - **Files Modified:**
    - `server/channels/app/slashcommands/slashcommands.go`

- [ ] **3.7 Redux Actions & Reducers (Frontend)**
  - [ ] Create summarize action creators
  - [ ] Implement summary reducer
  - [ ] Add summary selectors
  - [ ] Add loading/error state management
  - **Files Created:**
    - `webapp/channels/src/actions/ai_summarizer.ts`
    - `webapp/channels/src/reducers/ai/summarizer.ts`
    - `webapp/channels/src/selectors/ai_summarizer.ts`

- [ ] **3.8 Summary API Client (Frontend)**
  - [ ] Add summarizeThread method to Client4
  - [ ] Add summarizeChannel method to Client4
  - [ ] Add TypeScript types for requests/responses
  - **Files Modified:**
    - `webapp/channels/src/client/ai.ts`

- [ ] **3.9 Summary Panel Component**
  - [ ] Create RHS summary panel component
  - [ ] Display summary content with Markdown formatting
  - [ ] Show metadata (message count, participants, time range)
  - [ ] Add copy, share, regenerate actions
  - **Files Created:**
    - `webapp/channels/src/components/ai/summary/summary_panel.tsx`
    - `webapp/channels/src/components/ai/summary/summary_content.tsx`
    - `webapp/channels/src/components/ai/summary/summary_metadata.tsx`

- [ ] **3.10 Post Dropdown Integration**
  - [ ] Add "Summarize Thread" to post dropdown menu
  - [ ] Implement menu item click handler
  - [ ] Trigger RHS panel with summary
  - **Files Modified:**
    - `webapp/channels/src/components/post_view/post_menu/post_menu.tsx`
    - `webapp/channels/src/components/post_view/post_menu/index.ts`

- [ ] **3.11 Channel Header Integration**
  - [ ] Add "Summarize Channel" to channel header menu
  - [ ] Create date range picker modal
  - [ ] Trigger summarization with selected range
  - **Files Created:**
    - `webapp/channels/src/components/ai/summary/date_range_modal.tsx`
  - **Files Modified:**
    - `webapp/channels/src/components/channel_header/channel_header.tsx`

- [ ] **3.12 Verification**
  - [ ] Test thread summarization via slash command
  - [ ] Test thread summarization via post menu
  - [ ] Test channel summarization with various time ranges
  - [ ] Verify caching works correctly
  - [ ] Test permission restrictions on private channels
  - [ ] Test with large message volumes (500+ messages)

---

## PR #4: Action Item Extractor Feature

**Branch:** `feature/ai-action-items`  
**Description:** Build the AI-powered action item detection workflow, including backend services, database storage, reminders, slash commands, and the dashboards that surface personal and team commitments.

### Tasks

- [ ] **4.1 Action Item Service (Backend)**
  - [ ] Create action item service in app layer
  - [ ] Implement CRUD operations (create, list, update, delete)
  - [ ] Add validation logic for action items
  - [ ] Implement permission checking helpers
  - **Files Created:**
    - `server/channels/app/ai_action_items.go`
    - `server/channels/app/ai_action_items_types.go`

- [ ] **4.2 Detection Engine & Prompts**
  - [ ] Implement AI detection heuristics for commitments
  - [ ] Create OpenAI prompts for action item extraction
  - [ ] Parse LLM responses into structured action items
  - [ ] Extract assignees, deadlines, and priorities
  - **Files Created:**
    - `server/channels/app/ai_action_item_detector.go`
  - **Files Modified:**
    - `server/channels/app/openai/prompts.go`

- [ ] **4.3 Store Methods Implementation**
  - [ ] Implement CreateActionItem method
  - [ ] Implement GetActionItem, GetActionItemsByUser, GetActionItemsByChannel
  - [ ] Implement UpdateActionItem, DeleteActionItem
  - [ ] Add query methods for overdue/due-soon items
  - **Files Modified:**
    - `server/channels/store/sqlstore/ai_action_item_store.go`

- [ ] **4.4 Post Hook Integration**
  - [ ] Hook into MessageHasBeenPosted in app layer
  - [ ] Call action item detector asynchronously
  - [ ] Prevent duplicate detection on message edits
  - [ ] Add logging and error handling
  - **Files Modified:**
    - `server/channels/app/post.go` (add AI detection call)

- [ ] **4.5 Reminder Background Job**
  - [ ] Create reminder scheduler job
  - [ ] Implement due-soon reminder logic (24 hours before)
  - [ ] Implement overdue reminder logic (daily)
  - [ ] Send DM notifications to assignees
  - [ ] Respect user preferences for reminders
  - **Files Created:**
    - `server/channels/jobs/ai_action_item_reminders.go`
    - `server/channels/jobs/ai_action_item_reminders_scheduler.go`
  - **Files Modified:**
    - `server/channels/jobs/jobs.go` (register reminder job)

- [ ] **4.6 Action Item API Endpoints**
  - [ ] Create POST /api/v4/ai/actionitems (create item)
  - [ ] Create GET /api/v4/ai/actionitems (list with filters)
  - [ ] Create PUT /api/v4/ai/actionitems/{id} (update/complete)
  - [ ] Create DELETE /api/v4/ai/actionitems/{id}
  - [ ] Implement permission checks and validation
  - **Files Created:**
    - `server/channels/api4/ai_action_items.go`
  - **Files Modified:**
    - `server/channels/api4/ai.go`

- [ ] **4.7 Slash Command Implementation**
  - [ ] Register `/actionitems` command
  - [ ] Implement subcommands (list, mine, team, complete, create)
  - [ ] Add autocomplete support
  - [ ] Return ephemeral responses with action summaries
  - **Files Created:**
    - `server/channels/app/slashcommands/command_ai_actionitems.go`
  - **Files Modified:**
    - `server/channels/app/slashcommands/slashcommands.go`

- [ ] **4.8 Redux Actions & Reducers (Frontend)**
  - [ ] Create action item action creators
  - [ ] Implement action item reducer
  - [ ] Add action item selectors
  - [ ] Handle optimistic updates for completion
  - **Files Created:**
    - `webapp/channels/src/actions/ai_action_items.ts`
    - `webapp/channels/src/reducers/ai/action_items.ts`
    - `webapp/channels/src/selectors/ai_action_items.ts`

- [ ] **4.9 Action Items API Client (Frontend)**
  - [ ] Add createActionItem, getActionItems methods to Client4
  - [ ] Add updateActionItem, deleteActionItem, completeActionItem methods
  - [ ] Add filtering and pagination support
  - **Files Modified:**
    - `webapp/channels/src/client/ai.ts`

- [ ] **4.10 Personal Dashboard Component**
  - [ ] Create action items dashboard component
  - [ ] Build grouped lists (Overdue, Due Soon, No Deadline, Completed)
  - [ ] Add filtering by channel/status
  - [ ] Implement quick actions (mark done, view post)
  - **Files Created:**
    - `webapp/channels/src/components/ai/action_items/dashboard.tsx`
    - `webapp/channels/src/components/ai/action_items/action_item_list.tsx`

- [ ] **4.11 Action Item Card Component**
  - [ ] Display action item details (description, assignee, deadline)
  - [ ] Show source post link and channel context
  - [ ] Add inline status update controls
  - [ ] Display priority and reminder indicators
  - **Files Created:**
    - `webapp/channels/src/components/ai/action_items/action_item_card.tsx`

- [ ] **4.12 Team Action Items View**
  - [ ] Create manager/team view component
  - [ ] Group items by assignee
  - [ ] Add team-level filters
  - [ ] Implement CSV export functionality
  - **Files Created:**
    - `webapp/channels/src/components/ai/action_items/team_view.tsx`

- [ ] **4.13 Create/Edit Modal**
  - [ ] Create modal for manual action item creation
  - [ ] Add assignee selector (user picker)
  - [ ] Add due date/time picker
  - [ ] Add priority selector
  - **Files Created:**
    - `webapp/channels/src/components/ai/action_items/create_modal.tsx`
    - `webapp/channels/src/components/ai/common/date_time_picker.tsx`

- [ ] **4.14 Verification**
  - [ ] Test automatic detection on various message types
  - [ ] Verify reminders trigger correctly
  - [ ] Test permission enforcement
  - [ ] Validate slash command functionality
  - [ ] Test optimistic UI updates

---

## PR #5: Message Formatting Assistant Feature

**Branch:** `feature/ai-formatting`  
**Description:** Deliver the AI message formatting assistant, including backend formatting services, API endpoints, slash commands, and composer integration with previews and formatting profiles.

### Tasks

- [ ] **5.1 Formatter Service (Backend)**
  - [ ] Create formatter service in app layer
  - [ ] Support multiple formatting actions (professional, concise, list, code, grammar)
  - [ ] Implement OpenAI formatting requests
  - [ ] Add diff generation for preview
  - **Files Created:**
    - `server/channels/app/ai_formatter.go`
    - `server/channels/app/ai_formatter_types.go`

- [ ] **5.2 Formatting Prompt Templates**
  - [ ] Define formatting prompt templates for each profile
  - [ ] Create professional, casual, technical, concise profiles
  - [ ] Map profile metadata (label, description) for API responses
  - **Files Modified:**
    - `server/channels/app/openai/prompts.go`
  - **Files Created:**
    - `server/channels/app/ai_formatter_profiles.go`

- [ ] **5.3 User Preferences Storage**
  - [ ] Implement formatter preference getters/setters
  - [ ] Store default profile, auto-suggest, preview settings
  - [ ] Add preference validation
  - **Files Modified:**
    - `server/channels/store/sqlstore/ai_preferences_store.go`
    - `server/channels/app/ai_formatter.go`

- [ ] **5.4 Formatter API Endpoints**
  - [ ] Create POST /api/v4/ai/format/preview (preview formatting)
  - [ ] Create POST /api/v4/ai/format/apply (apply formatting)
  - [ ] Create GET /api/v4/ai/format/profiles (list available profiles)
  - [ ] Support custom instructions in request
  - **Files Created:**
    - `server/channels/api4/ai_formatter.go`
  - **Files Modified:**
    - `server/channels/api4/ai.go`

- [ ] **5.5 Slash Command Implementation**
  - [ ] Register `/format` command
  - [ ] Support action arguments (professional, concise, etc.)
  - [ ] Return ephemeral preview with apply/reject buttons
  - **Files Created:**
    - `server/channels/app/slashcommands/command_ai_format.go`
  - **Files Modified:**
    - `server/channels/app/slashcommands/slashcommands.go`

- [ ] **5.6 Post Hook Integration (Optional)**
  - [ ] Add optional auto-suggestion in MessageWillBePosted
  - [ ] Check user preferences before triggering
  - [ ] Ensure performance budget is maintained
  - **Files Modified:**
    - `server/channels/app/post.go` (optional integration)

- [ ] **5.7 Redux Actions & Reducers (Frontend)**
  - [ ] Create formatter action creators
  - [ ] Implement formatter reducer
  - [ ] Add formatter selectors
  - [ ] Manage formatting state (idle, formatting, preview)
  - **Files Created:**
    - `webapp/channels/src/actions/ai_formatter.ts`
    - `webapp/channels/src/reducers/ai/formatter.ts`
    - `webapp/channels/src/selectors/ai_formatter.ts`

- [ ] **5.8 Formatter API Client (Frontend)**
  - [ ] Add formatPreview, formatApply methods to Client4
  - [ ] Add getFormattingProfiles method
  - [ ] Handle large text payloads
  - **Files Modified:**
    - `webapp/channels/src/client/ai.ts`

- [ ] **5.9 Formatting Menu Component**
  - [ ] Create composer toolbar formatting button
  - [ ] Add dropdown menu with formatting profiles
  - [ ] Wire to formatting actions
  - [ ] Disable while formatting in progress
  - **Files Created:**
    - `webapp/channels/src/components/ai/formatter/formatting_menu.tsx`
  - **Files Modified:**
    - `webapp/channels/src/components/advanced_text_editor/advanced_text_editor.tsx`

- [ ] **5.10 Preview Modal Component**
  - [ ] Create side-by-side preview modal
  - [ ] Show original vs formatted text
  - [ ] Add diff highlighting
  - [ ] Provide apply, copy, dismiss actions
  - **Files Created:**
    - `webapp/channels/src/components/ai/formatter/preview_modal.tsx`
    - `webapp/channels/src/components/ai/formatter/diff_view.tsx`

- [ ] **5.11 Profile Selector Component**
  - [ ] Create profile selector UI
  - [ ] Display available formatting profiles
  - [ ] Add custom instruction input
  - [ ] Integrate with user preferences
  - **Files Created:**
    - `webapp/channels/src/components/ai/formatter/profile_selector.tsx`

- [ ] **5.12 Composer Integration & Styling**
  - [ ] Integrate formatting menu into text composer
  - [ ] Add inline suggestions bar (optional)
  - [ ] Apply formatted text to composer on accept
  - [ ] Add formatter-specific CSS
  - **Files Created:**
    - `webapp/channels/src/components/ai/formatter/formatter.scss`
  - **Files Modified:**
    - `webapp/channels/src/components/advanced_text_editor/advanced_text_editor.tsx`

- [ ] **5.13 Verification**
  - [ ] Test `/format` command with all profiles
  - [ ] Test composer integration
  - [ ] Verify preview modal displays correctly
  - [ ] Test with long messages (>1000 chars)
  - [ ] Verify user preferences are respected

---

## PR #6: Link & Article Summarizer Feature

**Branch:** `feature/ai-link-summarizer`  
**Description:** Implement AI-powered link and article summarization including URL detection, content fetching, AI summarization, caching, and rich preview cards in the UI.

### Tasks

- [ ] **6.1 Link Summarizer Service (Backend)**
  - [ ] Create link summarizer service in app layer
  - [ ] Implement URL detection regex patterns
  - [ ] Add content type detection (article, docs, GitHub, etc.)
  - [ ] Define summary generation methods
  - **Files Created:**
    - `server/channels/app/ai_link_summarizer.go`
    - `server/channels/app/ai_link_summarizer_types.go`

- [ ] **6.2 Content Fetcher**
  - [ ] Implement HTTP content fetcher with timeout
  - [ ] Add User-Agent and request headers
  - [ ] Implement robots.txt checking
  - [ ] Add content size limits (5MB max)
  - [ ] Handle redirects and error responses
  - **Files Created:**
    - `server/channels/app/ai_content_fetcher.go`

- [ ] **6.3 Content Extractor**
  - [ ] Implement HTML parsing and main content extraction
  - [ ] Add readability algorithm for article detection
  - [ ] Extract title, description, and metadata
  - [ ] Clean and normalize extracted text
  - [ ] Implement special extractors for GitHub, StackOverflow
  - **Files Created:**
    - `server/channels/app/ai_content_extractor.go`
    - `server/channels/app/ai_content_extractors/`

- [ ] **6.4 Link Summary Prompts**
  - [ ] Create link summarization prompt templates
  - [ ] Add key points extraction prompt
  - [ ] Implement reading time estimation
  - [ ] Add content type-specific prompts
  - **Files Modified:**
    - `server/channels/app/openai/prompts.go`
    - `server/channels/app/openai/prompt_templates.go`

- [ ] **6.5 Link Summary Store**
  - [ ] Implement CreateLinkSummary method
  - [ ] Implement GetLinkSummaryByUrl method (with URL hash)
  - [ ] Add cache expiration logic (7-day TTL)
  - [ ] Implement cleanup job for expired summaries
  - **Files Created:**
    - `server/channels/store/sqlstore/ai_link_summary_store.go`
  - **Files Modified:**
    - `server/channels/store/ai_store.go`

- [ ] **6.6 Link Summarizer API Endpoints**
  - [ ] Create POST /api/v4/ai/links/summarize
  - [ ] Create GET /api/v4/ai/links/summary?url=
  - [ ] Implement request validation
  - [ ] Add rate limiting for external fetches
  - [ ] Add permission checking
  - **Files Created:**
    - `server/channels/api4/ai_link_summarizer.go`
  - **Files Modified:**
    - `server/channels/api4/ai.go`

- [ ] **6.7 Post Hook for Auto-Detection**
  - [ ] Hook into MessageHasBeenPosted
  - [ ] Detect URLs in message content
  - [ ] Trigger async link summarization
  - [ ] Update post with summary attachment
  - **Files Modified:**
    - `server/channels/app/post.go`

- [ ] **6.8 Redux Actions & Reducers (Frontend)**
  - [ ] Create link summarizer action creators
  - [ ] Implement link summary reducer
  - [ ] Add link summary selectors
  - [ ] Handle loading and caching states
  - **Files Created:**
    - `webapp/channels/src/actions/ai_link_summarizer.ts`
    - `webapp/channels/src/reducers/ai/link_summaries.ts`
    - `webapp/channels/src/selectors/ai_link_summarizer.ts`

- [ ] **6.9 Link Summary API Client (Frontend)**
  - [ ] Add summarizeLink method to Client4
  - [ ] Add getLinkSummary method
  - [ ] Add TypeScript types for requests/responses
  - **Files Modified:**
    - `webapp/channels/src/client/ai.ts`

- [ ] **6.10 Link Summary Card Component**
  - [ ] Create rich preview card component
  - [ ] Display title, domain, favicon
  - [ ] Show AI-generated summary
  - [ ] Display key points as bullet list
  - [ ] Add reading time indicator
  - [ ] Add expand/collapse functionality
  - **Files Created:**
    - `webapp/channels/src/components/ai/link_summary/summary_card.tsx`
    - `webapp/channels/src/components/ai/link_summary/key_points.tsx`

- [ ] **6.11 Link Summary Actions**
  - [ ] Add "Read Full Article" button (opens new tab)
  - [ ] Add "Copy Summary" button
  - [ ] Add "Refresh Summary" button
  - [ ] Add collapse/expand toggle
  - **Files Modified:**
    - `webapp/channels/src/components/ai/link_summary/summary_card.tsx`

- [ ] **6.12 Post Integration**
  - [ ] Detect links in rendered posts
  - [ ] Show summary card below post with link
  - [ ] Handle multiple links in single post
  - [ ] Add loading state while fetching summary
  - **Files Created:**
    - `webapp/channels/src/components/ai/link_summary/post_link_summary.tsx`
  - **Files Modified:**
    - `webapp/channels/src/components/post_view/post_body/post_body.tsx`

- [ ] **6.13 Manual Summarize Action**
  - [ ] Add "Summarize Link" to link hover menu
  - [ ] Trigger summarization for unsummarized links
  - [ ] Show summary in modal or inline
  - **Files Created:**
    - `webapp/channels/src/components/ai/link_summary/summarize_link_button.tsx`

- [ ] **6.14 User Preferences**
  - [ ] Add auto-summarize toggle
  - [ ] Add default expanded/collapsed preference
  - [ ] Add summary length preference
  - [ ] Save preferences to AIPreferences table
  - **Files Modified:**
    - `webapp/channels/src/components/user_settings/advanced/`
    - `server/channels/store/sqlstore/ai_preferences_store.go`

- [ ] **6.15 Styling**
  - [ ] Create link summary card styles
  - [ ] Add loading skeleton animation
  - [ ] Ensure responsive design
  - [ ] Match Mattermost design system
  - **Files Created:**
    - `webapp/channels/src/components/ai/link_summary/link_summary.scss`

- [ ] **6.16 Verification**
  - [ ] Test URL detection in various message formats
  - [ ] Test content fetching for different site types
  - [ ] Verify caching works correctly (7-day TTL)
  - [ ] Test with various content types (articles, docs, GitHub)
  - [ ] Test error handling for unreachable URLs
  - [ ] Test UI summarize button functionality
  - [ ] Verify rate limiting prevents abuse

---

## PR #7: Testing, Documentation & Polish

**Branch:** `feature/ai-testing-docs`  
**Description:** Add comprehensive tests for backend and frontend, complete documentation, fix bugs, optimize performance, and prepare the feature for code review and merge.

### Tasks

- [ ] **7.1 Backend Unit Tests**
  - [ ] Write tests for AI summarizer service
  - [ ] Write tests for action item service and detector
  - [ ] Write tests for formatter service
  - [ ] Write tests for link summarizer service
  - [ ] Write tests for OpenAI client (mocked)
  - [ ] Write tests for store methods
  - **Files Created:**
    - `server/channels/app/ai_summarizer_test.go`
    - `server/channels/app/ai_action_items_test.go`
    - `server/channels/app/ai_formatter_test.go`
    - `server/channels/app/ai_link_summarizer_test.go`
    - `server/channels/app/openai/client_test.go`
    - `server/channels/store/sqlstore/ai_*_store_test.go`

- [ ] **7.2 API Endpoint Tests**
  - [ ] Write integration tests for all AI API endpoints
  - [ ] Test authentication and permissions
  - [ ] Test error responses and validation
  - [ ] Test edge cases (empty channels, large datasets)
  - **Files Created:**
    - `server/channels/api4/ai_test.go`
    - `server/channels/api4/ai_summarizer_test.go`
    - `server/channels/api4/ai_action_items_test.go`
    - `server/channels/api4/ai_formatter_test.go`
    - `server/channels/api4/ai_link_summarizer_test.go`

- [ ] **7.3 Frontend Component Tests**
  - [ ] Write tests for AI components (Jest + React Testing Library)
  - [ ] Write tests for Redux actions and reducers
  - [ ] Write tests for selectors
  - [ ] Test hooks and custom utilities
  - **Files Created:**
    - `webapp/channels/src/components/ai/**/*.test.tsx`
    - `webapp/channels/src/actions/ai_*.test.ts`
    - `webapp/channels/src/reducers/ai/*.test.ts`

- [ ] **7.4 E2E Tests (Optional)**
  - [ ] Write Cypress/Playwright tests for AI workflows
  - [ ] Test summarization end-to-end
  - [ ] Test action item creation and completion
  - [ ] Test formatting workflow
  - [ ] Test link summarization workflow
  - **Files Created:**
    - `e2e-tests/cypress/integration/ai/*.spec.js`

- [ ] **7.5 README & Feature Documentation**
  - [ ] Update main README with AI features section
  - [ ] Document each feature with examples
  - [ ] Add configuration guide
  - [ ] Add troubleshooting section
  - **Files Modified:**
    - `README.md`
  - **Files Created:**
    - `memory-bank/AI_FEATURES.md`
    - `memory-bank/AI_CONFIGURATION.md`

- [ ] **7.6 API Documentation**
  - [ ] Document all AI REST endpoints
  - [ ] Add request/response examples
  - [ ] Document slash commands
  - [ ] Add OpenAPI/Swagger specs
  - **Files Created:**
    - `api/v4/source/ai.yaml` (OpenAPI spec)
    - `memory-bank/AI_API.md`

- [ ] **7.7 Code Comments & Inline Documentation**
  - [ ] Add GoDoc comments to all exported functions
  - [ ] Add JSDoc comments to TypeScript code
  - [ ] Document complex algorithms
  - [ ] Add usage examples in comments
  - **Files Modified:**
    - All AI-related Go and TypeScript files

- [ ] **7.8 Database Migration Documentation**
  - [ ] Document database schema changes
  - [ ] Add migration rollback instructions
  - [ ] Document indexing strategy
  - **Files Created:**
    - `memory-bank/AI_DATABASE_SCHEMA.md`

- [ ] **7.9 Bug Fixes & Code Polish**
  - [ ] Fix any linting errors
  - [ ] Improve error messages and logging
  - [ ] Add missing loading states
  - [ ] Improve UI/UX responsiveness
  - [ ] Refactor duplicated code
  - **Files Modified:**
    - Various files as needed

- [ ] **7.10 Performance Optimization**
  - [ ] Verify database indexes are optimal
  - [ ] Optimize API response times
  - [ ] Implement query result caching where appropriate
  - [ ] Optimize frontend bundle size
  - [ ] Profile and optimize expensive operations
  - **Files Modified:**
    - Various files as needed

- [ ] **7.11 Security Review**
  - [ ] Audit all endpoints for authentication
  - [ ] Review SQL queries for injection risks
  - [ ] Verify XSS protection in UI components
  - [ ] Review OpenAI API key storage and handling
  - [ ] Test permission enforcement
  - **Files Modified:**
    - Various files as needed

- [ ] **7.12 Accessibility Review**
  - [ ] Ensure ARIA labels on AI components
  - [ ] Test keyboard navigation
  - [ ] Test screen reader compatibility
  - [ ] Verify color contrast ratios
  - **Files Modified:**
    - AI frontend components

- [ ] **7.13 Final Integration Testing**
  - [ ] Test full workflow: install → configure → use features
  - [ ] Test with real OpenAI API
  - [ ] Test with large datasets (1000+ messages)
  - [ ] Test concurrent user scenarios
  - [ ] Performance test under load
  - **Testing Tasks:**
    - Manual QA checklist
    - Load testing scripts

- [ ] **7.14 Demo & Presentation Materials**
  - [ ] Create demo video script
  - [ ] Record feature demonstration video (5-7 minutes)
  - [ ] Take screenshots of all features
  - [ ] Prepare presentation slides
  - **Files Created:**
    - `memory-bank/DEMO_SCRIPT.md`
    - `assets/ai-features/*.png`
    - `assets/demo-video.mp4`

---

## Progress Summary

| PR | Title | Status | Tasks |
|----|-------|--------|-------|
| #1 | Project Initialization & Core Infrastructure | ⬜ Not Started | 0/9 |
| #2 | AI Prompt System & API Foundation | ⬜ Not Started | 0/8 |
| #3 | AI Message Summarization Feature | ⬜ Not Started | 0/12 |
| #4 | Action Item Extractor Feature | ⬜ Not Started | 0/14 |
| #5 | Message Formatting Assistant Feature | ⬜ Not Started | 0/13 |
| #6 | Link & Article Summarizer Feature | ⬜ Not Started | 0/16 |
| #7 | Testing, Documentation & Polish | ⬜ Not Started | 0/14 |

**Total Tasks:** 87  
**Completed:** 0  
**Progress:** 0%

---

## Development Approach

This project follows a **brownfield development** approach, integrating AI features directly into the Mattermost core codebase rather than building a standalone plugin. This demonstrates:

- Working with an existing large-scale codebase
- Understanding and extending existing architecture (api4, app, store, Redux)
- Following established patterns and conventions
- Integrating with existing database schemas and data flows
- Extending existing UI components and state management

---

## Status Legend

- ⬜ Not Started
- 🟡 In Progress
- ✅ Completed
- ❌ Blocked

---

## Notes

- Each PR should be self-contained and deployable
- Run tests before each PR merge
- Update CHANGELOG.md with each PR
- Keep commits atomic and well-documented
- Request code review before merging to main

---

*Last Updated: December 2, 2024*