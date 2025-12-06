# Mattermost Local Development Guide

## ⚠️ Windows Users - Read This First!

**If you're on Windows**, the `make` commands in this guide require a Unix-like shell. **You have two options:**

1. **Git Bash** (Quick) - Use Git Bash instead of PowerShell
   - Open Git Bash: Start Menu → Git Bash
   - Navigate: `cd "/c/Gauntlet AI/Week8/mattermost/server"`
   - Run: `make run`

2. **WSL2** (Recommended for development) - Install Windows Subsystem for Linux
   - PowerShell as Admin: `wsl --install`
   - Restart computer
   - Use Ubuntu terminal

**📖 Full Windows Guide:** See [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md) for detailed instructions.

---

## Prerequisites

Before starting, ensure you have the following installed:

### Required Software
1. **Go** (version 1.24.6 or compatible)
   - Download: https://golang.org/dl/
   - Verify: `go version`

2. **Node.js** (version 18.10.0 or higher)
   - Download: https://nodejs.org/
   - Verify: `node --version`

3. **NPM** (version 9.0.0 or 10.0.0+)
   - Comes with Node.js
   - Verify: `npm --version`

4. **Docker Desktop** (for Windows)
   - Download: https://www.docker.com/products/docker-desktop
   - Required for: PostgreSQL, MinIO, Redis, and other development services
   - Verify: `docker --version` and `docker compose version`

5. **Make** (GNU Make)
   - On Windows, this typically comes with Git for Windows
   - Or install via Chocolatey: `choco install make`
   - Verify: `make --version`

6. **Git**
   - Likely already installed
   - Verify: `git --version`

## Quick Start (Recommended)

The fastest way to get Mattermost running locally:

### Option 1: Run Everything Together

```powershell
# Navigate to the server directory
cd server

# This single command will:
# - Start Docker dependencies
# - Build the Go server
# - Build the React webapp
# - Start both server and webapp
make run
```

The application will be available at: **http://localhost:8065**

### Option 2: Run Server and Webapp Separately

This is better for active development:

```powershell
# Terminal 1 - Start Docker dependencies and server
cd server
make run-server

# Terminal 2 - Start webapp (in a new terminal)
cd server
make run-client
```

## Detailed Setup Steps

### Step 1: Start Docker Dependencies

```powershell
cd server
make start-docker
```

This starts:
- **PostgreSQL** (port 5432): Database
- **MinIO** (port 9000): Object storage (S3-compatible)
- **Inbucket** (port 9001): Email testing server
- **Redis** (port 6379): Caching
- **Prometheus** (port 9090): Metrics
- **Grafana** (port 3000): Monitoring dashboards
- **Loki** (port 3100): Log aggregation
- **Promtail** (port 3180): Log collection

**First-time note**: This will download Docker images, which may take several minutes.

### Step 2: Build and Run the Server

```powershell
cd server
make run-server
```

This will:
1. Set up Go workspace
2. Download Go dependencies
3. Compile prepackaged binaries (mmctl)
4. Download prepackaged plugins
5. Build and start the Mattermost server

The server will listen on **port 8065**.

**What to expect**:
- First build takes 5-10 minutes (downloading dependencies)
- Subsequent builds are much faster (cached)
- Database migrations run automatically
- Server logs will appear in the terminal

### Step 3: Build and Run the Webapp

In a **new terminal**:

```powershell
cd server
make run-client
```

This will:
1. Navigate to the webapp directory
2. Install NPM dependencies
3. Start webpack in watch mode
4. Build the React application

**What to expect**:
- First build takes several minutes (downloading NPM packages)
- Webpack will watch for file changes and rebuild automatically
- Build output appears in the terminal

### Step 4: Access Mattermost

Open your browser and navigate to:
```
http://localhost:8065
```

You should see the Mattermost login/signup page.

## Initial Setup

### Create Your First Admin Account

1. Go to http://localhost:8065
2. Click "Create an account" or "Sign up"
3. Fill in the account creation form
4. The first account created becomes the System Admin

### Create a Team

1. After signup, you'll be prompted to create a team
2. Enter a team name
3. Start using Mattermost!

## Common Development Commands

### Server Commands

```powershell
cd server

# Start server with all dependencies
make run-server

# Stop the server
make stop-server

# Restart the server
make restart-server

# Run server tests
make test-server

# Run server in debug mode
make debug-server

# Clean build artifacts
make clean

# Stop Docker services
make stop-docker

# Clean Docker containers and volumes
make clean-docker
```

### Webapp Commands

```powershell
cd server

# Run webapp in watch mode (preferred for development)
# - Creates/refreshes a symlink: server/client -> webapp/channels/dist
# - Webpack watches files and rebuilds automatically
make run-client

# Stop webapp
make stop-client

# Restart webapp
make restart-client

# Run webapp tests
make test-client

# Clean webapp build
cd ../webapp && make clean

# One-off: build static bundle for server-only runs (no run-client/dev-server)
cd ../webapp
npm run build --workspace=channels
cd ..
Remove-Item -Recurse -Force .\server\client\*
Copy-Item -Recurse -Force .\webapp\channels\dist\* .\server\client\
```

### Combined Commands

```powershell
cd server

# Run both server and webapp
make run

# Stop both
make stop

# Restart both
make restart

# Run all tests
make test
```

## Development Workflow

### Making Changes

1. **Backend (Go) changes**:
   - Edit files in `server/`
   - Server auto-restarts on changes (if using `make run-server`)
   - Or manually restart: `make restart-server`

2. **Frontend (React/TypeScript) changes**:
   - Edit files in `webapp/`
   - Webpack rebuilds automatically
   - Refresh browser to see changes
   - Hot module replacement enabled for faster development

### Running Tests

```powershell
# Server tests
cd server
make test-server

# Webapp tests
cd server
make test-client

# Run specific test
cd server
go test ./channels/app -run TestSpecificFunction

# Run webapp tests in watch mode
cd webapp/channels
npm run test:watch
```

### Linting and Code Quality

```powershell
# Server linting
cd server
make check-style

# Webapp linting
cd webapp
npm run check

# Auto-fix issues
npm run fix
```

## Troubleshooting

### Docker Issues

**Problem**: Docker services won't start
```powershell
# Check Docker is running
docker ps

# Restart Docker Desktop
# Then try again
cd server
make start-docker
```

**Problem**: Port conflicts
```powershell
# Stop other services using the same ports
# Or stop and clean Docker services
cd server
make clean-docker
make start-docker
```

### Build Issues

**Problem**: Go build fails
```powershell
# Clean and rebuild
cd server
make clean
make run-server
```

**Problem**: NPM install fails
```powershell
# Clean and reinstall
cd webapp
make clean
npm install
```

**Problem**: Missing Go workspace
```powershell
cd server
make setup-go-work
```

### Database Issues

**Problem**: Database connection errors
```powershell
# Verify PostgreSQL is running
docker ps | findstr postgres

# Check logs
docker logs mattermost-postgres

# Restart PostgreSQL
make stop-docker
make start-docker
```

### Port Already in Use

**Problem**: Port 8065 already in use
```powershell
# Find what's using the port
netstat -ano | findstr :8065

# Kill the process or change Mattermost port in config
```

## Configuration

### Default Configuration
Located at: `server/config/config.json`

Key settings:
- **Site URL**: http://localhost:8065
- **Database**: PostgreSQL at localhost:5432
- **File Storage**: MinIO at localhost:9000
- **Email**: Inbucket (all emails caught locally)

### Environment Variables

You can override settings with environment variables:

```powershell
# Example: Change database
$env:MM_SQLSETTINGS_DATASOURCE = "postgres://user:pass@localhost/dbname?sslmode=disable"

# Example: Enable local mode
$env:MM_SERVICESETTINGS_ENABLELOCALMODE = "true"
```

### Using mmctl (CLI Tool)

```powershell
# After server is running
cd server

# Use mmctl with local mode (no login required)
bin/mmctl user create --email test@example.com --username testuser --password Password1! --local

# List users
bin/mmctl user list --local

# See all commands
bin/mmctl --help
```

## Advanced Development

### Using pgvector (for AI features)

```powershell
cd server
make run-pgvector
```

### Running in HA (High Availability) Mode

```powershell
cd server
make run-haserver
```

This starts a 3-node cluster with HAProxy.

### Debug Mode

```powershell
cd server
make debug-server
# Debugger listens on port 2345
```

### Adding Sample Data

```powershell
cd server
make test-data
```

Default accounts:
- **Admin**: sysadmin / Sys@dmin-sample1
- **User**: user-1 / SampleUs@r-1

## Stopping Everything

```powershell
cd server

# Stop server and webapp
make stop

# Also stop Docker services
make stop-docker
```

## Next Steps

Once you have Mattermost running:
1. Explore the codebase structure
2. Read the developer documentation: https://developers.mattermost.com/
3. Try making a small change
4. Run tests to verify your changes
5. Check out the plugin development guide for extensibility

## Getting Help

- **Developer Documentation**: https://developers.mattermost.com/
- **API Documentation**: https://api.mattermost.com/
- **Community Server**: https://community.mattermost.com/
- **GitHub Issues**: https://github.com/mattermost/mattermost-server/issues
- **Developer Discussion**: Join ~Contributors channel on community server

