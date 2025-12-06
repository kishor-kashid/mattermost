// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export type AIActionItemStatus = 'open' | 'in_progress' | 'completed' | 'dismissed';
export type AIActionItemPriority = 'low' | 'medium' | 'high' | 'urgent';

export type AISummaryType = 'thread' | 'channel';

export interface AIActionItem {
    id: string;
    channel_id: string;
    post_id?: string;
    created_by: string;
    assignee_id?: string;
    description: string;
    due_date?: number;
    priority: AIActionItemPriority;
    status: AIActionItemStatus;
    completed_at?: number;
    create_at: number;
    update_at: number;
    delete_at: number;
}

export interface ActionItemCreateRequest {
    description: string;
    assignee_id: string;
    channel_id: string;
    post_id?: string;
    due_date?: Date;
    priority?: AIActionItemPriority;
    status?: AIActionItemStatus;
}

export interface ActionItemUpdateRequest {
    description?: string;
    assignee_id?: string;
    due_date?: Date;
    priority?: AIActionItemPriority;
    status?: AIActionItemStatus;
    completed_at?: Date;
}

export interface ActionItemFilters {
    userId?: string;
    channelId?: string;
    status?: AIActionItemStatus;
    priority?: AIActionItemPriority;
    dueBefore?: Date;
    dueAfter?: Date;
    assignedBy?: string;
    includeCompleted?: boolean;
    page?: number;
    perPage?: number;
}

export interface ActionItemStats {
    total: number;
    overdue: number;
    dueToday: number;
    dueSoon: number;
    noDueDate: number;
    completed: number;
    byPriority: Record<string, number>;
    byStatus: Record<string, number>;
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
    custom_instructions?: string;
}

export interface FormatMessageResponse {
    formatted_text: string;
    profile: string;
    diff?: TextDiff;
    processing_ms: number;
}

export interface TextDiff {
    original: string;
    formatted: string;
    changes?: TextChange[];
}

export interface TextChange {
    type: 'insert' | 'delete' | 'replace';
    start: number;
    end: number;
    new_text?: string;
    old_text?: string;
}

export interface FormattingProfileInfo {
    id: string;
    label: string;
    description: string;
}

export interface AIConfig {
    enable: boolean;
    enable_summarization: boolean;
    enable_analytics: boolean;
    enable_action_items: boolean;
    enable_formatting: boolean;
    max_message_limit: number;
}

