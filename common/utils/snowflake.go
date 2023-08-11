package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

// GenSnowFlakeId
// @Description: 生成雪花算法id，默认nodeId为1
// @return int64
// @return error
func GenSnowFlakeId() (int64, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		logx.Errorw("生成雪花id失败：", logx.LogField{
			Key:   "error:",
			Value: err,
		})
		return 0, errors.Wrap(err, "雪花工具包错误")
	}
	return node.Generate().Int64(), nil
}
