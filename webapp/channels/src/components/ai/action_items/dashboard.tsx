// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import type {AIActionItem} from 'types/ai';
import type {GlobalState} from 'types/store';

import {getActionItems, completeActionItem, deleteActionItem} from 'actions/ai_action_items';
import {getOverdueActionItems, getDueSoonActionItems, getActiveActionItems, getCompletedActionItems} from 'selectors/ai_action_items';

import ActionItemCard from './action_item_card';
import CreateActionItemModal from './create_modal';

import './dashboard.scss';

const ActionItemsDashboard: React.FC = () => {
    const dispatch = useDispatch();
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [filter, setFilter] = useState('all');
    const [showCompleted, setShowCompleted] = useState(false);

    const overdueItems = useSelector(getOverdueActionItems);
    const dueSoonItems = useSelector(getDueSoonActionItems);
    const activeItems = useSelector(getActiveActionItems);
    const completedItems = useSelector(getCompletedActionItems);
    const loading = useSelector((state: GlobalState) => state.ai.actionItems.loading);
    const allItems = useSelector((state: GlobalState) => state.ai.actionItems.items);
    const allItemsArray = Object.values(allItems);

    console.log('[ActionItemsDashboard] Rendering:', {
        overdueCount: overdueItems.length,
        dueSoonCount: dueSoonItems.length,
        activeCount: activeItems.length,
        completedCount: completedItems.length,
        totalItems: Object.keys(allItems).length,
        allItemsArray: allItemsArray.length,
        loading,
        firstItem: allItemsArray[0],
    });

    useEffect(() => {
        console.log('[ActionItemsDashboard] Component mounted - fetching action items');
        console.log('[ActionItemsDashboard] showCompleted:', showCompleted);
        console.log('[ActionItemsDashboard] Current items in Redux before fetch:', Object.keys(allItems).length);
        dispatch(getActionItems({includeCompleted: showCompleted}));
    }, [dispatch, showCompleted]);

    // Log whenever allItems changes
    useEffect(() => {
        console.log('[ActionItemsDashboard] allItems changed:', allItems);
        console.log('[ActionItemsDashboard] allItemsArray:', allItemsArray);
    }, [allItems, allItemsArray]);

    const handleComplete = (itemId: string) => {
        dispatch(completeActionItem(itemId));
    };

    const handleDelete = (itemId: string) => {
        if (confirm('Are you sure you want to delete this action item?')) {
            dispatch(deleteActionItem(itemId));
        }
    };

    const handleEdit = (item: AIActionItem) => {
        // TODO: Implement edit modal
        console.log('Edit item:', item);
    };

    const renderActionItems = (items: AIActionItem[], title: string, icon: string) => {
        console.log(`[renderActionItems] ${title}:`, items.length, items);
        
        if (items.length === 0) {
            return null;
        }

        return (
            <div className='action-items-section'>
                <h3 className='section-title'>
                    <span className='icon'>{icon}</span>
                    {title}
                    <span className='count'>{items.length}</span>
                </h3>
                <div className='action-items-list'>
                    {items.map((item) => (
                        <ActionItemCard
                            key={item.id}
                            item={item}
                            onComplete={handleComplete}
                            onEdit={handleEdit}
                            onDelete={handleDelete}
                        />
                    ))}
                </div>
            </div>
        );
    };

    return (
        <div className='action-items-dashboard'>
            <div className='dashboard-header'>
                <h2>
                    <FormattedMessage
                        id='action_items.dashboard.title'
                        defaultMessage='Action Items'
                    />
                </h2>
                <div className='dashboard-actions'>
                    <label className='checkbox-label'>
                        <input
                            type='checkbox'
                            checked={showCompleted}
                            onChange={(e) => setShowCompleted(e.target.checked)}
                        />
                        <FormattedMessage
                            id='action_items.show_completed'
                            defaultMessage='Show completed'
                        />
                    </label>
                    <button
                        className='btn btn-primary'
                        onClick={() => setShowCreateModal(true)}
                    >
                        <i className='icon icon-plus'/>
                        <FormattedMessage
                            id='action_items.create_new'
                            defaultMessage='Create New'
                        />
                    </button>
                </div>
            </div>

            {loading ? (
                <div className='loading-spinner'>
                    <i className='icon icon-loading icon-spin'/>
                    <FormattedMessage
                        id='action_items.loading'
                        defaultMessage='Loading action items...'
                    />
                </div>
            ) : (
                <div className='dashboard-content'>
                    {/* Debug info - ALWAYS show during debugging */}
                    <div style={{padding: '10px', background: '#fffbe6', border: '1px solid #ffe58f', marginBottom: '10px', fontSize: '12px', borderRadius: '4px'}}>
                        <strong>🔍 Debug Info:</strong>
                        <div>Total items in Redux: <strong>{allItemsArray.length}</strong></div>
                        <div>Redux items object keys: {Object.keys(allItems).length}</div>
                        <div>Overdue: {overdueItems.length}</div>
                        <div>Due Soon: {dueSoonItems.length}</div>
                        <div>Active: {activeItems.length}</div>
                        <div>Completed: {completedItems.length}</div>
                        {allItemsArray.length > 0 ? (
                            <>
                                <div>First item ID: {allItemsArray[0]?.id}</div>
                                <div>First item status: {allItemsArray[0]?.status}</div>
                                <div>First item description: {allItemsArray[0]?.description?.substring(0, 50)}...</div>
                                <div>First item assignee_id: {allItemsArray[0]?.assignee_id}</div>
                                <div>First item created_by: {allItemsArray[0]?.created_by}</div>
                            </>
                        ) : (
                            <div style={{color: 'red'}}>⚠️ No items in Redux state!</div>
                        )}
                    </div>
                    
                    {renderActionItems(overdueItems, 'Overdue', '🔴')}
                    {renderActionItems(dueSoonItems, 'Due Soon', '⏰')}
                    {renderActionItems(
                        activeItems.filter((item) => !overdueItems.includes(item) && !dueSoonItems.includes(item)),
                        'Active',
                        '📋'
                    )}
                    {showCompleted && renderActionItems(completedItems, 'Completed', '✅')}

                    {/* Show all items if none are categorized */}
                    {allItemsArray.length > 0 && activeItems.length === 0 && completedItems.length === 0 && (
                        <div className='action-items-section'>
                            <h3 className='section-title'>
                                <span className='icon'>📋</span>
                                All Items (Uncategorized)
                                <span className='count'>{allItemsArray.length}</span>
                            </h3>
                            <div className='action-items-list'>
                                {allItemsArray.map((item) => (
                                    <ActionItemCard
                                        key={item.id}
                                        item={item}
                                        onComplete={handleComplete}
                                        onEdit={handleEdit}
                                        onDelete={handleDelete}
                                    />
                                ))}
                            </div>
                        </div>
                    )}

                    {allItemsArray.length === 0 && (
                        <div className='empty-state'>
                            <i className='icon icon-checkbox-marked-circle-outline'/>
                            <p>
                                <FormattedMessage
                                    id='action_items.no_items'
                                    defaultMessage='No action items yet'
                                />
                            </p>
                            <p className='subtitle'>
                                <FormattedMessage
                                    id='action_items.no_items_subtitle'
                                    defaultMessage='Action items will appear here when detected in messages or created manually'
                                />
                            </p>
                        </div>
                    )}
                </div>
            )}

            {showCreateModal && (
                <CreateActionItemModal
                    onClose={() => setShowCreateModal(false)}
                />
            )}
        </div>
    );
};

export default ActionItemsDashboard;

