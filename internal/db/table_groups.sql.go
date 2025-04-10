// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: table_groups.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addMember = `-- name: AddMember :exec
INSERT INTO group_members(user_id, group_id)
VALUES($1, $2)
`

type AddMemberParams struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (q *Queries) AddMember(ctx context.Context, arg AddMemberParams) error {
	_, err := q.db.ExecContext(ctx, addMember, arg.UserID, arg.GroupID)
	return err
}

const checkMemeber = `-- name: CheckMemeber :one
SELECT id, group_id, user_id, added_at FROM group_members
WHERE group_id=$1 AND user_id=$2
`

type CheckMemeberParams struct {
	GroupID uuid.UUID
	UserID  uuid.UUID
}

func (q *Queries) CheckMemeber(ctx context.Context, arg CheckMemeberParams) (GroupMember, error) {
	row := q.db.QueryRowContext(ctx, checkMemeber, arg.GroupID, arg.UserID)
	var i GroupMember
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.UserID,
		&i.AddedAt,
	)
	return i, err
}

const createGroup = `-- name: CreateGroup :one
INSERT INTO groups(name, description, created_by)
VALUES($1, $2, $3)
RETURNING id, name, description, created_by, created_at, updated_at, image_url
`

type CreateGroupParams struct {
	Name        string
	Description sql.NullString
	CreatedBy   uuid.UUID
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, createGroup, arg.Name, arg.Description, arg.CreatedBy)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ImageUrl,
	)
	return i, err
}

const deleteGroup = `-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id=$1
`

func (q *Queries) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteGroup, id)
	return err
}

const deleteGroupMember = `-- name: DeleteGroupMember :exec
DELETE FROM group_members
WHERE user_id=$1 AND group_id=$2
`

type DeleteGroupMemberParams struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (q *Queries) DeleteGroupMember(ctx context.Context, arg DeleteGroupMemberParams) error {
	_, err := q.db.ExecContext(ctx, deleteGroupMember, arg.UserID, arg.GroupID)
	return err
}

const getGroupByID = `-- name: GetGroupByID :one
SELECT id, name, description, created_by, created_at, updated_at, image_url FROM groups
WHERE id=$1
`

func (q *Queries) GetGroupByID(ctx context.Context, id uuid.UUID) (Group, error) {
	row := q.db.QueryRowContext(ctx, getGroupByID, id)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ImageUrl,
	)
	return i, err
}

const getGroupMembers = `-- name: GetGroupMembers :many
SELECT users.id, users.name, users.email, users.image_url, group_members.added_at,
CASE 
    WHEN g.created_by = users.id THEN TRUE
    ELSE FALSE
END AS is_admin
FROM users
INNER JOIN group_members ON users.id = group_members.user_id
INNER JOIN groups g ON group_members.group_id = g.id
WHERE group_id=$1
`

type GetGroupMembersRow struct {
	ID       uuid.UUID
	Name     string
	Email    string
	ImageUrl sql.NullString
	AddedAt  sql.NullTime
	IsAdmin  bool
}

func (q *Queries) GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]GetGroupMembersRow, error) {
	rows, err := q.db.QueryContext(ctx, getGroupMembers, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGroupMembersRow
	for rows.Next() {
		var i GetGroupMembersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.ImageUrl,
			&i.AddedAt,
			&i.IsAdmin,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserAllGroups = `-- name: GetUserAllGroups :many
SELECT groups.id, name, description, created_at FROM groups
INNER JOIN group_members ON groups.id = group_members.group_id
WHERE user_id=$1
`

type GetUserAllGroupsRow struct {
	ID          uuid.UUID
	Name        string
	Description sql.NullString
	CreatedAt   sql.NullTime
}

func (q *Queries) GetUserAllGroups(ctx context.Context, userID uuid.UUID) ([]GetUserAllGroupsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserAllGroups, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserAllGroupsRow
	for rows.Next() {
		var i GetUserAllGroupsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserGroups = `-- name: GetUserGroups :many
SELECT id, name, description, created_by, created_at, updated_at, image_url FROM groups
WHERE created_by=$1
`

func (q *Queries) GetUserGroups(ctx context.Context, createdBy uuid.UUID) ([]Group, error) {
	rows, err := q.db.QueryContext(ctx, getUserGroups, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Group
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGroup = `-- name: UpdateGroup :exec
UPDATE groups
SET name=$1, description=$2, updated_at=CURRENT_TIMESTAMP
WHERE id=$3
`

type UpdateGroupParams struct {
	Name        string
	Description sql.NullString
	ID          uuid.UUID
}

func (q *Queries) UpdateGroup(ctx context.Context, arg UpdateGroupParams) error {
	_, err := q.db.ExecContext(ctx, updateGroup, arg.Name, arg.Description, arg.ID)
	return err
}
