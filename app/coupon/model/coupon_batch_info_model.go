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
	CouponBatchInfoModel interface {
		Insert(ctx context.Context, data *CouponBatchInfo) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*CouponBatchInfo, error)
		FindByBatchToken(ctx context.Context, batchToken string) (*CouponBatchInfo, error)
		FindAll(ctx context.Context) ([]*CouponBatchInfo, error)
		Update(ctx context.Context, data *CouponBatchInfo) error
		UpdateCurrentCount(ctx context.Context, batchToken string, count int) (int64, error)
		UpdateCurrentCountWithLimit(ctx context.Context, batchToken string, count int) (int64, error)
		Delete(ctx context.Context, id int64) error
	}

	defaultCouponBatchInfoModel struct {
		conn  sqlx.SqlConn
		table string
	}

	CouponBatchInfo struct {
		Id                  int64     `db:"id"`
		AppId               int64     `db:"app_id"`
		Name                string    `db:"name"`
		BatchToken          string    `db:"batch_token"`
		MaxCount            int64     `db:"max_count"`
		CurrentCount        int64     `db:"current_count"`
		StartTime           int64     `db:"start_time"`
		ExpireTime          int64     `db:"expire_time"`
		ExpireDay           int64     `db:"expire_day"`
		LimitCount          int64     `db:"limit_count"`
		FullAmount          float64   `db:"full_amount"`
		Amount              float64   `db:"amount"`
		State               int8      `db:"state"`
		CouponType          int8      `db:"coupon_type"`
		PlatformLimit       string    `db:"platform_limit"`
		ProductLimitMonth   int8      `db:"product_limit_month"`
		ProductLimitRenewal int8      `db:"product_limit_renewal"`
		Ver                 int64     `db:"ver"`
		Ctime               time.Time `db:"ctime"`
		Mtime               time.Time `db:"mtime"`
	}
)

func NewCouponBatchInfoModel(conn sqlx.SqlConn, c cache.CacheConf) CouponBatchInfoModel {
	return &defaultCouponBatchInfoModel{
		conn:  conn,
		table: "coupon_batch_info",
	}
}

func (m *defaultCouponBatchInfoModel) Insert(ctx context.Context, data *CouponBatchInfo) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (app_id, name, batch_token, max_count, current_count, start_time, expire_time, expire_day, limit_count, full_amount, amount, state, coupon_type, platform_limit, product_limit_month, product_limit_renewal, ver) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.AppId, data.Name, data.BatchToken, data.MaxCount, data.CurrentCount, data.StartTime, data.ExpireTime, data.ExpireDay, data.LimitCount, data.FullAmount, data.Amount, data.State, data.CouponType, data.PlatformLimit, data.ProductLimitMonth, data.ProductLimitRenewal, data.Ver)
}

func (m *defaultCouponBatchInfoModel) FindOne(ctx context.Context, id int64) (*CouponBatchInfo, error) {
	var resp CouponBatchInfo
	query := fmt.Sprintf("SELECT id, app_id, name, batch_token, max_count, current_count, start_time, expire_time, expire_day, limit_count, full_amount, amount, state, coupon_type, platform_limit, product_limit_month, product_limit_renewal, ver, ctime, mtime FROM %s WHERE id = ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultCouponBatchInfoModel) FindByBatchToken(ctx context.Context, batchToken string) (*CouponBatchInfo, error) {
	var resp CouponBatchInfo
	query := fmt.Sprintf("SELECT id, app_id, name, batch_token, max_count, current_count, start_time, expire_time, expire_day, limit_count, full_amount, amount, state, coupon_type, platform_limit, product_limit_month, product_limit_renewal, ver, ctime, mtime FROM %s WHERE batch_token = ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, batchToken)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}

func (m *defaultCouponBatchInfoModel) FindAll(ctx context.Context) ([]*CouponBatchInfo, error) {
	var resp []*CouponBatchInfo
	query := fmt.Sprintf("SELECT id, app_id, name, batch_token, max_count, current_count, start_time, expire_time, expire_day, limit_count, full_amount, amount, state, coupon_type, platform_limit, product_limit_month, product_limit_renewal, ver, ctime, mtime FROM %s", m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultCouponBatchInfoModel) Update(ctx context.Context, data *CouponBatchInfo) error {
	query := fmt.Sprintf("UPDATE %s SET app_id = ?, name = ?, max_count = ?, start_time = ?, expire_time = ?, expire_day = ?, limit_count = ?, full_amount = ?, amount = ?, state = ?, coupon_type = ?, platform_limit = ?, product_limit_month = ?, product_limit_renewal = ? WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.AppId, data.Name, data.MaxCount, data.StartTime, data.ExpireTime, data.ExpireDay, data.LimitCount, data.FullAmount, data.Amount, data.State, data.CouponType, data.PlatformLimit, data.ProductLimitMonth, data.ProductLimitRenewal, data.Id)
	return err
}

func (m *defaultCouponBatchInfoModel) UpdateCurrentCount(ctx context.Context, batchToken string, count int) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET current_count = current_count + ? WHERE batch_token = ?", m.table)
	result, err := m.conn.ExecCtx(ctx, query, count, batchToken)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultCouponBatchInfoModel) UpdateCurrentCountWithLimit(ctx context.Context, batchToken string, count int) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET current_count = current_count + ? WHERE batch_token = ? AND current_count + ? <= max_count", m.table)
	result, err := m.conn.ExecCtx(ctx, query, count, batchToken, count)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultCouponBatchInfoModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
