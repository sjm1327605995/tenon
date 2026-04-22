package hooks

import (
	"fmt"
	"sync/atomic"
)

var globalIdCounter int64

// UseId 生成唯一的 ID（React 18+）
// 用于生成稳定的、唯一的标识符
func UseId() string {
	hook := getWorkInProgressHook()

	if hook.MemoizedState == nil {
		id := generateUniqueID()
		hook.MemoizedState = id
		return id
	}

	return hook.MemoizedState.(string)
}

func generateUniqueID() string {
	counter := atomic.AddInt64(&globalIdCounter, 1)
	return fmt.Sprintf(":r%d:", counter)
}
