package model

// 优惠券使用状态
const (
	UseFaild   int8 = 0
	UseSuccess int8 = 1
)

// 优惠券状态
const (
	NotUsed int32 = 0 // 未使用
	InUse   int32 = 1 // 使用中
	Used    int32 = 2 // 已使用
	Expire  int32 = 3 // 已过期
	Block   int32 = 4 // 已冻结
)

// 支付状态
const (
	WaitPay    int8 = 0
	InPay      int8 = 1
	PaySuccess int8 = 2
	PayFaild   int8 = 3
)

// 最大发放数量
const MaxSalaryCount = 100

// 余额变更类型
const (
	VipSalary         int8 = 1 // VIP发放
	SystemAdminSalary int8 = 2 // 系统发放
	Consume           int8 = 3 // 消费
	ConsumeFaildBack  int8 = 4 // 消费失败退回
)

// 优惠券类型
const (
	CouponVideo     int8 = 1 // 观影券
	CouponCartoon   int8 = 2 // 漫画券
	CouponAllowance int8 = 3 // 代金券
)

// 代金券来源
const (
	AllowanceNone            int8 = 0
	AllowanceSystemAdmin     int8 = 1 // 系统发放
	AllowanceBusinessReceive int8 = 2 // 业务领取
	AllowanceBusinessNewYear int8 = 3 // 新年活动
	AllowanceCodeOpen        int8 = 4 // 兑换码
)

// 批次状态
const (
	BatchStateNormal int8 = 0 // 正常
	BatchStateBlock  int8 = 1 // 冻结
)

// 代金券不可用原因
const (
	CouponHadBlock             = "代金券已被冻结"
	CouponFullAmountDissatisfy = "未达到满额条件"
	CouponNotInUsableTime      = "当前不在有效期内"
	CouponInUse                = "已绑定在其他未支付订单,点击解锁"
	CouponPlatformExplain      = "当前平台不可使用"
	CouponProductExplain       = "当前商品不可使用"
)

// 代金券提示
const (
	CouponTipNotUse      = "不使用代金券"
	CouponTipChooseOther = "选中其他商品有惊喜"
	CouponTipUse         = "抵扣%.2f元"
	CouponTipInUse       = "有代金券被锁定"
)

// 代金券可用状态
const (
	AllowanceDisables int8 = 0
	AllowanceUsable   int8 = 1
)

// 代金券变更类型
const (
	AllowanceSalary         int8 = 1 // 发放
	AllowanceConsume        int8 = 2 // 消费
	AllowanceCancel         int8 = 3 // 取消
	AllowanceConsumeSuccess int8 = 4 // 消费成功
	AllowanceConsumeFaild   int8 = 5 // 消费失败
	AllowanceReceive        int8 = 6 // 领取
)

// 支付通知状态
const (
	AllowanceUseFaild   int8 = 0
	AllowanceUseSuccess int8 = 1
)

// 兑换码状态
const (
	CodeStateNotUse int32 = 1 // 未使用
	CodeStateUsed   int32 = 2 // 已使用
	CodeStateBlock  int32 = 3 // 已冻结
)

// 设备类型
const (
	DeviceIOS         int = 1
	DeviceIPAD        int = 2
	DevicePC          int = 3
	DeviceANDROID     int = 4
	DeviceIPADHD      int = 5
	DeviceIOSBLUE     int = 6
	DeviceANDROIDBLUE int = 7
	DevicePUBLIC      int = 8
)

// 平台名称映射
var PlatformByName = map[string]int{
	"ios":       DeviceIOS,
	"ios_b":     DeviceIOS,
	"ipad":      DeviceIPAD,
	"ipadhd":    DeviceIPAD,
	"pc":        DevicePC,
	"public":    DevicePC,
	"android":   DeviceANDROID,
	"android_b": DeviceANDROID,
}

// 商品月份限制
const (
	ProdLimMonthNone int8 = 0
	ProdLimMonth1    int8 = 1
	ProdLimMonth3    int8 = 3
	ProdLimMonth12   int8 = 12
)

// 续费限制
const (
	ProdLimRenewalAll     int8 = 0 // 不限
	ProdLimRenewalAuto    int8 = 1 // 自动续期
	ProdLimRenewalNotAuto int8 = 2 // 非自动续期
)

// 分表函数
func HitCouponInfo(mid int64) int64 {
	return mid % 100
}

func HitCouponChangeLog(mid int64) int64 {
	return mid % 100
}

func HitAllowanceInfo(mid int64) int64 {
	return mid % 10
}

func HitAllowanceChangeLog(mid int64) int64 {
	return mid % 10
}

func HitBalanceInfo(mid int64) int64 {
	return mid % 10
}

func HitBalanceChangeLog(mid int64) int64 {
	return mid % 10
}
