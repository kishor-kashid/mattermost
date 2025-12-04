import React from 'react';

import {useSummary} from '../../hooks/useSummary';
import {SummarizeRequest} from '../../services/summarizerApi';
import {SummaryBridge, SummaryTarget} from '../../summarizer/bridge';
import ErrorMessage from '../common/ErrorMessage';
import LoadingSpinner from '../common/LoadingSpinner';
import SummaryContent from './SummaryContent';
import SummaryOptions from './SummaryOptions';

import '../../styles/summary.css';

type Props = {
    bridge: SummaryBridge;
    onClose?: () => void;
};

export const SummaryPanel = ({bridge, onClose}: Props) => {
    const [target, setTarget] = React.useState<SummaryTarget | null>(null);
    const [request, setRequest] = React.useState<SummarizeRequest | null>(null);
    const [rangePreset, setRangePreset] = React.useState('24h');

    React.useEffect(() => bridge.subscribe(setTarget), [bridge]);

    React.useEffect(() => {
        if (!target) {
            setRequest(null);
            return;
        }

        if (target.type === 'thread') {
            const root = target.rootPostId ?? target.postId;
            if (!root) {
                setRequest(null);
                return;
            }
            setRequest({
                type: 'thread',
                channel_id: target.channelId,
                root_post_id: root,
                post_id: target.postId,
            });
            return;
        }

        if (!target.channelId) {
            setRequest(null);
            return;
        }

        const preset = target.timeRange ?? rangePreset;
        setRangePreset(preset);
        setRequest({
            type: 'channel',
            channel_id: target.channelId,
            time_range: preset,
            since: target.since,
            until: target.until,
        });
    }, [rangePreset, target]);

    const {summary, isLoading, error, refresh} = useSummary({
        request,
        enabled: Boolean(request),
    });

    const handleCopy = React.useCallback(() => {
        if (!summary) {
            return;
        }
        void navigator.clipboard?.writeText(summary.summary);
    }, [summary]);

    const handleShare = React.useCallback(() => {
        if (!summary) {
            return;
        }

        const text = `${summary.title}\n\n${summary.summary}`;
        if (navigator.share) {
            void navigator.share({title: summary.title, text}).catch(() => {});
            return;
        }

        void navigator.clipboard?.writeText(text);
    }, [summary]);

    const handlePresetRange = (range: string) => {
        if (request?.type !== 'channel') {
            return;
        }
        setRangePreset(range);
        setRequest({
            ...request,
            time_range: range,
            since: undefined,
            until: undefined,
        });
    };

    const handleCustomRange = (since: number, until: number) => {
        if (request?.type !== 'channel') {
            return;
        }
        setRangePreset('custom');
        setRequest({
            ...request,
            time_range: undefined,
            since,
            until,
        });
    };

    const headerTitle = summary?.title ?? (target ? (target.type === 'thread' ? 'Thread Summary' : 'Channel Summary') : 'AI Summary');
    const meta = summary ? `${summary.message_count} messages • ${summary.participant_count} participants` : '';

    return (
        <div className='ai-suite-summary'>
            <div className='ai-suite-summary__header'>
                <div>
                    <h2>{headerTitle}</h2>
                    {summary && (
                        <div className='ai-suite-summary__meta'>
                            <span>{meta}</span>
                            <span>{summary.range.label}</span>
                            {summary.cached && (
                                <span className='ai-suite-summary__badge'>Cached</span>
                            )}
                        </div>
                    )}
                </div>
                <button
                    type='button'
                    className='ai-suite-summary__close'
                    onClick={onClose}
                >
                    ×
                </button>
            </div>

            {!target && (
                <div className='ai-suite-summary__empty'>
                    <p>Select “Summarize Thread” or “Summarize Channel” to get started.</p>
                </div>
            )}

            {target && (
                <>
                    {isLoading && !summary && <LoadingSpinner label='Generating summary'/>}

                    {error && (
                        <ErrorMessage
                            message={error}
                            onRetry={refresh}
                        />
                    )}

                    {summary && !error && (
                        <>
                            <SummaryContent summary={summary.summary}/>
                            {summary.limit_reached && (
                                <div className='ai-suite-summary__notice'>
                                    ⚠️ Limited to {summary.context.message_limit} messages
                                </div>
                            )}
                        </>
                    )}

                    <SummaryOptions
                        isChannel={target.type === 'channel'}
                        activeRange={rangePreset}
                        isLoading={isLoading}
                        onPresetRange={handlePresetRange}
                        onCustomRange={handleCustomRange}
                        onCopy={handleCopy}
                        onShare={handleShare}
                        onRegenerate={refresh}
                    />
                </>
            )}
        </div>
    );
};

export default SummaryPanel;


