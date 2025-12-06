// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import type {TextChange} from 'types/ai';

import './diff_view.scss';

interface Props {
    original: string;
    formatted: string;
    diff?: {
        original: string;
        formatted: string;
        changes?: TextChange[];
    };
}

const DiffView: React.FC<Props> = ({original, formatted, diff}) => {
    if (!diff || !diff.changes || diff.changes.length === 0) {
        // Simple character-by-character diff if no structured diff available
        return (
            <div className='ai-diff-view'>
                <div className='ai-diff-content'>
                    {renderSimpleDiff(original, formatted)}
                </div>
            </div>
        );
    }

    // Render structured diff
    const lines = original.split('\n');
    const formattedLines = formatted.split('\n');
    let originalIndex = 0;
    let formattedIndex = 0;

    return (
        <div className='ai-diff-view'>
            <div className='ai-diff-content'>
                {diff.changes.map((change, index) => {
                    const beforeText = original.substring(originalIndex, change.start);
                    const afterText = original.substring(change.end);
                    originalIndex = change.end;

                    return (
                        <React.Fragment key={index}>
                            {beforeText && (
                                <span className='ai-diff-unchanged'>{beforeText}</span>
                            )}
                            {change.type === 'delete' && (
                                <span className='ai-diff-deleted'>{change.old_text}</span>
                            )}
                            {change.type === 'insert' && (
                                <span className='ai-diff-inserted'>{change.new_text}</span>
                            )}
                            {change.type === 'replace' && (
                                <>
                                    <span className='ai-diff-deleted'>{change.old_text}</span>
                                    <span className='ai-diff-inserted'>{change.new_text}</span>
                                </>
                            )}
                        </React.Fragment>
                    );
                })}
                {originalIndex < original.length && (
                    <span className='ai-diff-unchanged'>{original.substring(originalIndex)}</span>
                )}
            </div>
        </div>
    );
};

function renderSimpleDiff(original: string, formatted: string): React.ReactNode {
    const maxLength = Math.max(original.length, formatted.length);
    const result: React.ReactNode[] = [];
    let i = 0;

    while (i < maxLength) {
        if (i >= original.length) {
            // Only in formatted
            result.push(
                <span key={i} className='ai-diff-inserted'>
                    {formatted.substring(i)}
                </span>
            );
            break;
        } else if (i >= formatted.length) {
            // Only in original
            result.push(
                <span key={i} className='ai-diff-deleted'>
                    {original.substring(i)}
                </span>
            );
            break;
        } else if (original[i] === formatted[i]) {
            // Find the next difference
            let j = i;
            while (j < maxLength && original[j] === formatted[j]) {
                j++;
            }
            result.push(
                <span key={i} className='ai-diff-unchanged'>
                    {original.substring(i, j)}
                </span>
            );
            i = j;
        } else {
            // Find the next match
            let j = i + 1;
            while (j < maxLength && original[j] !== formatted[j]) {
                j++;
            }
            result.push(
                <span key={i} className='ai-diff-deleted'>
                    {original.substring(i, j)}
                </span>
            );
            result.push(
                <span key={`${i}-new`} className='ai-diff-inserted'>
                    {formatted.substring(i, j)}
                </span>
            );
            i = j;
        }
    }

    return <>{result}</>;
}

export default DiffView;

