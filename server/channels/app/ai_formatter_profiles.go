// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// GetFormattingProfileMetadata returns metadata for all available formatting profiles
func GetFormattingProfileMetadata() map[openai.FormattingProfile]FormattingProfileInfo {
	return map[openai.FormattingProfile]FormattingProfileInfo{
		openai.FormattingProfessional: {
			Id:          string(openai.FormattingProfessional),
			Label:       "Professional",
			Description: "Improve grammar, clarity, and structure for professional business communication",
		},
		openai.FormattingCasual: {
			Id:          string(openai.FormattingCasual),
			Label:       "Casual",
			Description: "Improve message quality while maintaining a casual, friendly tone",
		},
		openai.FormattingTechnical: {
			Id:          string(openai.FormattingTechnical),
			Label:       "Technical",
			Description: "Improve technical communication with proper terminology and formatting",
		},
		openai.FormattingConcise: {
			Id:          string(openai.FormattingConcise),
			Label:       "Concise",
			Description: "Make messages more concise while preserving meaning",
		},
	}
}

// IsValidFormattingProfile checks if a profile ID is valid
func IsValidFormattingProfile(profileID string) bool {
	_, exists := GetFormattingProfileMetadata()[openai.FormattingProfile(profileID)]
	return exists
}

// GetDefaultFormattingProfile returns the default formatting profile
func GetDefaultFormattingProfile() openai.FormattingProfile {
	return openai.FormattingProfessional
}

// FormattingProfile is an alias for openai.FormattingProfile to use in app package
type FormattingProfile = openai.FormattingProfile

