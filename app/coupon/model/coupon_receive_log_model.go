package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	CouponReceiveLogModel interface {
		Insert(ctx context.Context, data *CouponReceiveLog) (sql.Result, error)
		FindByAppkeyAndOrderNo(ctx context.Context, appkey, orderNo string, couponType int8) (*CouponReceiveLog, error)
	}

	defaultCouponReceiveLogModel struct {
		conn  sqlx.SqlConn
		table string
	}

	CouponReceiveLog struct {
		Id          int64     `db:"id"`
		Appkey      string    `db:"appkey"`
		OrderNo     string    `db:"order_no"`
		Mid         int64     `db:"mid"`
		CouponToken string    `db:"coupon_token"`
		CouponType  int8      `db:"coupon_type"`
		Ctime       time.Time `db:"ctime"`
		Mtime       time.Time `db:"mtime"`
	}
)

func NewCouponReceiveLogModel(conn sqlx.SqlConn) CouponReceiveLogModel {
	return &defaultCouponReceiveLogModel{
		conn:  conn,
		table: "coupon_receive_log",
	}
}

func (m *defaultCouponReceiveLogModel) Insert(ctx context.Context, data *CouponReceiveLog) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (appkey, order_no, mid, coupon_token, coupon_type) VALUES (?, ?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.Appkey, data.OrderNo, data.Mid, data.CouponToken, data.CouponType)
}

func (m *defaultCouponReceiveLogModel) FindByAppkeyAndOrderNo(ctx context.Context, appkey, orderNo string, couponType int8) (*CouponReceiveLog, error) {
	var resp CouponReceiveLog
	query := fmt.Sprintf("SELECT id, appkey, order_no, mid, coupon_token, coupon_type, ctime, mtime FROM %s WHERE appkey = ? AND order_no = ? AND coupon_type = ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, appkey, orderNo, couponType)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}
