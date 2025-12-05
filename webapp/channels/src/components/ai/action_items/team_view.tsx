// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useMemo} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import type {AIActionItem} from 'types/ai';
import type {GlobalState} from 'types/store';

import {getActionItems} from 'actions/ai_action_items';
import {getActiveActionItems} from 'selectors/ai_action_items';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import ActionItemCard from './action_item_card';

import './team_view.scss';

const TeamActionItemsView: React.FC = () => {
    const dispatch = useDispatch();
    const currentChannelId = useSelector(getCurrentChannelId);
    const allItems = useSelector(getActiveActionItems);
    const loading = useSelector((state: GlobalState) => state.ai.actionItems.loading);

    useEffect(() => {
        if (currentChannelId) {
            dispatch(getActionItems({channelId: currentChannelId}));
        }
    }, [dispatch, currentChannelId]);

    // Group items by assignee
    const itemsByAssignee = useMemo(() => {
        const grouped: Record<string, AIActionItem[]> = {};
        
        allItems.forEach((item) => {
            const assigneeId = item.assignee_id || 'unassigned';
            if (!grouped[assigneeId]) {
                grouped[assigneeId] = [];
            }
            grouped[assigneeId].push(item);
        });

        return grouped;
    }, [allItems]);

    const exportToCSV = () => {
        const headers = ['Description', 'Assignee ID', 'Priority', 'Status', 'Due Date', 'Created At'];
        const rows = allItems.map((item) => [
            item.description,
            item.assignee_id || '',
            item.priority,
            item.status,
            item.due_date ? new Date(item.due_date).toISOString() : '',
            new Date(item.create_at).toISOString(),
        ]);

        const csvContent = [
            headers.join(','),
            ...rows.map((row) => row.map((cell) => `"${cell}"`).join(',')),
        ].join('\n');

        const blob = new Blob([csvContent], {type: 'text/csv;charset=utf-8;'});
        const link = document.createElement('a');
        const url = URL.createObjectURL(blob);
        
        link.setAttribute('href', url);
        link.setAttribute('download', `action_items_${Date.now()}.csv`);
        link.style.visibility = 'hidden';
        
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    if (loading) {
        return (
            <div className='team-action-items-view'>
                <div className='loading-spinner'>
                    <i className='icon icon-loading icon-spin'/>
                    <FormattedMessage
                        id='action_items.loading'
                        defaultMessage='Loading action items...'
                    />
                </div>
            </div>
        );
    }

    return (
        <div className='team-action-items-view'>
            <div className='team-view-header'>
                <h2>
                    <FormattedMessage
                        id='action_items.team_view.title'
                        defaultMessage='Team Action Items'
                    />
                </h2>
                <button
                    className='btn btn-secondary'
                    onClick={exportToCSV}
                    disabled={allItems.length === 0}
                >
                    <i className='icon icon-download-outline'/>
                    <FormattedMessage
                        id='action_items.export_csv'
                        defaultMessage='Export CSV'
                    />
                </button>
            </div>

            {allItems.length === 0 ? (
                <div className='empty-state'>
                    <i className='icon icon-account-multiple-outline'/>
                    <p>
                        <FormattedMessage
                            id='action_items.team_view.no_items'
                            defaultMessage='No team action items found'
                        />
                    </p>
                </div>
            ) : (
                <div className='team-view-content'>
                    {Object.entries(itemsByAssignee).map(([assigneeId, items]) => (
                        <div
                            key={assigneeId}
                            className='assignee-section'
                        >
                            <h3 className='assignee-title'>
                                <i className='icon icon-account-outline'/>
                                <span>
                                    {assigneeId === 'unassigned' ? 'Unassigned' : `User ${assigneeId.substring(0, 8)}`}
                                </span>
                                <span className='item-count'>{items.length}</span>
                            </h3>
                            <div className='assignee-items'>
                                {items.map((item) => (
                                    <ActionItemCard
                                        key={item.id}
                                        item={item}
                                    />
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default TeamActionItemsView;

