package tool

import (
	"strconv"
	"strings"
)

// JoinInts 将 int64 切片转换为逗号分隔的字符串（参考主项目 xstr.JoinInts）
// 用于拼接 SQL 的 IN 语句：WHERE id IN (1,2,3)
func JoinInts(ints []int64) string {
	if len(ints) == 0 {
		return ""
	}

	strs := make([]string, 0, len(ints))
	for _, v := range ints {
		strs = append(strs, strconv.FormatInt(v, 10))
	}
	return strings.Join(strs, ",")
}

// SplitInts 将逗号分隔的字符串转换为 int64 切片（参考主项目 xstr.SplitInts）
func SplitInts(s string) ([]int64, error) {
	if s == "" {
		return []int64{}, nil
	}

	strs := strings.Split(s, ",")
	ints := make([]int64, 0, len(strs))
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		ints = append(ints, v)
	}
	return ints, nil
}

// ContainsInt 判断切片中是否包含某个元素
func ContainsInt(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// UniqueInts 对 int64 切片去重
func UniqueInts(ints []int64) []int64 {
	if len(ints) == 0 {
		return ints
	}

	seen := make(map[int64]struct{}, len(ints))
	result := make([]int64, 0, len(ints))

	for _, v := range ints {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
