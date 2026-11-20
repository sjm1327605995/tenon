# tenon/pkg/decl

用 Go 结构体声明 UI，具名字段可直接补全。页面返回 `decl.View`，入口和子节点的转换由框架自动完成。

```go
package main

import (
    "fmt"
    d "github.com/sjm1327605995/tenon/pkg/decl"
    ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func Profile() d.View {
    name, setName := ui.UseState("")
    return d.Column{
        Width: 320, Gap: 12, Padding: 24,
        Children: []d.View{
            d.Text{Value: "个人资料", FontSize: 22},
            d.Input{Value: name, OnChange: setName, Placeholder: "姓名"},
            d.Button{Text: "保存", OnClick: func() { fmt.Println(name) }},
        },
    }
}

func main() { d.Run(Profile) }
```

运行结构体声明示例：`go run ./example/struct`。业务组件可内嵌 `Component`，实现 `Render() View`，通过 `p.Context()` 获取主题、生命周期与 UI 调度，详见 `go run ./example/component`。

## 组件与状态

- 对象组件入口：`d.Run(&Profile{})`；`d.Render(component)` 接收 `IComponent` 并返回普通节点。业务代码不显式转换声明。
- 函数页面仍可用 `d.Run(Profile)`；测试入口接受同样的根类型，例如 `d.Mount(&Profile{}, 640, 480)`。
- `Component` 是内嵌基类。嵌套函数组件使用 `d.Func(Profile)`；对象组件可直接放入子节点。
- 带参数组件：`d.Use(Row, props)`；比较 props 后跳过重渲染：`d.Memo(Row, props)`。这些函数也直接返回声明。
- `Children: []d.View{...}` 可混用结构体和现有 `*ui.Node`，例如 `ui.Text(...)`、`ui.Use(...)`、`ui.For(...)`；nil 子节点会被忽略。
- 静态声明结构体描述本轮 UI。业务对象组件使用指针：持久实例的字段可通过 `Context().Update` 或 `SetState` 通知变化；响应式数据优先使用 [`value.Bind`](../value)，调用 `Set` 自动更新控件。
- 同类型、同 key 的新组件描述保留 Context 中的 `d.State(ctx, initial)`；普通结构体字段视为当前实例/props，不自动复制旧实例的字段状态。一个对象实例只挂载到一个位置。
- 组件转换保留原函数身份和 key；有状态列表用 `d.Use(Row, item).Key(item.ID)`，无需手动转换行组件。

## 声明类型

| 声明 | 常用字段 |
| --- | --- |
| `Box` / `Column` / `Row` / `ScrollView` | `Width`、`Height`、`Padding`、`Gap`、`Background`、`Radius`、`Disabled`、`OnClick`、`Children` |
| `Text` | `Value`、`FontSize`、`Color`、`Bold` |
| `Button` | `Text`、`Variant`、`Size`、`Disabled`、`OnClick`、`Children` |
| `Input` | `Value`、`OnChange`、`Placeholder`、`OnSubmit`、`Disabled` |
| `Textarea` | `Value`、`OnChange`、`Placeholder`、`Rows`、`Disabled` |
| `Checkbox` / `Switch` | `Checked`、`OnChange`、`Disabled` |
| `Fragment` | `Key`、`Children`；不增加布局容器 |
| `ThemeProvider` | `Theme`、`Children` |

`Row` 默认垂直居中，与 `ui.HStack` 一致。控件复用 shadcn 的主题、禁用外观、动画和交互。声明方法只负责适配现有 Node，不新增渲染引擎。

## 样式与默认值

尺寸字段的零值表示自动/默认，颜色零值表示保持继承或组件默认值。`Style: []ui.StyleOpt{...}` 在具名样式字段之后应用，可表达百分比之外的高级布局、动画，以及显式零尺寸，例如 `Style: []ui.StyleOpt{ui.Width(0)}`。容器还可用 `Attrs` 附加 `UseMeasure`、`ArrowNav` 等现有属性。

结构体字段负责描述，方法负责内部转换，因此没有 `Width` 字段与链式 `Width()` 方法重名的问题。

热重载视图可以直接声明 `func View() d.View` 并返回结构体。解释器原有的泛型限制仍适用，Hooks 和泛型组件放在编译后的宿主中。

## Context 能力

`Context()` 由内嵌基类提供，挂载时绑定。它支持 `Theme()`、`Viewport()`、`Effect(...)`、`Post(...)`、`Invalidate()`，并实现标准 `context.Context`；卸载时取消生命周期，忽略迟到的更新。UI 状态和依赖读取在 UI 线程执行，后台工作通过 `Post` 回到 UI 线程。

文本和输入采用 Bind 时，无需手写更新回调：`shadcn.TextField(name)`、`ui.Text(name)`。可选结构体声明也支持 `Text{Bind: name}`、`Input{Bind: name}`、`Textarea{Bind: name}`。
