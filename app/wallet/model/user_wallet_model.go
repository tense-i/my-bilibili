package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserWalletModel = (*customUserWalletModel)(nil)

type (
	// UserWalletModel 用户钱包模型接口
	UserWalletModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *UserWallet) (sql.Result, error)
		FindOne(ctx context.Context, uid int64) (*UserWallet, error)
		FindOneForUpdate(ctx context.Context, session sqlx.Session, uid int64) (*UserWallet, error)
		UpdateRecharge(ctx context.Context, session sqlx.Session, uid int64, coinType int32, amount int64) error
		UpdatePay(ctx context.Context, session sqlx.Session, uid int64, coinType int32, amount int64) error
		UpdateExchange(ctx context.Context, session sqlx.Session, uid int64, srcType int32, srcAmount int64, destType int32, destAmount int64) error
		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
		Delete(ctx context.Context, uid int64) error
	}

	customUserWalletModel struct {
		*defaultUserWalletModel
	}

	defaultUserWalletModel struct {
		sqlConn sqlx.SqlConn
		table   string
	}

	// UserWallet 用户钱包结构体
	UserWallet struct {
		Uid              int64          `db:"uid"`
		Gold             int64          `db:"gold"`
		IapGold          int64          `db:"iap_gold"`
		Silver           int64          `db:"silver"`
		GoldRechargeCnt  int64          `db:"gold_recharge_cnt"`
		GoldPayCnt       int64          `db:"gold_pay_cnt"`
		SilverPayCnt     int64          `db:"silver_pay_cnt"`
		CostBase         int64          `db:"cost_base"`
		SnapshotTime     sql.NullString `db:"snapshot_time"`
		SnapshotGold     int64          `db:"snapshot_gold"`
		SnapshotIapGold  int64          `db:"snapshot_iap_gold"`
		SnapshotSilver   int64          `db:"snapshot_silver"`
		Reserved1        int64          `db:"reserved1"`
		Reserved2        string         `db:"reserved2"`
	}
)

// NewUserWalletModel 创建用户钱包模型
func NewUserWalletModel(conn sqlx.SqlConn, c cache.CacheConf) UserWalletModel {
	return &customUserWalletModel{
		defaultUserWalletModel: &defaultUserWalletModel{
			sqlConn: conn,
			table:   "user_wallet",
		},
	}
}

// getTableName 获取分表名称
func (m *defaultUserWalletModel) getTableName(uid int64) string {
	return fmt.Sprintf("%s_%d", m.table, GetWalletTableIndex(uid))
}

// Insert 插入钱包记录
func (m *defaultUserWalletModel) Insert(ctx context.Context, session sqlx.Session, data *UserWallet) (sql.Result, error) {
	tableName := m.getTableName(data.Uid)
	query := fmt.Sprintf(`
		INSERT INTO %s (uid, gold, iap_gold, silver, gold_recharge_cnt, gold_pay_cnt, silver_pay_cnt, cost_base)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	if session != nil {
		return session.ExecCtx(ctx, query, data.Uid, data.Gold, data.IapGold, data.Silver,
			data.GoldRechargeCnt, data.GoldPayCnt, data.SilverPayCnt, data.CostBase)
	}
	return m.sqlConn.ExecCtx(ctx, query, data.Uid, data.Gold, data.IapGold, data.Silver,
		data.GoldRechargeCnt, data.GoldPayCnt, data.SilverPayCnt, data.CostBase)
}

// FindOne 查询单条记录
func (m *defaultUserWalletModel) FindOne(ctx context.Context, uid int64) (*UserWallet, error) {
	tableName := m.getTableName(uid)
	query := fmt.Sprintf(`SELECT * FROM %s WHERE uid = ? LIMIT 1`, tableName)

	var resp UserWallet
	err := m.sqlConn.QueryRowCtx(ctx, &resp, query, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindOneForUpdate 查询单条记录（带FOR UPDATE锁）
func (m *defaultUserWalletModel) FindOneForUpdate(ctx context.Context, session sqlx.Session, uid int64) (*UserWallet, error) {
	tableName := m.getTableName(uid)
	query := fmt.Sprintf(`SELECT * FROM %s WHERE uid = ? FOR UPDATE`, tableName)

	var resp UserWallet
	err := session.QueryRowCtx(ctx, &resp, query, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound, sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// UpdateRecharge 更新充值
func (m *defaultUserWalletModel) UpdateRecharge(ctx context.Context, session sqlx.Session, uid int64, coinType int32, amount int64) error {
	tableName := m.getTableName(uid)
	coinField := getCoinField(coinType)

	query := fmt.Sprintf(`
		UPDATE %s SET 
			%s = %s + ?,
			gold_recharge_cnt = gold_recharge_cnt + ?,
			snapshot_time = NOW()
		WHERE uid = ?
	`, tableName, coinField, coinField)

	_, err := session.ExecCtx(ctx, query, amount, amount, uid)
	return err
}

// UpdatePay 更新消费
func (m *defaultUserWalletModel) UpdatePay(ctx context.Context, session sqlx.Session, uid int64, coinType int32, amount int64) error {
	tableName := m.getTableName(uid)
	coinField := getCoinField(coinType)
	cntField := getPayCntField(coinType)

	query := fmt.Sprintf(`
		UPDATE %s SET 
			%s = %s - ?,
			%s = %s + ?
		WHERE uid = ?
	`, tableName, coinField, coinField, cntField, cntField)

	_, err := session.ExecCtx(ctx, query, amount, amount, uid)
	return err
}

// UpdateExchange 更新兑换（原子性：一条SQL完成）
func (m *defaultUserWalletModel) UpdateExchange(ctx context.Context, session sqlx.Session, uid int64,
	srcType int32, srcAmount int64, destType int32, destAmount int64) error {

	tableName := m.getTableName(uid)
	srcField := getCoinField(srcType)
	destField := getCoinField(destType)
	srcCntField := getPayCntField(srcType)

	query := fmt.Sprintf(`
		UPDATE %s SET 
			%s = %s - ?,
			%s = %s + ?,
			%s = %s + ?
		WHERE uid = ?
	`, tableName, srcField, srcField, destField, destField, srcCntField, srcCntField)

	_, err := session.ExecCtx(ctx, query, srcAmount, destAmount, srcAmount, uid)
	return err
}

// Trans 执行事务
func (m *defaultUserWalletModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.sqlConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// Delete 删除记录
func (m *defaultUserWalletModel) Delete(ctx context.Context, uid int64) error {
	tableName := m.getTableName(uid)
	query := fmt.Sprintf(`DELETE FROM %s WHERE uid = ?`, tableName)
	_, err := m.sqlConn.ExecCtx(ctx, query, uid)
	return err
}

// getCoinField 获取币种字段名
func getCoinField(coinType int32) string {
	switch coinType {
	case CoinTypeGold:
		return "gold"
	case CoinTypeIapGold:
		return "iap_gold"
	case CoinTypeSilver:
		return "silver"
	default:
		return "gold"
	}
}

// getPayCntField 获取消费统计字段名
func getPayCntField(coinType int32) string {
	switch coinType {
	case CoinTypeGold, CoinTypeIapGold:
		return "gold_pay_cnt"
	case CoinTypeSilver:
		return "silver_pay_cnt"
	default:
		return "gold_pay_cnt"
	}
}

// GetCoinByType 获取指定币种余额
func GetCoinByType(wallet *UserWallet, coinType int32) int64 {
	switch coinType {
	case CoinTypeGold:
		return wallet.Gold
	case CoinTypeIapGold:
		return wallet.IapGold
	case CoinTypeSilver:
		return wallet.Silver
	default:
		return 0
	}
}

// AddCoin 增加币种余额（内存操作）
func AddCoin(wallet *UserWallet, coinType int32, amount int64) {
	switch coinType {
	case CoinTypeGold:
		wallet.Gold += amount
		wallet.GoldRechargeCnt += amount
	case CoinTypeIapGold:
		wallet.IapGold += amount
		wallet.GoldRechargeCnt += amount
	case CoinTypeSilver:
		wallet.Silver += amount
	}
}

// SubCoin 减少币种余额（内存操作）
func SubCoin(wallet *UserWallet, coinType int32, amount int64) {
	switch coinType {
	case CoinTypeGold:
		wallet.Gold -= amount
		wallet.GoldPayCnt += amount
	case CoinTypeIapGold:
		wallet.IapGold -= amount
		wallet.GoldPayCnt += amount
	case CoinTypeSilver:
		wallet.Silver -= amount
		wallet.SilverPayCnt += amount
	}
}
