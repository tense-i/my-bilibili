package model

// OArchive 热度排行榜表结构（对应 academy_archive 表）
type OArchive struct {
	ID       int64 `db:"id"`
	OID      int64 `db:"oid"`      // 对象ID（视频ID）
	Business int   `db:"business"` // 业务类型：1-视频，2-专栏
}

