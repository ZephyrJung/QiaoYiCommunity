package utils

import (
	"fmt"
	"sync"
	"time"
	"github.com/bwmarrin/snowflake"
)

var (
	sfNode *snowflake.Node
	sfOnce sync.Once
)

func InitSnowflake(nodeID int64) error {
	var err error
	sfOnce.Do(func() {
		sfNode, err = snowflake.NewNode(nodeID)
	})
	return err
}

func GenerateID() int64 {
	if sfNode == nil {
		panic("snowflake not initialized")
	}
	return sfNode.Generate().Int64()
}

func GenerateIDString() string {
	return fmt.Sprintf("%d", GenerateID())
}

func ParseID(id int64) time.Time {
	ms := snowflake.ParseInt64(id).Time()
	return time.Unix(0, ms*int64(time.Millisecond))
}
