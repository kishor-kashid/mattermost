# Mattermost AI Productivity Suite Plugin

The AI Productivity Suite brings AI-assisted summarization, channel analytics, action item extraction, and message formatting workflows to Mattermost. PR&nbsp;1 establishes the plugin scaffold so we can iterate feature-by-feature in subsequent pull requests.

## Project Layout

```
server/plugins/ai-suite/
├── plugin.json         # Plugin manifest
├── Makefile            # Build, bundle, deploy helpers
├── server/             # Go backend implementation
│   ├── plugin.go       # Mattermost plugin entry point
│   ├── configuration.go
│   └── constants.go
├── webapp/             # React/TypeScript web bundle
│   ├── src/
│   ├── tsconfig.json
│   └── webpack.config.js
└── docs/
    └── SETUP.md        # Development workflow documentation
```

## Requirements

- Go 1.24.6+
- Node.js 18.10+
- NPM 9+
- GNU Make
- A local Mattermost server (via `server/Makefile` in the monorepo)

## Common Commands

Run all commands from `server/plugins/ai-suite`.

| Command | Description |
| --- | --- |
| `make server` | Cross-compile the Go plugin for Linux, macOS, and Windows |
| `make webapp` | Install dependencies and produce the webapp bundle |
| `make bundle` | Build everything and emit `dist/ai-suite-<version>.tar.gz` |
| `make deploy` | Unpack the bundle into `$(MM_SERVER_PATH)/plugins` |
| `make watch` | Run the webapp build in watch mode for faster iteration |
| `make clean` | Remove build artifacts and installed node modules |

## Development Flow

1. Ensure the Mattermost server is running from the repository root (`cd server && make run`).
2. In another terminal, build the plugin bundle:
   ```bash
   cd server/plugins/ai-suite
   make bundle
   ```
3. Deploy into your local server (set `MM_SERVER_PATH` to the path of `server`):
   ```bash
   export MM_SERVER_PATH="$(pwd)/../.."
   make deploy
   ```
4. Enable the plugin via **System Console → Plugins → Management**.
5. Use `make watch` to keep the webapp rebuilding while developing UI features.

See `docs/SETUP.md` for deeper instructions, troubleshooting tips, and verification steps.

