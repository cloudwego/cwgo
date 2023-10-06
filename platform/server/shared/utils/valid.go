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

func ValidStatus(status string) bool {
	if status != consts.RepositoryStatusActive && status != consts.RepositoryStatusDisactive {
		return false
	}

	return true
}
