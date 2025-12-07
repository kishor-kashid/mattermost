import React, {useMemo, useState, useCallback} from 'react';
import {shallowEqual, useDispatch, useSelector} from 'react-redux';

import {summarizeLink} from 'actions/ai_link_summarizer';
import {getLinkSummaryForUrl, getLinkSummaryError, isLinkSummaryLoading} from 'selectors/ai_link_summarizer';
import type {GlobalState} from 'types/store';

import {LinkSummaryCard} from './summary_card';
import './link_summary.scss';

type Props = {
    urls: string[];
};

const normalizeUrl = (url: string) => url.trim().toLowerCase();

export const PostLinkSummary = ({urls}: Props) => {
    const uniqueUrls = useMemo(() => {
        const seen = new Set<string>();
        return urls.map((u) => normalizeUrl(u)).filter((u) => {
            if (!u || seen.has(u)) {
                return false;
            }
            seen.add(u);
            return true;
        });
    }, [urls]);

    const dispatch = useDispatch();

    const summaries = useSelector((state: GlobalState) => uniqueUrls.map((u) => ({
        url: u,
        summary: getLinkSummaryForUrl(state, u),
        loading: isLinkSummaryLoading(state, u),
        error: getLinkSummaryError(state, u),
    })), shallowEqual);

    // Debug logging
    console.log('[PostLinkSummary] Rendering with summaries:', summaries);

    // Track which URLs the user has explicitly asked to summarize
    const [requested, setRequested] = useState<Set<string>>(new Set(uniqueUrls.filter((u) => {
        // If a summary already exists (e.g., cached), consider it requested so it renders immediately
        const match = summaries.find((s) => s.url === u);
        return Boolean(match?.summary);
    })));

    const requestSummary = useCallback((url: string) => {
        setRequested((prev) => {
            const next = new Set(prev);
            next.add(url);
            return next;
        });
        dispatch(summarizeLink(url, false));
    }, [dispatch]);

    if (uniqueUrls.length === 0) {
        return null;
    }

    return (
        <div className='ai-link-summary__list'>
            {summaries.map(({url, summary, loading, error}) => {
                const isRequested = requested.has(url);
                
                // If requested but no loading/summary/error yet, treat as loading (race condition fix)
                const effectiveLoading = loading || (isRequested && !summary && !error);
                
                console.log('[PostLinkSummary] Render decision for:', url, {
                    summary: !!summary,
                    loading,
                    effectiveLoading,
                    error,
                    requested: isRequested,
                });

                // Determine which branch to render
                const showButton = !summary && !effectiveLoading && !isRequested;
                const showError = !summary && !effectiveLoading && error;
                const showCard = isRequested || summary;

                return (
                    <div
                        className='ai-link-summary__container'
                        key={url}
                    >
                        {showButton ? (
                            <button
                                type='button'
                                className='ai-link-summary__refresh'
                                onClick={() => requestSummary(url)}
                            >
                                Summarize link
                            </button>
                        ) : showError ? (
                            <div className='ai-link-summary__error'>
                                <div>{error}</div>
                                <button
                                    type='button'
                                    className='ai-link-summary__refresh'
                                    onClick={() => requestSummary(url)}
                                >
                                    Retry
                                </button>
                            </div>
                        ) : showCard ? (
                            <LinkSummaryCard
                                summary={summary || undefined}
                                loading={effectiveLoading}
                                onRefresh={() => dispatch(summarizeLink(url, true))}
                            />
                        ) : null}
                    </div>
                );
            })}
        </div>
    );
};

