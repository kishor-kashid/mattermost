// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"database/sql"

	"github.com/lib/pq"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/store"
)

type SqlAILinkSummaryStore struct {
	*SqlStore
}

func newSqlAILinkSummaryStore(sqlStore *SqlStore) store.AILinkSummaryStore {
	return &SqlAILinkSummaryStore{
		SqlStore: sqlStore,
	}
}

func (s *SqlAILinkSummaryStore) Save(summary *model.AILinkSummary) (*model.AILinkSummary, error) {
	summary.PreSave()

	if err := summary.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Insert("AILinkSummaries").
		Columns(
			"Id", "URL", "URLHash", "Title", "Description", "Summary",
			"KeyPoints", "ContentType", "ReadingTime", "Domain", "FaviconURL",
			"CreateAt", "ExpiresAt",
		).
		Values(
			summary.Id, summary.Url, summary.UrlHash, summary.Title, summary.Description, summary.Summary,
			pq.Array(summary.KeyPoints), summary.ContentType, summary.ReadingTime, summary.Domain, summary.FaviconURL,
			summary.CreateAt, summary.ExpiresAt,
		)

	if _, err := s.GetMaster().ExecBuilder(query); err != nil {
		return nil, errors.Wrap(err, "failed to save AILinkSummary")
	}

	return summary, nil
}

func (s *SqlAILinkSummaryStore) Get(id string) (*model.AILinkSummary, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AILinkSummaries").
		Where(sq.Eq{"Id": id})

	var summary model.AILinkSummary
	err := s.GetReplica().GetBuilder(&summary, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AILinkSummary", id)
		}
		return nil, errors.Wrapf(err, "failed to find AILinkSummary with id=%s", id)
	}

	return &summary, nil
}

func (s *SqlAILinkSummaryStore) GetByURLHash(urlHash string) (*model.AILinkSummary, error) {
	currentTime := model.GetMillis()

	query := s.getQueryBuilder().
		Select("*").
		From("AILinkSummaries").
		Where(sq.And{
			sq.Eq{"URLHash": urlHash},
			sq.Gt{"ExpiresAt": currentTime},
		}).
		OrderBy("CreateAt DESC").
		Limit(1)

	var summary model.AILinkSummary
	err := s.GetReplica().GetBuilder(&summary, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AILinkSummary", urlHash)
		}
		return nil, errors.Wrapf(err, "failed to find AILinkSummary for urlHash=%s", urlHash)
	}

	return &summary, nil
}

func (s *SqlAILinkSummaryStore) DeleteExpired(currentTime int64) (int64, error) {
	query := s.getQueryBuilder().
		Delete("AILinkSummaries").
		Where(sq.LtOrEq{"ExpiresAt": currentTime})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete expired AILinkSummaries")
	}

	rows, _ := result.RowsAffected()
	return rows, nil
}

func (s *SqlAILinkSummaryStore) Delete(id string) error {
	query := s.getQueryBuilder().
		Delete("AILinkSummaries").
		Where(sq.Eq{"Id": id})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return errors.Wrap(err, "failed to delete AILinkSummary")
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return store.NewErrNotFound("AILinkSummary", id)
	}

	return nil
}
