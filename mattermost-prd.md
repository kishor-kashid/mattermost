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
   - 5.2 [Channel Analytics Dashboard](#52-channel-analytics-dashboard)
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
- **Lack of Insights:** Teams have no visibility into communication patterns and channel health
- **Lost Action Items:** Tasks and commitments mentioned in conversations get forgotten
- **Communication Quality:** Messages lack proper formatting and clarity

### 1.3 Solution

The Mattermost AI Productivity Suite addresses these challenges through four integrated features:

1. **AI Message Summarization** - Intelligent summaries of threads and channels
2. **Channel Analytics Dashboard** - Visual insights into communication patterns
3. **Action Item Extractor** - Automatic detection and tracking of tasks and commitments
4. **Message Formatting Assistant** - AI-powered message improvement and formatting

### 1.4 Key Benefits

| Benefit | Impact |
|---------|--------|
| Reduced time catching up on channels | 60% faster information absorption |
| Data-driven team insights | Improved communication health visibility |
| Never miss action items | 100% task capture from conversations |
| Professional communication | Improved message clarity and formatting |

---

## 2. Project Overview

### 2.1 Project Name

**Mattermost AI Productivity Suite** (Internal codename: "MattermostAI")

### 2.2 Project Type

Mattermost Plugin (Server + Web App components)

### 2.3 Technology Stack

| Component | Technology |
|-----------|------------|
| Backend | Go (Golang) |
| Frontend | React, TypeScript |
| Database | PostgreSQL (Mattermost's existing DB) |
| AI/LLM | OpenAI GPT-4 API |
| Build System | Make, Webpack |
| Plugin Framework | Mattermost Plugin SDK |

### 2.4 Repository Structure

```
mattermost/
├── server/                 # Go backend (cmd, channels, enterprise, plugins, etc.)
│   ├── cmd/                # Binary entrypoints (`mattermost`, `platform`, tooling)
│   ├── channels/           # Core server app (api4, app, store, jobs, wsapi, web)
│   ├── plugins/            # Packaged plugins + local dev artifacts
│   ├── config/, data/, logs/ # Default configs, sample data, local runtime output
│   └── scripts/, tests/, public/ # Build helpers, QA tooling, static assets
├── webapp/                 # React/TypeScript clients
│   ├── channels/           # Primary Mattermost web client (src/, build/, dist/)
│   ├── platform/           # Shared platform components
│   ├── scripts/            # Frontend build/test helpers
│   └── patches/            # Yarn/NPM patch files
├── api/                    # REST/v4 API reference & OpenAPI specs
├── e2e-tests/              # Cypress/Playwright automation suites
├── tools/                  # Build, localization, packaging utilities
├── .github/                # Actions workflows, issue templates
├── memory-bank/            # Project knowledge base for this effort
└── Root docs & configs     # PRD, task list, README, licenses, etc.
```

This is the upstream Mattermost monorepo; our AI Productivity Suite plugin will live inside `server/plugins` (Go backend) and `webapp` (React client) while still following the Mattermost plugin conventions outlined elsewhere in this PRD.

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
| G2 | Provide communication insights | Dashboard loads in <2 seconds with 30-day data |
| G3 | Capture all action items | Extract 95%+ of commitments from conversations |
| G4 | Improve message quality | 80% of formatted messages rated as clearer |

### 3.2 Secondary Goals

- Demonstrate proficiency in Go and React development
- Learn Mattermost plugin architecture
- Build production-ready, deployable software
- Create comprehensive documentation

### 3.3 Non-Goals (Out of Scope)

- Mobile app modifications
- Mattermost core server changes
- Multi-LLM support (OpenAI only for v1.0)
- Real-time collaborative features
- Integration with external task management tools

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

#### Channel Analytics Dashboard

| ID | User Story | Priority |
|----|------------|----------|
| US-2.1 | As a user, I want to see message volume over time so that I understand channel activity patterns | P0 |
| US-2.2 | As a user, I want to see who the most active participants are so that I identify key contributors | P0 |
| US-2.3 | As a user, I want to see peak activity hours so that I know when to post for maximum visibility | P1 |
| US-2.4 | As a manager, I want to see response time metrics so that I can assess team responsiveness | P1 |
| US-2.5 | As a user, I want to export analytics data so that I can create custom reports | P2 |

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

### 5.2 Channel Analytics Dashboard

#### 5.2.1 Overview

The Channel Analytics Dashboard provides visual insights into communication patterns, helping users and managers understand channel activity, engagement, and team collaboration health.

#### 5.2.2 Functional Requirements

#### 5.3.1 Overview

The Channel Analytics Dashboard provides visual insights into communication patterns, helping users and managers understand channel activity, engagement, and team collaboration health.

#### 5.3.2 Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-3.1 | System shall display message volume over time (line chart) | P0 |
| FR-3.2 | System shall display top contributors by message count (bar chart) | P0 |
| FR-3.3 | System shall display activity by hour of day (heatmap) | P0 |
| FR-3.4 | System shall display activity by day of week (bar chart) | P0 |
| FR-3.5 | System shall calculate and display average response time | P1 |
| FR-3.6 | System shall display thread engagement rate (% of messages that get replies) | P1 |
| FR-3.7 | System shall display emoji reaction breakdown (pie chart) | P2 |
| FR-3.8 | System shall support date range filtering (7d, 30d, 90d, custom) | P0 |
| FR-3.9 | System shall support channel comparison (up to 3 channels) | P2 |
| FR-3.10 | System shall export data as CSV | P1 |
| FR-3.11 | System shall display file sharing statistics | P2 |
| FR-3.12 | System shall calculate channel health score | P2 |

#### 5.3.3 Metrics Definitions

| Metric | Definition | Calculation |
|--------|------------|-------------|
| Message Volume | Total messages in time period | COUNT(messages) |
| Active Users | Users who posted at least once | COUNT(DISTINCT user_id) |
| Avg Messages/Day | Average daily message count | SUM(messages) / days |
| Response Time | Avg time between message and first reply | AVG(first_reply_time - message_time) |
| Thread Rate | Percentage of messages that start threads | threads / total_messages * 100 |
| Engagement Rate | Percentage of messages with reactions or replies | engaged_messages / total_messages * 100 |
| Peak Hour | Hour with most messages | MODE(hour_of_day) |
| Top Contributor | User with most messages | MAX(user_message_count) |

#### 5.3.4 Dashboard Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  Channel Analytics: #engineering          [7d ▼] [Export CSV]   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌─────────┐│
│  │   Messages   │ │ Active Users │ │ Avg Response │ │ Threads ││
│  │     847      │ │      23      │ │    12 min    │ │   34%   ││
│  │   ↑ 12%      │ │    ↓ 3%     │ │    ↓ 8%     │ │  ↑ 5%   ││
│  └──────────────┘ └──────────────┘ └──────────────┘ └─────────┘│
│                                                                  │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  Message Volume Over Time                    ││
│  │  ▲                                                          ││
│  │  │    ╭─╮                        ╭───╮                      ││
│  │  │   ╭╯ ╰╮    ╭──╮             ╭╯   ╰╮                     ││
│  │  │  ╭╯   ╰────╯  ╰─╮  ╭───╮  ╭╯     ╰───╮                 ││
│  │  │──╯              ╰──╯   ╰──╯          ╰────              ││
│  │  └──────────────────────────────────────────────────▶       ││
│  │    Mon   Tue   Wed   Thu   Fri   Sat   Sun                  ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│  ┌────────────────────────┐  ┌─────────────────────────────────┐│
│  │   Top Contributors     │  │      Activity by Hour           ││
│  │                        │  │                                  ││
│  │  @sarah    ████████ 89 │  │  ░░▓▓████████████▓▓░░░░░░░░░░  ││
│  │  @john     ██████   67 │  │  12am    6am    12pm   6pm 12am ││
│  │  @mike     █████    52 │  │                                  ││
│  │  @lisa     ████     41 │  │  Peak: 2pm-4pm (34% of traffic) ││
│  │  @alex     ███      28 │  │                                  ││
│  └────────────────────────┘  └─────────────────────────────────┘│
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### 5.3.5 User Interface

**Access Points:**
- Channel header menu → "View Analytics"
- Slash command: `/analytics [channel-name]`
- Main menu → "Channel Analytics"

**Dashboard Features:**
- Responsive layout (desktop and tablet)
- Interactive charts (hover for details)
- Date range selector
- Channel selector dropdown
- Export button (CSV format)
- Refresh button

**Permissions:**
- All users can view analytics for public channels they belong to
- Private channel analytics visible only to members
- System admins can view all channel analytics

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
│  │                    AI Productivity Suite Plugin                 │ │
│  │                                                                 │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐│ │
│  │  │ Summarizer  │ │  Analytics  │ │ Action Item │ │ Formatter ││ │
│  │  │   Service   │ │  Collector  │ │  Extractor  │ │  Service  ││ │
│  │  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └─────┬─────┘│ │
│  │         │               │               │              │       │ │
│  │  ┌──────┴───────────────┴───────────────┴──────────────┴─────┐│ │
│  │  │                    Core Services Layer                     ││ │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ ││ │
│  │  │  │  OpenAI      │  │   Message    │  │   Notification   │ ││ │
│  │  │  │  Client      │  │   Processor  │  │   Manager        │ ││ │
│  │  │  └──────────────┘  └──────────────┘  └──────────────────┘ ││ │
│  │  └───────────────────────────────────────────────────────────┘│ │
│  │                              │                                 │ │
│  │  ┌───────────────────────────┴───────────────────────────────┐│ │
│  │  │                    Data Access Layer                       ││ │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐││ │
│  │  │  │   Plugin    │  │  Mattermost │  │   Key-Value Store   │││ │
│  │  │  │   Store     │  │    API      │  │   (Plugin KV)       │││ │
│  │  │  └─────────────┘  └─────────────┘  └─────────────────────┘││ │
│  │  └───────────────────────────────────────────────────────────┘│ │
│  │                                                                │ │
│  │  ┌───────────────────────────────────────────────────────────┐│ │
│  │  │                    REST API Endpoints                      ││ │
│  │  │  POST /summarize  |  GET /analytics  |  POST /format     ││ │
│  │  │  GET /actionitems  |  PUT /actionitem  |  POST /complete ││ │
│  │  └───────────────────────────────────────────────────────────┘│ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│                         PostgreSQL Database                          │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────────────┐│
│  │   Posts    │ │  Channels  │ │   Users    │ │  PluginKeyValue    ││
│  └────────────┘ └────────────┘ └────────────┘ └────────────────────┘│
└──────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │    OpenAI API       │
                         │  (GPT-4 / GPT-3.5)  │
                         └─────────────────────┘
```

### 6.2 Component Descriptions

| Component | Responsibility |
|-----------|----------------|
| Summarizer Service | Handles message retrieval, formatting, OpenAI calls, and summary generation |
| Analytics Collector | Aggregates message data, calculates metrics, prepares visualization data |
| Action Item Extractor | Detects commitments and tasks in conversations, manages action item tracking |
| Formatter Service | Improves message quality through AI-powered formatting and grammar checking |
| OpenAI Client | Wrapper for OpenAI API with rate limiting, retries, and error handling |
| Message Processor | Common utilities for message formatting, filtering, and parsing |
| Plugin Store | Persistent storage for plugin-specific data (action items, preferences, cache) |

### 6.3 Plugin Hooks Used

| Hook | Purpose |
|------|---------|
| `MessageHasBeenPosted` | Trigger analytics data collection and action item extraction |
| `MessageWillBePosted` | Apply formatting assistance if user requested |
| `OnActivate` | Initialize services, load action items |
| `OnDeactivate` | Cleanup, persist state |
| `ServeHTTP` | Handle custom REST API endpoints |

### 6.4 Data Flow Diagrams

**AI Summarization Flow:**
```
User Request → Validate Permissions → Fetch Messages → Format for LLM
     → Call OpenAI API → Parse Response → Cache Result → Return to User
```

**Analytics Flow:**
```
Message Posted → Extract Metrics → Update Aggregates → Store in KV
     → User Requests Dashboard → Query Aggregates → Generate Charts → Render
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

Used for: Schedule Message, Settings, Date Range Picker

#### 7.2.3 Channel Header Button

Dropdown menu with plugin actions:
- Summarize Channel
- View Analytics
- Schedule Message

#### 7.2.4 Main Menu Items

Under plugin section:
- Scheduled Messages
- Notification Settings
- Channel Analytics

### 7.3 Interaction Patterns

| Action | Trigger | Response |
|--------|---------|----------|
| Summarize Thread | Right-click → Summarize | RHS opens with loading → Summary |
| Schedule Message | Send button dropdown | Modal opens |
| View Analytics | Channel menu → Analytics | New tab/modal with dashboard |
| Check Notifications | Bell icon | RHS opens with notification list |

---

## 8. Data Model

### 8.1 Plugin Key-Value Store Schema

**Action Items:**
```
Key: "actionitem:{action_item_id}"
Value: {
  id: string
  description: string
  assignee_id: string
  creator_id: string
  channel_id: string
  message_id: string  // Link to original message
  thread_id: string (optional)
  deadline: timestamp (optional)
  created_at: timestamp
  completed_at: timestamp (optional)
  status: "active" | "completed" | "dismissed"
  priority: "high" | "medium" | "low"
  reminder_sent: boolean
}
```

**User Preferences:**
```
Key: "preferences:{user_id}"
Value: {
  summary_length: "brief" | "standard" | "detailed"
  action_item_reminders: boolean
  reminder_time: string  // e.g., "09:00"
  formatting_auto_suggest: boolean
  formatting_default_profile: "professional" | "casual" | "technical" | "concise"
  formatting_show_preview: boolean
}
```

**Plugin Configuration (System Settings):**
```
{
  openai_api_key: string (encrypted)
  openai_model: "gpt-4" | "gpt-3.5-turbo"
  max_message_limit: number (default: 500)
  api_rate_limit: number (default: 60)
  enable_summarization: boolean (default: true)
  enable_analytics: boolean (default: true)
  enable_action_items: boolean (default: true)
  enable_formatting: boolean (default: true)
  action_item_detection_channels: string[] (default: all channels)
  formatting_available_to_all: boolean (default: true)
}
```

**Summary Cache:**
```
Key: "summary:{channel_id}:{date_range_hash}"
Value: {
  summary: string
  generated_at: timestamp
  message_count: number
  expires_at: timestamp
}
```

**Analytics Aggregates:**
```
Key: "analytics:{channel_id}:{date}"
Value: {
  message_count: number
  unique_users: number
  thread_count: number
  reaction_count: number
  file_count: number
  hourly_distribution: number[24]
  top_contributors: {user_id: string, count: number}[]
  avg_response_time_ms: number
}
```

### 8.2 Indexes and Queries

| Query | Access Pattern |
|-------|----------------|
| Get user's action items | Scan keys matching `actionitem:*` filtered by assignee_id |
| Get channel action items | Scan keys matching `actionitem:*` filtered by channel_id |
| Get channel analytics for date range | Fetch `analytics:{channel_id}:{date}` for each date |
| Check summary cache | Direct key lookup with hash of parameters |
| Get user preferences | Direct key lookup `preferences:{user_id}` |

---

## 9. API Specifications

### 9.1 REST Endpoints

#### POST /plugins/ai-suite/api/v1/summarize

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

#### GET /plugins/ai-suite/api/v1/analytics/{channel_id}

**Query Parameters:**
- `start_date`: ISO8601 date
- `end_date`: ISO8601 date

**Response:**
```json
{
  "channel_id": "string",
  "period": {
    "start": "ISO8601",
    "end": "ISO8601"
  },
  "metrics": {
    "total_messages": "number",
    "active_users": "number",
    "avg_messages_per_day": "number",
    "avg_response_time_seconds": "number",
    "thread_percentage": "number",
    "engagement_rate": "number"
  },
  "time_series": [{
    "date": "ISO8601 date",
    "messages": "number",
    "users": "number"
  }],
  "hourly_distribution": "number[24]",
  "top_contributors": [{
    "user_id": "string",
    "username": "string",
    "message_count": "number"
  }]
}
```

#### POST /plugins/ai-suite/api/v1/schedule

**Request:**
```json
{
  "channel_id": "string",
  "message": "string",
  "scheduled_at": "ISO8601 timestamp",
  "timezone": "string (IANA timezone)",
  "file_ids": ["string"]
}
```

**Response:**
```json
{
  "id": "string",
  "status": "scheduled",
  "scheduled_at": "ISO8601 timestamp",
  "message_preview": "string"
}
```

#### GET /plugins/ai-suite/api/v1/schedule

**Response:**
```json
{
  "scheduled_messages": [{
    "id": "string",
    "channel_id": "string",
    "channel_name": "string",
    "message_preview": "string",
    "scheduled_at": "ISO8601 timestamp",
    "status": "pending" | "sent" | "cancelled",
    "created_at": "ISO8601 timestamp"
  }]
}
```

#### DELETE /plugins/ai-suite/api/v1/schedule/{message_id}

**Response:**
```json
{
  "success": true,
  "message": "Scheduled message cancelled"
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
| `/analytics` | Open analytics dashboard for current channel |
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
| Scheduled messages | Only visible to message author |

### 10.3 Data Retention

| Data Type | Retention Period |
|-----------|------------------|
| Summary cache | 24 hours |
| Analytics aggregates | 90 days |
| Scheduled messages (sent) | 7 days |
| Notification history | 7 days |
| User preferences | Until user deletes account |

### 10.4 Permissions Matrix

| Action | User | Channel Admin | System Admin |
|--------|------|---------------|--------------|
| Summarize public channel | ✅ (if member) | ✅ | ✅ |
| Summarize private channel | ✅ (if member) | ✅ (if member) | ✅ |
| View public channel analytics | ✅ (if member) | ✅ | ✅ |
| View private channel analytics | ✅ (if member) | ✅ (if member) | ✅ |
| Schedule messages | ✅ | ✅ | ✅ |
| Configure plugin settings | ❌ | ❌ | ✅ |

---

## 11. Performance Requirements

### 11.1 Response Time SLAs

| Operation | Target | Maximum |
|-----------|--------|---------|
| Thread summarization (< 50 messages) | 3 seconds | 10 seconds |
| Channel summarization (< 500 messages) | 5 seconds | 15 seconds |
| Analytics dashboard load | 1 second | 3 seconds |
| Schedule message | 200 ms | 1 second |
| Notification classification | 50 ms | 200 ms |

### 11.2 Throughput Requirements

| Metric | Requirement |
|--------|-------------|
| Concurrent summarization requests | 10 per server |
| Messages processed for notifications | 1000/second |
| Scheduled messages per user | 50 max |
| Analytics data points | 90 days of daily aggregates |

### 11.3 Resource Limits

| Resource | Limit |
|----------|-------|
| OpenAI API calls per minute | 60 (configurable) |
| Maximum messages per summary | 500 (configurable, range: 100-1000) |
| Summary cache size | 1000 entries |
| Notification queue per user | 500 items |
| Maximum message length for summarization | 100,000 characters total |

---

## 12. Dependencies

### 12.1 External Services

| Service | Purpose | Fallback |
|---------|---------|----------|
| OpenAI API | LLM for summarization and classification | Feature disabled with error message |

### 12.2 Mattermost Version Compatibility

| Mattermost Version | Plugin Compatibility |
|--------------------|---------------------|
| 9.0+ | Full support |
| 8.x | Limited (no RHS updates) |
| 7.x | Not supported |

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
| Scheduled messages lost on server restart | Low | High | Persist to database, recover on startup |
| Notification classification inaccurate | Medium | Low | Allow user feedback, manual override |
| Plugin conflicts with other plugins | Low | Medium | Namespace all resources, test with popular plugins |

---

## 14. Success Metrics

### 14.1 Adoption Metrics

| Metric | Target (30 days post-launch) |
|--------|------------------------------|
| Plugin installations | 100+ |
| Daily active users | 50+ |
| Summaries generated | 500+ |
| Action items tracked | 1000+ |
| Messages formatted | 300+ |

### 14.2 Engagement Metrics

| Metric | Target |
|--------|--------|
| Summaries per user per week | 5+ |
| Analytics dashboard views per week | 3+ |
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
| 6 | Feature: Analytics + Message Formatter | Analytics dashboard, formatting UI and AI integration |
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
- Action Item Extractor
- Semantic Search with vector embeddings
- Recurring scheduled messages
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
