# Tenon G3N backend

独立的 G3N 游戏后端模块。Tenon 主模块仍默认使用 Gio，不会引入 G3N、OpenGL 或 OpenAL 依赖。

## 使用

```go
package main

import (
    "time"

    "github.com/sjm1327605995/tenon/backend/g3n"
    "github.com/sjm1327605995/tenon/pkg/ui"
    "github.com/sjm1327605995/tenon/pkg/value"
)

func main() {
    g3n.Use(g3n.Options{
        Launch: func(game *g3n.Game) {
            // 使用 G3N 原生 scene、mesh、model、light、physics、audio API。
        },
        Update: func(game *g3n.Game, dt time.Duration) {
            // 游戏逐帧逻辑。
        },
    })
    ui.Run(ui.Text("Tenon UI on G3N"))
}
```

导入模块后由 G3N 接管 `ui.Run`。游戏场景和 Tenon UI 使用独立场景根，刷新 UI 不会删除游戏对象。字体复用 Tenon 内置 OPPOSans。

本仓库已配置 Go workspace，可直接运行：

```sh
cd backend/g3n
go build ./...
go run ./example/basic
```

示例创建一个原生 G3N 球体、环境光和点光，再叠加 Tenon UI。适配层不需要资源包、额外下载或链接参数；G3N 自身的平台要求仍以其[官方说明](https://github.com/g3n/engine#dependencies)为准（C 编译器、OpenGL，以及 Unix 的 OpenAL/Vorbis 开发库）。

## 声明式 3D 组件

`GameScene3D` 管理组件挂载、逐帧更新、资源共享和关闭清理。业务使用 Tenon 的 `Mesh3D`、`PerspectiveCamera3D`、`Light3D`、`Material3DProps` 和 `Transform3DProps` 声明场景；`Scene3D` 是 G3N 适配器。旧的 `GameScene`/`Mesh` API 保留用于原生扩展：

```go
spin := value.New(ui.V3(0, 1, 0))
scene := g3n.Scene3D(ui.GameScene3D(
    ui.PerspectiveCamera3D(ui.Camera3DProps{FOV: 60}),
    ui.Mesh3D(ui.Mesh3DProps{
        Shape:    ui.Sphere3D(1, 32, 32),
        Material: ui.Material3DProps{Color: value.Static(ui.Color{R: 41, G: 158, B: 235, A: 255})},
        Transform: ui.Transform3DProps{Spin: spin},
    }),
    ui.AmbientLight3D(ui.Light3DProps{Intensity: value.Static(float32(0.8))}),
))
g3n.Use(g3n.Options{Launch: scene.Launch, Update: scene.Frame})
ui.Run(ui.Text("Declarative 3D"))
```

- 相同 Shape 和 Material 在单个 Scene 内共享原生资源。
- `Position`、`Rotation`、`Scale` 通过 `value.Value[Vec3]` 订阅更新，不重建 Mesh。
- `Spin` 按秒增量旋转，修改绑定值即可暂停或改变速度。
- `Scene.Close` 退订并只移除自己的场景根，不影响其他原生对象；后端退出时自动调用。共享缓存与 Mesh 分别持有资源引用，`GeometryShape` 传入的应用几何不会被场景销毁。
- `MeshComponent.Native`、`CameraComponent.Native` 和 `LightComponent.Native` 保留原生出口。
- 游戏组件实现 `HandlePointer(PointerEvent)` 即可接收场景指针输入；后端会转换场景视口坐标，并阻止交互 UI 的点击穿透。

## 统一场景格式与插件

`pkg/ui` 的 `GameScene3D` 是后端中立的场景描述：Mesh、Camera、Light 和 Transform 不携带 G3N 类型。`Scene3DPluginRegistry` 负责把这些声明解析为原生组件，`Scene3DRuntime` 统一逐帧更新和逆序关闭。注册表按场景创建，不是全局单例。`Close` 不会丢弃组件定义，调用 `Activate` 可以在同一运行时重新挂载，因此各后端共享 `Mount → Update → Close → Mount` 生命周期。它与 UI `BackendRuntime` 的终态 `Close` 不同，详见[运行时生命周期契约](../../docs/RUNTIME_LIFECYCLE.md)。

内建节点使用 `Mesh3D`、`PerspectiveCamera3D`、`AmbientLight3D` 等构造函数；第三方节点使用 `Plugin3D(name, payload)` 或精确版本的 `Plugin3DVersion(name, version, payload)`。二者共用同一条适配与生命周期路径：

```go
registry := g3n.NewScene3DPluginRegistry()
registry.RegisterPluginVersion("particles", 1, &ParticleAdapter{})
if !registry.SupportsPlugin("particles", 1) {
    return fmt.Errorf("particles v1 is unavailable")
}
scene := g3n.GameSceneRuntime(registry.MustBuild(ui.GameScene3D(
    ui.PerspectiveCamera3D(ui.Camera3DProps{FOV: 60}),
    ui.Plugin3DVersion("particles", 1, ParticleSpec{Count: 1000}),
))
```

适配器可同时实现 `Scene3DValidator`，在创建原生资源前校验 payload：

```go
func (*ParticleAdapter) Validate3D(node *ui.GameNode3D) error {
    spec, ok := node.Native.(ParticleSpec)
    if !ok || spec.Count <= 0 {
        return fmt.Errorf("particles: payload must contain a positive Count")
    }
    return nil
}
func (a *ParticleAdapter) Adapt3D(node *ui.GameNode3D) (ui.Game3DComponent, error) {
    return &MyParticleComponent{Spec: node.Native.(ParticleSpec)}, nil
}
```

版本号必须是正整数，声明与注册精确匹配；版本化节点不会回退到未版本化适配器。旧的 `RegisterPlugin` 和 `Plugin3D` 作为 version 0 保留，独立升级现有代码不受影响。`SupportsPlugin(name, version)` 只查询注册能力，不执行插件代码。

插件组件只需实现 `Update(time.Duration)` 和 `Close()`；需要挂载原生资源时实现 `Mount(*SceneContext) error`，与内建组件一致。未知版本、校验失败或适配失败都会在构建时返回错误，已成功创建的组件会先清理。`GameScene` 和 `Native3D` 保留给直接使用 G3N 原生 API 的扩展。完整的原生版本化、带校验插件示例位于本模块的 `example/plugin`，可在启用 G3N 环境后运行 `go run ./example/plugin`。仅验证后端中立公开契约时，也可从仓库根目录运行 `go run ./example/plugin`，该示例不依赖 CGO，覆盖 payload 校验与 Close/Activate 生命周期。

## 当前边界

- 已实现：共享 G3N 窗口与场景、原生游戏 Launch/Update、Tenon 文本/矩形/图片、基础描边/渐变、鼠标/键盘/字符输入和系统剪贴板。
- 指针反馈与 Gio 共用命中语义：输入框使用 IBeam、可点击元素使用 Hand、禁用或空白区域使用 Arrow，仅在角色变化时调用原生窗口。
- 文本光标与下划线使用连续线条，不再按间隔绘制点阵。
- SVG 填充与描边与 Gio 共用同一份路径解析，通过 `x/image/vector` 抗锯齿光栅化并按外观缓存；`example/uiparity` 可验证图标、文字样式、输入框和按钮。
- 内置字体与 Gio 使用相同的 Semibold(600) 合成粗体阈值和约 12° 合成斜体；测量仍使用原始字体，不因合成样式改变布局。
- 2D 图层缩放、共享中心旋转、平移和整组透明度已通过几何回归。
- 尚未达到 Gio 像素级一致：圆角裁剪和精确 3D 投影目前降级或忽略。
- G3N v0.2.0 技术栈较旧（OpenGL 3.3），选择它是看重完整游戏 API 和简单导入，而不是现代渲染特性。
- 主模块与适配模块版本需在发布时同步更新。

调用方可以在选择后端后查询高级能力；未声明能力的第三方后端会保守返回全 false。能力表示完整支持，而不是存在一个可能降级的绘制入口：

```go
capabilities := ui.CurrentBackendCapabilities()
if !capabilities.SVGPaths {
    // 使用位图或后端自定义矢量出口。
}
```

当前可查询字段为 `SVGPaths`、`RoundedClipping`、`LayerTransforms`、`ExactProjection` 和 `Accessibility`。Gio 声明支持这些能力，其中无障碍是基础文本、按钮和输入框语义；G3N 声明 `SVGPaths` 和 `LayerTransforms`，其余尚未达到一致的能力仍返回 false。

## 原生窗口 smoke

在具备 CGO、G3N 原生库和显示环境的机器上，可用 `TENON_SMOKE_FRAMES=3 go run ./example/plugin` 运行三帧后自动退出。无显示器的 Linux 使用 `xvfb-run -a`。`TENON_SMOKE_FRAMES=12 go run ./example/lifecycle` 会先产生一次故意 Mount 失败，再恢复并执行两轮 Close/Remount，退出时断言挂载尝试次数。`TENON_UI_SMOKE=1 go run ./example/uiparity` 会通过原生事件队列点击按钮、Tab 聚焦、输入中文、Ctrl+A 全选删除并重新输入，断言状态后退出；它同时验证 SVG 图标、文字样式、命中、焦点导航与绑定更新。G3N 也支持 `TENON_CAPTURE=out.png TENON_CAPTURE_FRAMES=3`，在最终渲染后通过 `ReadPixels` 写入 PNG，可与根目录 `example/uiparity` 的 Gio 输出直接比较。Ubuntu CI 已配置这三个真实窗口验证；它们仍须实际运行成功后才算平台验证通过。
