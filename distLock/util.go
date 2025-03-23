package distLock

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateOwner 生成唯一锁ID
func GenerateOwner() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String())
}
