package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	CouponAllowanceChangeLogModel interface {
		Insert(ctx context.Context, data *CouponAllowanceChangeLog) (sql.Result, error)
		FindByToken(ctx context.Context, mid int64, couponToken string) ([]*CouponAllowanceChangeLog, error)
	}

	defaultCouponAllowanceChangeLogModel struct {
		conn        sqlx.SqlConn
		tablePrefix string
	}

	CouponAllowanceChangeLog struct {
		Id          int64     `db:"id"`
		CouponToken string    `db:"coupon_token"`
		OrderNo     string    `db:"order_no"`
		Mid         int64     `db:"mid"`
		State       int8      `db:"state"`
		ChangeType  int8      `db:"change_type"`
		Ctime       time.Time `db:"ctime"`
		Mtime       time.Time `db:"mtime"`
	}
)

func NewCouponAllowanceChangeLogModel(conn sqlx.SqlConn) CouponAllowanceChangeLogModel {
	return &defaultCouponAllowanceChangeLogModel{
		conn:        conn,
		tablePrefix: "coupon_allowance_change_log_",
	}
}

func (m *defaultCouponAllowanceChangeLogModel) tableName(mid int64) string {
	return fmt.Sprintf("%s0%d", m.tablePrefix, HitAllowanceChangeLog(mid))
}

func (m *defaultCouponAllowanceChangeLogModel) Insert(ctx context.Context, data *CouponAllowanceChangeLog) (sql.Result, error) {
	table := m.tableName(data.Mid)
	query := fmt.Sprintf("INSERT INTO %s (coupon_token, order_no, mid, state, change_type, ctime) VALUES (?, ?, ?, ?, ?, ?)", table)
	return m.conn.ExecCtx(ctx, query, data.CouponToken, data.OrderNo, data.Mid, data.State, data.ChangeType, data.Ctime)
}

func (m *defaultCouponAllowanceChangeLogModel) FindByToken(ctx context.Context, mid int64, couponToken string) ([]*CouponAllowanceChangeLog, error) {
	table := m.tableName(mid)
	var resp []*CouponAllowanceChangeLog
	query := fmt.Sprintf("SELECT id, coupon_token, order_no, mid, state, change_type, ctime, mtime FROM %s WHERE coupon_token = ?", table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, couponToken)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
