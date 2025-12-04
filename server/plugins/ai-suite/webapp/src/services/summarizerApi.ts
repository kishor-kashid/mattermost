import {apiClient} from './api';

export type SummaryType = 'thread' | 'channel';

export type SummaryParticipant = {
    id: string;
    username?: string;
    display_name?: string;
};

export type SummaryRange = {
    since: number;
    until: number;
    label: string;
};

export type SummaryContext = {
    type_label: string;
    message_limit: number;
    timeframe: string;
};

export type SummaryResponse = {
    id: string;
    type: SummaryType;
    channel_id: string;
    channel_name: string;
    root_post_id?: string;
    title: string;
    summary: string;
    message_count: number;
    participant_count: number;
    participants: SummaryParticipant[];
    generated_at: number;
    range: SummaryRange;
    context: SummaryContext;
    usage?: {
        prompt_tokens: number;
        completion_tokens: number;
        total_tokens: number;
    };
    limit_reached: boolean;
    cached: boolean;
};

export type SummarizeRequest = {
    type: SummaryType;
    channel_id?: string;
    root_post_id?: string;
    post_id?: string;
    time_range?: string;
    since?: number;
    until?: number;
    force?: boolean;
};

export const fetchSummary = (payload: SummarizeRequest): Promise<SummaryResponse> =>
    apiClient.post<SummaryResponse>('/summarize', payload);


