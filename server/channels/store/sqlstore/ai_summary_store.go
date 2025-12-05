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

type SqlAISummaryStore struct {
	*SqlStore
}

func newSqlAISummaryStore(sqlStore *SqlStore) store.AISummaryStore {
	return &SqlAISummaryStore{
		SqlStore: sqlStore,
	}
}

func (s *SqlAISummaryStore) Save(summary *model.AISummary) (*model.AISummary, error) {
	summary.PreSave()

	if err := summary.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Insert("AISummaries").
		Columns(
			"Id", "ChannelId", "PostId", "SummaryType", "Summary",
			"MessageCount", "StartTime", "EndTime", "CreateAt", "ExpiresAt",
		).
		Values(
			summary.Id, summary.ChannelId, summary.PostId, summary.SummaryType, summary.Summary,
			summary.MessageCount, summary.StartTime, summary.EndTime, summary.CreateAt, summary.ExpiresAt,
		)

	if _, err := s.GetMaster().ExecBuilder(query); err != nil {
		return nil, errors.Wrap(err, "failed to save AISummary")
	}

	return summary, nil
}

func (s *SqlAISummaryStore) Get(id string) (*model.AISummary, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AISummaries").
		Where(sq.Eq{"Id": id})

	var summary model.AISummary
	err := s.GetReplica().GetBuilder(&summary, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AISummary", id)
		}
		return nil, errors.Wrapf(err, "failed to find AISummary with id=%s", id)
	}

	return &summary, nil
}

func (s *SqlAISummaryStore) GetByChannel(channelId string, offset, limit int) ([]*model.AISummary, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AISummaries").
		Where(sq.Eq{"ChannelId": channelId}).
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	var summaries []*model.AISummary
	err := s.GetReplica().SelectBuilder(&summaries, query)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find AISummaries for channelId=%s", channelId)
	}

	return summaries, nil
}

func (s *SqlAISummaryStore) GetCachedSummary(channelId, summaryType string, startTime, endTime int64) (*model.AISummary, error) {
	currentTime := model.GetMillis()

	query := s.getQueryBuilder().
		Select("*").
		From("AISummaries").
		Where(sq.And{
			sq.Eq{"ChannelId": channelId},
			sq.Eq{"SummaryType": summaryType},
			sq.GtOrEq{"StartTime": startTime},
			sq.LtOrEq{"EndTime": endTime},
			sq.Gt{"ExpiresAt": currentTime},
		}).
		OrderBy("CreateAt DESC").
		Limit(1)

	var summary model.AISummary
	err := s.GetReplica().GetBuilder(&summary, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AISummary", channelId)
		}
		return nil, errors.Wrapf(err, "failed to find cached AISummary for channelId=%s", channelId)
	}

	return &summary, nil
}

func (s *SqlAISummaryStore) DeleteExpired(currentTime int64) (int64, error) {
	query := s.getQueryBuilder().
		Delete("AISummaries").
		Where(sq.LtOrEq{"ExpiresAt": currentTime})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete expired AISummaries")
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (s *SqlAISummaryStore) Delete(id string) error {
	query := s.getQueryBuilder().
		Delete("AISummaries").
		Where(sq.Eq{"Id": id})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return errors.Wrap(err, "failed to delete AISummary")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return store.NewErrNotFound("AISummary", id)
	}

	return nil
}

