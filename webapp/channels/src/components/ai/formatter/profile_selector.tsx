// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {useSelector} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import type {FormattingProfileInfo} from 'types/ai';
import type {GlobalState} from 'types/store';

import {getFormattingProfiles} from 'selectors/ai_formatter';

import './profile_selector.scss';

interface Props {
    selectedProfile?: string;
    onProfileChange: (profileId: string) => void;
    showCustomInstructions?: boolean;
    onCustomInstructionsChange?: (instructions: string) => void;
}

const ProfileSelector: React.FC<Props> = ({
    selectedProfile,
    onProfileChange,
    showCustomInstructions = false,
    onCustomInstructionsChange,
}) => {
    const profiles = useSelector((state: GlobalState) => getFormattingProfiles(state));
    const [customInstructions, setCustomInstructions] = useState('');

    const handleCustomInstructionsChange = (value: string) => {
        setCustomInstructions(value);
        if (onCustomInstructionsChange) {
            onCustomInstructionsChange(value);
        }
    };

    if (profiles.length === 0) {
        return null;
    }

    return (
        <div className='ai-profile-selector'>
            <div className='ai-profile-selector-label'>
                <FormattedMessage
                    id='ai.formatter.profile.label'
                    defaultMessage='Formatting Profile'
                />
            </div>
            <div className='ai-profile-selector-options'>
                {profiles.map((profile: FormattingProfileInfo) => (
                    <label
                        key={profile.id}
                        className={`ai-profile-option ${selectedProfile === profile.id ? 'selected' : ''}`}
                    >
                        <input
                            type='radio'
                            name='formatting-profile'
                            value={profile.id}
                            checked={selectedProfile === profile.id}
                            onChange={() => onProfileChange(profile.id)}
                        />
                        <div className='ai-profile-option-content'>
                            <div className='ai-profile-option-label'>{profile.label}</div>
                            <div className='ai-profile-option-description'>{profile.description}</div>
                        </div>
                    </label>
                ))}
            </div>
            {showCustomInstructions && (
                <div className='ai-profile-custom-instructions'>
                    <label>
                        <FormattedMessage
                            id='ai.formatter.customInstructions.label'
                            defaultMessage='Custom Instructions (optional)'
                        />
                    </label>
                    <textarea
                        value={customInstructions}
                        onChange={(e) => handleCustomInstructionsChange(e.target.value)}
                        placeholder='Add any specific formatting instructions...'
                        rows={3}
                    />
                </div>
            )}
        </div>
    );
};

export default ProfileSelector;

