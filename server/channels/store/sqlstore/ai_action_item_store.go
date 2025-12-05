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

type SqlAIActionItemStore struct {
	*SqlStore
}

func newSqlAIActionItemStore(sqlStore *SqlStore) store.AIActionItemStore {
	return &SqlAIActionItemStore{
		SqlStore: sqlStore,
	}
}

func (s *SqlAIActionItemStore) Save(actionItem *model.AIActionItem) (*model.AIActionItem, error) {
	actionItem.PreSave()

	if err := actionItem.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Insert("AIActionItems").
		Columns(
			"Id", "ChannelId", "PostId", "UserId", "AssigneeId",
			"Description", "Deadline", "Status", "ReminderSent",
			"CreateAt", "UpdateAt", "DeleteAt",
		).
		Values(
			actionItem.Id, actionItem.ChannelId, actionItem.PostId, actionItem.UserId, actionItem.AssigneeId,
			actionItem.Description, actionItem.Deadline, actionItem.Status, actionItem.ReminderSent,
			actionItem.CreateAt, actionItem.UpdateAt, actionItem.DeleteAt,
		)

	if _, err := s.GetMaster().ExecBuilder(query); err != nil {
		return nil, errors.Wrap(err, "failed to save AIActionItem")
	}

	return actionItem, nil
}

func (s *SqlAIActionItemStore) Get(id string) (*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AIActionItems").
		Where(sq.Eq{"Id": id, "DeleteAt": 0})

	var actionItem model.AIActionItem
	err := s.GetReplica().GetBuilder(&actionItem, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("AIActionItem", id)
		}
		return nil, errors.Wrapf(err, "failed to find AIActionItem with id=%s", id)
	}

	return &actionItem, nil
}

func (s *SqlAIActionItemStore) GetByChannel(channelId string, offset, limit int) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AIActionItems").
		Where(sq.Eq{"ChannelId": channelId, "DeleteAt": 0}).
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find AIActionItems for channelId=%s", channelId)
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) GetByUser(userId string, offset, limit int) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AIActionItems").
		Where(sq.Eq{"UserId": userId, "DeleteAt": 0}).
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find AIActionItems for userId=%s", userId)
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) GetByAssignee(assigneeId string, offset, limit int) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AIActionItems").
		Where(sq.Eq{"AssigneeId": assigneeId, "DeleteAt": 0}).
		OrderBy("CreateAt DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find AIActionItems for assigneeId=%s", assigneeId)
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) GetPendingReminders(currentTime int64) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("AIActionItems").
		Where(sq.And{
			sq.Eq{"Status": model.AIActionItemStatusPending},
			sq.Eq{"ReminderSent": false},
			sq.Eq{"DeleteAt": 0},
			sq.NotEq{"Deadline": nil},
			sq.LtOrEq{"Deadline": currentTime},
		})

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find pending reminder AIActionItems")
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) Update(actionItem *model.AIActionItem) (*model.AIActionItem, error) {
	actionItem.PreUpdate()

	if err := actionItem.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Update("AIActionItems").
		Set("AssigneeId", actionItem.AssigneeId).
		Set("Description", actionItem.Description).
		Set("Deadline", actionItem.Deadline).
		Set("Status", actionItem.Status).
		Set("ReminderSent", actionItem.ReminderSent).
		Set("UpdateAt", actionItem.UpdateAt).
		Where(sq.Eq{"Id": actionItem.Id})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update AIActionItem")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, store.NewErrNotFound("AIActionItem", actionItem.Id)
	}

	return actionItem, nil
}

func (s *SqlAIActionItemStore) Delete(id string, deleteAt int64) error {
	query := s.getQueryBuilder().
		Update("AIActionItems").
		Set("DeleteAt", deleteAt).
		Where(sq.Eq{"Id": id})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return errors.Wrap(err, "failed to delete AIActionItem")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return store.NewErrNotFound("AIActionItem", id)
	}

	return nil
}

func (s *SqlAIActionItemStore) PermanentDelete(id string) error {
	query := s.getQueryBuilder().
		Delete("AIActionItems").
		Where(sq.Eq{"Id": id})

	result, err := s.GetMaster().ExecBuilder(query)
	if err != nil {
		return errors.Wrap(err, "failed to permanently delete AIActionItem")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return store.NewErrNotFound("AIActionItem", id)
	}

	return nil
}

