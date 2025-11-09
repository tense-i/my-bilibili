package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CoinExchangeRecordModel = (*customCoinExchangeRecordModel)(nil)

type (
	// CoinExchangeRecordModel 兑换记录模型接口
	CoinExchangeRecordModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *CoinExchangeRecord) (sql.Result, error)
		FindByUid(ctx context.Context, uid int64, offset, limit int64) ([]*CoinExchangeRecord, error)
	}

	customCoinExchangeRecordModel struct {
		*defaultCoinExchangeRecordModel
	}

	defaultCoinExchangeRecordModel struct {
		sqlConn sqlx.SqlConn
		table   string
	}

	// CoinExchangeRecord 兑换记录结构体
	CoinExchangeRecord struct {
		Id            int64     `db:"id"`
		Uid           int64     `db:"uid"`
		TransactionId string    `db:"transaction_id"`
		ExtendTid     string    `db:"extend_tid"`
		SrcCoinType   int32     `db:"src_coin_type"`
		SrcCoinNum    int64     `db:"src_coin_num"`
		DestCoinType  int32     `db:"dest_coin_type"`
		DestCoinNum   int64     `db:"dest_coin_num"`
		ExchangeRate  float64   `db:"exchange_rate"`
		Status        int32     `db:"status"`
		Ctime         time.Time `db:"ctime"`
		Mtime         time.Time `db:"mtime"`
	}
)

// NewCoinExchangeRecordModel 创建兑换记录模型
func NewCoinExchangeRecordModel(conn sqlx.SqlConn) CoinExchangeRecordModel {
	return &customCoinExchangeRecordModel{
		defaultCoinExchangeRecordModel: &defaultCoinExchangeRecordModel{
			sqlConn: conn,
			table:   "coin_exchange_record",
		},
	}
}

// Insert 插入兑换记录
func (m *defaultCoinExchangeRecordModel) Insert(ctx context.Context, session sqlx.Session, data *CoinExchangeRecord) (sql.Result, error) {
	query := `
		INSERT INTO coin_exchange_record (
			uid, transaction_id, extend_tid, src_coin_type, src_coin_num,
			dest_coin_type, dest_coin_num, exchange_rate, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if session != nil {
		return session.ExecCtx(ctx, query,
			data.Uid, data.TransactionId, data.ExtendTid, data.SrcCoinType,
			data.SrcCoinNum, data.DestCoinType, data.DestCoinNum,
			data.ExchangeRate, data.Status)
	}

	return m.sqlConn.ExecCtx(ctx, query,
		data.Uid, data.TransactionId, data.ExtendTid, data.SrcCoinType,
		data.SrcCoinNum, data.DestCoinType, data.DestCoinNum,
		data.ExchangeRate, data.Status)
}

// FindByUid 根据用户ID查询兑换记录
func (m *defaultCoinExchangeRecordModel) FindByUid(ctx context.Context, uid int64, offset, limit int64) ([]*CoinExchangeRecord, error) {
	query := `
		SELECT * FROM coin_exchange_record 
		WHERE uid = ? 
		ORDER BY ctime DESC 
		LIMIT ? OFFSET ?
	`

	var resp []*CoinExchangeRecord
	err := m.sqlConn.QueryRowsCtx(ctx, &resp, query, uid, limit, offset)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
