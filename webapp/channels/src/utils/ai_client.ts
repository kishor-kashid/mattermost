// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AIActionItem, AIAnalytics, AIPreferences, AISummary, FormatMessageRequest, FormatMessageResponse, SummarizeRequest} from 'types/ai';

// AI Client utilities for making API calls to the AI endpoints
// These will be integrated with Client4 once the API endpoints are implemented

export class AIClient {
    // Action Items
    static async createActionItem(actionItem: Partial<AIActionItem>): Promise<AIActionItem> {
        // TODO: Implement API call to /api/v4/ai/actionitems
        throw new Error('Not implemented');
    }

    static async getActionItem(id: string): Promise<AIActionItem> {
        // TODO: Implement API call to /api/v4/ai/actionitems/{id}
        throw new Error('Not implemented');
    }

    static async getActionItemsByChannel(channelId: string, page = 0, perPage = 60): Promise<AIActionItem[]> {
        // TODO: Implement API call to /api/v4/ai/actionitems/channel/{channelId}
        throw new Error('Not implemented');
    }

    static async getActionItemsByUser(userId: string, page = 0, perPage = 60): Promise<AIActionItem[]> {
        // TODO: Implement API call to /api/v4/ai/actionitems/user/{userId}
        throw new Error('Not implemented');
    }

    static async updateActionItem(id: string, actionItem: Partial<AIActionItem>): Promise<AIActionItem> {
        // TODO: Implement API call to PATCH /api/v4/ai/actionitems/{id}
        throw new Error('Not implemented');
    }

    static async deleteActionItem(id: string): Promise<void> {
        // TODO: Implement API call to DELETE /api/v4/ai/actionitems/{id}
        throw new Error('Not implemented');
    }

    // Summaries
    static async createSummary(request: SummarizeRequest): Promise<AISummary> {
        // TODO: Implement API call to POST /api/v4/ai/summarize
        throw new Error('Not implemented');
    }

    static async getSummary(id: string): Promise<AISummary> {
        // TODO: Implement API call to GET /api/v4/ai/summaries/{id}
        throw new Error('Not implemented');
    }

    static async getSummariesByChannel(channelId: string, page = 0, perPage = 60): Promise<AISummary[]> {
        // TODO: Implement API call to GET /api/v4/ai/summaries/channel/{channelId}
        throw new Error('Not implemented');
    }

    // Analytics
    static async getChannelAnalytics(channelId: string, startDate: string, endDate: string): Promise<AIAnalytics[]> {
        // TODO: Implement API call to GET /api/v4/ai/analytics/{channelId}
        throw new Error('Not implemented');
    }

    // Message Formatting
    static async formatMessage(request: FormatMessageRequest): Promise<FormatMessageResponse> {
        // TODO: Implement API call to POST /api/v4/ai/format
        throw new Error('Not implemented');
    }

    // Preferences
    static async getPreferences(userId: string): Promise<AIPreferences> {
        // TODO: Implement API call to GET /api/v4/ai/preferences/{userId}
        throw new Error('Not implemented');
    }

    static async updatePreferences(userId: string, preferences: Partial<AIPreferences>): Promise<AIPreferences> {
        // TODO: Implement API call to PUT /api/v4/ai/preferences/{userId}
        throw new Error('Not implemented');
    }
}

