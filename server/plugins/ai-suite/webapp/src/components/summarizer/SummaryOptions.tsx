import React from 'react';

import Modal from '../common/Modal';

type Props = {
    isChannel: boolean;
    activeRange?: string;
    isLoading: boolean;
    onPresetRange?: (range: string) => void;
    onCustomRange?: (since: number, until: number) => void;
    onCopy: () => void;
    onShare: () => void;
    onRegenerate: () => void;
};

const RANGE_OPTIONS = [
    {label: '24h', value: '24h'},
    {label: '3d', value: '3d'},
    {label: '7d', value: '7d'},
];

export const SummaryOptions = ({
    isChannel,
    activeRange = '24h',
    isLoading,
    onPresetRange,
    onCustomRange,
    onCopy,
    onShare,
    onRegenerate,
}: Props) => {
    const [showModal, setShowModal] = React.useState(false);
    const [customStart, setCustomStart] = React.useState('');
    const [customEnd, setCustomEnd] = React.useState('');

    const applyCustomRange = () => {
        if (!customStart || !customEnd || !onCustomRange) {
            return;
        }
        const since = Date.parse(customStart);
        const until = Date.parse(customEnd);
        if (Number.isNaN(since) || Number.isNaN(until) || since >= until) {
            return;
        }

        onCustomRange(since, until);
        setShowModal(false);
    };

    return (
        <div className='ai-suite-summary__controls'>
            {isChannel && (
                <div className='ai-suite-summary__ranges'>
                    {RANGE_OPTIONS.map((option) => (
                        <button
                            key={option.value}
                            type='button'
                            className={option.value === activeRange ? 'ai-suite-summary__range ai-suite-summary__range--active' : 'ai-suite-summary__range'}
                            onClick={() => onPresetRange?.(option.value)}
                        >
                            {option.label}
                        </button>
                    ))}
                    <button
                        type='button'
                        className='ai-suite-summary__range'
                        onClick={() => setShowModal(true)}
                    >
                        Custom
                    </button>
                </div>
            )}

            <div className='ai-suite-summary__actions'>
                <button
                    type='button'
                    onClick={onCopy}
                >
                    Copy
                </button>
                <button
                    type='button'
                    onClick={onShare}
                >
                    Share
                </button>
                <button
                    type='button'
                    onClick={onRegenerate}
                    disabled={isLoading}
                >
                    {isLoading ? 'Generating…' : 'Regenerate'}
                </button>
            </div>

            {showModal && (
                <Modal
                    title='Select range'
                    onClose={() => setShowModal(false)}
                    footer={(
                        <>
                            <button
                                type='button'
                                className='ai-suite-button ai-suite-button--ghost'
                                onClick={() => setShowModal(false)}
                            >
                                Cancel
                            </button>
                            <button
                                type='button'
                                className='ai-suite-button'
                                onClick={applyCustomRange}
                                disabled={!customStart || !customEnd}
                            >
                                Apply
                            </button>
                        </>
                    )}
                >
                    <div className='ai-suite-summary__modal-fields'>
                        <label>
                            Start
                            <input
                                type='datetime-local'
                                value={customStart}
                                onChange={(event) => setCustomStart(event.target.value)}
                            />
                        </label>
                        <label>
                            End
                            <input
                                type='datetime-local'
                                value={customEnd}
                                onChange={(event) => setCustomEnd(event.target.value)}
                            />
                        </label>
                    </div>
                </Modal>
            )}
        </div>
    );
};

export default SummaryOptions;


