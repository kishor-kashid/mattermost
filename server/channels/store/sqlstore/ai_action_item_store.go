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
		Insert("aiactionitems").
		Columns(
			"id", "channelid", "postid", "createdby", "assigneeid",
			"description", "duedate", "priority", "status", "completedat",
			"createdat", "updatedat", "deletedat",
		).
		Values(
			actionItem.Id, actionItem.ChannelId, actionItem.PostId, actionItem.CreatedBy, actionItem.AssigneeId,
			actionItem.Description, actionItem.DueDate, actionItem.Priority, actionItem.Status, actionItem.CompletedAt,
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
		From("aiactionitems").
		Where(sq.Eq{"id": id, "deletedat": 0})

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

func (s *SqlAIActionItemStore) GetByChannel(channelId string, includeCompleted bool, offset, limit int) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("aiactionitems").
		Where(sq.Eq{"channelid": channelId, "deletedat": 0})
	
	if !includeCompleted {
		query = query.Where(sq.NotEq{"status": "completed"})
	}
	
	query = query.OrderBy("createdat DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find AIActionItems for channelId=%s", channelId)
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) GetByUser(userId string, includeCompleted bool, offset, limit int) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("aiactionitems").
		Where(sq.Or{
			sq.Eq{"assigneeid": userId},
			sq.Eq{"createdby": userId},
		}).
		Where(sq.Eq{"deletedat": 0})
	
	if !includeCompleted {
		query = query.Where(sq.NotEq{"status": "completed"})
	}
	
	query = query.OrderBy("createdat DESC").
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
		From("aiactionitems").
		Where(sq.Eq{"assigneeid": assigneeId, "deletedat": 0}).
		OrderBy("createdat DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find AIActionItems for assigneeId=%s", assigneeId)
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) GetOverdue(currentTime int64) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("aiactionitems").
		Where(sq.And{
			sq.NotEq{"status": "completed"},
			sq.NotEq{"status": "dismissed"},
			sq.Eq{"deletedat": 0},
			sq.NotEq{"duedate": 0},
			sq.Lt{"duedate": currentTime},
		}).
		OrderBy("duedate ASC")

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find overdue AIActionItems")
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) GetDueSoon(startTime, endTime int64) ([]*model.AIActionItem, error) {
	query := s.getQueryBuilder().
		Select("*").
		From("aiactionitems").
		Where(sq.And{
			sq.NotEq{"status": "completed"},
			sq.NotEq{"status": "dismissed"},
			sq.Eq{"deletedat": 0},
			sq.NotEq{"duedate": 0},
			sq.GtOrEq{"duedate": startTime},
			sq.LtOrEq{"duedate": endTime},
		}).
		OrderBy("duedate ASC")

	var actionItems []*model.AIActionItem
	err := s.GetReplica().SelectBuilder(&actionItems, query)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find due soon AIActionItems")
	}

	return actionItems, nil
}

func (s *SqlAIActionItemStore) Update(actionItem *model.AIActionItem) (*model.AIActionItem, error) {
	actionItem.PreUpdate()

	if err := actionItem.IsValid(); err != nil {
		return nil, err
	}

	query := s.getQueryBuilder().
		Update("aiactionitems").
		Set("assigneeid", actionItem.AssigneeId).
		Set("description", actionItem.Description).
		Set("duedate", actionItem.DueDate).
		Set("priority", actionItem.Priority).
		Set("status", actionItem.Status).
		Set("completedat", actionItem.CompletedAt).
		Set("updatedat", actionItem.UpdateAt).
		Where(sq.Eq{"id": actionItem.Id})

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
		Update("aiactionitems").
		Set("deletedat", deleteAt).
		Where(sq.Eq{"id": id})

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
		Delete("aiactionitems").
		Where(sq.Eq{"id": id})

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
