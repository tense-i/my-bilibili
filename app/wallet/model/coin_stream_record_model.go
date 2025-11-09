package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CoinStreamRecordModel = (*customCoinStreamRecordModel)(nil)

type (
	// CoinStreamRecordModel 流水记录模型接口
	CoinStreamRecordModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *CoinStreamRecord) (sql.Result, error)
		FindByUid(ctx context.Context, uid int64, offset, limit int64) ([]*CoinStreamRecord, error)
		Count(ctx context.Context, uid int64) (int64, error)
	}

	customCoinStreamRecordModel struct {
		*defaultCoinStreamRecordModel
	}

	defaultCoinStreamRecordModel struct {
		sqlConn sqlx.SqlConn
		table   string
	}

	// CoinStreamRecord 流水记录结构体
	CoinStreamRecord struct {
		Id            int64     `db:"id"`
		Uid           int64     `db:"uid"`
		TransactionId string    `db:"transaction_id"`
		ExtendTid     string    `db:"extend_tid"`
		CoinType      int32     `db:"coin_type"`
		DeltaCoinNum  int64     `db:"delta_coin_num"`
		OrgCoinNum    int64     `db:"org_coin_num"`
		OpResult      int32     `db:"op_result"`
		OpReason      int32     `db:"op_reason"`
		OpType        int32     `db:"op_type"`
		OpTime        time.Time `db:"op_time"`
		BizCode       string    `db:"biz_code"`
		Area          int64     `db:"area"`
		Source        string    `db:"source"`
		Metadata      string    `db:"metadata"`
		BizSource     string    `db:"biz_source"`
		Platform      int32     `db:"platform"`
		Reserved1     int64     `db:"reserved1"`
		Version       int64     `db:"version"`
	}
)

// NewCoinStreamRecordModel 创建流水记录模型
func NewCoinStreamRecordModel(conn sqlx.SqlConn) CoinStreamRecordModel {
	return &customCoinStreamRecordModel{
		defaultCoinStreamRecordModel: &defaultCoinStreamRecordModel{
			sqlConn: conn,
			table:   "coin_stream_record",
		},
	}
}

// getTableName 获取分表名称
func (m *defaultCoinStreamRecordModel) getTableName(transactionId string) string {
	return fmt.Sprintf("%s_%d", m.table, GetStreamTableIndex(transactionId))
}

// Insert 插入流水记录
func (m *defaultCoinStreamRecordModel) Insert(ctx context.Context, session sqlx.Session, data *CoinStreamRecord) (sql.Result, error) {
	tableName := m.getTableName(data.TransactionId)
	query := fmt.Sprintf(`
		INSERT INTO %s (
			uid, transaction_id, extend_tid, coin_type, delta_coin_num, org_coin_num,
			op_result, op_reason, op_type, op_time, biz_code, area, source, metadata,
			biz_source, platform, reserved1, version
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	if session != nil {
		return session.ExecCtx(ctx, query,
			data.Uid, data.TransactionId, data.ExtendTid, data.CoinType,
			data.DeltaCoinNum, data.OrgCoinNum, data.OpResult, data.OpReason,
			data.OpType, data.OpTime, data.BizCode, data.Area, data.Source,
			data.Metadata, data.BizSource, data.Platform, data.Reserved1, data.Version)
	}

	return m.sqlConn.ExecCtx(ctx, query,
		data.Uid, data.TransactionId, data.ExtendTid, data.CoinType,
		data.DeltaCoinNum, data.OrgCoinNum, data.OpResult, data.OpReason,
		data.OpType, data.OpTime, data.BizCode, data.Area, data.Source,
		data.Metadata, data.BizSource, data.Platform, data.Reserved1, data.Version)
}

// FindByUid 根据用户ID查询流水（需要查询所有分表）
func (m *defaultCoinStreamRecordModel) FindByUid(ctx context.Context, uid int64, offset, limit int64) ([]*CoinStreamRecord, error) {
	var result []*CoinStreamRecord

	// 查询所有10张分表
	for i := int64(0); i < 10; i++ {
		tableName := fmt.Sprintf("%s_%d", m.table, i)
		query := fmt.Sprintf(`
			SELECT * FROM %s 
			WHERE uid = ? 
			ORDER BY op_time DESC 
			LIMIT ? OFFSET ?
		`, tableName)

		var records []*CoinStreamRecord
		err := m.sqlConn.QueryRowsCtx(ctx, &records, query, uid, limit, offset)
		if err != nil && err != sqlx.ErrNotFound {
			return nil, err
		}
		result = append(result, records...)
	}

	return result, nil
}

// Count 统计用户流水数量（需要查询所有分表）
func (m *defaultCoinStreamRecordModel) Count(ctx context.Context, uid int64) (int64, error) {
	var total int64

	// 查询所有10张分表
	for i := int64(0); i < 10; i++ {
		tableName := fmt.Sprintf("%s_%d", m.table, i)
		query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE uid = ?`, tableName)

		var count int64
		err := m.sqlConn.QueryRowCtx(ctx, &count, query, uid)
		if err != nil && err != sqlx.ErrNotFound {
			return 0, err
		}
		total += count
	}

	return total, nil
}
