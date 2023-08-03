package dal

import (
	"cwgo/example/hex/biz/dal/mysql"
	"cwgo/example/hex/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
