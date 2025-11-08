package model

// 用户画像
type UserProfile struct {
	MID   int64  `json:"mid"`
	BUVID string `json:"buvid"`

	// 兴趣标签（标签名 -> 权重）
	Tags  map[string]float64 `json:"tags"`
	Zones map[int32]float64  `json:"zones"` // 分区ID -> 权重

	// 偏好UP主（UP主MID -> 权重）
	UPs map[int64]float64 `json:"ups"`

	// 行为统计
	TotalActions int64 `json:"total_actions"`
	PlayCount    int64 `json:"play_count"`
	LikeCount    int64 `json:"like_count"`
	CoinCount    int64 `json:"coin_count"`
	FavCount     int64 `json:"fav_count"`
	ShareCount   int64 `json:"share_count"`

	// 更新时间
	UpdatedAt int64 `json:"updated_at"`

	// 偏好UP主（UP主MID -> 时间戳）
	PrefUps map[int64]int64 `json:"pref_ups"`

	// 关注的UP主
	FollowUps map[int64]int64 `json:"follow_ups"`

	// 黑名单UP主
	BlackUps map[int64]bool `json:"black_ups"`

	// 实时行为数据（视频AVID -> 时间戳）
	LikeVideos map[int64]int64 `json:"like_videos"` // 点赞视频
	PosVideos  map[int64]int64 `json:"pos_videos"`  // 正反馈视频
	NegVideos  map[int64]int64 `json:"neg_videos"`  // 负反馈视频

	// 实时行为标签（标签ID -> 次数）
	LikeTagIDs map[int64]int64 `json:"like_tag_ids"` // 点赞视频的标签
	PosTagIDs  map[int64]int64 `json:"pos_tag_ids"`  // 正反馈视频的标签
	NegTagIDs  map[int64]int64 `json:"neg_tag_ids"`  // 负反馈视频的标签

	// 点赞UP主（UP主MID -> 时间戳）
	LikeUPs map[int64]int64 `json:"like_ups"`

	// 观看历史（用于去重）
	LastRecords []int64 `json:"last_records"`
}

// 推荐记录
// 注意：字段顺序必须与SQL查询顺序一致，以便 sqlx 正确扫描
type RecommendRecord struct {
	// 数据库字段（按SQL查询顺序）
	AVID            int64  `json:"avid" db:"avid"`
	CID             int64  `json:"cid" db:"cid"`
	UPMID           int64  `json:"up_mid" db:"mid"`
	Title           string `json:"title" db:"title"`
	Cover           string `json:"cover" db:"cover"`
	Duration        int32  `json:"duration" db:"duration"`
	ZoneID          int32  `json:"zone_id" db:"zone_id"`
	PubTime         int64  `json:"pub_time" db:"pub_time"`
	State           int8   `json:"state" db:"state"`
	PlayHive        int64  `json:"play_hive" db:"play_hive"`
	LikesHive       int64  `json:"likes_hive" db:"likes_hive"`
	FavHive         int64  `json:"fav_hive" db:"fav_hive"`
	ReplyHive       int64  `json:"reply_hive" db:"reply_hive"`
	ShareHive       int64  `json:"share_hive" db:"share_hive"`
	CoinHive        int64  `json:"coin_hive" db:"coin_hive"`
	PlayMonth       int64  `json:"play_month" db:"play_month"`
	LikesMonth      int64  `json:"likes_month" db:"likes_month"`
	ReplyMonth      int64  `json:"reply_month" db:"reply_month"`
	ShareMonth      int64  `json:"share_month" db:"share_month"`
	PlayMonthFinish int64  `json:"play_month_finish" db:"play_month_finish"`

	// 非数据库字段
	ZoneName string   `json:"zone_name" db:"-"`
	UPName   string   `json:"up_name" db:"-"`
	Tags     []string `json:"tags" db:"-"`
	TagIDs   []int64  `json:"tag_ids" db:"-"`

	// 统计数据（非数据库字段）
	Play  int64 `json:"play" db:"-"`
	Like  int64 `json:"like" db:"-"`
	Coin  int64 `json:"coin" db:"-"`
	Fav   int64 `json:"fav" db:"-"`
	Share int64 `json:"share" db:"-"`
	Reply int64 `json:"reply" db:"-"`

	// 推荐相关（非数据库字段）
	Score       float64           `json:"score" db:"-"`        // 排序分数
	Reason      string            `json:"reason" db:"-"`       // 推荐理由
	RecallTypes string            `json:"recall_types" db:"-"` // 召回类型（多个用|分隔）
	RecallTags  string            `json:"recall_tags" db:"-"`  // 召回标签
	Extra       map[string]string `json:"extra" db:"-"`        // 额外信息
}

// 视频状态常量
const (
	StateNormal      = 1 // 正常
	StateReview      = 3 // 回查可放出
	StateHighQuality = 4 // 优质
	StateSelection   = 5 // 精选
)

// 行为类型
const (
	BehaviorPlay   = 1 // 播放
	BehaviorLike   = 2 // 点赞
	BehaviorFav    = 3 // 收藏
	BehaviorShare  = 4 // 分享
	BehaviorFollow = 5 // 关注
)

// Redis Key 常量
const (
	RedisKeyUserProfile      = "RECOMMEND:USER_PROFILE:%d" // mid
	RedisKeyRecommendHistory = "RECOMMEND:HISTORY:%d"      // mid
)
