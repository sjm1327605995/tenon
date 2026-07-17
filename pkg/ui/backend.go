package ui

import "image"

// 渲染后端句柄：这些接口是「UI 引擎」与「具体渲染后端」（当前为 gio）之间的中立边界。
// 引擎侧（renderNode 字段、painter 接口、布局/折行/富文本代码）只依赖这些接口，不出现任何
// gio 的具体类型；后端定义自己的实现类型（gioFont/gioImage/gioPath），并在其 painter 内部
// 类型断言回去。保留这层中立边界，将来若要再换/加后端仍然只改后端侧。
//
// 与 gio 生态的边界（约定，改动前先想清楚）：
//
//   - Tenon 的引擎是 gio 的**对等层**，不是 gio widgets 的使用者：我们有自己的
//     fiber/hooks 与 yoga 布局，因此 gio 的 layout / widget / material 用不上（范式冲突）。
//   - gio 的**低层能力必须走它的惯例，不要自己造**：op/clip/paint、text.Shaper、
//     io/key 的 IME 协议、io/clipboard、io/pointer 的指针与光标语义。
//     有现成的就照 gio 自己的实现（如 widget.Editor）当参考。
//   - gio 的**单位与语义原样透传，禁止自创换算**。反例：曾把 gio 的像素滚动量除以 24
//     伪造成「档位」再让引擎乘回 24*uiScale，白白多乘一次 uiScale，高 DPI 上滚动速度失真。
//     后端与引擎统一用物理像素。
//
// 这样新增/替换一个后端只需：实现 painter 接口 + 提供 fontFace/bitmap/vecPath 的具体类型
// + 通过下面的构造钩子登记如何新建它们 + 提供 run 循环，而无需改动引擎核心。

// fontFace 是后端拥有的字体句柄：布局用它测量文本，painter 用它绘制字形。
type fontFace interface {
	// Measure 返回字符串 s 在给定行高下的推进宽度（物理像素）。
	Measure(s string, lineH float64) float32
	// Metrics 返回基线以上高度（ascent，物理像素）与是否需要伪粗/伪斜合成。
	Metrics() (ascent float32, fauxBold, fauxItalic bool)
}

// bitmap 是后端拥有的位图句柄（图片元素与解码后的图像）。
type bitmap interface {
	// Size 返回像素宽高。
	Size() (w, h int)
}

// vecPath 是后端拥有的矢量路径句柄（图标）。仅作为不透明句柄在引擎与 painter 间传递。
type vecPath interface{ isVecPath() }

// 后端构造钩子：由当前激活的渲染后端在 init 时登记。引擎在需要新建句柄/启动窗口时经此调用，
// 从而不硬编码任何具体后端类型。px 均为物理像素（已乘 uiScale）。
var (
	backendNewFont    func(px float32, weight int, italic bool) fontFace // 取字体句柄；失败可返回 nil
	backendNewBitmap  func(img image.Image) bitmap                       // 解码后的图像 -> 位图句柄
	backendNewVecPath func(svgPath string, scale float32) vecPath        // SVG 路径 d -> 矢量句柄；无内容返回 nil
	backendRun        func(root *Node, cfg windowConfig)                 // 启动窗口与渲染/事件循环（阻塞）
)
