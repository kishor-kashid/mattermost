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

## AI Message Summarization (PR #3)

PR #3 delivers the full end-to-end summarizer feature set:

- **Thread summaries** via `/summarize thread` or the new post dropdown action “Summarize Thread”. Results appear in a dedicated RHS panel with copy/share/regenerate controls.
- **Channel summaries** via `/summarize channel 24h|3d|7d|custom` or the channel header action “Summarize Channel”, including quick range chips and a custom date/time modal.
- **Backend services** including `server/summarizer` (message fetching, OpenAI prompt/response handling, 24h KV cache) and a `POST /plugins/com.mattermost.ai-suite/api/v1/summarize` endpoint.
- **Webapp UI** that registers a right-hand sidebar, post dropdown item, and channel header menu integration, plus React hooks for fetching/caching summaries.

### Verifying locally

1. `GOWORK=off make bundle && GOWORK=off make deploy` from `server/plugins/ai-suite`.
2. Set `MM_SERVER_PATH` to your `server/` directory before running `make deploy`.
3. Reload the plugin in **System Console → Plugins → Management** after deployment, then hard-refresh the browser.
4. Run `/summarize thread` inside a thread, or use the new UI actions, and confirm the RHS panel renders summaries with metadata and actions.

