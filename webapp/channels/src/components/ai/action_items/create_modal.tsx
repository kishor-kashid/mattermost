// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import type {ActionItemCreateRequest, AIActionItemPriority, AIActionItemStatus} from 'types/ai';
import type {GlobalState} from 'types/store';

import {createActionItem} from 'actions/ai_action_items';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/users';

import DateTimePicker from '../common/date_time_picker';

import './create_modal.scss';

interface Props {
    onClose: () => void;
    channelId?: string;
    postId?: string;
}

const CreateActionItemModal: React.FC<Props> = ({onClose, channelId, postId}) => {
    const dispatch = useDispatch();
    const currentUserId = useSelector(getCurrentUserId);
    const currentChannelId = useSelector(getCurrentChannelId);
    const defaultChannelId = channelId || currentChannelId;

    const [description, setDescription] = useState('');
    const [assigneeId, setAssigneeId] = useState(currentUserId);
    const [dueDate, setDueDate] = useState<Date | null>(null);
    const [priority, setPriority] = useState<AIActionItemPriority>('medium');
    const [submitting, setSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!description.trim()) {
            return;
        }

        setSubmitting(true);

        const request: ActionItemCreateRequest = {
            description: description.trim(),
            assignee_id: assigneeId,
            channel_id: defaultChannelId,
            post_id: postId,
            priority,
            due_date: dueDate || undefined,
        };

        const result = await dispatch(createActionItem(request));
        
        setSubmitting(false);

        if ('data' in result) {
            onClose();
        }
    };

    return (
        <div
            className='create-action-item-modal-overlay'
            onClick={onClose}
        >
            <div
                className='create-action-item-modal'
                onClick={(e) => e.stopPropagation()}
            >
                <div className='modal-header'>
                    <h3>
                        <FormattedMessage
                            id='action_items.create_modal.title'
                            defaultMessage='Create Action Item'
                        />
                    </h3>
                    <button
                        className='close-button'
                        onClick={onClose}
                    >
                        <i className='icon icon-close'/>
                    </button>
                </div>

                <form
                    className='modal-body'
                    onSubmit={handleSubmit}
                >
                    <div className='form-group'>
                        <label htmlFor='description'>
                            <FormattedMessage
                                id='action_items.create_modal.description'
                                defaultMessage='Description *'
                            />
                        </label>
                        <textarea
                            id='description'
                            className='form-control'
                            rows={3}
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder='Describe what needs to be done...'
                            required
                            autoFocus
                        />
                    </div>

                    <div className='form-group'>
                        <label htmlFor='assignee'>
                            <FormattedMessage
                                id='action_items.create_modal.assignee'
                                defaultMessage='Assignee'
                            />
                        </label>
                        <input
                            id='assignee'
                            type='text'
                            className='form-control'
                            value={assigneeId}
                            onChange={(e) => setAssigneeId(e.target.value)}
                            placeholder='User ID'
                        />
                        <small className='form-text'>
                            <FormattedMessage
                                id='action_items.create_modal.assignee_hint'
                                defaultMessage='Enter user ID (defaults to you)'
                            />
                        </small>
                    </div>

                    <div className='form-row'>
                        <div className='form-group flex-1'>
                            <label htmlFor='priority'>
                                <FormattedMessage
                                    id='action_items.create_modal.priority'
                                    defaultMessage='Priority'
                                />
                            </label>
                            <select
                                id='priority'
                                className='form-control'
                                value={priority}
                                onChange={(e) => setPriority(e.target.value as AIActionItemPriority)}
                            >
                                <option value='low'>🟢 Low</option>
                                <option value='medium'>🟡 Medium</option>
                                <option value='high'>🔴 High</option>
                                <option value='urgent'>🔥 Urgent</option>
                            </select>
                        </div>

                        <div className='form-group flex-1'>
                            <label htmlFor='due-date'>
                                <FormattedMessage
                                    id='action_items.create_modal.due_date'
                                    defaultMessage='Due Date'
                                />
                            </label>
                            <DateTimePicker
                                value={dueDate}
                                onChange={setDueDate}
                            />
                        </div>
                    </div>

                    <div className='modal-footer'>
                        <button
                            type='button'
                            className='btn btn-secondary'
                            onClick={onClose}
                            disabled={submitting}
                        >
                            <FormattedMessage
                                id='action_items.create_modal.cancel'
                                defaultMessage='Cancel'
                            />
                        </button>
                        <button
                            type='submit'
                            className='btn btn-primary'
                            disabled={!description.trim() || submitting}
                        >
                            {submitting ? (
                                <>
                                    <i className='icon icon-loading icon-spin'/>
                                    <FormattedMessage
                                        id='action_items.create_modal.creating'
                                        defaultMessage='Creating...'
                                    />
                                </>
                            ) : (
                                <>
                                    <i className='icon icon-plus'/>
                                    <FormattedMessage
                                        id='action_items.create_modal.create'
                                        defaultMessage='Create'
                                    />
                                </>
                            )}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default CreateActionItemModal;

