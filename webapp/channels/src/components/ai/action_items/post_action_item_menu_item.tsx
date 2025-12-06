// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import ReactDOM from 'react-dom';
import {Provider, useSelector} from 'react-redux';
import {IntlProvider} from 'react-intl';
import {FormattedMessage} from 'react-intl';

import {CheckCircleOutlineIcon} from '@mattermost/compass-icons/components';

import * as Menu from 'components/menu';
import store from 'stores/redux_store';
import {getCurrentLocale, getTranslations} from 'selectors/i18n';

import CreateActionItemModal from './create_modal';

interface Props {
    postId: string;
    channelId: string;
    onClose?: () => void;
}

// Global modal container - persists across component lifecycles
let currentModal: {element: HTMLElement; cleanup: () => void} | null = null;

const showCreateModal = (postId: string, channelId: string) => {
    console.log('[PostActionItemMenuItem] showCreateModal called', {postId, channelId});
    
    // Clean up any existing modal
    if (currentModal) {
        currentModal.cleanup();
        currentModal = null;
    }

    // Create a container for the modal
    const modalContainer = document.createElement('div');
    modalContainer.id = `action-item-modal-${postId}`;
    document.body.appendChild(modalContainer);
    console.log('[PostActionItemMenuItem] Modal container created and appended to body');

    const cleanup = () => {
        console.log('[PostActionItemMenuItem] Cleaning up modal');
        ReactDOM.unmountComponentAtNode(modalContainer);
        if (modalContainer.parentNode) {
            modalContainer.parentNode.removeChild(modalContainer);
        }
        currentModal = null;
    };

    // Get current locale and translations from store
    const state = store.getState();
    const locale = getCurrentLocale(state);
    const translations = getTranslations(state, locale);

    // Render the modal wrapped in Redux Provider and IntlProvider
    ReactDOM.render(
        <Provider store={store}>
            <IntlProvider
                locale={locale}
                messages={translations}
            >
                <CreateActionItemModal
                    onClose={cleanup}
                    channelId={channelId}
                    postId={postId}
                />
            </IntlProvider>
        </Provider>,
        modalContainer,
    );

    currentModal = {element: modalContainer, cleanup};
    console.log('[PostActionItemMenuItem] Modal rendered');
};

/**
 * Menu item component for creating action items from posts
 * This can be added to the post dot menu
 */
const PostActionItemMenuItem: React.FC<Props> = ({postId, channelId, onClose}) => {
    const handleClick = () => {
        console.log('[PostActionItemMenuItem] Clicked for post:', postId);
        
        // Open the modal using the global function
        showCreateModal(postId, channelId);
        
        // Close the menu immediately
        if (onClose) {
            onClose();
        }
    };

    return (
        <Menu.Item
            id={`create_action_item_${postId}`}
            data-testid={`create_action_item_${postId}`}
            leadingElement={<CheckCircleOutlineIcon size={18}/>}
            labels={
                <FormattedMessage
                    id='ai.action_items.createFromPost'
                    defaultMessage='Create Action Item'
                />
            }
            onClick={handleClick}
        />
    );
};

export default PostActionItemMenuItem;
