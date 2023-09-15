package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func GetCurrentUnixTimestamp() int64 {
	return time.Now().Unix()
}

func FromInt64TimeStampToProtobufTimestamp(timestamp int64) *timestamppb.Timestamp {
	t := time.Unix(timestamp, 0)
	return timestamppb.New(t)
}

func FromUnixTimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}
