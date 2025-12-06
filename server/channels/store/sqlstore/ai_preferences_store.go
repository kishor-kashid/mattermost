// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"database/sql"

	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/store"
)

type SqlAIPreferencesStore struct {
	*SqlStore
}

func newSqlAIPreferencesStore(sqlStore *SqlStore) store.AIPreferencesStore {
	return &SqlAIPreferencesStore{
		SqlStore: sqlStore,
	}
}

func (s *SqlAIPreferencesStore) Save(preferences *model.AIPreferences) (*model.AIPreferences, error) {
	preferences.PreSave()

	if err := preferences.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Insert("AIPreferences").
		Columns(
			"Id", "UserId", "EnableSummarization", "EnableAnalytics", "EnableActionItems",
			"EnableFormatting", "DefaultModel", "FormattingProfile", "CreateAt", "UpdateAt",
		).
		Values(
			preferences.Id, preferences.UserId, preferences.EnableSummarization, preferences.EnableAnalytics, preferences.EnableActionItems,
			preferences.EnableFormatting, preferences.DefaultModel, preferences.FormattingProfile, preferences.CreateAt, preferences.UpdateAt,
		)

	if _, err := s.GetMaster().ExecBuilder(query); err != nil {
		return nil, errors.Wrap(err, "failed to save AIPreferences")
	}

	return preferences, nil
}

func (s *SqlAIPreferencesStore) Get(id string) (*model.AIPreferences, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AIPreferences").
		Where(sq.Eq{"Id": id})

	var preferences model.AIPreferences
	err := s.GetReplica().GetBuilder(&preferences, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AIPreferences", id)
		}
		return nil, errors.Wrapf(err, "failed to find AIPreferences with id=%s", id)
	}

	return &preferences, nil
}

func (s *SqlAIPreferencesStore) GetByUser(userId string) (*model.AIPreferences, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AIPreferences").
		Where(sq.Eq{"UserId": userId})

	var preferences model.AIPreferences
	err := s.GetReplica().GetBuilder(&preferences, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AIPreferences", userId)
		}
		return nil, errors.Wrapf(err, "failed to find AIPreferences for userId=%s", userId)
	}

	return &preferences, nil
}

func (s *SqlAIPreferencesStore) Update(preferences *model.AIPreferences) (*model.AIPreferences, error) {
	preferences.PreUpdate()

	if err := preferences.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Update("AIPreferences").
		Set("EnableSummarization", preferences.EnableSummarization).
		Set("EnableAnalytics", preferences.EnableAnalytics).
		Set("EnableActionItems", preferences.EnableActionItems).
		Set("EnableFormatting", preferences.EnableFormatting).
		Set("DefaultModel", preferences.DefaultModel).
		Set("FormattingProfile", preferences.FormattingProfile).
		Set("UpdateAt", preferences.UpdateAt).
		Where(sq.Eq{"Id": preferences.Id})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update AIPreferences")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, store.NewErrNotFound("AIPreferences", preferences.Id)
	}

	return preferences, nil
}

func (s *SqlAIPreferencesStore) Delete(userId string) error {
	query := s.getQueryBuilder().
		Delete("AIPreferences").
		Where(sq.Eq{"UserId": userId})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return errors.Wrap(err, "failed to delete AIPreferences")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return store.NewErrNotFound("AIPreferences", userId)
	}

	return nil
}

// GetFormatterPreferences retrieves formatting preferences for a user
func (s *SqlAIPreferencesStore) GetFormatterPreferences(userId string) (string, bool, error) {
	preferences, err := s.GetByUser(userId)
	if err != nil {
		// Return defaults if not found
		return "professional", false, nil
	}

	defaultProfile := preferences.FormattingProfile
	if defaultProfile == "" {
		defaultProfile = "professional"
	}

	// Auto-suggest is not stored in current schema, default to false
	autoSuggest := false

	return defaultProfile, autoSuggest, nil
}

// SetFormatterPreferences updates formatting preferences for a user
func (s *SqlAIPreferencesStore) SetFormatterPreferences(userId string, defaultProfile string, autoSuggest bool) error {
	preferences, err := s.GetByUser(userId)
	if err != nil {
		// Create new preferences if not found
		preferences = &model.AIPreferences{
			UserId:            userId,
			FormattingProfile: defaultProfile,
		}
		preferences.PreSave()
		_, err = s.Save(preferences)
		return err
	}

	// Update existing preferences
	preferences.FormattingProfile = defaultProfile
	preferences.PreUpdate()
	_, err = s.Update(preferences)
	return err
}

