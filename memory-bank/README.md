# Memory Bank - Mattermost Project

This directory contains comprehensive documentation about the Mattermost codebase and our **AI Productivity Suite** development project.

## 🎯 Current Project: AI Productivity Suite (Native Features)

**Status**: PR #6 core functionality complete - All 4 AI features working!  
**Timeline**: 6-7 days, 87 tasks, 7 PRs  
**Progress**: 6 of 7 PRs complete (~86%)  
**Features**: 4 AI-powered productivity enhancements - ALL FUNCTIONAL ✅

### AI Features
1. **AI Message Summarization** ✅ - GPT-powered summaries (max 500 messages)
2. **Link & Article Summarizer** ✅ - AI summaries of shared URLs and articles
3. **Action Item Extractor** ✅ - Auto-detect and track commitments
4. **Message Formatting Assistant** ✅ - AI grammar and tone improvements

### Development Progress
- [x] **PR #1**: Core Infrastructure (database, store, OpenAI client) ✅ Dec 4
- [x] **PR #2**: API Foundation (routes, prompts, Redux, UI components) ✅ Dec 5
- [x] **PR #3**: AI Message Summarization ✅ Dec 5
- [x] **PR #4**: Action Item Extractor ✅ Dec 5-6
- [x] **PR #5**: Message Formatting Assistant ✅ Dec 5-6
- [x] **PR #6**: Link & Article Summarizer ✅ Dec 7 (Core working)
- [ ] PR #7: Testing, Documentation & Polish

📄 **Full Details**: See `mattermost-prd.md` and `mattermost-task-list.md` in project root

## 📚 Documentation Files

### Core Documentation

1. **[projectbrief.md](projectbrief.md)**
   - High-level project overview
   - Core technologies (Go + React)
   - Repository structure
   - Key features and goals
   - **Start here** for project overview

2. **[productContext.md](productContext.md)**
   - Why Mattermost exists
   - Problems it solves
   - Primary use cases
   - User experience goals
   - **Read this** to understand the product

3. **[systemPatterns.md](systemPatterns.md)**
   - Architecture overview
   - Key technical decisions
   - Design patterns in use
   - Component relationships
   - **Essential** for understanding architecture

4. **[techContext.md](techContext.md)**
   - Complete technology stack
   - Development setup requirements
   - Technical constraints
   - Dependencies management
   - **Reference** for technical details

5. **[activeContext.md](activeContext.md)**
   - Current work focus
   - Recent changes
   - Next steps
   - Development workflow
   - **Check this** to see current state

6. **[progress.md](progress.md)**
   - What works
   - Current development tasks
   - Testing status
   - Performance notes
   - **Track progress** here

### Guides

7. **[LOCAL_DEVELOPMENT_GUIDE.md](LOCAL_DEVELOPMENT_GUIDE.md)**
   - **Complete setup instructions**
   - Prerequisites and installation
   - Quick start commands
   - Troubleshooting guide
   - **ESSENTIAL** for running locally

## 🚀 Quick Start

### To Run Mattermost Locally

1. **Read the setup guide**:
   ```
   Open: memory-bank/LOCAL_DEVELOPMENT_GUIDE.md
   ```

2. **Ensure prerequisites**:
   - Go 1.24.6+
   - Node.js 18.10.0+
   - Docker Desktop
   - Make (GNU Make)

3. **Run with one command**:
   ```powershell
   cd server
   make run
   ```

4. **Access the application**:
   ```
   http://localhost:8065
   ```

### To Test AI Features

1. **Set OpenAI API Key**:
   ```bash
   export MM_AISETTINGS_OPENAIAPIKEY="sk-your-key-here"
   ```

2. **Test Message Summarization**:
   - Type `/summarize channel` in any channel

3. **Test Link Summarizer**:
   - Post a URL in a message
   - Click "Summarize link" button

4. **Test Action Items**:
   - Click the checkmark icon in channel header
   - Or post a message like "I'll review this by tomorrow"

5. **Test Message Formatting**:
   - Type a message in the composer
   - Click the AI formatting button in toolbar

## 📖 Reading Order

### For Complete Understanding
1. Start → `projectbrief.md` (5 min read)
2. Then → `productContext.md` (5 min read)
3. Deep dive → `systemPatterns.md` (10 min read)
4. Reference → `techContext.md` (10 min read)
5. Setup → `LOCAL_DEVELOPMENT_GUIDE.md` (follow along)

### For Quick Setup
1. `LOCAL_DEVELOPMENT_GUIDE.md` - Complete setup guide
2. `activeContext.md` - Current workflow
3. Start coding!

## 🎯 Common Questions

**Q: What is Mattermost?**
A: See `projectbrief.md` and `productContext.md`

**Q: How do I run it locally?**
A: Follow `LOCAL_DEVELOPMENT_GUIDE.md`

**Q: What's the architecture?**
A: Read `systemPatterns.md`

**Q: What technologies are used?**
A: Check `techContext.md`

**Q: What's the current status?**
A: See `activeContext.md` and `progress.md`

**Q: How do I test AI features?**
A: See Quick Start section above

## 🏗️ Project Structure Summary

```
mattermost/
├── server/              # Go backend
│   ├── channels/       # Core functionality
│   │   ├── api4/      # REST API (ai_*.go for AI endpoints)
│   │   ├── app/       # Business logic (ai_*.go for AI services)
│   │   ├── store/     # Data layer
│   │   └── wsapi/     # WebSocket API
│   └── cmd/           # CLI tools
│
├── webapp/             # React frontend
│   ├── channels/      # Main web app
│   │   └── src/
│   │       ├── components/ai/  # AI UI components
│   │       ├── actions/        # Redux actions (ai_*.ts)
│   │       ├── reducers/ai/    # AI reducers
│   │       └── selectors/      # AI selectors (ai_*.ts)
│   └── platform/      # Shared libraries
│
└── memory-bank/       # This documentation
```

## 🛠️ Essential Commands

```powershell
# Run everything
cd server && make run

# Access app
# → http://localhost:8065

# Stop everything
cd server && make stop

# Run tests
cd server && make test-server    # Backend
cd server && make test-client     # Frontend

# Get help
cd server && make help
```

## 📊 Codebase Scale

- **Backend**: 2,000+ Go files
- **Frontend**: 7,933 files (3,479 TypeScript/TSX)
- **Tests**: Extensive coverage (Go + Jest + Cypress + Playwright)
- **Database Migrations**: 292+ SQL files
- **API Endpoints**: 148 handlers in api4/
- **Business Logic**: 441 files in app/
- **AI Feature Files**: 30+ new files for AI Productivity Suite

## 🔄 Development Workflow

```
1. Start Docker → make start-docker
2. Run Server → make run-server (Terminal 1)
3. Run Webapp → make run-client (Terminal 2)
4. Edit Code → Auto-reload
5. Test → make test-server / make test-client
6. Commit → Git workflow
```

## 🎓 Learning Path

### Day 1: Setup & Overview
- [ ] Read `projectbrief.md`
- [ ] Read `productContext.md`
- [ ] Follow `LOCAL_DEVELOPMENT_GUIDE.md`
- [ ] Get app running locally
- [ ] Create test account
- [ ] Explore the UI

### Day 2: Architecture Deep Dive
- [ ] Read `systemPatterns.md`
- [ ] Read `techContext.md`
- [ ] Explore `server/channels/app/`
- [ ] Explore `webapp/channels/src/components/`
- [ ] Review API documentation

### Day 3: Development
- [ ] Make a small change
- [ ] Run tests
- [ ] Debug an issue
- [ ] Read plugin documentation
- [ ] Explore extensions

## 📚 External Resources

- **Developer Docs**: https://developers.mattermost.com/
- **API Reference**: https://api.mattermost.com/
- **Community**: https://community.mattermost.com/
- **GitHub**: https://github.com/mattermost/mattermost
- **Contributing Guide**: `../CONTRIBUTING.md`

## 🔍 Finding Your Way

### Backend Development
Start in: `server/channels/app/`
- `user.go` - User management
- `channel.go` - Channel operations
- `post.go` - Post/message handling
- `team.go` - Team management
- `ai_*.go` - AI feature services

### Frontend Development
Start in: `webapp/channels/src/components/`
- Browse component folders
- Check Redux actions/reducers
- Review selectors for data access
- `ai/` folder - AI UI components

### API Development
Check: `server/channels/api4/`
- Each file handles specific API endpoints
- `api4.go` - Main API initialization
- `ai.go` - AI endpoint registration
- `ai_*.go` - AI feature endpoints

## 📝 Keep This Updated

This Memory Bank should be updated when:
- Major architectural changes occur
- New patterns are established
- Setup process changes
- Important decisions are made
- After completing significant features

To update: Edit the relevant .md files in this directory

## 🆘 Getting Help

If stuck:
1. Check troubleshooting in `LOCAL_DEVELOPMENT_GUIDE.md`
2. Review relevant documentation files
3. Check official docs at https://developers.mattermost.com/
4. Ask in ~Contributors channel on community server
5. Check GitHub issues

---

**Last Updated**: December 7, 2024  
**Status**: ✅ All 4 AI features functional (Summarization, Action Items, Formatting, Link Summarizer)  
**Version**: AI Productivity Suite - 86% complete (6/7 PRs done)
