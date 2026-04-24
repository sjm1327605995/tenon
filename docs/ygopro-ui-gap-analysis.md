# YGOPro UI → Tenon v2 能力缺口分析

## 一、YGOPro 的 UI 架构概览

YGOPro 基于 **Irrlicht 3D 引擎**，UI 分为两层：

1. **3D 渲染层** (`drawing.cpp`)：场地、卡牌、特效
2. **2D GUI 层** (`game.cpp` Irrlicht GUI)：菜单、对话框、按钮、文本

### 核心渲染特征

| 特征 | 实现方式 | Tenon 替代方案 |
|------|---------|---------------|
| 3D 卡牌 | `irr::core::matrix4` 变换，支持位置/旋转/透明度 | ebiten `GeoM` 仿射变换模拟 |
| 场地格子 | 预计算 `S3DVertex[4]` 四边形，3D 空间定位 | 2D 倾斜矩形 + 透视近似 |
| 卡牌翻转 | Y轴旋转矩阵，正反面切换 | `scaleX` 从 1→0→-1 模拟 |
| 选择框 | `GL_LINE_STIPPLE` 虚线动画 | 自定义虚线描边 + 动画 |
| 阴影文字 | 偏移绘制两次 | 自定义 ShadowText 组件 |
| 渐变填充 | `draw2DRectangle` 四色渐变 | ebiten 顶点色或着色器 |

---

## 二、组件缺口对照表

### Irrlicht GUI 组件 → Tenon v2 映射

| YGOPro 使用的组件 | 用途 | tenon 已有 | 状态 |
|------------------|------|-----------|------|
| `IGUIButton` | 各种按钮 | `Button`, `FloatButton` | ✅ |
| `IGUIStaticText` | 静态文字 | `Text` | ✅ |
| `IGUIEditBox` | 昵称、IP、聊天输入 | `TextInput` | ✅ |
| `IGUICheckBox` | 设置选项 | `Checkbox` | ✅ |
| `IGUIScrollBar` | 滚动条 | 内嵌于 `ScrollView` | ⚠️ 缺独立滚动条 |
| `IGUIListBox` | 主机列表、聊天记录、卡组列表 | **无** | ❌ 需 `ListView` |
| `IGUIComboBox` | 规则选择、卡组分类 | **无** | ❌ 需 `Dropdown` |
| `IGUITabControl` | 设置标签页 | **无** | ❌ 需 `Tabs` |
| `IGUIWindow` | 对话框容器 | `Window`, `Modal` | ✅ |
| `IGUIImage` | 卡牌图片 | `Image` | ✅ |
| `IGUIProgressBar` | LP 血条 | `ProgressBar` | ✅ |
| `IGUISlider` | 音量调节 | `Slider` | ✅ |
| `IGUIContextMenu` | 右键菜单 | `Menu` | ✅ |

**缺失的核心组件：ListView、Dropdown、Tabs、独立 ScrollBar**

---

## 三、关键功能缺口

### 3.1 🔴 高优先级（必须实现，否则无法渲染游戏核心）

#### 1. 元素变换系统（Transform）

YGOPro 的卡牌使用 3D 矩阵变换：
```cpp
pcard->mTransform.setTranslation(pcard->curPos);
pcard->mTransform.setRotationRadians(pcard->curRot);
driver->setTransform(ETS_WORLD, pcard->mTransform);
```

Tenon 当前 `Element` 只有 `LayoutBounds`（x, y, w, h），**没有旋转、缩放、倾斜**。

**需要增加：**
```go
type Transform struct {
    Rotation  float32 // 角度，度
    ScaleX    float32 // 默认 1
    ScaleY    float32 // 默认 1
    SkewX     float32 // X轴倾斜，度
    SkewY     float32 // Y轴倾斜，度
    OriginX   float32 // 变换中心 X（相对自身 0-1）
    OriginY   float32 // 变换中心 Y（相对自身 0-1）
    Alpha     float32 // 透明度 0-1
}
```

绘制时通过 ebiten `GeoM` 应用：
```go
// 近似 3D 翻转（Y轴旋转）
func simulateYRotation(angleDeg float32) (scaleX, skewY float32) {
    rad := angleDeg * math.Pi / 180
    scaleX = float32(math.Cos(float64(rad)))
    // 加上透视收缩效果
    return
}
```

**影响范围：** 卡牌翻转、场地倾斜、卡牌弹出动画、攻击指示器旋转。

#### 2. 卡牌组件（Card）

YGOPro 卡牌不是简单图片，而是有状态的 3D 对象：
- 正反面（正面卡图 / 背面卡背）
- 位置状态（手牌、场地、墓地等）
- 覆盖标记（装备、目标、连锁、无效）
- 选择高亮（虚线边框）
- 攻击力/防御力/等级文字叠加

**需要实现：**
```go
type Card struct {
    core.BaseElement
    frontImage   *ebiten.Image  // 正面卡图
    backImage    *ebiten.Image  // 卡背
    isFaceUp     bool
    flipProgress float32        // 翻转动画 0-1
    
    // 状态标记
    isSelectable bool
    isSelected   bool
    isHighlighted bool
    cmdFlag      int            // 可执行命令标记
    
    // 覆盖图标
    equipIcon    bool
    targetIcon   bool
    negateIcon   bool
    attackIcon   bool
}
```

#### 3. 阴影文字（ShadowText）

YGOPro 大量使用：
```cpp
DrawShadowText(lpcFont, text, position, padding, color, shadowcolor, hcenter, vcenter, clip);
```

实现简单但关键（LP 数字、玩家名、卡牌数量等都依赖）：
```go
type ShadowText struct {
    core.BaseElement
    content    string
    fontSize   float64
    color      color.Color
    shadowColor color.Color
    shadowOffsetX, shadowOffsetY float32
}
// Draw: 先画阴影偏移文本，再画正文
```

#### 4. 列表组件（ListView）

YGOPro 大量使用列表：
- `lstHostList` - 主机列表
- `lstReplayList` - 录像列表
- `lstChatLog` - 聊天记录
- `lstBotList` - 机器人列表
- `lstSinglePlayList` - 单人模式列表
- `lstCategories` / `lstDecks` - 卡组管理

需要支持：
- 垂直滚动列表
- 项选择（单选/多选）
- 项模板（文字项、图文项）
- 虚拟滚动（大量项时优化）

#### 5. 下拉选择（Dropdown / Select）

YGOPro 中的下拉框：
- `cbHostLFlist` - 禁卡表选择
- `cbRule` - 规则选择
- `cbMatchMode` - 对战模式
- `cbDuelRule` - 决斗规则
- `cbDBCategory` / `cbDBDecks` - 卡组分类

---

### 3.2 🟡 中优先级（影响游戏体验）

#### 6. 标签页（Tabs）

YGOPro 设置页面：
- `tabHelper` - 辅助设置
- `tabSystem` - 系统设置

需要 Tab + TabPanel 组合组件。

#### 7. 渐变支持

YGOPro LP 条：
```cpp
driver->draw2DRectangle(0xa0000000, Resize(327, 8, 630, 51));
// 以及多色渐变 LP 条
driver->draw2DImage(imageManager.tLPBar, ...);
```

需要：
- `LinearGradient` 线性渐变填充
- 可指定起点/终点颜色
- 方向：水平/垂直

#### 8. 虚线边框/选择框动画

YGOPro 的选择框有流动的虚线效果：
```cpp
if(linePattern < 15) { /* 动画第一段 */ }
else { /* 动画第二段 */ }
linePattern = (linePattern + 1) % 30;
```

实现方案：
- `drawDashedRect(screen, bounds, dashLen, gapLen, color, phase)`
- phase 每帧递增实现流动效果

#### 9. 网格布局/网格组件

卡组编辑界面需要卡牌网格：
- 主卡组（40-60 张）
- 额外卡组（0-15 张）
- 副卡组（0-15 张）

Yoga 的 Flex 布局可以部分模拟，但固定列数的网格用 `Grid` 更方便。

#### 10. 拖拽支持（Drag & Drop）

卡组编辑时拖拽卡牌排序。
游戏内拖拽卡牌到场地。

需要在 Engine 层面支持：
- `EventDragStart`, `EventDragMove`, `EventDragEnd`
- 拖拽时的半透明预览
- 放置目标高亮

---

### 3.3 🟢 低优先级（锦上添花）

#### 11. 圆角图片/剪切蒙版

卡牌图片目前用矩形显示，实际 YGOPro 卡牌有微小圆角。
可以用 ebiten `DrawTriangles` + alpha 蒙版实现。

#### 12. Toast / Snackbar 通知

YGOPro 的自动关闭消息：
```cpp
PopupElement(element, hideframe); // 自动在 hideframe 后关闭
```

#### 13. 分隔线（Divider）

已有 `Divider`，但 YGOPro 中更多使用。

#### 14. 粒子/特效系统

YGOPro 的召唤特效、连锁标记旋转动画等。
可以用简单的 Sprite 动画 + Tween 模拟。

---

## 四、"界面倾斜实现 3D" 的技术方案

### 4.1 方案：2D 仿射变换模拟 3D（推荐）

ebiten 的 `GeoM` 是 3×3 2D 仿射变换矩阵，支持：
- `Translate` 平移
- `Rotate` 旋转
- `Scale` 缩放
- `Concat` 矩阵组合

**卡牌翻转（模拟 Y 轴旋转）：**
```go
func drawCardFlip(screen *ebiten.Image, img *ebiten.Image, bounds LayoutBounds, flipAngle float32) {
    // flipAngle: 0=正面, 90=侧立, 180=背面
    rad := float64(flipAngle) * math.Pi / 180
    scaleX := float32(math.Cos(rad))
    
    op := &ebiten.DrawImageOptions{}
    cx := float64(bounds.X + bounds.Width/2)
    cy := float64(bounds.Y + bounds.Height/2)
    
    op.GeoM.Translate(-cx, -cy)      // 移到中心
    op.GeoM.Scale(float64(scaleX), 1) // X轴收缩模拟旋转
    op.GeoM.Translate(cx, cy)         // 移回
    
    if math.Abs(scaleX) > 0 {
        // 选择正面或背面
        if flipAngle < 90 || flipAngle > 270 {
            screen.DrawImage(frontImg, op)
        } else {
            screen.DrawImage(backImg, op)
        }
    }
}
```

**场地倾斜（模拟透视）：**
```go
// 底部向上收缩模拟远处
skewY := -0.3  // 向上倾斜
scaleY := 0.8  // 远处变小

op.GeoM.Skew(0, float64(skewY))
op.GeoM.Scale(1, float64(scaleY))
```

**局限性：**
- 不能真正处理 Z 轴深度排序（远处卡牌应在近处后面）
- 不能实现复杂透视（如卡牌站立时的顶部/底部不同收缩）
- 对于游戏王这种俯视角+卡牌平铺的场景，**足够了**

### 4.2 更优方案：DrawTriangles 透视贴图

如果需要更好的 3D 效果，可以用 ebiten 的 `DrawTriangles`：

```go
// 定义四边形的 4 个顶点（可任意变形）
vertices := []ebiten.Vertex{
    {DstX: x1, DstY: y1, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
    {DstX: x2, DstY: y2, SrcX: w, SrcY: 0, ...},
    {DstX: x3, DstY: y3, SrcX: 0, SrcY: h, ...},
    {DstX: x4, DstY: y4, SrcX: w, SrcY: h, ...},
}
indices := []uint16{0, 1, 2, 2, 1, 3}
screen.DrawTriangles(vertices, indices, img, nil)
```

这允许任意四边形变形，可以实现真正的透视效果（类似 YGOPro 的 `Draw2DImageQuad`）。

**推荐实现优先级：**
1. 先用 `GeoM` 仿射变换实现基础效果（开发快，够用）
2. 后期对卡牌组件升级到 `DrawTriangles` 实现完美透视

---

## 五、实施路线图

### Phase 1: 核心基础设施（阻塞其他所有工作）

1. **Transform 系统**
   - `BaseElement` 增加 `transform Transform`
   - `drawElement` 绘制前应用 `GeoM` 变换
   - `hitTest` 需要考虑变换后的命中区域

2. **ShadowText 组件**
   - 简单但大量使用，优先实现

3. **Card 组件**
   - 正面/背面切换
   - 翻转动画（Tween + Transform）
   - 状态图标叠加

### Phase 2: 通用 UI 组件

4. **ListView**
5. **Dropdown**
6. **Tabs**
7. **ScrollBar**（独立可视化）

### Phase 3: 游戏专用效果

8. **渐变填充**（View 增加 gradient 支持）
9. **虚线边框动画**
10. **拖拽系统**
11. **场地网格组件**（基于 Card 的网格布局）

### Phase 4:  polish

12. 圆角图片剪切
13. Toast 通知
14. 粒子特效系统

---

## 六、YGOPro 各界面的 Tenon 实现难度评估

| 界面 | 难度 | 关键依赖 |
|------|------|---------|
| 主菜单 | ⭐ 低 | Button, Window, Text |
| LAN 对战 | ⭐⭐ 中 | TextInput, ListView, Button |
| 创建主机 | ⭐⭐ 中 | Dropdown, TextInput, Checkbox |
| 卡组编辑 | ⭐⭐⭐ 高 | Grid, DragDrop, Image, ScrollView |
| **决斗场地** | ⭐⭐⭐⭐ 很高 | Card, Transform, 场地布局, 动画 |
| 卡牌选择 | ⭐⭐ 中 | Grid, Button(Image), Modal |
| 设置页面 | ⭐⭐ 中 | Tabs, Checkbox, Slider, Dropdown |
| 聊天框 | ⭐⭐ 中 | ListView, TextInput |

---

## 七、结论

Tenon v2 的架构（Widget→Element→Yoga 布局→脏标记刷新）**完全适合** YGOPro 的 UI 需求。主要缺口是：

1. **变换系统** — 必须实现，否则无法做卡牌动画和场地倾斜
2. **Card 组件** — 游戏核心
3. **ListView / Dropdown** — 通用 UI，大量界面依赖
4. **ShadowText / 渐变 / 虚线** — 视觉效果，提升还原度

建议先实现 **Transform + Card + ShadowText**，然后做一个简化版决斗场地 demo 验证方案可行性。
