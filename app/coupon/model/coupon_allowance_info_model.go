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
	CouponAllowanceInfoModel interface {
		Insert(ctx context.Context, data *CouponAllowanceInfo) (sql.Result, error)
		BatchInsert(ctx context.Context, mid int64, list []*CouponAllowanceInfo) (int64, error)
		FindByToken(ctx context.Context, mid int64, couponToken string) (*CouponAllowanceInfo, error)
		FindByOrderNo(ctx context.Context, mid int64, orderNo string) (*CouponAllowanceInfo, error)
		FindByMidAndState(ctx context.Context, mid int64, state int32, now int64) ([]*CouponAllowanceInfo, error)
		FindUsable(ctx context.Context, mid int64, state int32, now int64) ([]*CouponAllowanceInfo, error)
		FindList(ctx context.Context, mid int64, state int32, now int64, stime time.Time) ([]*CouponAllowanceInfo, error)
		UpdateState(ctx context.Context, id, mid int64, state int32, orderNo, remark string, ver int64) (int64, error)
		UpdateStateToUsed(ctx context.Context, id, mid int64, orderNo string, ver int64) (int64, error)
		UpdateStateToNotUsed(ctx context.Context, id, mid int64, ver int64) (int64, error)
		CountByBatchToken(ctx context.Context, mid int64, batchToken string) (int64, error)
	}

	defaultCouponAllowanceInfoModel struct {
		conn        sqlx.SqlConn
		tablePrefix string
	}

	CouponAllowanceInfo struct {
		Id          int64     `db:"id"`
		CouponToken string    `db:"coupon_token"`
		Mid         int64     `db:"mid"`
		State       int32     `db:"state"`
		StartTime   int64     `db:"start_time"`
		ExpireTime  int64     `db:"expire_time"`
		Origin      int8      `db:"origin"`
		Ver         int64     `db:"ver"`
		BatchToken  string    `db:"batch_token"`
		OrderNO     string    `db:"order_no"`
		Amount      float64   `db:"amount"`
		FullAmount  float64   `db:"full_amount"`
		AppId       int64     `db:"app_id"`
		Remark      string    `db:"remark"`
		Ctime       time.Time `db:"ctime"`
		Mtime       time.Time `db:"mtime"`
	}
)

func NewCouponAllowanceInfoModel(conn sqlx.SqlConn, c cache.CacheConf) CouponAllowanceInfoModel {
	return &defaultCouponAllowanceInfoModel{
		conn:        conn,
		tablePrefix: "coupon_allowance_info_",
	}
}

func (m *defaultCouponAllowanceInfoModel) tableName(mid int64) string {
	return fmt.Sprintf("%s0%d", m.tablePrefix, HitAllowanceInfo(mid))
}

func (m *defaultCouponAllowanceInfoModel) Insert(ctx context.Context, data *CouponAllowanceInfo) (sql.Result, error) {
	table := m.tableName(data.Mid)
	query := fmt.Sprintf("INSERT INTO %s (coupon_token, mid, state, start_time, expire_time, origin, batch_token, amount, full_amount, app_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", table)
	return m.conn.ExecCtx(ctx, query, data.CouponToken, data.Mid, data.State, data.StartTime, data.ExpireTime, data.Origin, data.BatchToken, data.Amount, data.FullAmount, data.AppId)
}

func (m *defaultCouponAllowanceInfoModel) BatchInsert(ctx context.Context, mid int64, list []*CouponAllowanceInfo) (int64, error) {
	if len(list) == 0 {
		return 0, nil
	}
	table := m.tableName(mid)
	query := fmt.Sprintf("INSERT INTO %s (coupon_token, mid, state, start_time, expire_time, origin, batch_token, amount, full_amount, app_id, ctime) VALUES ", table)

	var values []interface{}
	for i, v := range list {
		if i > 0 {
			query += ","
		}
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		values = append(values, v.CouponToken, v.Mid, v.State, v.StartTime, v.ExpireTime, v.Origin, v.BatchToken, v.Amount, v.FullAmount, v.AppId, v.Ctime)
	}

	result, err := m.conn.ExecCtx(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultCouponAllowanceInfoModel) FindByToken(ctx context.Context, mid int64, couponToken string) (*CouponAllowanceInfo, error) {
	table := m.tableName(mid)
	var resp CouponAllowanceInfo
	query := fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, ver, batch_token, order_no, amount, full_amount, app_id, remark, ctime, mtime FROM %s WHERE coupon_token = ?", table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, couponToken)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultCouponAllowanceInfoModel) FindByOrderNo(ctx context.Context, mid int64, orderNo string) (*CouponAllowanceInfo, error) {
	table := m.tableName(mid)
	var resp CouponAllowanceInfo
	query := fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, ver, batch_token, order_no, amount, full_amount, app_id, remark, ctime, mtime FROM %s WHERE order_no = ?", table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderNo)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultCouponAllowanceInfoModel) FindByMidAndState(ctx context.Context, mid int64, state int32, now int64) ([]*CouponAllowanceInfo, error) {
	table := m.tableName(mid)
	var resp []*CouponAllowanceInfo
	query := fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, ver, batch_token, order_no, amount, full_amount, app_id, remark, ctime, mtime FROM %s WHERE mid = ? AND expire_time > ? AND state = ?", table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, mid, now, state)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultCouponAllowanceInfoModel) FindUsable(ctx context.Context, mid int64, state int32, now int64) ([]*CouponAllowanceInfo, error) {
	table := m.tableName(mid)
	var resp []*CouponAllowanceInfo
	query := fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, ver, batch_token, order_no, amount, full_amount, app_id, remark, ctime, mtime FROM %s WHERE mid = ? AND expire_time > ? AND start_time < ? AND state = ?", table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, mid, now, now, state)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultCouponAllowanceInfoModel) FindList(ctx context.Context, mid int64, state int32, now int64, stime time.Time) ([]*CouponAllowanceInfo, error) {
	table := m.tableName(mid)
	var resp []*CouponAllowanceInfo
	var query string

	switch state {
	case NotUsed:
		query = fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, ver, batch_token, order_no, amount, full_amount, app_id, remark, ctime, mtime FROM %s WHERE mid = ? AND (state = 0 OR state = 1) AND expire_time > ? AND start_time < ? AND ctime > ? ORDER BY id DESC", table)
		err := m.conn.QueryRowsCtx(ctx, &resp, query, mid, now, now, stime)
		return resp, err
	case Used:
		query = fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, ver, batch_token, order_no, amount, full_amount, app_id, remark, ctime, mtime FROM %s WHERE mid = ? AND state = 2 AND ctime > ? ORDER BY id DESC", table)
		err := m.conn.QueryRowsCtx(ctx, &resp, query, mid, stime)
		return resp, err
	case Expire:
		query = fmt.Sprintf("SELECT id, coupon_token, mid, state, start_time, expire_time, origin, ver, batch_token, order_no, amount, full_amount, app_id, remark, ctime, mtime FROM %s WHERE mid = ? AND state <> 2 AND expire_time < ? AND ctime > ? ORDER BY id DESC", table)
		err := m.conn.QueryRowsCtx(ctx, &resp, query, mid, now, stime)
		return resp, err
	}
	return nil, nil
}

func (m *defaultCouponAllowanceInfoModel) UpdateState(ctx context.Context, id, mid int64, state int32, orderNo, remark string, ver int64) (int64, error) {
	table := m.tableName(mid)
	query := fmt.Sprintf("UPDATE %s SET state = ?, order_no = ?, remark = ?, ver = ver + 1 WHERE id = ? AND ver = ?", table)
	result, err := m.conn.ExecCtx(ctx, query, state, orderNo, remark, id, ver)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultCouponAllowanceInfoModel) UpdateStateToUsed(ctx context.Context, id, mid int64, orderNo string, ver int64) (int64, error) {
	table := m.tableName(mid)
	query := fmt.Sprintf("UPDATE %s SET state = ?, order_no = ?, ver = ver + 1 WHERE id = ? AND ver = ? AND state = ?", table)
	result, err := m.conn.ExecCtx(ctx, query, Used, orderNo, id, ver, InUse)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultCouponAllowanceInfoModel) UpdateStateToNotUsed(ctx context.Context, id, mid int64, ver int64) (int64, error) {
	table := m.tableName(mid)
	query := fmt.Sprintf("UPDATE %s SET state = ?, order_no = '', ver = ver + 1 WHERE id = ? AND ver = ? AND state = ?", table)
	result, err := m.conn.ExecCtx(ctx, query, NotUsed, id, ver, InUse)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultCouponAllowanceInfoModel) CountByBatchToken(ctx context.Context, mid int64, batchToken string) (int64, error) {
	table := m.tableName(mid)
	var count int64
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE mid = ? AND batch_token = ?", table)
	err := m.conn.QueryRowCtx(ctx, &count, query, mid, batchToken)
	if err != nil {
		return 0, err
	}
	return count, nil
}
