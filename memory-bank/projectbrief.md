# Project Brief: Mattermost

## Overview
Mattermost is an open-core, self-hosted collaboration platform that offers chat, workflow automation, voice calling, screen sharing, and AI integration. This repository is the primary source for core development on the Mattermost platform.

## Current Project: AI Productivity Suite (Native Feature Integration)

We are developing **AI Productivity Suite** features integrated directly into the Mattermost core codebase to enhance team productivity and communication efficiency. This demonstrates **brownfield development** - extending an existing large-scale production codebase rather than building from scratch.

## Core Technologies
- **Backend**: Go (server-side)
- **Frontend**: React + TypeScript (web application)
- **Database**: PostgreSQL (primary)
- **Architecture**: Monorepo structure with separate server and webapp directories

## Repository Structure
```
mattermost/
├── server/          # Go backend server
├── webapp/          # React/TypeScript frontend
├── api/             # API documentation
├── e2e-tests/       # End-to-end testing (Cypress & Playwright)
└── tools/           # Build and development tools
```

## Key Features
- Real-time messaging and collaboration
- Team and channel management
- File sharing and storage (MinIO integration)
- Plugin system for extensibility
- Enterprise features (LDAP, SAML, clustering)
- Search capabilities (Elasticsearch/OpenSearch)
- Mobile and desktop applications

## Project Goals
- Provide a self-hosted alternative to proprietary team collaboration tools
- Maintain high security and compliance standards
- Support extensibility through plugins and APIs
- Enable deployment at scale with high availability

## Build Information
- Released monthly on the 16th
- MIT license for core features
- Enterprise edition available with additional features
- Go version: 1.24.6
- Node version: >=18.10.0
- NPM version: ^9.0.0 || ^10.0.0

## AI Productivity Suite Development

### Core Features (4 Total)
1. **AI Message Summarization** - GPT-powered summaries of threads and channels (max 500 messages, configurable)
2. **Channel Analytics Dashboard** - Visual insights into communication patterns and metrics
3. **Action Item Extractor** - Auto-detect tasks and commitments, track with reminders
4. **Message Formatting Assistant** - AI-powered grammar, tone, and structure improvements

### Technology Stack (Native Integration)
- **Backend**: Go (Mattermost core: api4, app, store layers)
- **Frontend**: React + TypeScript + Redux (Mattermost channels webapp)
- **Database**: PostgreSQL (new AI tables: AIActionItems, AISummaries, AIAnalytics, AIPreferences)
- **AI/LLM**: OpenAI GPT-4 / GPT-3.5-turbo API
- **Background Jobs**: Native Mattermost jobs framework
- **Build**: Standard Mattermost build system (Make + Webpack)

### Development Approach
- **Type**: Brownfield Development (extending existing codebase)
- **Integration**: Native Mattermost core features, not plugins
- **Code Location**: 
  - Backend: `server/channels/app/ai_*.go`, `server/channels/api4/ai_*.go`
  - Frontend: `webapp/channels/src/components/ai/`, `webapp/channels/src/actions/ai_*.ts`
- **Database**: Migrations for new AI tables with proper indexes

### Timeline
- **Total Duration**: 6-7 days
- **Scope**: 87 tasks across 7 PRs
- **Status**: PR #1 Complete (Infrastructure) - December 4, 2024
- **Progress**: 1 of 7 PRs completed (14% complete)

