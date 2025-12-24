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
	CouponInfoModel interface {
		Insert(ctx context.Context, data *CouponInfo) (sql.Result, error)
		FindByToken(ctx context.Context, mid int64, couponToken string) (*CouponInfo, error)
		FindByMidAndType(ctx context.Context, mid int64, couponType int8) ([]*CouponInfo, error)
		FindByMidAndState(ctx context.Context, mid int64, state int32, pn, ps int32) ([]*CouponInfo, error)
		CountByMidAndState(ctx context.Context, mid int64, state int32) (int64, error)
		CountByMidAndType(ctx context.Context, mid int64, couponType int8) (int32, error)
		UpdateState(ctx context.Context, id, mid int64, state int32, orderNo string, oid int64, remark string, ver int64) (int64, error)
	}

	defaultCouponInfoModel struct {
		conn        sqlx.SqlConn
		tablePrefix string
	}

	CouponInfo struct {
		Id          int64     `db:"id"`
		CouponToken string    `db:"coupon_token"`
		Mid         int64     `db:"mid"`
		State       int32     `db:"state"`
		StartTime   int64     `db:"start_time"`
		ExpireTime  int64     `db:"expire_time"`
		Origin      int8      `db:"origin"`
		CouponType  int8      `db:"coupon_type"`
		OrderNO     string    `db:"order_no"`
		Oid         int64     `db:"oid"`
		BatchToken  string    `db:"batch_token"`
		Ver         int64     `db:"ver"`
		Remark      string    `db:"remark"`
		Ctime       time.Time `db:"ctime"`
		Mtime       time.Time `db:"mtime"`
	}
)

func NewCouponInfoModel(conn sqlx.SqlConn, c cache.CacheConf) CouponInfoModel {
	return &defaultCouponInfoModel{
		conn:        conn,
		tablePrefix: "coupon_info_",
	}
}

func (m *defaultCouponInfoModel) tableName(mid int64) string {
	return fmt.Sprintf("%s0%d", m.tablePrefix, HitCouponInfo(mid))
}

func (m *defaultCouponInfoModel) Insert(ctx context.Context, data *CouponInfo) (sql.Result, error) {
	table := m.tableName(data.Mid)
	query := fmt.Sprintf("INSERT INTO %s (coupon_token, mid, state, start_time, expire_time, origin, coupon_type, batch_token) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", table)
	return m.conn.ExecCtx(ctx, query, data.CouponToken, data.Mid, data.State, data.StartTime, data.ExpireTime, data.Origin, data.CouponType, data.BatchToken)
}

func (m *defaultCouponInfoModel) FindByToken(ctx context.Context, mid int64, couponToken string) (*CouponInfo, error) {
	table := m.tableName(mid)
	var resp CouponInfo
	query := fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, coupon_type, order_no, oid, batch_token, ver, remark, ctime, mtime FROM %s WHERE coupon_token = ?", table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, couponToken)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultCouponInfoModel) FindByMidAndType(ctx context.Context, mid int64, couponType int8) ([]*CouponInfo, error) {
	table := m.tableName(mid)
	var resp []*CouponInfo
	now := time.Now().Unix()
	query := fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, coupon_type, order_no, oid, batch_token, ver, remark, ctime, mtime FROM %s WHERE mid = ? AND coupon_type = ? AND state = 0 AND expire_time > ? ORDER BY id DESC", table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, mid, couponType, now)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultCouponInfoModel) FindByMidAndState(ctx context.Context, mid int64, state int32, pn, ps int32) ([]*CouponInfo, error) {
	table := m.tableName(mid)
	var resp []*CouponInfo
	offset := (pn - 1) * ps
	query := fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, coupon_type, order_no, oid, batch_token, ver, remark, ctime, mtime FROM %s WHERE mid = ? AND state = ? ORDER BY id DESC LIMIT ?, ?", table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, mid, state, offset, ps)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultCouponInfoModel) CountByMidAndState(ctx context.Context, mid int64, state int32) (int64, error) {
	table := m.tableName(mid)
	var count int64
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE mid = ? AND state = ?", table)
	err := m.conn.QueryRowCtx(ctx, &count, query, mid, state)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *defaultCouponInfoModel) CountByMidAndType(ctx context.Context, mid int64, couponType int8) (int32, error) {
	table := m.tableName(mid)
	var count int32
	now := time.Now().Unix()
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE mid = ? AND coupon_type = ? AND state = 0 AND expire_time > ?", table)
	err := m.conn.QueryRowCtx(ctx, &count, query, mid, couponType, now)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *defaultCouponInfoModel) UpdateState(ctx context.Context, id, mid int64, state int32, orderNo string, oid int64, remark string, ver int64) (int64, error) {
	table := m.tableName(mid)
	query := fmt.Sprintf("UPDATE %s SET state = ?, order_no = ?, oid = ?, remark = ?, ver = ver + 1 WHERE id = ? AND ver = ?", table)
	result, err := m.conn.ExecCtx(ctx, query, state, orderNo, oid, remark, id, ver)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
