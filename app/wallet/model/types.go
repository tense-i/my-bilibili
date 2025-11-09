package model

import (
	"errors"
	"hash/fnv"
)

// 错误定义
var (
	ErrNotFound        = errors.New("record not found")
	ErrInvalidParam    = errors.New("invalid parameter")
	ErrCoinNotEnough   = errors.New("coin not enough")
	ErrLockFailed      = errors.New("lock failed")
	ErrTransactionFail = errors.New("transaction failed")
)

// 币种类型常量
const (
	CoinTypeGold    = 1 // 金瓜子（Android/PC/H5）
	CoinTypeIapGold = 2 // IAP金瓜子（iOS）
	CoinTypeSilver  = 3 // 银瓜子
)

// 操作类型常量
const (
	OpTypeRecharge = 1 // 充值
	OpTypePay      = 2 // 消费
	OpTypeExchange = 3 // 兑换
)

// 操作结果常量
const (
	OpResultAddSucc   = 1  // 增加成功
	OpResultSubSucc   = 2  // 减少成功
	OpResultAddFailed = -1 // 增加失败
	OpResultSubFailed = -2 // 减少失败
)

// 失败原因常量
const (
	OpReasonSuccess      = 0 // 成功
	OpReasonNotEnough    = 1 // 余额不足
	OpReasonInvalidParam = 2 // 参数错误
	OpReasonLockFailed   = 3 // 锁失败
)

// 平台类型常量
const (
	PlatformIOS     = 1 // iOS
	PlatformAndroid = 2 // Android
	PlatformPC      = 3 // PC
	PlatformH5      = 4 // H5
)

// GetSysCoinType 根据币种名称和平台获取系统币种类型
func GetSysCoinType(coinType string, platform string) string {
	// iOS平台的gold自动转为iap_gold
	if coinType == "gold" && platform == "ios" {
		return "iap_gold"
	}
	return coinType
}

// GetCoinTypeNumber 获取币种类型编号
func GetCoinTypeNumber(coinType string) int32 {
	switch coinType {
	case "gold":
		return CoinTypeGold
	case "iap_gold":
		return CoinTypeIapGold
	case "silver":
		return CoinTypeSilver
	default:
		return 0
	}
}

// GetCoinTypeName 获取币种类型名称
func GetCoinTypeName(coinType int32) string {
	switch coinType {
	case CoinTypeGold:
		return "gold"
	case CoinTypeIapGold:
		return "iap_gold"
	case CoinTypeSilver:
		return "silver"
	default:
		return ""
	}
}

// GetPlatformNumber 获取平台编号
func GetPlatformNumber(platform string) int32 {
	switch platform {
	case "ios":
		return PlatformIOS
	case "android":
		return PlatformAndroid
	case "pc":
		return PlatformPC
	case "h5":
		return PlatformH5
	default:
		return PlatformAndroid
	}
}

// GetWalletTableIndex 获取用户钱包表索引（按uid取模）
func GetWalletTableIndex(uid int64) int64 {
	return uid % 10
}

// GetStreamTableIndex 获取流水记录表索引（按transaction_id hash）
func GetStreamTableIndex(transactionId string) int64 {
	h := fnv.New32a()
	h.Write([]byte(transactionId))
	return int64(h.Sum32() % 10)
}
