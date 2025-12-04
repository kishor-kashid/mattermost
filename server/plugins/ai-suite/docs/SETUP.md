# AI Productivity Suite – Local Setup

This guide walks through building, installing, and verifying the Mattermost AI Productivity Suite plugin during PR&nbsp;1.

## 1. Prerequisites

- Go **1.24.6** or later (`go version`)
- Node.js **18.10.0** or later (`node --version`)
- NPM **9.x** or **10.x**
- GNU Make
- A running Mattermost developer server (`cd server && make run`)

> 💡 On Windows, run the plugin commands from **Git Bash** or **WSL** since the Makefile expects a POSIX shell.

## 2. Environment Variables

Set `MM_SERVER_PATH` so the `make deploy` target knows where to place the bundle:

```bash
export MM_SERVER_PATH="$REPO_ROOT/server"
```

You can also create a `.env` file alongside the Makefile that exports the same variable for convenience.

## 3. Build Steps

```bash
cd server/plugins/ai-suite

# 1. Build Go binaries for all major platforms
make server

# 2. Build the React/TypeScript webapp bundle
make webapp

# 3. Create a distributable .tar.gz bundle
make bundle
```

Each command is idempotent and can be re-run as you iterate.

## 4. Deploy to Local Mattermost

```bash
export MM_SERVER_PATH="$REPO_ROOT/server"
make deploy
```

Then open **System Console → Plugins → Plugin Management** and enable **AI Productivity Suite**.

## 5. Watch Mode (Optional)

During UI development, keep webpack running:

```bash
make watch
```

Webpack runs in development mode with inline source maps and automatically rebuilds on save.

## 6. Verification Checklist

After deployment, confirm the scaffold works:

1. Plugin appears in **System Console** with ID `com.mattermost.ai-suite`.
2. Enabling/disabling the plugin succeeds without server errors.
3. The placeholder channel header button displays “AI Suite” and triggers an alert when clicked.
4. Server logs show:
   ```
   AI Productivity Suite plugin activated version=0.1.0
   ```

If any step fails, re-run `make bundle` and review the Mattermost server logs for plugin-specific errors.

