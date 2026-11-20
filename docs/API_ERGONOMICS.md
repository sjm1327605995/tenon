# Go UI 编写体验调研与简化方案

调研日期：2026-11-19。目标是减少 Tenon 的样板代码，同时保留 Go 自动补全、编译期参数检查、gofmt、现有组件身份和 Hooks 状态。

## 调研结果

以下结论来自各项目官方示例和文档；适用性是针对 Tenon 的设计判断。HTML 项目只作为 API 写法参考。

| 项目 | 写法与来源 | 对 Tenon 的启发 |
| --- | --- | --- |
| Fyne | [`NewButton(label, callback)` + `NewVBox(children...)`](https://github.com/fyne-io/fyne#example)，[类型化数据绑定](https://docs.fyne.io/binding/simple/) | 高频控件直接接收最常用参数，减少空配置结构体；输入值和更新回调应容易成对使用。 |
| go-app | [`Div().Body(...)`、`Input().Value(...).Placeholder(...)`](https://github.com/maxence-charriere/go-app#components) | 链式属性便于 IDE 发现，也能把布局结构和修饰区分开。 |
| gomponents | [纯 Go 函数、`Map`、`Group`、`If`](https://github.com/maragudk/gomponents#features) | 保留函数组合；列表辅助返回可直接嵌入的节点，减少手动拼接切片。 |
| Walk | [`MainWindow{Children: []Widget{...}}`](https://github.com/lxn/walk#example) | 复杂配置适合具名字段，但所有常用控件都写完整结构体会增加样板。 |
| Gio | [`Layout` 回调、`layout.Context` 和约束](https://gioui.org/doc/architecture/layout) | 适合底层布局控制。Tenon 已有 Fiber 和 Yoga，应用层继续隐藏这些参数。 |
| templ | [`.templ` 组件与 `templ generate`](https://templ.guide/quick-start/creating-a-simple-templ-component/) | 标记语言可读，但需要额外生成和编辑器工具；用于 Tenon 原生控件还需要另建语法与编译器。 |

## 两种声明形式，共用同一运行时

业务组件可内嵌 `Component`，实现 `Render() View` 并通过基类的 `Context()` 获取运行时能力。UI 树推荐继续使用函数和变参组织，避免反复书写 `Children: []View{...}`；结构体布局声明保留为可选形式。函数组件使用 `decl.Func`/`decl.Use`，对象组件由 `decl.Render(IComponent)` 或 `decl.Run(instance)` 自动挂载。

## 数据使用 Value / Bind

数据协议位于 [`pkg/value`](../pkg/value/README.md)：`Value[T]` 负责读取和订阅，`Bind[T]` 增加 `Set`。静态值和内置 Bind 都是实现，自定义数据源可独立实现接口。控件自动订阅，在卸载或更换数据源时退订，常规数据修改不需要 `Context.Update`。

```go
name := value.New("")
return ui.VStack(12,
    ui.Text("个人资料"),
    shadcn.TextField(name).Placeholder("姓名"),
    ui.Text(name),
).Width(320).Padding(24)
```

校验、转换、变化和订阅钩子通过 `value.Hooks[T]` 配置。派生值使用 `value.Map`，布局/条件依赖值时使用 `ui.UseValue`。现有字符串、setter 与 Signal API 保留。

## 可选的结构体布局声明

```go
func Profile() d.View {
    name, setName := ui.UseState("")
    return d.Column{
        Width: 320, Gap: 12, Padding: 24,
        Children: []d.View{
            d.Input{Value: name, OnChange: setName, Placeholder: "姓名"},
            d.Button{Text: "保存", OnClick: func() { fmt.Println(name) }},
        },
    }
}

func main() { d.Run(Profile) }
```

具名字段有完整 Go 类型检查；`[]d.View` 比函数变参多一些标点和行数，但配置集中、便于补全。页面状态仍在函数中的 Hooks 里，结构体只描述本轮 UI。`Children` 可以混入已有 Node，转换保留原组件函数身份和 key，也不增加布局容器。详见 [结构体声明说明](../pkg/decl/README.md)。

## 函数与链式声明

采用一套渐进 API：所有结果仍是 `*ui.Node`，可以和旧构造函数、原有属性函数混用。

```go
func Profile() *ui.Node {
    name, setName := ui.UseState("")
    return ui.VStack(12,
        shadcn.TextField(name, setName).Placeholder("姓名"),
        shadcn.TextButton("保存", func() { fmt.Println(name) }),
    ).Width(320).Padding(24)
}

func main() { ui.Run(ui.View(Profile)) }
```

原来的等价写法需要无参数组件占位、完整控件 props 和混在子节点中的样式：

```go
func Profile(_ struct{}) *ui.Node {
    name, setName := ui.UseState("")
    return ui.VStack(12,
        ui.Style(ui.Width(320), ui.Padding(24)),
        shadcn.Input(shadcn.InputProps{
            Value: name, OnChange: setName, Placeholder: "姓名",
        }),
        shadcn.Button(shadcn.ButtonProps{
            OnClick: func() { fmt.Println(name) },
        }, ui.Text("保存")),
    )
}

func main() { ui.Run(ui.Use(Profile, struct{}{})) }
```

这段表单省去了两个完整 props 字面量和两个 `struct{}` 占位，原来的输入受控行为和按钮样式保持一致。更少代码是否更符合具体团队习惯，仍应通过实际页面试用判断。

## 函数与链式 API

| 场景 | 写法 |
| --- | --- |
| 无 props 组件 | `ui.View(Profile)`，组件签名为 `func() *ui.Node` |
| 有 props 组件 | 继续使用 `ui.Use(Row, props)`、`ui.Memo(Row, props)` |
| 常用修饰 | `.Width(...)`、`.Height(...)`、`.Padding(...)`、`.Gap(...)`、`.Bg(...)`、`.Radius(...)`、`.FontSize(...)`、`.TextColor(...)` |
| 其他样式/属性 | `.Style(ui.Grow(1), ui.Clip)`、`.With(ui.OnHover(...))` |
| 点击与提交 | `.OnClick(func() {...})`、`.OnSubmit(func(value string) {...})` |
| 受控原生输入 | `ui.Input().Bind(name, setName).Placeholder("姓名")` |
| 主题化常用控件 | `shadcn.TextField(name, setName)`、`shadcn.TextButton("保存", save)` |
| 列表 | `ui.For(items, func(item Item) *ui.Node { return ui.Use(Row, item).Key(item.ID) })` |
| 条件渲染 | 继续使用 `ui.If(condition, node)`；条件为假时的子节点不会挂载，但参数表达式仍按 Go 规则求值。 |

复杂按钮继续用具名配置，例如 `ButtonProps{Variant: Outline, Size: SizeSm, Disabled: saving}`。不为每个复杂组件再复制一套 builder。

## 语义约束

- 修饰返回新节点，必须使用返回值；共享基础节点可以安全派生多种外观。`Keyed` 也改为返回副本。
- 修饰函数组件时直接修饰它返回的根节点，不额外插入 `Div`，组件函数身份和本地状态保留。
- Fragment、Provider、Portal 的单个根可透明传递修饰；多个根应明确用 `Div`/`VStack` 包含，再修饰容器。
- 后写的样式或事件覆盖前面的同名设置。`.Bind` 是 `Value + OnChange` 的简写，仍是受控数据流，不会观察普通 Go 变量。
- `.Disabled(...)` 提供引擎禁用行为。shadcn 的统一禁用外观仍使用组件自身的 `Disabled` props。
- 每次创建新的组件修饰都会视为新的输入，保守刷新，以更新样式和捕获状态的回调；稳定复用同一个修饰后的 Memo 节点仍可跳过重渲染。
- `For` 仅做映射；有状态行用 `ui.Use(Row, item).Key(item.ID)`，Hooks 留在行组件里，不在列表回调中直接调用。
- 本轮未引入新依赖、代码生成器或字符串样式 DSL。热重载的原有限制仍在：解释源码不支持泛型，因此 `For` 等泛型辅助在编译后的宿主中使用。

## 验证与试用

`go run ./example/struct` 展示结构体页面、嵌套组件和自动转换；`go run ./example/simple` 展示函数与链式写法。`example/accordion` 的 Input/Textarea 预览也采用了简写。

回归测试覆盖修饰分支不互相污染、组件状态保留、最新回调、受控输入、Memo、Context、列表重排，以及新旧 API 混用。`go run _debug_snapshot/test_v3.go` 包含新写法在状态变化后的 bounds 检查。
