# value

可扩展的 UI 数据协议，核心只有两个名字：

```go
type Value[T any] interface {
    Get() T
    Subscribe(func(T)) func()
}

type Bind[T any] interface {
    Value[T]
    Set(T)
}
```

`Static` 创建静态值，`New` 创建 Bind，`Map` 创建只读派生值。自定义数据源实现这些方法即可接入，不需要继承内置类型。

## 自动绑定

```go
name := value.New("Alice")

ui.VStack(12,
    ui.Text(value.Static("个人资料")),
    shadcn.TextField(name),
    ui.Text(name),
    ui.Text(value.Map(name, func(s string) string { return "你好，" + s })),
)

name.Set("张三")
```

输入框的编辑会自动写回 `name`。`Set` 会同步输入框、文本与派生标签，无需 `Context.Update` 或额外的 `OnChange`。文本与输入框直接更新对应渲染节点，不重渲染整个组件。

`ui.Text` 接收普通字符串或 `Value[string]`；其他类型可先用 `Map` 格式化。可编辑输入接收 `Bind[string]`。原生输入可写 `ui.Input(ui.BindInput(name))`；主题化输入可写 `shadcn.TextField(name)`，或 `shadcn.Input(InputProps{Bind: name})`。`TextareaProps.Bind` 同样支持双向绑定。

布局、条件分支等依赖值时，在组件中用 `ui.UseValue(source)` 读取并订阅，只重渲染该订阅组件；普通 `Get()` 读取快照。

## 钩子

```go
name := value.New("", value.Hooks[string]{
    BeforeSet: func(old, next string) (string, error) {
        if len(next) > 100 {
            return old, errors.New("姓名过长")
        }
        return next, nil
    },
    OnChange:      func(old, next string) { log.Println(old, next) },
    OnError:       func(err error) { log.Println(err) },
    OnSubscribe:   func() { /* 一个消费者开始订阅 */ },
    OnUnsubscribe: func() { /* 一个消费者结束订阅 */ },
})
```

`BeforeSet` 可转换或拒绝写入。转换后的值与旧值相等时不通知；拒绝时保留旧值并调用 `OnError`。输入框会回读规范值，保证校验失败或转换后显示值与数据一致。`OnChange` 在提交之后、订阅通知之前调用。钩子可省略，扩展数据源也不必实现它们。

通知按订阅顺序执行。通知中的重入 `Set` 排队处理，当前通知的消费者看到同一个快照。退订函数可重复调用；框架在值被替换、绑定被移除、节点卸载或窗口关闭时自动退订。无头测试通过 `Harness.Close()` 释放整棵树。

## 生命周期与兼容

- 值可以在组件外创建、共享，卸载控件只结束订阅，不销毁数据对象。
- 切片、映射等引用数据用副本替换后再 `Set`；`Get` 返回当前值，不会追踪普通 Go 字段或容器的原地修改。
- 将持久 Bind 放在组件字段中，并在创建组件时初始化；不要每次 Render 都重新创建它。
- 操作和通知由调用方串行化。UI 绑定值在 UI 线程中读写；后台任务通过 `ui.Post` 或组件 Context 的 `Post` 调用 `Set`。
- 自定义 `Subscribe` 可以同步通知初值；框架会订阅后再读快照。取消后迟到的通知会被忽略。
- 原有 `ui.Signal` 已满足此协议，继续支持 `Get/Set/Subscribe`、`ui.Bind` 和 `ui.BindInput`。
- 新协议位于独立的 `value` 包，保留 `ui.Value(...)` 和 `ui.Bind(...)` 原有属性函数的名称。

完整示例：`go run ./example/component`。
