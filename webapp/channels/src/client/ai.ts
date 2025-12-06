// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AIActionItem, AIAnalytics, AIPreferences, AISummary, FormatMessageRequest, FormatMessageResponse, FormattingProfileInfo, SummarizeRequest} from 'types/ai';

import {Client4} from '@mattermost/client';

// AI Client for making API calls to AI endpoints
export class AIClient {
    private baseRoute = '/api/v4/ai';

    // System Endpoints

    async healthCheck(): Promise<any> {
        return Client4.doFetch(
            `${this.baseRoute}/health`,
            {method: 'get'},
        );
    }

    async validateConfig(apiKey: string, model: string): Promise<any> {
        return Client4.doFetch(
            `${this.baseRoute}/config/validate`,
            {
                method: 'post',
                body: JSON.stringify({openai_api_key: apiKey, model}),
            },
        );
    }

    async testConnection(testPrompt?: string): Promise<any> {
        return Client4.doFetch(
            `${this.baseRoute}/test`,
            {
                method: 'post',
                body: JSON.stringify({test_prompt: testPrompt || 'Hello'}),
            },
        );
    }

    // Summaries

    async createSummary(request: SummarizeRequest): Promise<AISummary> {
        return Client4.doFetch(
            `${this.baseRoute}/summarize`,
            {
                method: 'post',
                body: JSON.stringify(request),
            },
        );
    }

    async summarizeThread(postId: string, level: string = 'standard', useCache: boolean = true): Promise<{summary: AISummary; from_cache: boolean; processing_ms: number}> {
        return Client4.doFetch(
            `${this.baseRoute}/summarize/thread/${postId}?level=${level}&use_cache=${useCache}`,
            {method: 'get'},
        );
    }

    async summarizeChannel(channelId: string, startTime?: number, endTime?: number, level: string = 'standard', useCache: boolean = true): Promise<{summary: AISummary; from_cache: boolean; processing_ms: number}> {
        const body: any = {
            summary_level: level,
            use_cache: useCache,
        };
        if (startTime) {
            body.start_time = startTime;
        }
        if (endTime) {
            body.end_time = endTime;
        }

        return Client4.doFetch(
            `${this.baseRoute}/summarize/channel/${channelId}`,
            {
                method: 'post',
                body: JSON.stringify(body),
            },
        );
    }

    async getSummary(id: string): Promise<AISummary> {
        return Client4.doFetch(
            `${this.baseRoute}/summaries/${id}`,
            {method: 'get'},
        );
    }

    async getSummariesByChannel(channelId: string, page = 0, perPage = 60): Promise<AISummary[]> {
        return Client4.doFetch(
            `${this.baseRoute}/summaries/channel/${channelId}?page=${page}&per_page=${perPage}`,
            {method: 'get'},
        );
    }

    // Action Items

    async createActionItem(actionItem: Partial<AIActionItem>): Promise<AIActionItem> {
        return Client4.doFetch(
            `${this.baseRoute}/actionitems`,
            {
                method: 'post',
                body: JSON.stringify(actionItem),
            },
        );
    }

    async getActionItem(id: string): Promise<AIActionItem> {
        return Client4.doFetch(
            `${this.baseRoute}/actionitems/${id}`,
            {method: 'get'},
        );
    }

    async getActionItems(filters?: any): Promise<AIActionItem[]> {
        const params = new URLSearchParams();
        if (filters) {
            if (filters.userId) params.append('user_id', filters.userId);
            if (filters.channelId) params.append('channel_id', filters.channelId);
            if (filters.status) params.append('status', filters.status);
            if (filters.priority) params.append('priority', filters.priority);
            if (filters.includeCompleted) params.append('include_completed', 'true');
            if (filters.page !== undefined) params.append('page', filters.page.toString());
            if (filters.perPage !== undefined) params.append('per_page', filters.perPage.toString());
        }
        
        const queryString = params.toString();
        return Client4.doFetch(
            `${this.baseRoute}/actionitems${queryString ? '?' + queryString : ''}`,
            {method: 'get'},
        );
    }

    async getActionItemsByChannel(channelId: string, page = 0, perPage = 60): Promise<AIActionItem[]> {
        return this.getActionItems({channelId, page, perPage});
    }

    async getActionItemsByUser(userId: string, page = 0, perPage = 60): Promise<AIActionItem[]> {
        return this.getActionItems({userId, page, perPage});
    }

    async updateActionItem(id: string, actionItem: Partial<AIActionItem>): Promise<AIActionItem> {
        return Client4.doFetch(
            `${this.baseRoute}/actionitems/${id}`,
            {
                method: 'put',
                body: JSON.stringify(actionItem),
            },
        );
    }

    async completeActionItem(id: string): Promise<AIActionItem> {
        return Client4.doFetch(
            `${this.baseRoute}/actionitems/${id}/complete`,
            {method: 'post'},
        );
    }

    async deleteActionItem(id: string): Promise<void> {
        return Client4.doFetch(
            `${this.baseRoute}/actionitems/${id}`,
            {method: 'delete'},
        );
    }

    async getActionItemStats(userId?: string): Promise<any> {
        const params = userId ? `?user_id=${userId}` : '';
        return Client4.doFetch(
            `${this.baseRoute}/actionitems/stats${params}`,
            {method: 'get'},
        );
    }

    // Analytics

    async getChannelAnalytics(channelId: string, startDate: string, endDate: string): Promise<AIAnalytics[]> {
        return Client4.doFetch(
            `${this.baseRoute}/analytics/${channelId}?start_date=${startDate}&end_date=${endDate}`,
            {method: 'get'},
        );
    }

    // Message Formatting

    async formatPreview(request: FormatMessageRequest): Promise<FormatMessageResponse> {
        return Client4.doFetch(
            `${this.baseRoute}/format/preview`,
            {
                method: 'post',
                body: JSON.stringify(request),
            },
        );
    }

    async formatApply(request: FormatMessageRequest): Promise<FormatMessageResponse> {
        return Client4.doFetch(
            `${this.baseRoute}/format/apply`,
            {
                method: 'post',
                body: JSON.stringify(request),
            },
        );
    }

    async getFormattingProfiles(): Promise<FormattingProfileInfo[]> {
        return Client4.doFetch(
            `${this.baseRoute}/format/profiles`,
            {method: 'get'},
        );
    }

    // Preferences

    async getPreferences(userId: string): Promise<AIPreferences> {
        return Client4.doFetch(
            `${this.baseRoute}/preferences/${userId}`,
            {method: 'get'},
        );
    }

    async updatePreferences(userId: string, preferences: Partial<AIPreferences>): Promise<AIPreferences> {
        return Client4.doFetch(
            `${this.baseRoute}/preferences/${userId}`,
            {
                method: 'put',
                body: JSON.stringify(preferences),
            },
        );
    }
}

// Export a singleton instance
export const aiClient = new AIClient();

