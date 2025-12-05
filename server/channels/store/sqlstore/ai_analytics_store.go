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

type SqlAIAnalyticsStore struct {
	*SqlStore
}

func newSqlAIAnalyticsStore(sqlStore *SqlStore) store.AIAnalyticsStore {
	return &SqlAIAnalyticsStore{
		SqlStore: sqlStore,
	}
}

func (s *SqlAIAnalyticsStore) Save(analytics *model.AIAnalytics) (*model.AIAnalytics, error) {
	analytics.PreSave()

	if err := analytics.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Insert("AIAnalytics").
		Columns(
			"Id", "ChannelId", "Date", "MessageCount", "UserCount",
			"AvgResponseTime", "CreateAt", "UpdateAt",
		).
		Values(
			analytics.Id, analytics.ChannelId, analytics.Date, analytics.MessageCount, analytics.UserCount,
			analytics.AvgResponseTime, analytics.CreateAt, analytics.UpdateAt,
		)

	if _, err := s.GetMaster().ExecBuilder(query); err != nil {
		return nil, errors.Wrap(err, "failed to save AIAnalytics")
	}

	return analytics, nil
}

func (s *SqlAIAnalyticsStore) Get(id string) (*model.AIAnalytics, error) {
	query := s.getQueryBuilder().
		Select("Id", "ChannelId", "Date", "MessageCount", "UserCount", "AvgResponseTime", "CreateAt", "UpdateAt").
		From("AIAnalytics").
		Where(sq.Eq{"Id": id})

	var analytics model.AIAnalytics
	err := s.GetReplica().GetBuilder(&analytics, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AIAnalytics", id)
		}
		return nil, errors.Wrapf(err, "failed to find AIAnalytics with id=%s", id)
	}

	// Initialize maps if nil
	if analytics.TopContributors == nil {
		analytics.TopContributors = make(map[string]interface{})
	}
	if analytics.HourlyDistribution == nil {
		analytics.HourlyDistribution = make(map[string]interface{})
	}

	return &analytics, nil
}

func (s *SqlAIAnalyticsStore) GetByChannel(channelId string, startDate, endDate string) ([]*model.AIAnalytics, error) {
	query := s.getQueryBuilder().
		Select("Id", "ChannelId", "Date", "MessageCount", "UserCount", "AvgResponseTime", "CreateAt", "UpdateAt").
		From("AIAnalytics").
		Where(sq.And{
			sq.Eq{"ChannelId": channelId},
			sq.GtOrEq{"Date": startDate},
			sq.LtOrEq{"Date": endDate},
		}).
		OrderBy("Date DESC")

	var analyticsList []*model.AIAnalytics
	err := s.GetReplica().SelectBuilder(&analyticsList, query)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find AIAnalytics for channelId=%s", channelId)
	}

	// Initialize maps for each analytics entry
	for _, analytics := range analyticsList {
		if analytics.TopContributors == nil {
			analytics.TopContributors = make(map[string]interface{})
		}
		if analytics.HourlyDistribution == nil {
			analytics.HourlyDistribution = make(map[string]interface{})
		}
	}

	return analyticsList, nil
}

func (s *SqlAIAnalyticsStore) GetByChannelAndDate(channelId, date string) (*model.AIAnalytics, error) {
	query := s.getQueryBuilder().
		Select("Id", "ChannelId", "Date", "MessageCount", "UserCount", "AvgResponseTime", "CreateAt", "UpdateAt").
		From("AIAnalytics").
		Where(sq.And{
			sq.Eq{"ChannelId": channelId},
			sq.Eq{"Date": date},
		})

	var analytics model.AIAnalytics
	err := s.GetReplica().GetBuilder(&analytics, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AIAnalytics", channelId+":"+date)
		}
		return nil, errors.Wrapf(err, "failed to find AIAnalytics for channelId=%s, date=%s", channelId, date)
	}

	// Initialize maps if nil
	if analytics.TopContributors == nil {
		analytics.TopContributors = make(map[string]interface{})
	}
	if analytics.HourlyDistribution == nil {
		analytics.HourlyDistribution = make(map[string]interface{})
	}

	return &analytics, nil
}

func (s *SqlAIAnalyticsStore) Update(analytics *model.AIAnalytics) (*model.AIAnalytics, error) {
	analytics.PreUpdate()

	if err := analytics.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Update("AIAnalytics").
		Set("MessageCount", analytics.MessageCount).
		Set("UserCount", analytics.UserCount).
		Set("AvgResponseTime", analytics.AvgResponseTime).
		Set("UpdateAt", analytics.UpdateAt).
		Where(sq.Eq{"Id": analytics.Id})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update AIAnalytics")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, store.NewErrNotFound("AIAnalytics", analytics.Id)
	}

	return analytics, nil
}

func (s *SqlAIAnalyticsStore) DeleteOlderThan(date string) (int64, error) {
	query := s.getQueryBuilder().
		Delete("AIAnalytics").
		Where(sq.Lt{"Date": date})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete old AIAnalytics")
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

