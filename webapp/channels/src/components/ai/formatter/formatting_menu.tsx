// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import type {FormattingProfileInfo} from 'types/ai';
import type {GlobalState} from 'types/store';

import {getFormattingProfiles, formatPreview} from 'actions/ai_formatter';
import {getFormattingProfiles as getProfilesSelector, isFormatting} from 'selectors/ai_formatter';

import './formatting_menu.scss';

interface Props {
    message: string;
    onFormat: (formattedText: string) => void;
    onShowPreview: () => void;
    disabled?: boolean;
}

const FormattingMenu: React.FC<Props> = ({message, onFormat, onShowPreview, disabled}) => {
    const dispatch = useDispatch();
    const profiles = useSelector((state: GlobalState) => getProfilesSelector(state));
    const formatting = useSelector((state: GlobalState) => isFormatting(state));
    const [showMenu, setShowMenu] = useState(false);
    const [selectedProfile, setSelectedProfile] = useState<string | null>(null);

    // Debug logging
    console.log('[FormattingMenu] State:', {
        disabled,
        formatting,
        profilesCount: profiles.length,
        messageLength: message.length,
        messageTrimmed: message.trim().length,
    });

    useEffect(() => {
        // Always try to load profiles on mount
        if (profiles.length === 0 && !formatting) {
            dispatch(getFormattingProfiles());
        }
    }, [profiles.length, formatting, dispatch]);

    const handleProfileClick = async (profile: FormattingProfileInfo) => {
        if (!message.trim() || formatting) {
            return;
        }

        setSelectedProfile(profile.id);
        setShowMenu(false);

        // Show preview modal
        const result = await dispatch(formatPreview(message, profile.id));
        if ('data' in result) {
            onShowPreview();
        }
    };

    return (
        <div className='ai-formatting-menu-container'>
            <button
                className='ai-formatting-button'
                onClick={() => {
                    // Load profiles if not loaded yet
                    if (profiles.length === 0) {
                        dispatch(getFormattingProfiles());
                    }
                    setShowMenu(!showMenu);
                }}
                disabled={disabled || formatting || !message.trim()}
                title={!message.trim() ? 'Type a message to format' : (profiles.length === 0 ? 'Loading formatting profiles...' : 'AI Formatting')}
                aria-label='AI Formatting'
            >
                <i className='icon icon-robot'/>
                {formatting && (
                    <span className='ai-formatting-spinner'>
                        <i className='icon icon-spinner icon-spin'/>
                    </span>
                )}
            </button>

            {showMenu && (
                <>
                    <div
                        className='ai-formatting-menu-overlay'
                        onClick={() => setShowMenu(false)}
                    />
                    <div className='ai-formatting-menu'>
                        <div className='ai-formatting-menu-header'>
                            <FormattedMessage
                                id='ai.formatter.menu.title'
                                defaultMessage='Format Message'
                            />
                        </div>
                        <div className='ai-formatting-menu-items'>
                            {profiles.length === 0 ? (
                                <div className='ai-formatting-menu-loading'>
                                    <i className='icon icon-spinner icon-spin'/>
                                    <span>Loading formatting profiles...</span>
                                </div>
                            ) : (
                                profiles.map((profile) => (
                                    <button
                                        key={profile.id}
                                        className='ai-formatting-menu-item'
                                        onClick={() => handleProfileClick(profile)}
                                        disabled={formatting}
                                    >
                                        <div className='ai-formatting-menu-item-label'>
                                            {profile.label}
                                        </div>
                                        <div className='ai-formatting-menu-item-description'>
                                            {profile.description}
                                        </div>
                                    </button>
                                ))
                            )}
                        </div>
                    </div>
                </>
            )}
        </div>
    );
};

export default FormattingMenu;

