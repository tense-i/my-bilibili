package model

// 召回项
type RecallItem struct {
	AVID       int64             `json:"avid"`
	Score      float64           `json:"score"`
	RecallType string            `json:"recall_type"`
	RecallTag  string            `json:"recall_tag"`
	Priority   int32             `json:"priority"`
	Extra      map[string]string `json:"extra"`
}

// 视频索引（正排索引）
type VideoIndex struct {
	AVID     int64  `json:"avid"`
	MID      int64  `json:"mid"` // UP主MID
	Title    string `json:"title"`
	ZoneID   int32  `json:"zone_id"`
	Duration int32  `json:"duration"`
	PubTime  int64  `json:"pub_time"`
	State    int32  `json:"state"`
	Tags     []Tag  `json:"tags"`
}

// 标签
type Tag struct {
	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
	TagType int32  `json:"tag_type"`
}

// 召回策略名称常量
const (
	HotRecall         = "HotRecall"
	SelectionRecall   = "SelectionRecall"
	LikeI2IRecall     = "LikeI2IRecall"
	LikeTagRecall     = "LikeTagRecall"
	LikeUPRecall      = "LikeUPRecall"
	PosI2IRecall      = "PosI2IRecall" // 正反馈I2I召回
	PosTagRecall      = "PosTagRecall" // 正反馈标签召回
	FollowRecall      = "FollowRecall"
	UserProfileRecall = "UserProfileRecall"
	UserProfileBili   = "UserProfileBili" // B站用户画像召回
	UserProfileBBQ    = "UserProfileBBQ"  // BBQ用户画像召回
	RandomRecall      = "RandomRecall"    // 随机召回
)

// 优先级常量
const (
	PriorityVeryHigh = 10000
	PriorityHigh     = 1000
	PriorityMid      = 100
	PriorityLow      = 10
)

// Redis Key 前缀
const (
	RecallKeyI2IPrefix   = "RECALL:I2I"
	RecallKeyTagIDPrefix = "RECALL:HOT_T"
	RecallKeyUpIDPrefix  = "RECALL:HOT_UP"
)

// Redis Key 常量
const (
	RedisKeyHotIndex       = "RECALL:HOT:INDEX"
	RedisKeySelectionIndex = "RECALL:SELECTION:INDEX"
	RedisKeyI2IIndex       = "RECALL:I2I:%d"        // avid
	RedisKeyTagIndex       = "RECALL:TAG:%s"        // tag_name
	RedisKeyTagIDIndex     = "RECALL:HOT_T:%d"      // tag_id
	RedisKeyTagIDIndexStr  = "RECALL:HOT_T:%s"      // tag_name
	RedisKeyTagNewPubIndex = "RECALL:T:%s"          // tag_name (新发布视频)
	RedisKeyUPIndex        = "RECALL:UP:%d"         // up_mid
	RedisKeyUPIndexHot     = "RECALL:HOT_UP:%d"     // up_mid (热门视频)
	RedisKeyBloomFilter    = "RECALL:BLOOM:%d"      // mid
	RedisKeyHotDefault     = "RECALL:HOT_DEFAULT:0" // 默认热门视频
	RedisKeyOpVideo        = "job:bbq:rec:op"       // 运营推荐视频
	RedisKeyUserActionLike = "user:action:%d:like"  // mid (点赞)
	RedisKeyUserActionPos  = "user:action:%d:pos"   // mid (正反馈)
	RedisKeyUserActionNeg  = "user:action:%d:neg"   // mid (负反馈)
)

// 视频信息
type VideoInfo struct {
	AVID     int64  `json:"avid"`
	MID      int64  `json:"mid"`
	Title    string `json:"title"`
	ZoneID   int32  `json:"zone_id"`
	Duration int32  `json:"duration"`
	PubTime  int64  `json:"pub_time"`
	State    int32  `json:"state"`

	// 全站统计
	PlayHive  int64 `json:"play_hive"`
	LikesHive int64 `json:"likes_hive"`
	FavHive   int64 `json:"fav_hive"`
	ReplyHive int64 `json:"reply_hive"`
	ShareHive int64 `json:"share_hive"`
	CoinHive  int64 `json:"coin_hive"`

	// 月度统计
	PlayMonth       int64 `json:"play_month"`
	LikesMonth      int64 `json:"likes_month"`
	ReplyMonth      int64 `json:"reply_month"`
	ShareMonth      int64 `json:"share_month"`
	PlayMonthFinish int64 `json:"play_month_finish"`

	// 标签
	Tags []VideoTag `json:"tags"`
}

// 视频标签
type VideoTag struct {
	AVID    int64  `json:"avid"`
	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
}
