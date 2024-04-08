// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, create_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, create_at, updated_at, user_id, feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreateAt  time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreateAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreateAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteFeedFollows = `-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2
`

type DeleteFeedFollowsParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteFeedFollows(ctx context.Context, arg DeleteFeedFollowsParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollows, arg.ID, arg.UserID)
	return err
}

const getFeedFollows = `-- name: GetFeedFollows :many
SELECT id, create_at, updated_at, user_id, feed_id FROM feed_follows ff
WHERE ff.user_id = $1
`

func (q *Queries) GetFeedFollows(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollows, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
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

const getPopularFeeds = `-- name: GetPopularFeeds :many
SELECT ff.feed_id, count(ff.feed_id) as num_follows, f.url, f.name
FROM feed_follows ff
JOIN feeds f ON ff.feed_id = f.id
GROUP BY ff.feed_id, f.url, f.name
ORDER BY num_follows DESC
OFFSET $1
LIMIT $2
`

type GetPopularFeedsParams struct {
	Offset int32
	Limit  int32
}

type GetPopularFeedsRow struct {
	FeedID     uuid.UUID
	NumFollows int64
	Url        string
	Name       string
}

func (q *Queries) GetPopularFeeds(ctx context.Context, arg GetPopularFeedsParams) ([]GetPopularFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPopularFeeds, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPopularFeedsRow
	for rows.Next() {
		var i GetPopularFeedsRow
		if err := rows.Scan(
			&i.FeedID,
			&i.NumFollows,
			&i.Url,
			&i.Name,
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
