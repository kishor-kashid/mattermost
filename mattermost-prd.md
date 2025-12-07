# Product Requirements Document (PRD)

## Mattermost AI Productivity Suite

**Version:** 1.0  
**Date:** December 2, 2024  
**Author:** [Your Name]  
**Status:** Draft

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Project Overview](#2-project-overview)
3. [Goals and Objectives](#3-goals-and-objectives)
4. [Target Users](#4-target-users)
5. [Feature Specifications](#5-feature-specifications)
   - 5.1 [AI Message Summarization](#51-ai-message-summarization)
   - 5.2 [Link & Article Summarizer](#52-link--article-summarizer)
   - 5.3 [Action Item Extractor](#53-action-item-extractor)
   - 5.4 [Message Formatting Assistant](#54-message-formatting-assistant)
6. [Technical Architecture](#6-technical-architecture)
7. [User Interface Design](#7-user-interface-design)
8. [Data Model](#8-data-model)
9. [API Specifications](#9-api-specifications)
10. [Security and Privacy](#10-security-and-privacy)
11. [Performance Requirements](#11-performance-requirements)
12. [Dependencies](#12-dependencies)
13. [Risks and Mitigations](#13-risks-and-mitigations)
14. [Success Metrics](#14-success-metrics)
15. [Timeline and Milestones](#15-timeline-and-milestones)
16. [Future Enhancements](#16-future-enhancements)
17. [Appendix](#17-appendix)

---

## 1. Executive Summary

### 1.1 Purpose

This document outlines the product requirements for the **Mattermost AI Productivity Suite**, a plugin that enhances the Mattermost collaboration platform with AI-powered features designed to reduce information overload, improve team communication efficiency, and provide actionable insights.

### 1.2 Problem Statement

Modern teams using Mattermost face several challenges:

- **Information Overload:** High-volume channels make it difficult to catch up on missed conversations
- **Link Fatigue:** Team members share articles, docs, and links that take time to read; many go unread
- **Lost Action Items:** Tasks and commitments mentioned in conversations get forgotten
- **Communication Quality:** Messages lack proper formatting and clarity

### 1.3 Solution

The Mattermost AI Productivity Suite addresses these challenges through four integrated **native features** added directly to the Mattermost core:

1. **AI Message Summarization** - Intelligent summaries of threads and channels
2. **Link & Article Summarizer** - AI-generated summaries of shared URLs and articles
3. **Action Item Extractor** - Automatic detection and tracking of tasks and commitments
4. **Message Formatting Assistant** - AI-powered message improvement and formatting

### 1.4 Key Benefits

| Benefit | Impact |
|---------|--------|
| Reduced time catching up on channels | 60% faster information absorption |
| Shared links become actionable | 80% faster knowledge extraction from articles |
| Never miss action items | 100% task capture from conversations |
| Professional communication | Improved message clarity and formatting |

---

## 2. Project Overview

### 2.1 Project Name

**Mattermost AI Productivity Suite** (Internal codename: "MattermostAI")

### 2.2 Project Type

Native Mattermost Feature Integration (Server Backend + Web App Frontend)

### 2.3 Technology Stack

| Component | Technology |
|-----------|------------|
| Backend | Go (Golang) - Mattermost Core App Layer |
| Frontend | React, TypeScript, Redux - Mattermost Channels Webapp |
| Database | PostgreSQL (Mattermost's existing DB with new tables) |
| AI/LLM | OpenAI GPT-4 API |
| Build System | Make, Webpack, Go Modules |
| Framework | Native Mattermost Core (api4, app, store layers) |

### 2.4 Repository Structure

```
mattermost/
├── server/                 # Go backend (cmd, channels, enterprise, plugins, etc.)
│   ├── cmd/                # Binary entrypoints (`mattermost`, `platform`, tooling)
│   ├── channels/           # Core server app (api4, app, store, jobs, wsapi, web)
│   │   ├── api4/           # REST API endpoints → AI endpoints here
│   │   ├── app/            # Business logic → AI services here
│   │   ├── store/          # Database layer → AI data models here
│   │   └── jobs/           # Background workers → AI reminders here
│   ├── config/, data/, logs/ # Default configs, sample data, local runtime output
│   └── scripts/, tests/, public/ # Build helpers, QA tooling, static assets
├── webapp/                 # React/TypeScript clients
│   ├── channels/           # Primary Mattermost web client
│   │   └── src/            # Frontend source code
│   │       ├── actions/    # Redux actions → AI actions here
│   │       ├── components/ # React components → AI UI here
│   │       ├── reducers/   # Redux reducers → AI state here
│   │       ├── selectors/  # Redux selectors → AI selectors here
│   │       └── utils/      # Utilities → AI helpers here
│   ├── platform/           # Shared platform packages
│   └── scripts/            # Frontend build helpers
├── api/                    # REST/v4 API reference & OpenAPI specs
├── e2e-tests/              # Cypress/Playwright automation suites
├── memory-bank/            # Project knowledge base for this effort
└── Root docs & configs     # PRD, task list, README, licenses, etc.
```

This is the upstream Mattermost monorepo; our AI Productivity Suite features will be integrated directly into the core codebase (`server/channels` for backend, `webapp/channels/src` for frontend) as native Mattermost functionality.

### 2.5 Forked Repository

**Base Repository:** https://github.com/mattermost/mattermost  
**Stars:** ~32,000  
**License:** MIT (open core)

---

## 3. Goals and Objectives

### 3.1 Primary Goals

| Goal | Description | Success Criteria |
|------|-------------|------------------|
| G1 | Reduce information overload | Users can summarize 50+ messages in <5 seconds |
| G2 | Make shared links actionable | Link summaries generated in <5 seconds with key insights |
| G3 | Capture all action items | Extract 95%+ of commitments from conversations |
| G4 | Improve message quality | 80% of formatted messages rated as clearer |

### 3.2 Secondary Goals

- Demonstrate proficiency in Go and React development
- Learn Mattermost core architecture (api4, app, store, Redux)
- **Practice brownfield development** - extending an existing large codebase
- Build production-ready, deployable software
- Create comprehensive documentation

### 3.3 Non-Goals (Out of Scope)

- Mobile app modifications
- Multi-LLM support (OpenAI only for v1.0)
- Real-time collaborative features
- Integration with external task management tools
- Plugin-based architecture (we're building native features)

---

## 4. Target Users

### 4.1 User Personas

#### Persona 1: Developer Dan

| Attribute | Description |
|-----------|-------------|
| Role | Software Engineer |
| Team Size | 15-person engineering team |
| Pain Points | Misses important messages in high-volume channels, spends 30+ minutes daily catching up |
| Goals | Quickly understand what happened while away, never miss critical updates |
| Technical Comfort | High |

#### Persona 2: Manager Maria

| Attribute | Description |
|-----------|-------------|
| Role | Engineering Manager |
| Team Size | Manages 3 teams (45 people) |
| Pain Points | No visibility into team communication health, cannot identify silos |
| Goals | Understand team dynamics, ensure healthy communication patterns |
| Technical Comfort | Medium |

#### Persona 3: Remote Rachel

| Attribute | Description |
|-----------|-------------|
| Role | Product Designer |
| Team Size | 10-person product team across 4 timezones |
| Pain Points | Messages sent at wrong times, colleagues miss updates |
| Goals | Coordinate effectively across timezones |
| Technical Comfort | Medium |

### 4.2 User Stories

#### AI Message Summarization

| ID | User Story | Priority |
|----|------------|----------|
| US-1.1 | As a user, I want to summarize a long thread so that I can quickly understand the discussion | P0 |
| US-1.2 | As a user, I want to summarize all messages in a channel from a time range so that I can catch up after being away | P0 |
| US-1.3 | As a user, I want to receive a daily digest of important channel activity so that I stay informed | P1 |
| US-1.4 | As a user, I want to customize the summary length and format so that it fits my preferences | P2 |

#### Link & Article Summarizer

| ID | User Story | Priority |
|----|------------|----------|
| US-2.1 | As a user, I want to see AI summaries of shared links so that I can quickly understand article content | P0 |
| US-2.2 | As a user, I want automatic link detection and summarization so that summaries appear without manual action | P0 |
| US-2.3 | As a user, I want to see key points extracted from articles so that I can decide whether to read the full content | P1 |
| US-2.4 | As a user, I want to manually trigger link summarization so that I can get summaries for older shared links | P1 |
| US-2.5 | As a user, I want link summaries to be cached so that repeated access is instant | P2 |

#### Action Item Extractor

| ID | User Story | Priority |
|----|------------|----------|
| US-3.1 | As a user, I want the system to automatically detect action items in conversations so that nothing gets forgotten | P0 |
| US-3.2 | As a user, I want to see who is assigned to each action item so that accountability is clear | P0 |
| US-3.3 | As a user, I want to view all my action items in a personal dashboard so that I can track what I need to do | P0 |
| US-3.4 | As a user, I want to mark action items as complete so that I can track progress | P0 |
| US-3.5 | As a user, I want to receive reminders for overdue action items so that I don't miss deadlines | P1 |
| US-3.6 | As a manager, I want to see action items for my team so that I can track team commitments | P1 |

#### Message Formatting Assistant

| ID | User Story | Priority |
|----|------------|----------|
| US-4.1 | As a user, I want AI to help format my message professionally so that my communication is clearer | P0 |
| US-4.2 | As a user, I want to convert plain text into proper lists and formatting so that messages are easier to read | P0 |
| US-4.3 | As a user, I want grammar and spelling suggestions so that my messages are error-free | P0 |
| US-4.4 | As a user, I want to make technical messages more concise so that they're easier to understand | P1 |
| US-4.5 | As a user, I want to preview the formatted version before sending so that I can review changes | P1 |

---

## 5. Feature Specifications

### 5.1 AI Message Summarization

#### 5.1.1 Overview

AI Message Summarization uses OpenAI's GPT-4 to generate concise, accurate summaries of Mattermost conversations. Users can summarize individual threads, entire channels over a time range, or receive automated daily digests.

#### 5.1.2 Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-1.1 | System shall summarize threads with all messages, capped at 500 messages maximum | P0 |
| FR-1.2 | System shall summarize channel messages within a specified time range (1 hour to 30 days), limited to 500 messages or messages in time range, whichever is fewer | P0 |
| FR-1.3 | System shall allow administrators to configure the maximum message limit for summarization (default: 500) | P1 |
| FR-1.4 | System shall provide slash command `/summarize` for on-demand summarization | P0 |
| FR-1.5 | System shall provide right-click context menu option for thread summarization | P0 |
| FR-1.6 | System shall display summaries in the right-hand sidebar (RHS) | P0 |
| FR-1.7 | System shall support summary lengths: brief (1-2 sentences), standard (paragraph), detailed (multiple paragraphs) | P1 |
| FR-1.8 | System shall extract and highlight key decisions made in the conversation | P1 |
| FR-1.9 | System shall identify and list participants mentioned in the summary | P1 |
| FR-1.10 | System shall support scheduled daily digest emails | P2 |
| FR-1.11 | System shall cache summaries for 24 hours to reduce API calls | P1 |

#### 5.1.3 Message Limits

**Default Limits:**

| Summarization Type | Default Behavior | Configurable |
|--------------------|------------------|--------------|
| Thread Summary | All messages in thread (max 500) | Yes |
| Channel Summary | Last 500 messages OR messages in time range, whichever is fewer | Yes |

**Configuration:**
- System administrators can configure the maximum message limit in plugin settings
- Default limit: 500 messages
- Recommended range: 100-1000 messages
- Higher limits increase API costs and response time

**Behavior:**
- If a thread has 300 messages, all 300 are summarized
- If a thread has 700 messages, only the first 500 are summarized
- If a channel time range contains 2000 messages, only the most recent 500 are summarized
- If a channel time range contains 200 messages, all 200 are summarized
- Users receive a notice if the message limit was reached: "⚠️ Summary limited to 500 most recent messages"

#### 5.1.4 Slash Command Specification

**Command:** `/summarize`

**Syntax:**
```
/summarize thread
/summarize channel [time-range]
/summarize channel today
/summarize channel 7d
/summarize channel 2024-11-01 2024-11-30
```

**Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| target | enum | Yes | - | `thread` or `channel` |
| time-range | string | No | 24h | Time range for channel summary |

**Response Format:**
```
📋 **Summary** (42 messages, 8 participants)

**Key Points:**
• The team decided to postpone the v2.0 release to December 15
• @john identified a critical bug in the authentication module
• @sarah will lead the bug fix effort with support from @mike

**Decisions Made:**
1. Release postponed to Dec 15
2. Bug fix prioritized over new features

**Action Items:**
• @john: Document the bug in Jira
• @sarah: Create fix branch by EOD
```

#### 5.1.5 User Interface

**Thread Summarization:**
- Right-click context menu on thread root post → "Summarize Thread"
- Summary appears in RHS panel
- Options: Copy, Share, Regenerate

**Channel Summarization:**
- Channel header dropdown → "Summarize Channel"
- Date range picker modal
- Summary appears in RHS panel

**Daily Digest:**
- Settings: Enable/disable, delivery time, channels to include
- Delivered as DM from bot account

#### 5.1.6 AI Prompt Engineering

**System Prompt:**
```
You are a professional workplace communication summarizer. Your task is to 
create concise, accurate summaries of team conversations from Mattermost.

Guidelines:
- Focus on decisions, action items, and key information
- Use bullet points for clarity
- Mention specific people when they have action items
- Highlight any deadlines or time-sensitive information
- Maintain professional tone
- Do not include speculation or information not in the messages
- If the conversation is unclear or lacks substance, say so
```

**User Prompt Template:**
```
Summarize the following conversation from the #{channel_name} channel.
Time range: {start_time} to {end_time}
Number of messages: {message_count}
Participants: {participant_list}

Provide:
1. A brief overview (2-3 sentences)
2. Key points as bullet points
3. Any decisions made
4. Action items with assignees (if mentioned)

Messages:
{formatted_messages}
```

---

### 5.2 Link & Article Summarizer

#### 5.2.1 Overview

The Link & Article Summarizer uses AI to generate concise summaries of URLs shared in Mattermost. When team members share articles, documentation, blog posts, or other web content, the system automatically extracts key information, enabling users to quickly understand the content without leaving the chat.

#### 5.2.2 Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-2.1 | System shall detect URLs in posted messages automatically | P0 |
| FR-2.2 | System shall fetch and parse web page content from detected URLs | P0 |
| FR-2.3 | System shall generate AI summaries of fetched content | P0 |
| FR-2.4 | System shall display summaries as rich preview cards below messages | P0 |
| FR-2.5 | System shall extract and display key points from articles | P1 |
| FR-2.6 | System shall estimate reading time for original content | P1 |
| FR-2.7 | System shall cache link summaries for 7 days | P1 |
| FR-2.8 | System shall support manual "Summarize Link" action for any URL | P1 |
| FR-2.9 | System shall handle various content types (articles, docs, GitHub READMEs) | P1 |
| FR-2.10 | System shall respect robots.txt and rate limit external fetches | P0 |
| FR-2.11 | System shall allow users to expand/collapse link summaries | P2 |

#### 5.2.3 Supported Content Types

| Content Type | Examples | Extraction Method |
|--------------|----------|-------------------|
| Web Articles | News sites, blogs, Medium | HTML parsing, main content extraction |
| Documentation | Confluence, Notion, GitBook | HTML parsing, structured content |
| GitHub READMEs | Repository README.md files | Markdown parsing via GitHub API |
| PDF Documents | Linked PDF files | PDF text extraction (first 10 pages) |
| Stack Overflow | Q&A threads | Question + accepted answer extraction |

#### 5.2.4 Link Summary Card Layout

```
┌────────────────────────────────────────────────────────────────┐
│  @bob shared a link:                                           │
│  ─────────────────────────────────────────────────────────────  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ 📄 "Microservices vs Monoliths in 2024"                  │  │
│  │    martinfowler.com · ⏱️ 12 min read                      │  │
│  │ ───────────────────────────────────────────────────────  │  │
│  │ 🤖 AI Summary:                                           │  │
│  │                                                          │  │
│  │ This article compares microservices and monolithic       │  │
│  │ architectures for modern applications. Key insights:     │  │
│  │                                                          │  │
│  │ • Microservices suit large teams with clear domain       │  │
│  │   boundaries and high deployment frequency               │  │
│  │ • Monoliths are better for startups and small teams      │  │
│  │   (<10 developers) due to lower operational overhead     │  │
│  │ • The key decision factor is deployment frequency and    │  │
│  │   team autonomy requirements                             │  │
│  │                                                          │  │
│  │ [Read Full Article ↗] [Copy Summary] [▼ Collapse]        │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
```

#### 5.2.5 User Interface

**Automatic Summarization:**
- Links detected in messages trigger automatic summarization
- Summary card appears below the message (expandable/collapsible)
- Loading state shown while fetching and summarizing

**Manual Summarization:**
- Click "Summarize" on any link preview
- Right-click link → "Summarize with AI"
- Hover over link → Click summarize icon

**Summary Card Features:**
- Site favicon and domain name
- Article title
- Estimated reading time
- AI-generated summary (3-5 key points)
- "Read Full Article" button (opens in new tab)
- Copy summary option
- Expand/collapse toggle

**User Preferences:**
- Enable/disable auto-summarization
- Default expanded/collapsed state
- Summary length preference (brief/standard)

#### 5.2.6 AI Prompt Engineering

**System Prompt:**
```
You are a professional content summarizer. Your task is to create concise, 
accurate summaries of web articles and documents shared in workplace chat.

Guidelines:
- Extract the main thesis or purpose of the content
- Identify 3-5 key points or takeaways
- Use bullet points for clarity
- Maintain factual accuracy - do not add information not in the source
- Keep summaries professional and neutral in tone
- If the content is technical, preserve important technical details
- If the content is unclear or too short, say so
- Do not include personal opinions or editorializing
```

**User Prompt Template:**
```
Summarize the following article/document.

Title: {title}
Source: {domain}
Content length: {word_count} words

Provide:
1. A brief overview (1-2 sentences)
2. 3-5 key points as bullet points
3. Any important conclusions or recommendations

Content:
{extracted_content}
```

#### 5.2.7 Content Extraction

**Extraction Pipeline:**
1. **URL Detection**: Regex pattern matching in message content
2. **Fetch**: HTTP GET with appropriate User-Agent, timeout, and size limits
3. **Parse**: HTML parsing to extract main content (removing ads, navigation, etc.)
4. **Clean**: Remove HTML tags, normalize whitespace, limit to 10,000 characters
5. **Summarize**: Send to OpenAI for AI summary generation
6. **Cache**: Store summary in AILinkSummaries table with 7-day TTL

**Extraction Libraries:**
- `golang.org/x/net/html` for HTML parsing
- `github.com/PuerkitoBio/goquery` for content extraction
- Readability algorithm for main content detection

#### 5.2.8 Error Handling

| Error | User Message | Action |
|-------|--------------|--------|
| URL unreachable | "Unable to fetch this link" | Show basic link preview only |
| Content too short | "Not enough content to summarize" | Show link title only |
| Paywall detected | "This content is behind a paywall" | Show link preview only |
| Rate limited | "Summarization temporarily unavailable" | Queue for later |
| Unsupported format | "This content type is not supported" | Show link preview only |

---

### 5.3 Action Item Extractor

#### 5.3.1 Overview

The Action Item Extractor automatically scans channels and threads for commitments, tasks, and action items mentioned in conversations. It uses AI to detect who promised to do what and by when, creating a centralized dashboard to track all commitments and prevent tasks from being forgotten.

#### 5.3.2 Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-3.1 | System shall automatically detect action items in messages using AI analysis | P0 |
| FR-3.2 | System shall extract assignee information (who will do it) from action items | P0 |
| FR-3.3 | System shall extract deadline information (when it's due) from action items when mentioned | P1 |
| FR-3.4 | System shall provide a personal action item dashboard for each user | P0 |
| FR-3.5 | System shall allow users to mark action items as complete | P0 |
| FR-3.6 | System shall allow users to edit extracted action items (assignee, deadline, description) | P1 |
| FR-3.7 | System shall send reminders for action items approaching their deadline | P1 |
| FR-3.8 | System shall send reminders for overdue action items | P1 |
| FR-3.9 | System shall link action items back to the original message/thread | P0 |
| FR-3.10 | System shall provide a team view for managers to see all team action items | P1 |
| FR-3.11 | System shall support slash command `/actionitems` to view and manage items | P0 |
| FR-3.12 | System shall allow manual creation of action items | P2 |

#### 5.3.3 Action Item Detection

**AI Detection Patterns:**

The system detects action items through pattern recognition and semantic analysis:

| Pattern Type | Examples | Detected Fields |
|--------------|----------|-----------------|
| Explicit assignment | "@john will update the docs by Friday" | Assignee: john, Deadline: Friday |
| Commitment | "I'll review the PR tomorrow" | Assignee: speaker, Deadline: tomorrow |
| Request | "@sarah can you send the report?" | Assignee: sarah |
| Group task | "Team needs to test this before launch" | Assignee: team/channel |
| Deadline mention | "This needs to be done by EOD" | Deadline: EOD |
| Multiple tasks | "I'll do X and Y, then Z" | Multiple items |

**Extraction Process:**
1. Message is posted in a channel
2. AI analyzes message content for commitment language
3. Extracts: task description, assignee(s), deadline (if mentioned), priority indicators
4. Creates action item entry linked to original message
5. Notifies assigned user(s)
6. Adds to personal dashboard

#### 5.3.4 Personal Dashboard

**Action Items Dashboard:**

```
┌─────────────────────────────────────────────────────────────────┐
│  My Action Items                      [All ▼] [Filter by: All]  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  🔴 Overdue (2)                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ ☐ Update documentation for API v2                         │ │
│  │    From: #engineering · @mike                             │ │
│  │    Due: Dec 1 (2 days ago) · [View Message] [Mark Done]   │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                  │
│  🟡 Due Soon (3)                                                │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ ☐ Review Sarah's PR #234                                  │ │
│  │    From: #code-review · @sarah                            │ │
│  │    Due: Tomorrow · [View Message] [Mark Done]             │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ⚪ No Deadline (5)                                             │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ ☐ Schedule team offsite meeting                           │ │
│  │    From: #general · @manager                              │ │
│  │    Created: Dec 2 · [View Message] [Mark Done] [Set Due]  │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ✅ Completed (12) [Show]                                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Dashboard Features:**
- Filter by status: All, Overdue, Due Today, This Week, No Deadline, Completed
- Filter by channel
- Sort by due date, created date, priority
- Bulk actions: mark multiple as done
- Search action items

#### 5.3.5 Team View (For Managers)

**Team Action Items View:**

Shows all action items assigned to team members with:
- Grouped by assignee
- Filter by team member
- Overview of team workload
- Identify blocked or overdue items
- Export capability

#### 5.3.6 Slash Command Specification

**Command:** `/actionitems`

**Syntax:**
```
/actionitems                     # Open personal dashboard
/actionitems list                # List all active items
/actionitems mine                # Show only my items
/actionitems team                # Show team items (managers)
/actionitems create [description] @[assignee] [deadline]
/actionitems complete [id]       # Mark item as complete
```

**Response Examples:**

List:
```
📋 Your Action Items (5 active)

🔴 Overdue:
1. Update API docs (Due: Dec 1) - #engineering
   [Mark Done] [View Message]

🟡 Due Soon:
2. Review PR #234 (Due: Tomorrow) - #code-review
3. Send weekly report (Due: Friday) - #management

⚪ No Deadline:
4. Schedule team meeting - #general
5. Update roadmap - #planning

[View Full Dashboard]
```

#### 5.3.7 Reminders

**Reminder Schedule:**

| Type | Trigger | Frequency |
|------|---------|-----------|
| Due Today | 9:00 AM on due date | Once |
| Due Tomorrow | 5:00 PM day before | Once |
| Overdue | 9:00 AM daily | Daily until completed |
| Weekly Summary | Configurable (default: Monday 9 AM) | Weekly |

**Reminder Format:**
```
⏰ Action Item Reminder

You have 3 action items due today:
1. Update API documentation - #engineering
2. Review budget proposal - @manager
3. Submit timesheet - #admin

[View All Action Items]
```

#### 5.3.8 Integration Points

- **Channel Messages**: Auto-detection in all channels
- **Thread Replies**: Detect commitments in thread conversations
- **Direct Messages**: Optional (can be enabled/disabled per user)
- **Mentions**: Special attention to messages that @mention users
- **Meeting Notes**: Extract action items from meeting summaries

---

### 5.4 Message Formatting Assistant

#### 5.4.1 Overview

The Message Formatting Assistant uses AI to help users improve their message quality by automatically formatting text, fixing grammar and spelling, making messages more professional, and enhancing clarity. Users can transform plain text into well-structured, properly formatted messages with a single click.

#### 5.4.2 Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-4.1 | System shall provide AI-powered message formatting suggestions | P0 |
| FR-4.2 | System shall convert plain text to structured formats (lists, headers, code blocks) | P0 |
| FR-4.3 | System shall provide grammar and spelling correction | P0 |
| FR-4.4 | System shall offer professional tone enhancement | P1 |
| FR-4.5 | System shall provide message conciseness suggestions | P1 |
| FR-4.6 | System shall show preview before applying changes | P0 |
| FR-4.7 | System shall support one-click formatting actions from message box | P0 |
| FR-4.8 | System shall maintain user's original meaning while formatting | P0 |
| FR-4.9 | System shall support technical content (code, terminal commands) | P1 |
| FR-4.10 | System shall allow custom formatting instructions | P2 |
| FR-4.11 | System shall learn user's preferred writing style | P2 |

#### 5.4.3 Formatting Actions

**Available Formatting Operations:**

| Action | Description | Example Input | Example Output |
|--------|-------------|---------------|----------------|
| Make Professional | Improve tone | "hey can u check this?" | "Hi team, could you please review this when you have a moment?" |
| Format as List | Convert to bullets | "we need to do A and B and C" | "We need to:
• A
• B
• C" |
| Fix Grammar | Correct errors | "Their going to the meeting" | "They're going to the meeting" |
| Make Concise | Reduce wordiness | "In my opinion, I think that we should probably..." | "We should..." |
| Add Code Blocks | Format code properly | "run npm install then start" | "\`\`\`bash
npm install
npm start
\`\`\`" |
| Improve Clarity | Simplify language | Complex technical explanation | Clear, understandable version |

#### 5.4.4 User Interface

**Message Composer Integration:**

```
┌─────────────────────────────────────────────────────────────────┐
│  Message to #engineering                                    🤖  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  hey team i think we should use postgres instead of mysql       │
│  because its better for our usecase and also we need to         │
│  consider the licensing                                         │
│                                                                  │
│  ────────────────────────────────────────────────────────────  │
│  ✨ AI Suggestions:                                            │
│  • 2 grammar issues detected                                   │
│  • Convert to professional tone                                │
│  • Structure as arguments                                      │
│                                                                  │
│  [✓ Make Professional] [✓ Format as List] [Fix Grammar]       │
│                                                                  │
├─────────────────────────────────────────────────────────────────┤
│  [Bold] [Italic] [Code] [Link] | [🎨 AI Format ▼] [@] [emoji] │
└─────────────────────────────────────────────────────────────────┘
```

**AI Format Dropdown Menu:**

```
🎨 AI Format
├── ✨ Make Professional
├── 📝 Format as List
├── ✔️ Fix Grammar & Spelling
├── 📉 Make Concise
├── 💻 Add Code Formatting
├── 📊 Improve Clarity
└── ⚙️ Custom Instruction...
```

**Preview Modal:**

```
┌─────────────────────────────────────────────────────────────────┐
│  AI Formatting Preview                              [✕ Cancel]  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Original:                                                       │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ hey team i think we should use postgres instead of mysql  │ │
│  │ because its better for our usecase and also we need to    │ │
│  │ consider the licensing                                    │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                  │
│  Formatted (Professional + List):                              │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ Hi team,                                                  │ │
│  │                                                           │ │
│  │ I recommend we use PostgreSQL instead of MySQL for       │ │
│  │ the following reasons:                                    │ │
│  │                                                           │ │
│  │ • Better suited for our use case                         │ │
│  │ • Licensing considerations favor PostgreSQL              │ │
│  │                                                           │ │
│  │ Would appreciate your thoughts.                           │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                  │
│  Changes:                                                        │
│  • Fixed: "its" → "it's" (not applicable in rewrite)            │
│  • Improved: Professional tone                                  │
│  • Structured: Bullet points for clarity                        │
│                                                                  │
│           [Use Original] [Apply Formatting]                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### 5.4.5 Formatting Profiles

**Quick Profiles:**

Users can save formatting preferences:

| Profile | Settings | Use Case |
|---------|----------|----------|
| Professional | Formal tone, proper grammar, structured | Client communication |
| Casual | Friendly tone, simple fixes | Internal team chat |
| Technical | Preserve technical terms, add code formatting | Development discussions |
| Concise | Brevity prioritized, remove fluff | Quick updates |

#### 5.4.6 AI Prompt Engineering

**System Prompt for Formatting:**
```
You are a professional communication assistant. Your task is to improve 
workplace messages while preserving the author's original intent and meaning.

Guidelines:
- Maintain the core message and intent
- Fix grammar and spelling errors
- Improve clarity and structure
- Adjust tone as requested (professional, concise, etc.)
- Preserve technical terms and proper nouns
- Never change factual content or commitments
- If adding structure, use appropriate Markdown
- Keep formatting compatible with Mattermost
```

**Action-Specific Prompts:**

*Make Professional:*
```
Transform this message to a professional workplace tone while keeping 
the same meaning. Fix grammar issues and improve structure.
```

*Format as List:*
```
Convert this message into a well-structured list using bullet points or 
numbered items where appropriate. Maintain all key information.
```

*Make Concise:*
```
Rewrite this message to be more concise while preserving all important 
information and meaning. Remove unnecessary words and redundancy.
```

#### 5.4.7 User Settings

**Formatting Preferences:**

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| Auto-suggest formatting | Boolean | true | Show suggestions while typing |
| Default profile | Enum | Professional | Default formatting style |
| Show preview before applying | Boolean | true | Always show preview modal |
| Grammar check level | Enum | Standard | Basic, Standard, Strict |
| Preserve technical terms | Boolean | true | Don't change code/technical words |

---


## 6. Technical Architecture

### 6.1 System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Mattermost Server                            │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │              server/channels/ (Core Application)                │ │
│  │                                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────┐ │ │
│  │  │                    api4/ (REST API Layer)                  │ │ │
│  │  │  POST /api/v4/ai/summarize  |  GET /api/v4/ai/analytics  │ │ │
│  │  │  POST /api/v4/ai/format  |  GET /api/v4/ai/actionitems   │ │ │
│  │  └────────────────────────────┬──────────────────────────────┘ │ │
│  │                               │                                 │ │
│  │  ┌────────────────────────────┴──────────────────────────────┐ │ │
│  │  │                 app/ (Business Logic Layer)                │ │ │
│  │  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────┐│ │ │
│  │  │  │ Summarizer  │ │    Link     │ │ Action Item │ │Format││ │ │
│  │  │  │   Service   │ │ Summarizer  │ │  Extractor  │ │ Svc  ││ │ │
│  │  │  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └───┬──┘│ │ │
│  │  │         │               │               │             │    │ │ │
│  │  │  ┌──────┴───────────────┴───────────────┴─────────────┴──┐│ │ │
│  │  │  │              Shared AI Services Layer                 ││ │ │
│  │  │  │  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐││ │ │
│  │  │  │  │  OpenAI      │  │   Message    │  │ Notification│││ │ │
│  │  │  │  │  Client      │  │   Processor  │  │  Manager    │││ │ │
│  │  │  │  └──────────────┘  └──────────────┘  └─────────────┘││ │ │
│  │  │  └───────────────────────────────────────────────────────┘│ │ │
│  │  └────────────────────────────┬──────────────────────────────┘ │ │
│  │                               │                                 │ │
│  │  ┌────────────────────────────┴──────────────────────────────┐ │ │
│  │  │               store/ (Data Access Layer)                   │ │ │
│  │  │  ┌───────────────┐  ┌──────────────┐  ┌────────────────┐ │ │ │
│  │  │  │  AI Summaries │  │ Action Items │  │ Link Summaries │ │ │ │
│  │  │  │  Store        │  │  Store       │  │  Store         │ │ │ │
│  │  │  └───────────────┘  └──────────────┘  └────────────────┘ │ │ │
│  │  └───────────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────────────┘ │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│                         PostgreSQL Database                          │
│  ┌────────────┐ ┌────────────┐ ┌───────────────┐ ┌───────────────┐ │
│  │   Posts    │ │  Channels  │ │ AISummaries   │ │ AIActionItems ││
│  │   Users    │ │ Reactions  │ │AILinkSummaries│ │ AIPreferences ││
│  └────────────┘ └────────────┘ └───────────────┘ └───────────────┘ │
└──────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │    OpenAI API       │
                         │  (GPT-4 / GPT-3.5)  │
                         └─────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                    webapp/channels/ (Frontend)                       │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                    src/components/ (UI Layer)                   │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐│ │
│  │  │  Summary    │ │  Analytics  │ │ Action Item │ │ Formatter ││ │
│  │  │  Panel      │ │  Dashboard  │ │  Dashboard  │ │  UI       ││ │
│  │  └─────────────┘ └─────────────┘ └─────────────┘ └───────────┘│ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │              src/actions/ (Redux Action Creators)               │ │
│  │  summarizeThread() | summarizeLink() | createActionItem()       │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │            src/reducers/ (Redux State Management)               │ │
│  │  aiSummaries | aiLinkSummaries | aiActionItems | aiFormatter   │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │               src/selectors/ (Data Selectors)                   │ │
│  │  getAISummary() | getLinkSummary() | getActionItems()          │ │
│  └────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────┘
```

### 6.2 Component Descriptions

**Backend (server/channels/)**

| Component | Layer | Responsibility |
|-----------|-------|----------------|
| API Handlers (`api4/ai_*.go`) | REST API | HTTP request handling, validation, permissions, response formatting |
| Summarizer Service (`app/ai_summarizer.go`) | Business Logic | Message retrieval, formatting, OpenAI calls, summary generation |
| Link Summarizer (`app/ai_link_summarizer.go`) | Business Logic | URL detection, content fetching, article summarization |
| Action Item Extractor (`app/ai_action_items.go`) | Business Logic | Detects commitments and tasks, manages action item lifecycle |
| Formatter Service (`app/ai_formatter.go`) | Business Logic | AI-powered formatting and grammar checking |
| OpenAI Client (`app/openai_client.go`) | Shared Service | Wrapper for OpenAI API with rate limiting, retries, error handling |
| AI Stores (`store/sqlstore/ai_*.go`) | Data Layer | Database operations for summaries, link summaries, action items, preferences |

**Frontend (webapp/channels/src/)**

| Component | Layer | Responsibility |
|-----------|-------|----------------|
| AI Components (`components/ai/`) | UI | React components for summaries, link summaries, action items, formatting |
| AI Actions (`actions/ai_*.ts`) | Redux Actions | Action creators for API calls and state updates |
| AI Reducers (`reducers/ai_*.ts`) | Redux State | State management for AI features |
| AI Selectors (`selectors/ai_*.ts`) | Redux Selectors | Memoized data selectors and computed state |

### 6.3 Integration Points

| Integration Point | Purpose |
|-------------------|---------|
| `app.MessageHasBeenPosted()` | Trigger analytics data collection and action item extraction |
| `app.MessageWillBePosted()` | Apply formatting assistance if user requested |
| `app.InitializeAIServices()` | Initialize AI services at server startup |
| `api4.InitAI()` | Register AI REST API endpoints |
| Redux Store Integration | Wire AI reducers and middleware into main store |

### 6.4 Data Flow Diagrams

**AI Summarization Flow:**
```
User Request → Validate Permissions → Fetch Messages → Format for LLM
     → Call OpenAI API → Parse Response → Cache Result → Return to User
```

**Link Summarization Flow:**
```
User Posts Link → Detect URL → Fetch Content → Extract Main Content
     → Send to OpenAI → Generate Summary → Cache in DB → Display Card
```

**Action Item Extraction Flow:**
```
Message Posted → AI Analyzes for Commitments → Extract Task Details
     → Identify Assignee & Deadline → Create Action Item → Notify User
     → Add to Dashboard → Schedule Reminders
```

**Message Formatting Flow:**
```
User Writes Message → Requests Formatting → Send to AI → Apply Rules
     → Generate Formatted Version → Show Preview → User Approves
     → Replace Original → Post Message
```

---

## 7. User Interface Design

### 7.1 Design Principles

| Principle | Description |
|-----------|-------------|
| Non-intrusive | Features enhance, not replace, native Mattermost UI |
| Consistent | Follow Mattermost design patterns and component library |
| Accessible | Support keyboard navigation and screen readers |
| Responsive | Work on desktop and tablet viewports |
| Fast | UI interactions respond within 100ms |

### 7.2 UI Components

#### 7.2.1 Right-Hand Sidebar (RHS) Panel

Used for: Summaries, Notification Center

```
┌────────────────────────────────┐
│  📋 Summary         [✕ Close] │
├────────────────────────────────┤
│                                │
│  Channel: #engineering         │
│  Period: Last 24 hours         │
│  Messages: 47                  │
│                                │
│  ─────────────────────────────│
│                                │
│  **Overview**                  │
│  The team discussed the...     │
│                                │
│  **Key Points**                │
│  • Point one                   │
│  • Point two                   │
│                                │
│  **Decisions**                 │
│  1. Decision one               │
│                                │
│  ─────────────────────────────│
│                                │
│  [Copy] [Share] [Regenerate]   │
│                                │
└────────────────────────────────┘
```

#### 7.2.2 Modal Dialogs

Used for: Action Item creation/editing, Formatting preview/confirm, Date Range Picker

#### 7.2.3 Channel Header Button

Dropdown menu with plugin actions:
- Summarize Channel
- Summarize Link
- Manage Action Items

#### 7.2.4 Main Menu Items

Under plugin section:
- Action Items Dashboard
- Formatting Preferences
- Link Summary Settings

### 7.3 Interaction Patterns

| Action | Trigger | Response |
|--------|---------|----------|
| Summarize Thread | Right-click → Summarize | RHS opens with loading → Summary |
| Manage Action Items | Channel menu → Action Items / `/actionitems` | RHS dashboard or modal opens |
| Apply Formatting | Composer toolbar → AI Format | Preview modal opens with before/after |
| Summarize Link | Click on link → Summarize | Summary card appears below message |

---

## 8. Data Model

### 8.1 Database Schema

**AIActionItems Table:**
```sql
CREATE TABLE AIActionItems (
  Id VARCHAR(26) PRIMARY KEY,
  Description TEXT NOT NULL,
  AssigneeId VARCHAR(26) NOT NULL,
  CreatorId VARCHAR(26) NOT NULL,
  ChannelId VARCHAR(26) NOT NULL,
  PostId VARCHAR(26),  -- Link to original message
  ThreadId VARCHAR(26),
  Deadline BIGINT,  -- Unix timestamp in milliseconds
  CreatedAt BIGINT NOT NULL,
  CompletedAt BIGINT,
  Status VARCHAR(32) NOT NULL,  -- 'active', 'completed', 'dismissed'
  Priority VARCHAR(32),  -- 'high', 'medium', 'low'
  ReminderSent BOOLEAN DEFAULT FALSE,
  DeleteAt BIGINT DEFAULT 0,
  
  INDEX idx_assignee (AssigneeId, Status, DeleteAt),
  INDEX idx_channel (ChannelId, Status, DeleteAt),
  INDEX idx_deadline (Deadline, Status, DeleteAt),
  FOREIGN KEY (AssigneeId) REFERENCES Users(Id),
  FOREIGN KEY (CreatorId) REFERENCES Users(Id),
  FOREIGN KEY (ChannelId) REFERENCES Channels(Id),
  FOREIGN KEY (PostId) REFERENCES Posts(Id)
);
```

**AISummaries Table (Cache):**
```sql
CREATE TABLE AISummaries (
  Id VARCHAR(26) PRIMARY KEY,
  ChannelId VARCHAR(26),
  ThreadId VARCHAR(26),
  Summary TEXT NOT NULL,
  KeyPoints TEXT,  -- JSON array
  Decisions TEXT,  -- JSON array
  MessageCount INT,
  ParticipantCount INT,
  StartTime BIGINT,
  EndTime BIGINT,
  CreatedAt BIGINT NOT NULL,
  ExpiresAt BIGINT NOT NULL,
  
  INDEX idx_channel_time (ChannelId, StartTime, EndTime),
  INDEX idx_thread (ThreadId),
  INDEX idx_expires (ExpiresAt)
);
```

**AILinkSummaries Table:**
```sql
CREATE TABLE AILinkSummaries (
  Id VARCHAR(26) PRIMARY KEY,
  Url TEXT NOT NULL,
  UrlHash VARCHAR(64) NOT NULL,  -- SHA-256 hash for indexing
  Title TEXT,
  Domain VARCHAR(255),
  Summary TEXT NOT NULL,
  KeyPoints TEXT,  -- JSON array of key points
  ReadingTimeMinutes INT,
  ContentLength INT,
  FetchedAt BIGINT NOT NULL,
  CreatedAt BIGINT NOT NULL,
  ExpiresAt BIGINT NOT NULL,  -- 7-day TTL
  
  UNIQUE INDEX idx_url_hash (UrlHash),
  INDEX idx_domain (Domain),
  INDEX idx_expires (ExpiresAt)
);
```

**AIPreferences Table:**
```sql
CREATE TABLE AIPreferences (
  UserId VARCHAR(26) PRIMARY KEY,
  SummaryLength VARCHAR(32) DEFAULT 'standard',  -- 'brief', 'standard', 'detailed'
  ActionItemReminders BOOLEAN DEFAULT TRUE,
  ReminderTime VARCHAR(5) DEFAULT '09:00',
  FormattingAutoSuggest BOOLEAN DEFAULT TRUE,
  FormattingDefaultProfile VARCHAR(32) DEFAULT 'professional',
  FormattingShowPreview BOOLEAN DEFAULT TRUE,
  CreatedAt BIGINT NOT NULL,
  UpdatedAt BIGINT NOT NULL,
  
  FOREIGN KEY (UserId) REFERENCES Users(Id)
);
```

**System Configuration (config.json):**
```json
{
  "AISettings": {
    "Enable": true,
    "OpenAIAPIKey": "encrypted_key",
    "OpenAIModel": "gpt-4",
    "MaxMessageLimit": 500,
    "APIRateLimit": 60,
    "EnableSummarization": true,
    "EnableLinkSummarizer": true,
    "EnableActionItems": true,
    "EnableFormatting": true,
    "ActionItemDetectionChannels": [],
    "FormattingAvailableToAll": true
  }
}
```

### 8.2 Query Patterns

| Query | SQL / Access Pattern |
|-------|----------------------|
| Get user's action items | `SELECT * FROM AIActionItems WHERE AssigneeId = ? AND Status = ? AND DeleteAt = 0 ORDER BY Deadline` |
| Get channel action items | `SELECT * FROM AIActionItems WHERE ChannelId = ? AND Status = ? AND DeleteAt = 0` |
| Get link summary by URL | `SELECT * FROM AILinkSummaries WHERE UrlHash = ? AND ExpiresAt > ?` |
| Check summary cache | `SELECT * FROM AISummaries WHERE ChannelId = ? AND StartTime = ? AND EndTime = ? AND ExpiresAt > ?` |
| Get user preferences | `SELECT * FROM AIPreferences WHERE UserId = ?` |
| Get overdue action items | `SELECT * FROM AIActionItems WHERE Deadline < ? AND Status = 'active' AND DeleteAt = 0` |

---

## 9. API Specifications

### 9.1 REST Endpoints

All AI endpoints are part of the native Mattermost API v4 under the `/api/v4/ai/` namespace.

#### POST /api/v4/ai/summarize

**Request:**
```json
{
  "type": "thread" | "channel",
  "channel_id": "string",
  "thread_id": "string (optional)",
  "start_time": "ISO8601 timestamp (optional)",
  "end_time": "ISO8601 timestamp (optional)",
  "length": "brief" | "standard" | "detailed"
}
```

**Response:**
```json
{
  "summary": "string",
  "key_points": ["string"],
  "decisions": ["string"],
  "action_items": [{
    "assignee": "string",
    "task": "string"
  }],
  "message_count": "number",
  "participant_count": "number",
  "generated_at": "ISO8601 timestamp"
}
```

#### POST /api/v4/ai/links/summarize

**Request:**
```json
{
  "url": "string",
  "force_refresh": "boolean (optional, default: false)"
}
```

**Response:**
```json
{
  "id": "string",
  "url": "string",
  "title": "string",
  "domain": "string",
  "summary": "string",
  "key_points": ["string"],
  "reading_time_minutes": "number",
  "content_length": "number",
  "fetched_at": "ISO8601 timestamp",
  "cached": "boolean"
}
```

#### GET /api/v4/ai/links/summary

**Query Parameters:**
- `url`: The URL to get summary for (required)

**Response:**
```json
{
  "id": "string",
  "url": "string",
  "title": "string",
  "domain": "string",
  "summary": "string",
  "key_points": ["string"],
  "reading_time_minutes": "number",
  "cached": "boolean",
  "expires_at": "ISO8601 timestamp"
}
```

#### POST /api/v4/ai/actionitems

**Request:**
```json
{
  "channel_id": "string",
  "description": "string",
  "assignee_id": "string",
  "due_at": "ISO8601 timestamp (optional)",
  "source_post_id": "string (optional)"
}
```

**Response:**
```json
{
  "id": "string",
  "channel_id": "string",
  "description": "string",
  "assignee_id": "string",
  "status": "active",
  "due_at": "ISO8601 timestamp (nullable)",
  "created_at": "ISO8601 timestamp"
}
```

#### GET /api/v4/ai/actionitems

**Query Parameters:**
- `channel_id`: filter to a single channel (optional)
- `assignee_id`: filter to a specific user (optional)
- `status`: `active`, `completed`, or `dismissed` (optional)

**Response:**
```json
{
  "action_items": [{
    "id": "string",
    "channel_id": "string",
    "channel_name": "string",
    "description": "string",
    "assignee_id": "string",
    "status": "active" | "completed" | "dismissed",
    "due_at": "ISO8601 timestamp (nullable)",
    "source_post_id": "string",
    "created_at": "ISO8601 timestamp",
    "completed_at": "ISO8601 timestamp (nullable)"
  }]
}
```

#### PUT /api/v4/ai/actionitems/{action_item_id}

**Request:**
```json
{
  "status": "active" | "completed" | "dismissed",
  "assignee_id": "string (optional)",
  "due_at": "ISO8601 timestamp (optional)"
}
```

**Response:**
```json
{
  "success": true,
  "action_item": {
    "id": "string",
    "status": "completed",
    "assignee_id": "string",
    "due_at": "ISO8601 timestamp (nullable)"
  }
}
```

### 9.2 Slash Commands

| Command | Description |
|---------|-------------|
| `/summarize thread` | Summarize current thread |
| `/summarize channel [time]` | Summarize channel for time period |
| `/actionitems` | Open personal action items dashboard |
| `/actionitems list` | List all active action items |
| `/actionitems create` | Manually create an action item |
| `/actionitems complete [id]` | Mark action item as complete |
| `/format [action]` | Apply formatting to current message |

---

## 10. Security and Privacy

### 10.1 Security Requirements

| ID | Requirement |
|----|-------------|
| SEC-1 | All API endpoints must verify user authentication |
| SEC-2 | Users can only access data for channels they are members of |
| SEC-3 | OpenAI API key must be stored encrypted in plugin settings |
| SEC-4 | No message content stored in logs |
| SEC-5 | Summary cache must respect channel membership changes |
| SEC-6 | Rate limiting on all API endpoints (100 req/min per user) |

### 10.2 Privacy Considerations

| Concern | Mitigation |
|---------|------------|
| Message content sent to OpenAI | Document in plugin description; admin can disable |
| Analytics tracking | Only aggregate data stored; no individual message tracking |
| Notification content | Stored temporarily; auto-deleted after 7 days |
| Action items | Only visible to channel members with access to the originating post; DM reminders sent only to assignee |

### 10.3 Data Retention

| Data Type | Retention Period |
|-----------|------------------|
| Summary cache | 24 hours |
| Link summaries cache | 7 days |
| Action items (active) | Until completed or dismissed |
| Notification history | 7 days |
| User preferences | Until user deletes account |

### 10.4 Permissions Matrix

| Action | User | Channel Admin | System Admin |
|--------|------|---------------|--------------|
| Summarize public channel | ✅ (if member) | ✅ | ✅ |
| Summarize private channel | ✅ (if member) | ✅ (if member) | ✅ |
| View public channel analytics | ✅ (if member) | ✅ | ✅ |
| View private channel analytics | ✅ (if member) | ✅ (if member) | ✅ |
| Manage action items | ✅ (if member) | ✅ | ✅ |
| Configure plugin settings | ❌ | ❌ | ✅ |

---

## 11. Performance Requirements

### 11.1 Response Time SLAs

| Operation | Target | Maximum |
|-----------|--------|---------|
| Thread summarization (< 50 messages) | 3 seconds | 10 seconds |
| Channel summarization (< 500 messages) | 5 seconds | 15 seconds |
| Link summarization (new URL) | 5 seconds | 15 seconds |
| Link summarization (cached) | 100 ms | 500 ms |
| Create/complete action item | 200 ms | 1 second |
| Notification classification | 50 ms | 200 ms |

### 11.2 Throughput Requirements

| Metric | Requirement |
|--------|-------------|
| Concurrent summarization requests | 10 per server |
| Messages processed for notifications | 1000/second |
| Active action items per user | 200 max |
| Link summary cache entries | 10,000 max |

### 11.3 Resource Limits

| Resource | Limit |
|----------|-------|
| OpenAI API calls per minute | 60 (configurable) |
| Maximum messages per summary | 500 (configurable, range: 100-1000) |
| Summary cache size | 1000 entries |
| Notification queue per user | 500 items |
| Maximum message length for summarization | 100,000 characters total |
| Maximum fetched page size | 5 MB |
| Maximum content for link summarization | 10,000 characters |

---

## 12. Dependencies

### 12.1 External Services

| Service | Purpose | Fallback |
|---------|---------|----------|
| OpenAI API | LLM for summarization and classification | Feature disabled with error message |

### 12.2 Mattermost Version Compatibility

| Mattermost Version | Feature Compatibility |
|--------------------|----------------------|
| 9.11+ (our fork) | Full support |
| 9.0-9.10 | May require minor adjustments |
| 8.x and below | Not compatible (requires core changes) |

### 12.3 Third-Party Libraries

**Go (Server):**
| Library | Version | Purpose |
|---------|---------|---------|
| sashabaranov/go-openai | latest | OpenAI API client |
| robfig/cron | v3 | Scheduled job execution |

**React (Web App):**
| Library | Version | Purpose |
|---------|---------|---------|
| recharts | ^2.x | Analytics charts |
| date-fns | ^2.x | Date manipulation |
| @mattermost/compass-components | latest | UI components |

---

## 13. Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| OpenAI API rate limits exceeded | Medium | High | Implement caching, queue requests, configurable limits |
| OpenAI API costs exceed budget | Medium | Medium | Set daily/monthly spending caps, optimize prompts |
| Large channels cause timeout | Medium | Medium | Paginate messages, set maximum message count |
| Reminder worker backlog (action items) | Low | High | Persist reminders in KV store, restart-safe scheduler with retries |
| Notification classification inaccurate | Medium | Low | Allow user feedback, manual override |
| Plugin conflicts with other plugins | Low | Medium | Namespace all resources, test with popular plugins |

---

## 14. Success Metrics

### 14.1 Adoption Metrics

| Metric | Target (30 days post-launch) |
|--------|------------------------------|
| Plugin installations | 100+ |
| Daily active users | 50+ |
| Message summaries generated | 500+ |
| Link summaries generated | 300+ |
| Action items tracked | 1000+ |
| Messages formatted | 300+ |

### 14.2 Engagement Metrics

| Metric | Target |
|--------|--------|
| Message summaries per user per week | 5+ |
| Link summaries per user per week | 10+ |
| Action items completed per user per week | 8+ |
| Messages formatted per user per week | 3+ |

### 14.3 Quality Metrics

| Metric | Target |
|--------|--------|
| Summary usefulness rating | 4.0/5.0 |
| Action item detection accuracy | 90% |
| Formatting improvement rating | 4.0/5.0 |
| Error rate | < 1% |

---

## 15. Timeline and Milestones

### 15.1 Development Timeline (6 Days)

| Day | Focus | Deliverables |
|-----|-------|--------------|
| 1-2 | Setup & Learning | Environment setup, Go basics, plugin architecture understanding, "Hello World" plugin |
| 3 | Core Infrastructure | OpenAI integration, plugin settings, basic API structure |
| 4 | Feature: Summarization | Slash command, thread/channel summarization, RHS display |
| 5 | Feature: Action Item Extractor | Auto-detection, dashboard, reminders |
| 6 | Feature: Link Summarizer + Formatter | Link summarization, formatting UI and AI integration |
| 7 | Polish & Documentation | Bug fixes, testing, README, demo video |

### 15.2 Milestones

| Milestone | Target Date | Criteria |
|-----------|-------------|----------|
| M1: Dev Environment Ready | Day 1 | Mattermost running locally, plugin compiles |
| M2: First Feature Complete | Day 4 | Summarization working end-to-end |
| M3: All Features Complete | Day 6 | All 4 features functional |
| M4: Production Ready | Day 7 | Documented, tested, deployable |

---

## 16. Future Enhancements

### 16.1 Version 1.1 (Potential)

- **Smart Notifications** - AI-powered notification prioritization and filtering
- **Channel Analytics Dashboard** - Visual insights into communication patterns
- Semantic Search with vector embeddings
- Recurring reminder rules for action items
- Multi-LLM support (Anthropic, local Ollama)

### 16.2 Version 1.2 (Potential)

- Real-time translation
- Voice message transcription
- Team health monitoring
- Mobile app support

---

## 17. Appendix

### 17.1 Glossary

| Term | Definition |
|------|------------|
| LLM | Large Language Model (e.g., GPT-4) |
| RHS | Right-Hand Sidebar in Mattermost UI |
| KV Store | Key-Value storage provided by Mattermost plugin API |
| Thread | A conversation started as a reply to a message |
| Channel | A Mattermost conversation space (public or private) |

### 17.2 References

- Mattermost Plugin Documentation: https://developers.mattermost.com/integrate/plugins/
- Mattermost Server Repository: https://github.com/mattermost/mattermost
- OpenAI API Documentation: https://platform.openai.com/docs
- Go Programming Language: https://go.dev/doc/

### 17.3 Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | Dec 2, 2024 | [Your Name] | Initial draft |

---

*End of Document*
