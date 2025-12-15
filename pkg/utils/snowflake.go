package utils

import (
	"errors"
	"sync"

	"github.com/sony/sonyflake"
)

var (
	sf   *sonyflake.Sonyflake
	once sync.Once
)

// Init 初始化 Sonyflake（只允许调用一次）
func InitSnowflake() error {
	var snowflakeInitErr error

	once.Do(func() {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{})
		if sf == nil {
			snowflakeInitErr = errors.New("sonyflake初始化失败!")
		}
	})

	return snowflakeInitErr
}

// NextID 生成全局唯一 ID
func NextUniqueID() (uint64, error) {
	if sf == nil {
		return 0, errors.New("sonyflake没有初始化!")
	}
	return sf.NextID()
}
