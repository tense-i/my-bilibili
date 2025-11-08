package xerr

// 错误消息映射
var msgMap = map[uint32]string{
	OK:                  "成功",
	SERVER_COMMON_ERROR: "服务器内部错误",
	REQUEST_PARAM_ERROR: "参数错误",
	DB_ERROR:            "数据库错误",
	DB_UPDATE_ERROR:     "数据库更新失败",

	VIDEO_NOT_FOUND:      "视频不存在",
	VIDEO_STAT_NOT_FOUND: "视频统计数据不存在",
	VIDEO_INFO_ERROR:     "获取视频信息失败",
	VIDEO_STAT_ERROR:     "获取视频统计数据失败",
	VIDEO_LIST_ERROR:     "获取视频列表失败",

	HOTRANK_NOT_FOUND:       "热门排行榜数据不存在",
	HOTRANK_CALCULATE_ERROR: "热度计算失败",
	HOTRANK_UPDATE_ERROR:    "热度更新失败",
	HOTRANK_QUERY_ERROR:     "热门排行榜查询失败",

	USER_NOT_FOUND:  "用户不存在",
	USER_AUTH_ERROR: "用户认证失败",
}

// MapErrMsg 根据错误码获取错误消息
func MapErrMsg(errCode uint32) string {
	if msg, ok := msgMap[errCode]; ok {
		return msg
	}
	return msgMap[SERVER_COMMON_ERROR]
}

// IsCodeErr 判断是否为自定义错误码
func IsCodeErr(errCode uint32) bool {
	_, ok := msgMap[errCode]
	return ok
}
