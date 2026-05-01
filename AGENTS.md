# Tenon — AI 开发约束与踩坑记录

> 本文档记录框架中**非显而易见的设计约束**和**已踩过的坑**，供 AI 编码时参考。

---

## 一、Go 嵌入类型的方法调用陷阱

### 现象

`RenderObjectElement.Mount()` 中调用 `r.CreateRenderObject()`，但子类（如 `ColumnElement`）覆盖的 `CreateRenderObject()` 不会被调用，而是调用了 `RenderObjectElement.CreateRenderObject()`（返回 nil）。

### 根因

Go 的结构体嵌入（embedding）在方法调用时不做动态分发。当 `ColumnElement.Mount` 调用 `e.RenderObjectElement.Mount(parent, slot)` 时，`RenderObjectElement.Mount` 方法内的接收者 `r` 是 `*RenderObjectElement` 类型，而不是 `*ColumnElement`。因此 `r.CreateRenderObject()` 调用的是 `RenderObjectElement.CreateRenderObject`，不是 `ColumnElement.CreateRenderObject`。

### 修复方案

在子类的 `Mount` 方法中，**先显式调用自己的 `CreateRenderObject()` 并赋值给 `e.RenderObject`，再调用父类的 `Mount`**：

```go
func (e *ColumnElement) Mount(parent ui.Element, slot int) {
    e.RenderObject = e.CreateRenderObject()  // ← 必须先做
    e.MultiChildRenderObjectElement.Mount(parent, slot)
    // ... 创建子元素
}
```

### AI 编码约束

1. **所有 RenderObjectElement 的子类**，如果覆盖了 `CreateRenderObject()`，则必须在自己的 `Mount()` 中先调用 `e.RenderObject = e.CreateRenderObject()`，再调用父类 `Mount`。
2. **不要在父类的 Mount 中依赖动态方法分发**。

---

## 二、Widget 字段名与链式方法名冲突

### 现象

```go
type TextWidget struct {
    FontSize float32
}

func (t TextWidget) FontSize(size float32) TextWidget {  // 编译错误！
    t.FontSize = size
    return t
}
```

Go 不允许结构体有同名的字段和方法。

### 修复方案

Widget 的**字段使用小写**（不可导出），链式方法使用大写（导出）：

```go
type TextWidget struct {
    ui.BaseWidget
    fontSize float32  // 小写，不可导出
}

func (t TextWidget) FontSize(size float32) TextWidget {  // 大写，方法
    t.fontSize = size
    return t
}
```

用户通过构造函数和链式方法创建 Widget，不直接使用结构体字面量。

---

## 三、Widget 值类型与接口实现

### 现象

```go
type Widget interface {
    CreateElement() Element
    GetKey() Key
}
```

如果 `BaseWidget.GetKey()` 使用指针接收者 `func (b *BaseWidget) GetKey() Key`，则内嵌 `BaseWidget` 的值类型（如 `TextWidget`）不会自动实现 `Widget` 接口。

### 修复方案

`BaseWidget` 的方法使用**值接收者**：

```go
func (b BaseWidget) GetKey() Key { ... }  // 值接收者
```

这样值类型和指针类型都实现接口。

---

## 四、Color 类型实现 color.Color 接口

### 现象

内嵌 `color.RGBA` 会导致 `RGBA()` 方法名与字段名提升冲突。

### 修复方案

不内嵌 `color.RGBA`，而是自己定义 `R, G, B, A uint8` 字段并显式实现 `RGBA()` 方法：

```go
type Color struct {
    R, G, B, A uint8
}

func (c *Color) RGBA() (r, g, b, a uint32) {
    // 实现 color.Color 接口
}
```

---

## 五、验证清单（修改核心引擎或组件后）

- [ ] `go build ./...` 通过
- [ ] `go run _debug_snapshot/test_v3.go` 无 zero-bounds 输出
- [ ] 如果是新组件，`test_v3.go` 中包含状态变化后的 bounds 断言
- [ ] Widget 链式方法名不与字段名冲突
- [ ] RenderObjectElement 子类的 Mount 中先调用 `e.RenderObject = e.CreateRenderObject()`
