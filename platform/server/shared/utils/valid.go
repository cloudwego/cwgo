package utils

import "github.com/cloudwego/cwgo/platform/server/shared/consts"

func ValidStrings(ss ...string) bool {
	for _, s := range ss {
		if len(s) == 0 {
			return false
		}
	}

	return true
}

func ValidOrder(order int32) bool {
	if order != consts.OrderNumInc && order != consts.OrderNumDec {
		return false
	}

	return true
}

func ValidOrderBy(orderBy string) bool {
	if orderBy != consts.OrderByCreateTime && orderBy != consts.OrderBySyncTime && orderBy != consts.OrderByUpdateTime {
		return false
	}

	return true
}

func ValidStatus(status string) bool {
	if status != consts.Active && status != consts.DisActive {
		return false
	}

	return true
}

func ValidRepoType(_type int32) bool {
	if _type != consts.GitLab && _type != consts.Github {
		return false
	}

	return true
}
