package dao

import (
	"context"
	"fmt"
	"time"

	"mybilibili/app/hotrank/cmd/job/internal/model"
	"mybilibili/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	// 查询SQL（参考主项目）
	_getArcsSQL = "SELECT id,oid,business FROM academy_archive WHERE state=? AND business=? AND id > ? ORDER BY id ASC LIMIT ?"
)

type AcademyDao struct {
	conn sqlx.SqlConn
}

func NewAcademyDao(conn sqlx.SqlConn) *AcademyDao {
	return &AcademyDao{
		conn: conn,
	}
}

// Archives 游标分页查询（完全参考主项目）
func (d *AcademyDao) Archives(ctx context.Context, id int64, bs, limit int) ([]*model.OArchive, error) {
	var res []*model.OArchive

	err := d.conn.QueryRowsCtx(ctx, &res, _getArcsSQL, 0, bs, id, limit)
	if err != nil {
		logx.Errorf("AcademyDao.Archives QueryRows error(%v)", err)
		return nil, err
	}

	return res, nil
}

// UPHotByAIDs 批量更新热度值（完全参考主项目的 CASE WHEN 实现）
func (d *AcademyDao) UPHotByAIDs(ctx context.Context, hots map[int64]int64) error {
	if len(hots) == 0 {
		return nil
	}

	var oids []int64

	// 构建 CASE WHEN SQL（参考主项目）
	sqlStr := "UPDATE academy_archive SET hot = CASE oid "
	for oid, hot := range hots {
		sqlStr += fmt.Sprintf("WHEN %d THEN %d ", oid, hot)
		oids = append(oids, oid)
	}
	sqlStr += fmt.Sprintf("END, mtime=? WHERE oid IN (%s)", tool.JoinInts(oids))

	_, err := d.conn.ExecCtx(ctx, sqlStr, time.Now())
	if err != nil {
		logx.Errorf("AcademyDao.UPHotByAIDs ExecCtx sql(%s) error(%v)", sqlStr, err)
		return err
	}

	logx.Infof("AcademyDao.UPHotByAIDs success: updated %d records", len(hots))
	return nil
}




