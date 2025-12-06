// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useCallback, useEffect, useMemo, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import type {GlobalState} from 'types/store';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {getFormatterState} from 'selectors/ai_formatter';
import {aiClient} from 'client/ai';
import AIActionTypes from 'utils/constants/ai';

import FormattingMenu from 'components/ai/formatter/formatting_menu';
import PreviewModal from 'components/ai/formatter/preview_modal';

import type {PostDraft} from 'types/store/draft';
import type TextboxClass from 'components/textbox/textbox';

const useFormatter = (
    draft: PostDraft,
    handleDraftChange: ((draft: PostDraft, options: {instant?: boolean; show?: boolean}) => void),
    textboxRef: React.RefObject<TextboxClass>, // kept for future enhancements (cursor/selection-aware prompts)
) => {
    const dispatch = useDispatch();
    const formatterState = useSelector((state: GlobalState) => getFormatterState(state));
    const config = useSelector((state: GlobalState) => getConfig(state));
    const aiSystemState = useSelector((state: GlobalState) => state.entities.ai?.system);
    const [showPreview, setShowPreview] = useState(false);

    // Trigger health check on mount if not already done
    useEffect(() => {
        // Only check if we haven't checked recently (within last 5 minutes)
        const lastCheck = aiSystemState?.health?.lastCheck;
        const fiveMinutesAgo = Date.now() - (5 * 60 * 1000);
        
        if (!lastCheck || lastCheck < fiveMinutesAgo) {
            aiClient.healthCheck().then((response) => {
                dispatch({
                    type: AIActionTypes.AI_HEALTH_CHECK_SUCCESS,
                    data: response,
                });
            }).catch((error) => {
                dispatch({
                    type: AIActionTypes.AI_HEALTH_CHECK_FAILURE,
                });
            });
        }
    }, [dispatch, aiSystemState?.health?.lastCheck]);

    // Check if AI formatting is enabled.
    // We rely primarily on config.AISettings so that the UI can render
    // even if the health check fails or the backend AI service is unavailable.
    const configAISettings = config?.AISettings as any;

    const isFormattingEnabled = useMemo(() => {
        if (!configAISettings) {
            // If there is no AISettings block in config, default to enabled.
            return true;
        }

        const enableFormatting = configAISettings.EnableFormatting;
        const enableGlobal = configAISettings.Enable;

        // Explicit per-feature disable takes precedence.
        if (enableFormatting === false || enableFormatting === 'false' || enableFormatting === 'False') {
            return false;
        }

        // Explicit per-feature enable.
        if (enableFormatting === true || enableFormatting === 'true' || enableFormatting === 'True') {
            return true;
        }

        // Fall back to the global AI enable flag.
        if (enableGlobal === false || enableGlobal === 'false' || enableGlobal === 'False') {
            return false;
        }

        if (enableGlobal === true || enableGlobal === 'true' || enableGlobal === 'True') {
            return true;
        }

        // Default: show the button if we can't determine it's disabled.
        return true;
    }, [configAISettings]);

    const handleFormat = useCallback((formattedText: string) => {
        const updatedDraft = {
            ...draft,
            message: formattedText,
        };
        handleDraftChange(updatedDraft, {instant: true});
    }, [draft, handleDraftChange]);

    const handleShowPreview = useCallback(() => {
        setShowPreview(true);
    }, []);

    const handleClosePreview = useCallback(() => {
        setShowPreview(false);
    }, []);

    const handleApplyFromPreview = useCallback((formattedText: string) => {
        handleFormat(formattedText);
        setShowPreview(false);
    }, [handleFormat]);

    if (!isFormattingEnabled) {
        return {
            additionalControl: null,
            previewModal: null,
        };
    }

    return {
        additionalControl: useMemo(() => (
            <FormattingMenu
                message={draft.message}
                onFormat={handleFormat}
                onShowPreview={handleShowPreview}
                disabled={formatterState.formatting || formatterState.loading}
            />
        ), [
            draft.message,
            handleFormat,
            handleShowPreview,
            formatterState.formatting,
            formatterState.loading,
        ]),
        previewModal: useMemo(() => (
            <PreviewModal
                show={showPreview && formatterState.preview !== null}
                originalMessage={draft.message}
                onClose={handleClosePreview}
                onApply={handleApplyFromPreview}
            />
        ), [
            showPreview,
            formatterState.preview,
            draft.message,
            handleClosePreview,
            handleApplyFromPreview,
        ]),
    };
};

export default useFormatter;

