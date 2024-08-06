package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FileModel = (*customFileModel)(nil)

var cacheFileMD5Prefix = "cache:files:md5:"

type (
	// FileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFileModel.
	FileModel interface {
		fileModel
		FindOneByMD5(ctx context.Context, md5 string) (*File, error)
		ListFiles(ctx context.Context, userId int64, limit, offset int) ([]*File, error)
		CountFiles(ctx context.Context, userId int64) (int64, error)
	}

	customFileModel struct {
		*defaultFileModel
	}
)

// NewFileModel returns a model for the database table.
func NewFileModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FileModel {
	return &customFileModel{
		defaultFileModel: newFileModel(conn, c, opts...),
	}
}

func (m *customFileModel) FindOneByMD5(ctx context.Context, md5 string) (*File, error) {
	fileIdKey := fmt.Sprintf("%s%v", cacheFileMD5Prefix, md5)

	var resp File
	err := m.QueryRowCtx(ctx, &resp, fileIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `checksum`= ? limit 1", fileRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, md5)
	})
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customFileModel) CountFiles(ctx context.Context, userId int64) (int64, error) {
	var cnt int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM file WHERE `owner_id` = %d AND `status`='active'", userId)

	err := m.QueryRowNoCacheCtx(ctx, &cnt, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 如果没有结果，返回0
			return 0, nil
		}
		// 返回其他错误
		return 0, err
	}
	return cnt, nil
}

func (m *customFileModel) ListFiles(ctx context.Context, userId int64, limit, offset int) ([]*File, error) {
	var files []*File
	query := fmt.Sprintf("SELECT %s FROM %s where `owner_id`=%d AND `status`='active' ORDER BY `updated_at` DESC  limit %d offset %d",
		fileRows, m.table, userId, limit, offset)
	err := m.QueryRowsNoCacheCtx(ctx, &files, query)
	switch {
	case err == nil:
		return files, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
