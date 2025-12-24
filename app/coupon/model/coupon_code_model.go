package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	CouponCodeModel interface {
		Insert(ctx context.Context, data *CouponCode) (sql.Result, error)
		FindByCode(ctx context.Context, code string) (*CouponCode, error)
		UpdateState(ctx context.Context, code string, state int32, mid int64, couponToken string, ver int64) (int64, error)
		CountByMidAndBatchToken(ctx context.Context, mid int64, batchToken string) (int64, error)
	}

	defaultCouponCodeModel struct {
		conn  sqlx.SqlConn
		table string
	}

	CouponCode struct {
		Id          int64     `db:"id"`
		BatchToken  string    `db:"batch_token"`
		Code        string    `db:"code"`
		State       int32     `db:"state"`
		Mid         int64     `db:"mid"`
		CouponType  int8      `db:"coupon_type"`
		CouponToken string    `db:"coupon_token"`
		Ver         int64     `db:"ver"`
		Ctime       time.Time `db:"ctime"`
		Mtime       time.Time `db:"mtime"`
	}
)

func NewCouponCodeModel(conn sqlx.SqlConn, c cache.CacheConf) CouponCodeModel {
	return &defaultCouponCodeModel{
		conn:  conn,
		table: "coupon_code",
	}
}

func (m *defaultCouponCodeModel) Insert(ctx context.Context, data *CouponCode) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (batch_token, code, state, mid, coupon_type, coupon_token, ver) VALUES (?, ?, ?, ?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.BatchToken, data.Code, data.State, data.Mid, data.CouponType, data.CouponToken, data.Ver)
}

func (m *defaultCouponCodeModel) FindByCode(ctx context.Context, code string) (*CouponCode, error) {
	var resp CouponCode
	query := fmt.Sprintf("SELECT id, batch_token, code, state, mid, coupon_type, coupon_token, ver, ctime, mtime FROM %s WHERE code = ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, code)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultCouponCodeModel) UpdateState(ctx context.Context, code string, state int32, mid int64, couponToken string, ver int64) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET state = ?, mid = ?, coupon_token = ? WHERE code = ? AND ver = ?", m.table)
	result, err := m.conn.ExecCtx(ctx, query, state, mid, couponToken, code, ver)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultCouponCodeModel) CountByMidAndBatchToken(ctx context.Context, mid int64, batchToken string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE mid = ? AND batch_token = ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &count, query, mid, batchToken)
	if err != nil {
		return 0, err
	}
	return count, nil
}
