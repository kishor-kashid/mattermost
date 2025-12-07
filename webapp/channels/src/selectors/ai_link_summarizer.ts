import type {GlobalState} from 'types/store';
import type {AILinkSummary} from 'types/ai';

const normalizeUrl = (url: string) => url.trim().toLowerCase();

export const getLinkSummaryState = (state: GlobalState) => state.ai.linkSummaries;

export function getLinkSummaryForUrl(state: GlobalState, url: string): AILinkSummary | null {
    const key = normalizeUrl(url);
    return state.ai.linkSummaries.byUrl[key]?.summary || null;
}

export function isLinkSummaryLoading(state: GlobalState, url: string): boolean {
    const key = normalizeUrl(url);
    return Boolean(state.ai.linkSummaries.loading[key]);
}

export function getLinkSummaryError(state: GlobalState, url: string): string | null {
    const key = normalizeUrl(url);
    return state.ai.linkSummaries.error[key] || null;
}

