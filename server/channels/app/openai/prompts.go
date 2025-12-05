// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package openai

import (
	"fmt"
	"strings"
)

// PromptType defines the type of AI prompt
type PromptType string

const (
	PromptTypeSummarize        PromptType = "summarize"
	PromptTypeActionItemExtract PromptType = "action_item_extract"
	PromptTypeFormatMessage     PromptType = "format_message"
)

// SummarizationLevel defines the detail level of summarization
type SummarizationLevel string

const (
	SummarizationBrief    SummarizationLevel = "brief"
	SummarizationStandard SummarizationLevel = "standard"
	SummarizationDetailed SummarizationLevel = "detailed"
)

// FormattingProfile defines the tone/style for message formatting
type FormattingProfile string

const (
	FormattingProfessional FormattingProfile = "professional"
	FormattingCasual       FormattingProfile = "casual"
	FormattingTechnical    FormattingProfile = "technical"
	FormattingConcise      FormattingProfile = "concise"
)

// PromptTemplate represents a prompt template with variable substitution
type PromptTemplate struct {
	System string
	User   string
}

// Substitute replaces variables in the template with provided values
func (pt *PromptTemplate) Substitute(variables map[string]string) (system, user string) {
	system = pt.System
	user = pt.User

	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		system = strings.ReplaceAll(system, placeholder, value)
		user = strings.ReplaceAll(user, placeholder, value)
	}

	return system, user
}

// GetSummarizationPrompt returns the prompt template for summarization
func GetSummarizationPrompt(level SummarizationLevel) *PromptTemplate {
	switch level {
	case SummarizationBrief:
		return summaryPromptBrief
	case SummarizationDetailed:
		return summaryPromptDetailed
	default:
		return summaryPromptStandard
	}
}

// GetActionItemExtractionPrompt returns the prompt template for action item extraction
func GetActionItemExtractionPrompt() *PromptTemplate {
	return actionItemExtractionPrompt
}

// GetMessageFormattingPrompt returns the prompt template for message formatting
func GetMessageFormattingPrompt(profile FormattingProfile) *PromptTemplate {
	switch profile {
	case FormattingProfessional:
		return messageFormattingProfessional
	case FormattingCasual:
		return messageFormattingCasual
	case FormattingTechnical:
		return messageFormattingTechnical
	case FormattingConcise:
		return messageFormattingConcise
	default:
		return messageFormattingProfessional
	}
}

// BuildSummarizationUserPrompt builds the user prompt for summarization
func BuildSummarizationUserPrompt(messages, participants string, messageCount int, isThread bool) string {
	contextType := "channel"
	if isThread {
		contextType = "thread"
	}

	variables := map[string]string{
		"context_type":  contextType,
		"messages":      messages,
		"participants":  participants,
		"message_count": fmt.Sprintf("%d", messageCount),
	}

	_, userPrompt := GetSummarizationPrompt(SummarizationStandard).Substitute(variables)
	return userPrompt
}

// BuildActionItemExtractionUserPrompt builds the user prompt for action item extraction
func BuildActionItemExtractionUserPrompt(message, author, channelName string) string {
	variables := map[string]string{
		"message":      message,
		"author":       author,
		"channel_name": channelName,
	}

	_, userPrompt := GetActionItemExtractionPrompt().Substitute(variables)
	return userPrompt
}

// BuildMessageFormattingUserPrompt builds the user prompt for message formatting
func BuildMessageFormattingUserPrompt(message string, profile FormattingProfile) string {
	variables := map[string]string{
		"message": message,
	}

	_, userPrompt := GetMessageFormattingPrompt(profile).Substitute(variables)
	return userPrompt
}

