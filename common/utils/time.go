package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func GetCurrentTimeStamp() int64 {
	return time.Now().Unix()
}

func FromInt64TimeStampToProtobufTimeStamp(timestamp int64) *timestamppb.Timestamp {
	t := time.Unix(timestamp, 0)
	return timestamppb.New(t)
}
