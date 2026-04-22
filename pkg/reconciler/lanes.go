package reconciler

// Lane 表示单个优先级车道
type Lane = int

// Lanes 表示多个优先级的位掩码
type Lanes = int

const (
	NoLane                Lane = 0
	SyncLane              Lane = 1
	InputContinuousLane   Lane = 1 << 1
	DefaultLane           Lane = 1 << 2
	TransitionLane1       Lane = 1 << 3
	TransitionLane2       Lane = 1 << 4
	IdleLane              Lane = 1 << 5
)

// 预定义优先级集合
const (
	SyncUpdateLanes      Lanes = SyncLane
	DefaultUpdateLanes   Lanes = DefaultLane | SyncLane | InputContinuousLane
	TransitionLanes      Lanes = TransitionLane1 | TransitionLane2
)

// MergeLanes 合并车道
func MergeLanes(a, b Lanes) Lanes {
	return a | b
}

// RemoveLanes 从集合中移除车道
func RemoveLanes(set Lanes, subset Lanes) Lanes {
	return set &^ subset
}

// IncludesLane 检查集合是否包含某个车道
func IncludesLane(set Lanes, lane Lane) bool {
	return (set & lane) != 0
}

// IncludesSomeLane 检查两个集合是否有交集
func IncludesSomeLane(a, b Lanes) bool {
	return (a & b) != 0
}

// GetHighestPriorityLane 获取最高优先级的车道
func GetHighestPriorityLane(lanes Lanes) Lane {
	return lanes & -lanes
}

// MarkRootUpdated 标记根节点更新（简化版）
func MarkRootUpdated(lanes Lanes, updateLane Lane) Lanes {
	return MergeLanes(lanes, updateLane)
}
