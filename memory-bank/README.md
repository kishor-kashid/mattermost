# Memory Bank - Mattermost Project

This directory contains comprehensive documentation about the Mattermost codebase and our **AI Productivity Suite Plugin** development project.

## 🎯 Current Project: AI Productivity Suite (Native Features)

**Status**: PR #5 complete and verified - Moving to PR #6  
**Timeline**: 6-7 days, 87 tasks, 7 PRs  
**Progress**: 5 of 7 PRs complete (~71%)  
**Features**: 4 AI-powered productivity enhancements

### AI Features
1. **AI Message Summarization** - GPT-powered summaries (max 500 messages)
2. **Link & Article Summarizer** - AI summaries of shared URLs and articles
3. **Action Item Extractor** - Auto-detect and track commitments
4. **Message Formatting Assistant** - AI grammar and tone improvements

### Development Progress
- [x] **PR #1**: Core Infrastructure (database, store, OpenAI client) ✅ Dec 4
- [x] **PR #2**: API Foundation (routes, prompts, Redux, UI components) ✅ Dec 5
- [x] **PR #3**: AI Message Summarization ✅ Dec 5
- [x] **PR #4**: Action Item Extractor ✅ Dec 5
- [x] **PR #5**: Message Formatting Assistant ✅ Dec 5 (Complete - all 13 tasks)
- [ ] PR #6: Link & Article Summarizer
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

## 🏗️ Project Structure Summary

```
mattermost/
├── server/              # Go backend
│   ├── channels/       # Core functionality
│   │   ├── api4/      # REST API
│   │   ├── app/       # Business logic (START HERE for backend)
│   │   ├── store/     # Data layer
│   │   └── wsapi/     # WebSocket API
│   └── cmd/           # CLI tools
│
├── webapp/             # React frontend
│   ├── channels/      # Main web app
│   │   └── src/
│   │       └── components/  # React components (START HERE for frontend)
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

### Frontend Development
Start in: `webapp/channels/src/components/`
- Browse component folders
- Check Redux actions/reducers
- Review selectors for data access

### API Development
Check: `server/channels/api4/`
- Each file handles specific API endpoints
- `api4.go` - Main API initialization
- `user.go` - User endpoints
- `channel.go` - Channel endpoints

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

**Last Updated**: December 5, 2024  
**Status**: ✅ PR #5 complete, 3 features fully implemented (Summarization, Action Items, Formatting)  
**Version**: Updated with AI Productivity Suite native integration project

