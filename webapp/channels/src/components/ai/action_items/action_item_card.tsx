// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage, FormattedDate} from 'react-intl';

import type {AIActionItem} from 'types/ai';

import './action_item_card.scss';

interface Props {
    item: AIActionItem;
    onComplete?: (itemId: string) => void;
    onEdit?: (item: AIActionItem) => void;
    onDelete?: (itemId: string) => void;
}

const ActionItemCard: React.FC<Props> = ({item, onComplete, onEdit, onDelete}) => {
    const getPriorityClass = (priority: string) => {
        return `priority-${priority}`;
    };

    const getPriorityIcon = (priority: string) => {
        switch (priority) {
        case 'urgent':
            return '🔥';
        case 'high':
            return '🔴';
        case 'medium':
            return '🟡';
        case 'low':
            return '🟢';
        default:
            return '⚪';
        }
    };

    const getStatusBadge = (status: string) => {
        const badges: Record<string, {text: string; class: string}> = {
            open: {text: 'Open', class: 'status-open'},
            in_progress: {text: 'In Progress', class: 'status-in-progress'},
            completed: {text: 'Completed', class: 'status-completed'},
            dismissed: {text: 'Dismissed', class: 'status-dismissed'},
        };
        return badges[status] || {text: status, class: ''};
    };

    const isOverdue = item.due_date && item.due_date < Date.now() && item.status !== 'completed';
    const isDueSoon = item.due_date && item.due_date > Date.now() && item.due_date < Date.now() + (7 * 24 * 60 * 60 * 1000);

    const statusBadge = getStatusBadge(item.status);

    return (
        <div className={`action-item-card ${getPriorityClass(item.priority)}`}>
            <div className='action-item-card__header'>
                <div className='action-item-card__priority'>
                    <span className='priority-icon'>{getPriorityIcon(item.priority)}</span>
                    <span className='priority-text'>{item.priority}</span>
                </div>
                <div className='action-item-card__status'>
                    <span className={`status-badge ${statusBadge.class}`}>
                        {statusBadge.text}
                    </span>
                </div>
            </div>

            <div className='action-item-card__body'>
                <p className='action-item-card__description'>
                    {item.description}
                </p>
            </div>

            <div className='action-item-card__footer'>
                <div className='action-item-card__meta'>
                    {item.due_date && (
                        <div className={`due-date ${isOverdue ? 'overdue' : ''} ${isDueSoon ? 'due-soon' : ''}`}>
                            {isOverdue && <span className='overdue-icon'>⏰</span>}
                            <span>
                                <FormattedMessage
                                    id='action_items.due'
                                    defaultMessage='Due: '
                                />
                                <FormattedDate
                                    value={new Date(item.due_date)}
                                    year='numeric'
                                    month='short'
                                    day='2-digit'
                                />
                            </span>
                        </div>
                    )}
                    {item.post_id && (
                        <a
                            href={`/_redirect/pl/${item.post_id}`}
                            className='source-link'
                        >
                            <i className='icon icon-message-text-outline'/>
                            <FormattedMessage
                                id='action_items.view_source'
                                defaultMessage='View in context'
                            />
                        </a>
                    )}
                </div>

                {item.status !== 'completed' && item.status !== 'dismissed' && (
                    <div className='action-item-card__actions'>
                        {onComplete && (
                            <button
                                className='btn btn-sm btn-primary'
                                onClick={() => onComplete(item.id)}
                            >
                                <i className='icon icon-check'/>
                                <FormattedMessage
                                    id='action_items.complete'
                                    defaultMessage='Complete'
                                />
                            </button>
                        )}
                        {onEdit && (
                            <button
                                className='btn btn-sm btn-secondary'
                                onClick={() => onEdit(item)}
                            >
                                <i className='icon icon-pencil-outline'/>
                            </button>
                        )}
                        {onDelete && (
                            <button
                                className='btn btn-sm btn-danger'
                                onClick={() => onDelete(item.id)}
                            >
                                <i className='icon icon-trash-can-outline'/>
                            </button>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default ActionItemCard;

