// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export type AIActionItemStatus = 'pending' | 'completed' | 'dismissed';

export type AISummaryType = 'thread' | 'channel';

export interface AIActionItem {
    id: string;
    channel_id: string;
    post_id?: string;
    user_id: string;
    assignee_id?: string;
    description: string;
    deadline?: number;
    status: AIActionItemStatus;
    reminder_sent: boolean;
    create_at: number;
    update_at: number;
    delete_at: number;
}

export interface AISummary {
    id: string;
    channel_id: string;
    post_id?: string;
    summary_type: AISummaryType;
    summary: string;
    message_count: number;
    start_time: number;
    end_time: number;
    create_at: number;
    expires_at: number;
}

export interface AIAnalytics {
    id: string;
    channel_id: string;
    date: string;
    message_count: number;
    user_count: number;
    avg_response_time: number;
    top_contributors: Record<string, any>;
    hourly_distribution: Record<string, any>;
    create_at: number;
    update_at: number;
}

export interface AIPreferences {
    id: string;
    user_id: string;
    enable_summarization: boolean;
    enable_analytics: boolean;
    enable_action_items: boolean;
    enable_formatting: boolean;
    default_model: string;
    formatting_profile: string;
    create_at: number;
    update_at: number;
}

export interface SummarizeRequest {
    channel_id: string;
    post_id?: string;
    summary_type: AISummaryType;
    message_count?: number;
}

export interface FormatMessageRequest {
    message: string;
    profile?: string;
}

export interface FormatMessageResponse {
    formatted_message: string;
}

export interface AIConfig {
    enable: boolean;
    enable_summarization: boolean;
    enable_analytics: boolean;
    enable_action_items: boolean;
    enable_formatting: boolean;
    max_message_limit: number;
}

