package engine

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

// Widget 是不可变的 UI 配置描述对象。
// 每次状态变化 rebuild 都会创建新的 Widget 实例。
// Widget 本身不包含任何状态，只描述「UI 应该长什么样」。
type Widget interface {
	// CreateElement 为该 Widget 创建对应的 Element。
	CreateElement() Element

	// Key 返回该 Widget 的标识键，用于 Element 复用时的匹配。
	// 返回 nil 表示无 Key。
	GetKey() Key
}

// Key 用于在 Widget diff 时识别和复用 Element。
type Key interface {
	Equals(other Key) bool
	String() string
}

// NilKey 表示没有 Key。
type NilKey struct{}

func (NilKey) Equals(other Key) bool { _, ok := other.(NilKey); return ok }
func (NilKey) String() string       { return "NilKey" }

// IsNilKey 判断 Key 是否为 NilKey。
func IsNilKey(k Key) bool {
	_, ok := k.(NilKey)
	return ok
}

// ValueKey 基于值相等的 Key，是最常用的 Key 类型。
type ValueKey[T comparable] struct {
	Value T
}

func NewValueKey[T comparable](v T) *ValueKey[T] {
	return &ValueKey[T]{Value: v}
}

// NewStringKey 创建字符串类型的 ValueKey（非泛型便捷函数）。
func NewStringKey(v string) *ValueKey[string] {
	return &ValueKey[string]{Value: v}
}

func (k *ValueKey[T]) Equals(other Key) bool {
	if o, ok := other.(*ValueKey[T]); ok {
		return k.Value == o.Value
	}
	return false
}

func (k *ValueKey[T]) String() string {
	return fmt.Sprintf("ValueKey(%v)", k.Value)
}

// LocalKey 基于内存地址的 Key，用于框架内部自动生成。
type LocalKey struct {
	id int
}

var localKeyCounter int64

func NewLocalKey() *LocalKey {
	id := atomic.AddInt64(&localKeyCounter, 1)
	return &LocalKey{id: int(id)}
}

func (k *LocalKey) Equals(other Key) bool {
	if o, ok := other.(*LocalKey); ok {
		return k.id == o.id
	}
	return false
}

func (k *LocalKey) String() string {
	return fmt.Sprintf("LocalKey(%d)", k.id)
}

// GlobalKey 是全局唯一的 Key，可跨树访问对应的 Element/State。
type GlobalKey struct {
	id int64
}

var globalKeyCounter int64

// NewGlobalKey 创建一个新的全局唯一 Key。
func NewGlobalKey() *GlobalKey {
	id := atomic.AddInt64(&globalKeyCounter, 1)
	return &GlobalKey{id: id}
}

func (k *GlobalKey) Equals(other Key) bool {
	if o, ok := other.(*GlobalKey); ok {
		return k.id == o.id
	}
	return false
}

func (k *GlobalKey) String() string {
	return fmt.Sprintf("GlobalKey(%d)", k.id)
}

// CurrentContext 返回与该 GlobalKey 关联的 BuildContext。
func (k *GlobalKey) CurrentContext() BuildContext {
	if el := getGlobalKeyElement(k); el != nil {
		if se, ok := el.(*StatefulElement); ok {
			return se.buildContext
		}
	}
	return nil
}

// CurrentWidget 返回与该 GlobalKey 关联的 Widget。
func (k *GlobalKey) CurrentWidget() Widget {
	if el := getGlobalKeyElement(k); el != nil {
		return el.GetWidget()
	}
	return nil
}

// CurrentState 返回与该 GlobalKey 关联的 State（如果是 StatefulElement）。
func (k *GlobalKey) CurrentState() State {
	if el := getGlobalKeyElement(k); el != nil {
		if se, ok := el.(*StatefulElement); ok {
			return se.state
		}
	}
	return nil
}

// BaseWidget 提供 Widget 的默认实现。
// 用户自定义 Widget 可内嵌此结构体。
type BaseWidget struct {
	key Key
}

func (b *BaseWidget) SetKey(k Key) {
	b.key = k
}

func (b BaseWidget) GetKey() Key {
	if b.key == nil {
		return NilKey{}
	}
	return b.key
}

// CanUpdate 判断两个 Widget 是否可以在 Element 更新时复用同一个 Element。
// 条件：1. runtimeType 相同；2. Key 相同（或都无 Key）。
func CanUpdate(oldWidget, newWidget Widget) bool {
	if oldWidget == nil || newWidget == nil {
		return false
	}
	// 比较类型
	if reflect.TypeOf(oldWidget) != reflect.TypeOf(newWidget) {
		return false
	}
	// 比较 Key
	return oldWidget.GetKey().Equals(newWidget.GetKey())
}

// BuildFunc 是用户提供的构建函数，每次 rebuild 被调用，产出 Widget 树。
type BuildFunc func() Widget
