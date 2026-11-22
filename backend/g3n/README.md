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

## 当前边界

- 已实现：共享 G3N 窗口与场景、原生游戏 Launch/Update、Tenon 文本/矩形/图片、基础描边/渐变、鼠标/键盘/字符输入和系统剪贴板。
- 尚未达到 Gio 像素级一致：SVG 路径、圆角裁剪、图层变换和精确投影目前降级或忽略。
- G3N v0.2.0 技术栈较旧（OpenGL 3.3），选择它是看重完整游戏 API 和简单导入，而不是现代渲染特性。
- 主模块与适配模块版本需在发布时同步更新。
