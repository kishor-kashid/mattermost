# Project Brief: Mattermost

## Overview
Mattermost is an open-core, self-hosted collaboration platform that offers chat, workflow automation, voice calling, screen sharing, and AI integration. This repository is the primary source for core development on the Mattermost platform.

## Current Project: AI Productivity Suite Plugin

We are developing a **Mattermost AI Productivity Suite** plugin that adds AI-powered features to enhance team productivity and communication efficiency.

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

## Plugin Development Project

### Plugin Features (4 Total)
1. **AI Message Summarization** - GPT-powered summaries of threads and channels (max 500 messages, configurable)
2. **Channel Analytics Dashboard** - Visual insights into communication patterns and metrics
3. **Action Item Extractor** - Auto-detect tasks and commitments, track with reminders
4. **Message Formatting Assistant** - AI-powered grammar, tone, and structure improvements

### Plugin Technology Stack
- **Backend**: Go (Mattermost Plugin SDK)
- **Frontend**: React + TypeScript (Mattermost webapp integration)
- **AI/LLM**: OpenAI GPT-4 / GPT-3.5-turbo API
- **Storage**: Mattermost Plugin Key-Value Store
- **Build**: Make + Webpack

### Timeline
- **Total Duration**: 6-7 days
- **Scope**: 86 tasks across 7 PRs
- **Status**: Planning phase complete, ready to begin development

