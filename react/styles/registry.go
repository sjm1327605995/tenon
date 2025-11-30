package styles

import (
	"reflect"
	"sync"
)

// StyleHandler 定义样式处理函数类型
type StyleHandler func(view interface{}, style interface{})

// StyleRegistry 样式注册器
type StyleRegistry struct {
	// handlers 存储元素类型到样式类型的处理函数映射
	// 格式: elementType -> styleType -> handler
	handlers map[string]map[string]StyleHandler
	mutex    sync.RWMutex
}

// NewStyleRegistry 创建一个新的样式注册器
func NewStyleRegistry() *StyleRegistry {
	return &StyleRegistry{
		handlers: make(map[string]map[string]StyleHandler),
	}
}

// Register 注册样式处理函数
func (r *StyleRegistry) Register(elementType string, styleType string, handler StyleHandler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.handlers[elementType]; !exists {
		r.handlers[elementType] = make(map[string]StyleHandler)
	}

	r.handlers[elementType][styleType] = handler
}

// Handle 处理指定元素和样式
func (r *StyleRegistry) Handle(view interface{}, elementType string, style interface{}) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 获取元素类型的处理器映射
	elementHandlers, exists := r.handlers[elementType]
	if !exists {
		return false
	}

	// 获取样式类型名称
	styleTypeName := typeName(style)

	// 获取对应的处理函数
	handler, exists := elementHandlers[styleTypeName]
	if !exists {
		return false
	}

	// 调用处理函数
	handler(view, style)
	return true
}

// 全局样式注册器实例
var globalRegistry = NewStyleRegistry()

// RegisterGlobalHandler 注册全局样式处理函数
func RegisterGlobalHandler(elementType string, styleType string, handler StyleHandler) {
	globalRegistry.Register(elementType, styleType, handler)
}

// HandleGlobalStyle 处理全局样式
func HandleGlobalStyle(view interface{}, elementType string, style interface{}) bool {
	return globalRegistry.Handle(view, elementType, style)
}

// typeName 获取值的类型名称
func typeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}
