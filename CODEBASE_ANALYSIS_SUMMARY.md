# Mattermost Codebase Analysis Summary

## 📋 Overview

**Mattermost** is an open-source, self-hosted team collaboration platform built with:
- **Backend**: Go (server-side API and business logic)
- **Frontend**: React + TypeScript (web application)
- **Database**: PostgreSQL
- **Architecture**: Monorepo with separate server and webapp components

---

## 🏗️ Repository Structure

```
mattermost/
├── server/                  # Go backend server
│   ├── channels/           # Core channel functionality
│   │   ├── api4/          # REST API v4 handlers
│   │   ├── app/           # Business logic
│   │   ├── store/         # Data access layer
│   │   ├── jobs/          # Background jobs
│   │   └── wsapi/         # WebSocket API
│   ├── cmd/               # Command-line tools
│   │   ├── mattermost/    # Main server binary
│   │   └── mmctl/         # CLI admin tool
│   ├── public/            # Public Go modules
│   ├── platform/          # Platform services
│   ├── config/            # Configuration management
│   └── build/             # Build scripts and Docker configs
│
├── webapp/                  # React/TypeScript frontend
│   ├── channels/           # Main web application
│   │   └── src/
│   │       ├── components/ # React components
│   │       ├── actions/    # Redux actions
│   │       ├── reducers/   # Redux reducers
│   │       └── selectors/  # Redux selectors
│   └── platform/           # Shared platform code
│       ├── client/         # API client library
│       ├── types/          # TypeScript definitions
│       └── components/     # Shared components
│
├── api/                     # API documentation (OpenAPI)
├── e2e-tests/              # End-to-end tests
│   ├── cypress/            # Cypress E2E tests
│   └── playwright/         # Playwright E2E tests
├── tools/                  # Development tools
└── memory-bank/            # Project documentation (NEW)
```

---

## 🔑 Key Components

### Backend (Go)

1. **API Layer** (`channels/api4/`)
   - RESTful API handlers
   - Version 4 of the API
   - HTTP endpoint implementations

2. **Business Logic** (`channels/app/`)
   - Core application logic
   - User management, channels, posts, teams
   - File handling, notifications, plugins
   - 441+ files covering all features

3. **Data Access** (`channels/store/`)
   - Database abstraction layer
   - PostgreSQL queries
   - Transaction management
   - Caching integration

4. **Jobs System** (`channels/jobs/`)
   - Background task processing
   - Scheduled jobs
   - Async operations

5. **WebSocket API** (`channels/wsapi/`)
   - Real-time communication
   - Live updates
   - Event broadcasting

### Frontend (React/TypeScript)

1. **Components** (`webapp/channels/src/components/`)
   - 2,394 TypeScript/TSX files
   - Reusable UI components
   - Container and presentational patterns

2. **State Management**
   - Redux for global state
   - Redux Thunk for async actions
   - Selectors for derived data
   - Redux Persist for offline support

3. **Platform Libraries**
   - `@mattermost/client`: API client
   - `@mattermost/types`: TypeScript definitions
   - `@mattermost/components`: Shared components

---

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.24.6
- **Web Framework**: Custom (Gorilla Mux)
- **Database**: PostgreSQL 14+
- **Migrations**: golang-migrate
- **WebSockets**: gorilla/websocket
- **Storage**: MinIO (S3-compatible)
- **Cache**: Redis
- **Search**: Elasticsearch/OpenSearch
- **Auth**: JWT, LDAP, SAML

### Frontend
- **Framework**: React 18.2.0
- **Language**: TypeScript 5.6.3
- **State**: Redux 5.0.1
- **Build**: Webpack 5.95.0
- **Styling**: SCSS + Styled Components
- **Testing**: Jest + React Testing Library
- **Icons**: Compass Icons
- **UI Libraries**: MUI, React Bootstrap

### Development Tools
- **Build**: Make (GNU Make)
- **Containers**: Docker + Docker Compose
- **E2E Testing**: Cypress + Playwright
- **Linting**: ESLint, golangci-lint
- **Package Management**: Go modules, NPM workspaces

---

## 🚀 Running Locally

### Quick Start (3 Commands)

```powershell
# 1. Navigate to server directory
cd server

# 2. Run everything (Docker + Server + Webapp)
make run

# 3. Open browser to http://localhost:8065
```

### Prerequisites
✅ Go 1.24.6+
✅ Node.js 18.10.0+
✅ NPM 9.0.0+
✅ Docker Desktop
✅ Make (GNU Make)

### Detailed Steps

**Step 1: Start Docker Dependencies**
```powershell
cd server
make start-docker
```

Starts: PostgreSQL, MinIO, Redis, Inbucket (email), Prometheus, Grafana

**Step 2: Start Server**
```powershell
cd server
make run-server
```

Builds and runs the Go backend on port 8065

**Step 3: Start Webapp**
```powershell
cd server
make run-client
```

Builds React app with webpack in watch mode

**Step 4: Access Application**
```
http://localhost:8065
```

### First-Time Setup
- First build downloads dependencies (5-10 minutes)
- Create admin account on first access
- Database migrations run automatically
- Sample data available: `make test-data`

---

## 📊 Codebase Statistics

### Backend (Server)
- **Total Files**: 2,000+ Go files
- **Core App Logic**: 441 files
- **API Handlers**: 148 files
- **Store Layer**: 295 files
- **Jobs**: 68 files
- **Tests**: Extensive coverage throughout

### Frontend (Webapp)
- **Total Files**: 7,933 files
- **TypeScript/TSX**: 3,479 files
- **Components**: 2,394+ React components
- **Images/Assets**: 3,379 PNG files
- **Tests**: Jest test files throughout

### Overall
- **Languages**: Go, TypeScript, JavaScript, YAML, SQL
- **Lines of Code**: 100,000+ (estimate)
- **Database Migrations**: 292+ SQL migration files
- **Prepackaged Plugins**: 13 official plugins

---

## 🏛️ Architecture Patterns

### Backend Patterns
- **Repository Pattern**: Store layer abstracts database
- **Dependency Injection**: Service-based architecture
- **Plugin System**: Hashicorp plugin architecture
- **Middleware Chain**: HTTP middleware stack
- **Event Bus**: Internal event system for plugins

### Frontend Patterns
- **Container/Presenter**: Separation of concerns
- **Redux Thunks**: Async action creators
- **Selectors**: Memoized data derivation
- **HOCs**: Higher-order components
- **Code Splitting**: Lazy-loaded modules

---

## 🔌 Key Features

### Core Functionality
- ✅ Real-time messaging (WebSocket)
- ✅ Channels and teams
- ✅ File sharing and storage
- ✅ Search (Elasticsearch/OpenSearch)
- ✅ Notifications (email, push, desktop)
- ✅ User and team management
- ✅ OAuth, LDAP, SAML authentication
- ✅ Mobile and desktop apps
- ✅ Extensive plugin system
- ✅ Webhooks and slash commands
- ✅ API (REST + WebSocket)

### Enterprise Features
- 🔒 High availability clustering
- 🔒 Advanced LDAP/AD integration
- 🔒 SAML 2.0 SSO
- 🔒 Compliance and data retention
- 🔒 Advanced metrics and monitoring
- 🔒 Guest access controls

---

## 🧪 Testing

### Test Infrastructure
- **Backend**: Go testing framework
- **Frontend**: Jest + React Testing Library
- **E2E**: Cypress (986 test files) + Playwright (117 specs)
- **Integration**: Docker-based test environment

### Running Tests
```powershell
# Server tests
cd server
make test-server

# Webapp tests
cd server
make test-client

# E2E tests
cd e2e-tests/cypress
npm test
```

---

## 📚 Documentation

### Created Documentation (in memory-bank/)
- ✅ `projectbrief.md` - Project overview and goals
- ✅ `productContext.md` - Why Mattermost exists, problems it solves
- ✅ `systemPatterns.md` - Architecture and design patterns
- ✅ `techContext.md` - Technologies and dependencies
- ✅ `activeContext.md` - Current work and next steps
- ✅ `progress.md` - Development status and tasks
- ✅ `LOCAL_DEVELOPMENT_GUIDE.md` - Complete setup guide

### Official Resources
- Developer Docs: https://developers.mattermost.com/
- API Docs: https://api.mattermost.com/
- Community: https://community.mattermost.com/
- GitHub: https://github.com/mattermost/mattermost

---

## 🎯 Common Development Tasks

### Making Changes
```powershell
# Backend changes (auto-restart)
cd server
make run-server
# Edit files in server/

# Frontend changes (auto-rebuild)
cd server
make run-client
# Edit files in webapp/
```

### Managing Services
```powershell
cd server

# Start everything
make run

# Stop everything
make stop

# Restart server only
make restart-server

# Clean build artifacts
make clean

# Clean Docker
make clean-docker
```

### Using mmctl CLI
```powershell
cd server

# Create user (server must be running)
bin/mmctl user create --email user@example.com --username myuser --password Password1! --local

# List teams
bin/mmctl team list --local
```

---

## 🔍 Important Directories

### Configuration
- `server/config/` - Configuration management
- `server/config/config.json` - Default config file
- `server/build/docker/` - Docker configurations

### Development
- `server/Makefile` - Build automation
- `webapp/Makefile` - Frontend build automation
- `server/docker-compose.yaml` - Docker services
- `server/build/docker-compose.common.yml` - Service definitions

### Source Code
- `server/channels/app/` - **Main business logic** (start here)
- `webapp/channels/src/components/` - **UI components** (start here)
- `server/public/model/` - Data models
- `webapp/platform/types/` - TypeScript types

---

## 🚦 Next Steps

1. **✅ Codebase Indexed**: Complete
2. **✅ Documentation Created**: Complete
3. **⏳ Run Locally**: Ready for you to execute
4. **⏳ Explore Features**: After setup
5. **⏳ Begin Development**: Based on your needs

---

## 💡 Quick Reference

| Task | Command |
|------|---------|
| Start everything | `cd server && make run` |
| Access app | http://localhost:8065 |
| Stop everything | `cd server && make stop` |
| Server tests | `cd server && make test-server` |
| Webapp tests | `cd server && make test-client` |
| Add sample data | `cd server && make test-data` |
| View all commands | `cd server && make help` |

---

## 📝 Notes

- First build takes 5-10 minutes (downloads dependencies)
- Subsequent builds are much faster (cached)
- Enterprise features require `../../enterprise` directory
- Database migrations run automatically
- All development dependencies run in Docker
- Changes auto-reload during development

---

**Status**: ✅ Codebase fully analyzed and documented
**Ready**: 🚀 For local development
**Next**: Follow `memory-bank/LOCAL_DEVELOPMENT_GUIDE.md` to start

