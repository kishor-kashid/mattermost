// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import type {GlobalState} from 'types/store';

import {clearFormatPreview} from 'actions/ai_formatter';
import {getFormatPreview, isFormatting} from 'selectors/ai_formatter';

import DiffView from './diff_view';

import './preview_modal.scss';

interface Props {
    show: boolean;
    originalMessage: string;
    onClose: () => void;
    onApply: (formattedText: string) => void;
}

const PreviewModal: React.FC<Props> = ({show, originalMessage, onClose, onApply}) => {
    const dispatch = useDispatch();
    const preview = useSelector((state: GlobalState) => getFormatPreview(state));
    const formatting = useSelector((state: GlobalState) => isFormatting(state));
    const [viewMode, setViewMode] = useState<'side-by-side' | 'diff'>('side-by-side');

    if (!show || !preview) {
        return null;
    }

    const handleApply = () => {
        if (!preview || formatting) {
            return;
        }

        // Use the formatted text from preview directly
        onApply(preview.formatted_text);
        dispatch(clearFormatPreview());
        onClose();
    };

    const handleCopy = () => {
        if (preview?.formatted_text) {
            navigator.clipboard.writeText(preview.formatted_text);
        }
    };

    const handleDismiss = () => {
        dispatch(clearFormatPreview());
        onClose();
    };

    return (
        <div
            className='ai-preview-modal-overlay'
            onClick={handleDismiss}
        >
            <div
                className='ai-preview-modal'
                onClick={(e) => e.stopPropagation()}
            >
                <div className='ai-preview-modal-header'>
                    <h3>
                        <FormattedMessage
                            id='ai.formatter.preview.title'
                            defaultMessage='Formatting Preview'
                        />
                    </h3>
                    <div className='ai-preview-modal-actions'>
                        <button
                            className='ai-preview-view-toggle'
                            onClick={() => setViewMode(viewMode === 'side-by-side' ? 'diff' : 'side-by-side')}
                        >
                            {viewMode === 'side-by-side' ? (
                                <FormattedMessage
                                    id='ai.formatter.preview.showDiff'
                                    defaultMessage='Show Diff'
                                />
                            ) : (
                                <FormattedMessage
                                    id='ai.formatter.preview.showSideBySide'
                                    defaultMessage='Show Side-by-Side'
                                />
                            )}
                        </button>
                        <button
                            className='ai-preview-close-button'
                            onClick={handleDismiss}
                        >
                            <i className='icon icon-close'/>
                        </button>
                    </div>
                </div>

                <div className='ai-preview-modal-body'>
                    {viewMode === 'side-by-side' ? (
                        <div className='ai-preview-side-by-side'>
                            <div className='ai-preview-original'>
                                <div className='ai-preview-section-header'>
                                    <FormattedMessage
                                        id='ai.formatter.preview.original'
                                        defaultMessage='Original'
                                    />
                                </div>
                                <div className='ai-preview-text'>
                                    {originalMessage}
                                </div>
                            </div>
                            <div className='ai-preview-formatted'>
                                <div className='ai-preview-section-header'>
                                    <FormattedMessage
                                        id='ai.formatter.preview.formatted'
                                        defaultMessage='Formatted'
                                    />
                                </div>
                                <div className='ai-preview-text'>
                                    {preview.formatted_text}
                                </div>
                            </div>
                        </div>
                    ) : (
                        <div className='ai-preview-diff'>
                            <DiffView
                                original={originalMessage}
                                formatted={preview.formatted_text}
                                diff={preview.diff}
                            />
                        </div>
                    )}
                </div>

                <div className='ai-preview-modal-footer'>
                    <button
                        className='ai-preview-button ai-preview-button-secondary'
                        onClick={handleCopy}
                    >
                        <i className='icon icon-content-copy'/>
                        <FormattedMessage
                            id='ai.formatter.preview.copy'
                            defaultMessage='Copy'
                        />
                    </button>
                    <div className='ai-preview-button-group'>
                        <button
                            className='ai-preview-button ai-preview-button-secondary'
                            onClick={handleDismiss}
                        >
                            <FormattedMessage
                                id='ai.formatter.preview.dismiss'
                                defaultMessage='Dismiss'
                            />
                        </button>
                        <button
                            className='ai-preview-button ai-preview-button-primary'
                            onClick={handleApply}
                            disabled={formatting}
                        >
                            <FormattedMessage
                                id='ai.formatter.preview.apply'
                                defaultMessage='Apply'
                            />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default PreviewModal;

